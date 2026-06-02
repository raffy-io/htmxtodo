package main

import (
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"

	"github.com/raffy-io/htmxtodo"

	"github.com/raffy-io/htmxtodo/internal/config"
	"github.com/raffy-io/htmxtodo/internal/connection"
	"github.com/raffy-io/htmxtodo/internal/db"
	"github.com/raffy-io/htmxtodo/internal/handlers"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
)


func main() {
	
	// ENV
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to initialize configuration: %v", err)
	}

	// initialize the database pool
	conn, err := connection.Connect(cfg.DBURL, cfg.AuthToken)
	if err != nil {
		log.Fatalf("Database initialization failed: %v", err)
	}
	defer conn.Close()

	fmt.Println("Successfully connected to Turso!")

	queries := db.New(conn)

	// handlers
	taskHandler := handlers.NewHandler(queries)
	
	//routes
	mux := http.NewServeMux()

	mux.HandleFunc("GET /" ,taskHandler.GetTasks)
	mux.HandleFunc("POST /addtask", taskHandler.AddTask)
	mux.HandleFunc("DELETE /deletetask/{id}", taskHandler.DeleteTask)


	// static assets
	staticFS,err := fs.Sub(htmxtodo.EmbeddedAssets,"static")
	if err != nil {
		log.Fatal(err)
	}
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.FS(staticFS))))

	// server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server is running on port %s...", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}