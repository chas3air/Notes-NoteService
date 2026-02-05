package notes

import (
	"encoding/json"
	"errors"
	"net/http"
	"notesservice/internal/handlers"
	"notesservice/internal/handlers/rest/utils"
	noteModels "notesservice/internal/models/domain"
	swaggerModels "notesservice/internal/models/swagger"
	"notesservice/internal/service"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Handler struct {
	log     *zap.Logger
	service handlers.Service
}

func NewHandler(log *zap.Logger, service handlers.Service) *Handler {
	return &Handler{
		log:     log,
		service: service,
	}
}

// GetNotes godoc
// @Summary      Get all notes
// @Description  Get a list of all notes with pagination
// @Tags         notes
// @Produce      json
// @Param        limit   query     int  false  "Limit (default 10)"
// @Param        offset  query     int  false  "Offset (default 0)"
// @Success      200  {array}   noteModels.Note
// @Failure      500  {string}  string "failed to get notes"
// @Router       /notes [get]
func (h *Handler) GetNotes(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.notes.GetNotes"
	log := h.log.With(zap.String("op", op))

	paginationData := utils.GetPagination(r)

	notes, err := h.service.GetNotes(r.Context(), paginationData.Offset, paginationData.Limit)
	if err != nil {
		log.Error("failed to get notes", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//
	var _ noteModels.Note

	if err := json.NewEncoder(w).Encode(notes); err != nil {
		log.Error("failed to encode notes response", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// GetNotesByUser godoc
// @Summary      Get notes by user
// @Description  Get a list of notes for a specific user via query parameter
// @Tags         notes
// @Produce      json
// @Param        id      query     string  true   "User UUID" format(uuid)
// @Param        limit   query     int     false  "Limit"
// @Param        offset  query     int     false  "Offset"
// @Success      200  {array}   noteModels.Note
// @Failure      400  {string}  string "missing or invalid id"
// @Failure      500  {string}  string "failed to get notes by user"
// @Router       /notes/by-user [get]
func (h *Handler) GetNotesByUser(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.notes.GetNotesByUser"
	log := h.log.With(zap.String("op", op))

	sUserId := r.URL.Query().Get("id")
	if sUserId == "" {
		log.Error("missing id parameter")
		http.Error(w, "missing id parameter", http.StatusBadRequest)
		return
	}

	userId, err := uuid.Parse(sUserId)
	if err != nil {
		log.Error("invalid id format", zap.Error(err))
		http.Error(w, "invalid id format", http.StatusBadRequest)
		return
	}

	paginationData := utils.GetPagination(r)

	notesByUser, err := h.service.GetNotesByUser(r.Context(), userId, paginationData.Offset, paginationData.Limit)
	if err != nil {
		log.Error("failed to get notes by user", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(notesByUser); err != nil {
		log.Error("failed to encode notes by user response", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// GetNoteById godoc
// @Summary      Get note details
// @Description  Get details of a single note by its UUID in query
// @Tags         notes
// @Produce      json
// @Param        id   query     string  true  "Note UUID" format(uuid)
// @Success      200  {object}  noteModels.Note
// @Failure      400  {string}  string "invalid id"
// @Failure      404  {string}  string "note not found"
// @Failure      500  {string}  string "failed to get note by id"
// @Router       /notes/detail [get]
func (h *Handler) GetNoteById(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.notes.GetNoteById"
	log := h.log.With(zap.String("op", op))

	sId := r.URL.Query().Get("id")
	if sId == "" {
		log.Error("missing id parameter")
		http.Error(w, "missing id parameter", http.StatusBadRequest)
		return
	}

	noteId, err := uuid.Parse(sId)
	if err != nil {
		log.Error("invalid id parameter", zap.Error(err))
		http.Error(w, "invalid id parameter", http.StatusBadRequest)
		return
	}

	note, err := h.service.GetNoteById(r.Context(), noteId)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			log.Warn("note not found", zap.String("id", noteId.String()))
			http.Error(w, "note not found", http.StatusNotFound)
			return
		}

		log.Error("failed to get note by id", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(note); err != nil {
		log.Error("failed to encode note response", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Insert godoc
// @Summary      Create note
// @Description  Insert a new note record
// @Tags         notes
// @Accept       json
// @Param        request body   swaggerModels.NoteRequest true "Inserted note data"
// @Success      201  {string}  string "Created"
// @Failure      400  {string}  string "invalid body"
// @Failure      409  {string}  string "note already exists"
// @Failure      500  {string}  string "failed to insert note"
// @Router       /notes [post]
func (h *Handler) Insert(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.notes.Insert"
	log := h.log.With(zap.String("op", op))

	request := swaggerModels.NoteRequest{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Error("failed to decode request body", zap.Error(err))
		http.Error(w, "failed to decode request body", http.StatusBadRequest)
		return
	}

	if err := h.service.Insert(r.Context(), request.Note); err != nil {
		if errors.Is(err, service.ErrAlreadyExists) {
			log.Warn("note already exists", zap.String("id", request.Note.Id.String()))
			http.Error(w, "note already exists", http.StatusConflict)
			return
		}
		log.Error("failed to insert note", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// Update godoc
// @Summary      Update note
// @Description  Update an existing note record
// @Tags         notes
// @Accept       json
// @Param        request body      swaggerModels.NoteRequest  true  "Updated note data"
// @Success      204  {string}  string "No Content"
// @Failure      404  {string}  string "not found"
// @Failure      500  {string}  string "failed to update note"
// @Router       /notes [put]
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.notes.Update"
	log := h.log.With(zap.String("op", op))

	request := swaggerModels.NoteRequest{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Error("failed to decode request body", zap.Error(err))
		http.Error(w, "failed to decode request body", http.StatusBadRequest)
		return
	}

	if err := h.service.Update(r.Context(), request.Note.Id, request.Note); err != nil {
		if errors.Is(err, service.ErrNotFound) {
			log.Warn("note not found for update", zap.String("id", request.Note.Id.String()))
			http.Error(w, "note not found for update", http.StatusNotFound)
			return
		}
		log.Error("failed to update note", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Delete godoc
// @Summary      Delete note
// @Description  Remove a note record by ID in body
// @Tags         notes
// @Accept       json
// @Param        request body      swaggerModels.DeleteRequest  true  "ID container: {'id': 'uuid'}"
// @Success      204  {string}  string "No Content"
// @Failure      404  {string}  string "not found"
// @Failure      500  {string}  string "failed to delete note"
// @Router       /notes [delete]
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.notes.Delete"
	log := h.log.With(zap.String("op", op))

	var request swaggerModels.DeleteRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Error("failed to decode request body", zap.Error(err))
		http.Error(w, "failed to decode request body", http.StatusBadRequest)
		return
	}

	noteId, err := uuid.Parse(request.Id)
	if err != nil {
		log.Error("invalid id parameter", zap.Error(err))
		http.Error(w, "invalid id parameter", http.StatusBadRequest)
		return
	}

	if err := h.service.Delete(r.Context(), noteId); err != nil {
		if errors.Is(err, service.ErrNotFound) {
			log.Warn("note not found for delete", zap.String("id", noteId.String()))
			http.Error(w, "note not found for delete", http.StatusNotFound)
			return
		}
		log.Error("failed to delete note", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
