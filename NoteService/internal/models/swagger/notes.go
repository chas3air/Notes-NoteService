package swagger

import "notesservice/internal/models/domain"

type NoteRequest struct {
	Note domain.Note `json:"note"`
}

type DeleteRequest struct {
	Id string `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
}
