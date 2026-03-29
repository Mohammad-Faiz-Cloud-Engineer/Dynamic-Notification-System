// Dynamic Notification System
// Author: Mohammad Faiz
// Copyright (c) 2024 Mohammad Faiz
// Licensed under the MIT License

package main

import (
	"dynamic-notification-system/config"
	"dynamic-notification-system/notifier"
	"dynamic-notification-system/plugins"
	"dynamic-notification-system/scheduler"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

const (
	serverPort         = ":8080"
	readTimeout        = 15 * time.Second
	writeTimeout       = 15 * time.Second
	idleTimeout        = 60 * time.Second
	maxHeaderBytes     = 1 << 20 // 1 MB
	shutdownTimeout    = 30 * time.Second
)

func main() {
	var err error

	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// Load plugins based on configuration
	notifiers, err := plugins.LoadPlugins(cfg.Channels)
	if err != nil {
		log.Fatalf("Error loading plugins: %v", err)
	}

	// Pass the loaded notifiers to the notifier package
	notifier.SetNotifiers(notifiers)

	// Initialize Scheduler if enabled
	if cfg.Scheduler {
		err = scheduler.Initialize(cfg, notifiers)
		if err != nil {
			log.Fatalf("Error initializing scheduler: %v", err)
		}
		defer scheduler.Shutdown()
	}

	r := mux.NewRouter()
	if cfg.Scheduler {
		// Scheduled notification endpoints
		r.HandleFunc("/schema/job", scheduler.GetJobSchema())
		r.HandleFunc("/jobs", scheduler.HandlePostJob).Methods("POST")
		r.HandleFunc("/jobs", scheduler.HandleGetJobs).Methods("GET")
	}
	// Instant notification endpoint
	r.HandleFunc("/notify", notifier.HandlePostJob).Methods("POST")

	// Configure HTTP server with timeouts
	srv := &http.Server{
		Addr:           serverPort,
		Handler:        r,
		ReadTimeout:    readTimeout,
		WriteTimeout:   writeTimeout,
		IdleTimeout:    idleTimeout,
		MaxHeaderBytes: maxHeaderBytes,
	}

	log.Printf("Server listening on port %s", serverPort)
	log.Fatal(srv.ListenAndServe())
}
