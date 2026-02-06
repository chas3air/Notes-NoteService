package profilers

import (
	"notesservice/internal/models/domain"

	notesv1 "github.com/chas3air/protos/gen/go/notes"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ProtoToNote(obj *notesv1.Note) domain.Note {
	return domain.Note{
		Id:        uuid.MustParse(obj.Id),
		UserId:    uuid.MustParse(obj.UserId),
		Title:     obj.Title,
		Content:   obj.Content,
		CreatedAt: obj.CreatedAt.AsTime(),
		IsPrivate: obj.IsPrivate,
	}
}

func NoteToProto(obj domain.Note) notesv1.Note {
	return notesv1.Note{
		Id:        obj.Id.String(),
		UserId:    obj.UserId.String(),
		Title:     obj.Title,
		Content:   obj.Content,
		CreatedAt: timestamppb.New(obj.CreatedAt),
		IsPrivate: obj.IsPrivate,
	}
}
