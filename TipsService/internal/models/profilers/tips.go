package profilers

import (
	"tipsservice/internal/models/domain"

	tipsv1 "github.com/chas3air/protos/gen/go/tips"
	"github.com/google/uuid"
)

func ProtoToTip(obj *tipsv1.Tip) domain.Tip {
	return domain.Tip{
		Id:        uuid.MustParse(obj.Id),
		UserId:    uuid.MustParse(obj.UserId),
		Title:     obj.Title,
		Content:   obj.Content,
		CreatedAt: obj.CreatedAt.AsTime(), // TODO: fix it later
		IsPrivate: obj.IsPrivate,
	}
}

func TipToProto(obj domain.Tip) tipsv1.Tip {
	return tipsv1.Tip{
		Id:        obj.Id.String(),
		UserId:    obj.UserId.String(),
		Title:     obj.Title,
		Content:   obj.Content,
		CreatedAt: nil, // TODO: fix it later
		IsPrivate: obj.IsPrivate,
	}
}
