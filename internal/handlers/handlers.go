package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/a-h/templ"
	"github.com/raffy-io/htmxtodo/internal/db"
	"github.com/raffy-io/htmxtodo/ui/components"
	"github.com/raffy-io/htmxtodo/ui/layout"
)

type TasksHandler struct {
	Queries *db.Queries
}

func (h *TasksHandler) GetTasks(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	data ,err := h.Queries.ListTasks(ctx)
	if err != nil {
		log.Printf("Data not found: %v\n", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	component := components.Tasks(data)
	pageLayout := layout.Base("Welcome",component)
	templ.Handler(pageLayout).ServeHTTP(w,r)

}

func (h *TasksHandler) AddTask(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		err := r.ParseForm()
		if err != nil{
			http.Error(w,"Bad Request",http.StatusBadRequest)
			return
		}
		data := r.FormValue("task")
		if data == "" {
			http.Error(w, "Task cannot be empty", http.StatusUnprocessableEntity)
        	return
		}

	    err = h.Queries.CreateTask(ctx,data)
		if err != nil {
        log.Printf("Failed to add task: %v", err)
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    	}

		newList,err := h.Queries.ListTasks(ctx)
		component := components.TasksList(newList)
		templ.Handler(component).ServeHTTP(w,r)
		
	
}

func (h *TasksHandler) DeleteTask(w http.ResponseWriter, r *http.Request){
	ctx := r.Context()

	err := r.ParseForm()
	if err != nil{
			http.Error(w,"Bad Request",http.StatusBadRequest)
			return
	}

	idValue := r.PathValue("id")
    if idValue == "" {
        http.Error(w, "Missing item ID", http.StatusBadRequest)
        return
    }

	id, err := strconv.Atoi(idValue)
    if err != nil {
        log.Printf("Invalid ID format: %v\n", err)
        http.Error(w, "Bad request", http.StatusBadRequest)
        return
    }


    err = h.Queries.DeleteTask(ctx, int64(id)) // Cast to int64 if required by your DB schema
    if err != nil {
        log.Printf("Failed to delete todo: %v\n", err)
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
    }

	newList,err := h.Queries.ListTasks(ctx)
	component := components.TasksList(newList)
	templ.Handler(component).ServeHTTP(w,r)

}