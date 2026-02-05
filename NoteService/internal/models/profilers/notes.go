package profilers

import (
	"notesservice/internal/models/domain"

	tipsv1 "github.com/chas3air/protos/gen/go/tips" // TODO: change it later
	"github.com/google/uuid"
)

func ProtoToNote(obj *tipsv1.Tip) domain.Note {
	return domain.Note{
		Id:        uuid.MustParse(obj.Id),
		UserId:    uuid.MustParse(obj.UserId),
		Title:     obj.Title,
		Content:   obj.Content,
		CreatedAt: obj.CreatedAt.AsTime(), // TODO: fix it later
		IsPrivate: obj.IsPrivate,
	}
}

func NoteToProto(obj domain.Note) tipsv1.Tip {
	return tipsv1.Tip{
		Id:        obj.Id.String(),
		UserId:    obj.UserId.String(),
		Title:     obj.Title,
		Content:   obj.Content,
		CreatedAt: nil, // TODO: fix it later
		IsPrivate: obj.IsPrivate,
	}
}
