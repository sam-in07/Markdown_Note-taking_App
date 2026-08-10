package handlers

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strings"

	"note-app-go/internal/services"
	"note-app-go/internal/utils"
)

// NoteHandler handles HTTP requests related to notes.
type NoteHandler struct {
	noteService *services.NoteService
}

// NewNoteHandler creates a new NoteHandler with the given NoteService.
func NewNoteHandler(service *services.NoteService) *NoteHandler {
	return &NoteHandler{
		noteService: service,
	}
}

// IndexHandler renders the main page.
func (h *NoteHandler) IndexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	tmpl, err := template.ParseFiles("templates/index.tmpl")
	if err != nil {
		log.Printf("Failed to parse index template: %s", err)
		http.Error(w, "Failed to load page", http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"PageTitle": "Markdown Note App",
		"Content":   template.HTML("<h1>Welcome to Markdown Note App</h1><p>Your notes will appear here.</p>"),
	}

	if err := tmpl.Execute(w, data); err != nil {
		log.Printf("Failed to render index template: %s", err)
		http.Error(w, "Failed to render page", http.StatusInternalServerError)
		return
	}
}

// UploadHandler handles the uploading of Markdown files.

// UploadHandler handles the uploading of Markdown files.

// UploadHandler handles the uploading of Markdown files.
func (h *NoteHandler) UploadHandler(w http.ResponseWriter, r *http.Request) {
	_, fileHeader, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Invalid file", http.StatusBadRequest)
		return
	}

	// Make sure the uploaded file is a Markdown file.
	if !strings.HasSuffix(strings.ToLower(fileHeader.Filename), ".md") {
		http.Error(w, "File must be a markdown file", http.StatusBadRequest)
		return
	}

	// Upload the Markdown file.
	_, err = h.noteService.UploadNote(fileHeader, h.noteService.UploadDir)
	if err != nil {
		log.Printf("Failed to upload file: %s", err)
		http.Error(w, "Failed to upload file", http.StatusInternalServerError)
		return
	}

	log.Printf("%s uploaded successfully", fileHeader.Filename)

	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(map[string]string{
		"message": "File uploaded successfully",
	}); err != nil {
		log.Printf("Failed to encode upload response: %s", err)
	}
}




// DeleteHandler handles deleting a Markdown note.
func (h *NoteHandler) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	filename := strings.TrimPrefix(r.URL.Path, "/api/notes/")

	if filename == "" {
		http.Error(w, "Invalid filename", http.StatusBadRequest)
		return
	}

	err := h.noteService.DeleteNote(filename)
	if err != nil {
		log.Printf("Failed to delete file: %s", err)
		http.Error(w, "Failed to delete file", http.StatusInternalServerError)
		return
	}

	log.Printf("%s deleted successfully", filename)

	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(map[string]string{
		"message": "File deleted successfully",
	}); err != nil {
		log.Printf("Failed to encode delete response: %s", err)
	}
}

// ListNotesHandler returns a list of all Markdown notes.
func (h *NoteHandler) ListNotesHandler(w http.ResponseWriter, r *http.Request) {
	notes, err := h.noteService.GetNoteList()
	if err != nil {
		log.Printf("Failed to list notes: %s", err)
		http.Error(w, "Failed to list notes", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(notes); err != nil {
		log.Printf("Failed to encode notes: %s", err)
	}
}

// RenderToHtmlHandler reads a Markdown file,
// converts it to HTML, and renders it in index.tmpl.
func (h *NoteHandler) RenderToHtmlHandler(w http.ResponseWriter, r *http.Request) {
	filename := strings.TrimPrefix(r.URL.Path, "/api/notes/")

	if filename == "" {
		http.Error(w, "Invalid filename", http.StatusBadRequest)
		return
	}

	log.Printf("Requested filename: %s", filename)

	// Get Markdown content from the uploads directory.
	content, err := h.noteService.GetNoteContent(filename)
	if err != nil {
		log.Printf("Failed to get file content: %s", err)
		http.Error(w, "Failed to get file content", http.StatusInternalServerError)
		return
	}

	log.Printf("Markdown content length: %d", len(content))
	log.Printf("Markdown content: %s", content)

	// Convert Markdown to HTML.
	htmlContent := utils.MdToHtml(content)

	log.Printf("HTML content length: %d", len(htmlContent))
	log.Printf("HTML content: %s", htmlContent)

	// Load the HTML template.
	tmpl, err := template.ParseFiles("templates/index.tmpl")
	if err != nil {
		log.Printf("Failed to parse template file: %s", err)
		http.Error(w, "Failed to parse template file", http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"PageTitle": filename,
		"Content":   template.HTML(string(htmlContent)),
	}

	// Render the template.
	if err := tmpl.Execute(w, data); err != nil {
		log.Printf("Failed to execute template: %s", err)
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		return
	}
}

// CheckGrammarHandler checks the grammar of a Markdown note.
func (h *NoteHandler) CheckGrammarHandler(w http.ResponseWriter, r *http.Request) {
	filename := strings.TrimPrefix(r.URL.Path, "/api/notes/check/")

	if filename == "" {
		http.Error(w, "Invalid filename", http.StatusBadRequest)
		return
	}

	content, err := h.noteService.GetNoteContent(filename)
	if err != nil {
		log.Printf("Failed to get file content: %s", err)
		http.Error(w, "Failed to get file content", http.StatusInternalServerError)
		return
	}

	// Check the grammar.
	issues, err := h.noteService.CheckGrammar(string(content))
	if err != nil {
		log.Printf("Failed to check grammar: %s", err)
		http.Error(w, "Failed to check grammar", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(issues); err != nil {
		log.Printf("Failed to encode grammar issues: %s", err)
	}
}
