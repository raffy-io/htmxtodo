package main

import (
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"

	"github.com/raffy-io/htmxtodo"

	"github.com/raffy-io/htmxtodo/internal/connection"
	"github.com/raffy-io/htmxtodo/internal/db"
	"github.com/raffy-io/htmxtodo/internal/handlers"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
)


func main() {
	// Database Connection

	// Turso URLs look like: libsql://your-db-name-username.turso.io
	dbURL := os.Getenv("DB_URL")
	authToken := os.Getenv("AUTH_TOKEN")

	// 1. Initialize the pool
	conn, err := connection.Connect(dbURL, authToken)
	if err != nil {
		log.Fatalf("Database initialization failed: %v", err)
	}
	
	// 2. THIS is where you defer the close. It stays alive for the duration of main().
	defer conn.Close()

	fmt.Println("Successfully connected to Turso!")

	queries := db.New(conn)

	tasksHandler := &handlers.TasksHandler{
		Queries: queries,
	}


	//routes
	mux := http.NewServeMux()

	mux.HandleFunc("GET /" ,tasksHandler.GetTasks)
	mux.HandleFunc("POST /addtask", tasksHandler.AddTask)
	mux.HandleFunc("DELETE /deletetask/{id}", tasksHandler.DeleteTask)


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