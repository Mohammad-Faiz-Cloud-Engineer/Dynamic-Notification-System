package scheduler

import (
	"dynamic-notification-system/config"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/alecthomas/jsonschema"
	"github.com/robfig/cron/v3"
)

func GetJobSchema() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reflector := jsonschema.Reflector{}
		schema := reflector.Reflect(&config.ScheduledJob{})

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(w).Encode(schema); err != nil {
			http.Error(w, fmt.Sprintf("Error encoding schema: %v", err), http.StatusInternalServerError)
		}
	}
}

func loadJobs(c *cron.Cron) {
	dbJobs, err := loadJobsFromDB(db)
	if err != nil {
		log.Printf("Error loading jobs from DB: %v", err)
		return
	}
	for _, job := range dbJobs {
		addCronJob(c, job, notifiers)
	}
}

func addCronJob(c *cron.Cron, job config.ScheduledJob, notifiers []config.Notifier) {
	jobCopy := job
	_, err := c.AddFunc(job.ScheduleExpression, func() {
		for _, notifier := range notifiers {
			if notifier.Type() == jobCopy.NotificationType {
				err := notifier.Notify(&jobCopy.Message)
				if err != nil {
					log.Printf("Error sending notification via %s: %v", notifier.Name(), err)
				}
				_, err = db.Exec("UPDATE scheduled_jobs SET last_run = NOW() WHERE id = ?", jobCopy.ID)
				if err != nil {
					log.Printf("Error updating last_run: %v", err)
				}
			}
		}
	})
	if err != nil {
		log.Printf("Error adding cron job %s: %v", job.Name, err)
	}
}

func HandlePostJob(w http.ResponseWriter, r *http.Request) {
	var job config.ScheduledJob

	err := json.NewDecoder(r.Body).Decode(&job)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := validateJob(&job); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := db.Exec("INSERT INTO scheduled_jobs (name, notification_type, recipient, message, schedule_expression) VALUES (?, ?, ?, ?, ?)",
		job.Name, job.NotificationType, job.Recipient, job.Message, job.ScheduleExpression)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error inserting job: %v", err), http.StatusInternalServerError)
		return
	}

	id, _ := result.LastInsertId()
	job.ID = int(id)

	addCronJob(cronInstance, job, notifiers)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(job)
}

func HandleGetJobs(w http.ResponseWriter, r *http.Request) {
	var jobs []config.ScheduledJob

	rows, err := db.Query("SELECT id, name, notification_type, recipient, message, schedule_expression FROM scheduled_jobs")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var job config.ScheduledJob
		err := rows.Scan(&job.ID, &job.Name, &job.NotificationType, &job.Recipient, &job.Message, &job.ScheduleExpression)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		jobs = append(jobs, job)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(jobs)
}

func validateJob(job *config.ScheduledJob) error {
	if job.Name == "" {
		return fmt.Errorf("job name is required")
	}
	if len(job.Name) > 255 {
		return fmt.Errorf("job name must not exceed 255 characters")
	}
	if job.ScheduleExpression == "" {
		return fmt.Errorf("schedule expression is required")
	}
	if job.NotificationType == "" {
		return fmt.Errorf("notification type is required")
	}
	if job.Recipient == "" {
		return fmt.Errorf("recipient is required")
	}
	if len(job.Recipient) > 255 {
		return fmt.Errorf("recipient must not exceed 255 characters")
	}
	return nil
}
