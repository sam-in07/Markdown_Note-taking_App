package routers

import (
	"note-app-go/internal/handlers"
	"note-app-go/internal/services"

	"github.com/gorilla/mux"
)

// NewRouter creates a new mux.Router and sets up the routes.
func NewRouter() *mux.Router {
	service := services.NewNoteService("./uploads")
	handler := handlers.NewNoteHandler(service)

	r := mux.NewRouter()

	// Frontend
	r.HandleFunc("/", handler.IndexHandler).Methods("GET")

	// API
	r.HandleFunc("/api/notes", handler.UploadHandler).Methods("POST")
	r.HandleFunc("/api/notes/{filename}", handler.DeleteHandler).Methods("DELETE")
	r.HandleFunc("/api/notes", handler.ListNotesHandler).Methods("GET")
	r.HandleFunc("/api/notes/{filename}", handler.RenderToHtmlHandler).Methods("GET")
	r.HandleFunc("/api/notes/check/{filename}", handler.CheckGrammarHandler).Methods("GET")

	return r
}
