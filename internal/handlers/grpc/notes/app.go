package notes

import (
	"context"
	"errors"
	"notesservice/internal/handlers"
	"notesservice/internal/models/profilers"
	"notesservice/internal/service"

	notesv1 "github.com/chas3air/protos/gen/go/notes"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

var (
	defaultLimit  = 10
	defaultOffset = 0
)

type serverAPI struct {
	notesv1.UnimplementedNotesServer

	log     *zap.Logger
	service handlers.Service
}

func Register(svc *grpc.Server, log *zap.Logger, service handlers.Service) {
	notesv1.RegisterNotesServer(svc, &serverAPI{
		log:     log,
		service: service,
	})
}

func (s *serverAPI) GetNotes(ctx context.Context, req *notesv1.GetNotesRequest) (*notesv1.GetNotesResponse, error) {
	const op = "handlers.grpc.GetNotes"
	log := s.log.With(zap.String("op", op))

	offset := req.GetOffset()
	limit := req.GetLimit()

	if limit <= 0 || offset < 0 {
		log.Error("invalid pagination parameters", zap.Int32("limit", limit), zap.Int32("offset", offset))
		return nil, status.Error(codes.InvalidArgument, "invalid pagination parameters")
	}

	notes, err := s.service.GetNotes(ctx, int(offset), int(limit))
	if err != nil {
		log.Error("failed to get notes", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to get notes")
	}

	protoList := profilers.NotesToProtoList(notes)

	return &notesv1.GetNotesResponse{
		Notes: protoList,
	}, nil
}

func (s *serverAPI) GetNotesByUser(ctx context.Context, req *notesv1.GetNotesByUserRequest) (*notesv1.GetNotesByUserResponse, error) {
	const op = "handlers.grpc.GetNotesByUser"
	log := s.log.With(zap.String("op", op))

	userIds := req.GetUserId()
	userId, err := uuid.Parse(userIds)
	if err != nil {
		log.Error("failed to parse user ID", zap.String("user_id", userIds), zap.Error(err))
		return nil, status.Error(codes.InvalidArgument, "invalid user ID format")
	}

	notes, err := s.service.GetNotesByUser(ctx, userId, defaultOffset, defaultLimit)
	if err != nil {
		log.Error("failed to get notes by user", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to get notes by user")
	}

	protoList := profilers.NotesToProtoList(notes)

	return &notesv1.GetNotesByUserResponse{
		Notes: protoList,
	}, nil
}

func (s *serverAPI) GetNotesById(ctx context.Context, req *notesv1.GetNotesByIdRequest) (*notesv1.GetNotesByIdResponse, error) {
	const op = "handlers.grpc.GetNotesById"
	log := s.log.With(zap.String("op", op))

	noteIds := req.GetId()
	noteId, err := uuid.Parse(noteIds)
	if err != nil {
		log.Error("failed to parse note ID", zap.String("note_id", noteIds), zap.Error(err))
		return nil, status.Error(codes.InvalidArgument, "invalid note ID format")
	}

	note, err := s.service.GetNoteById(ctx, noteId)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			log.Error("note not found", zap.String("note_id", noteIds))
			return nil, status.Error(codes.NotFound, "note not found")
		}

		log.Error("failed to get note by ID", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to get note by ID")
	}

	protoNote := profilers.NoteToProto(note)

	return &notesv1.GetNotesByIdResponse{
		Note: &protoNote,
	}, nil
}

func (s *serverAPI) Insert(ctx context.Context, req *notesv1.InsertRequest) (*emptypb.Empty, error) {
	const op = "handlers.grpc.Insert"
	log := s.log.With(zap.String("op", op))

	notepb := req.GetNote()
	if notepb == nil {
		log.Error("note data is nil")
		return nil, status.Error(codes.InvalidArgument, "note data is required")
	}

	note := profilers.ProtoToNote(notepb)

	if err := s.service.Insert(ctx, note); err != nil {
		if errors.Is(err, service.ErrAlreadyExists) {
			log.Error("note already exists", zap.String("id", note.Id.String()))
			return nil, status.Error(codes.AlreadyExists, "note already exists")
		}

		log.Error("failed to insert note", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to insert note")
	}

	return &emptypb.Empty{}, nil
}

func (s *serverAPI) Update(ctx context.Context, req *notesv1.UpdateRequest) (*emptypb.Empty, error) {
	const op = "handlers.grpc.Update"
	log := s.log.With(zap.String("op", op))

	notepb := req.GetNote()
	if notepb == nil {
		log.Error("note data is nil")
		return nil, status.Error(codes.InvalidArgument, "note data is required")
	}

	note := profilers.ProtoToNote(notepb)

	if err := s.service.Update(ctx, note.Id, note); err != nil {
		if errors.Is(err, service.ErrNotFound) {
			log.Error("note not found", zap.String("id", note.Id.String()))
			return nil, status.Error(codes.NotFound, "note not found")
		}

		log.Error("failed to update note", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to update note")
	}

	return &emptypb.Empty{}, nil
}

func (s *serverAPI) Delete(ctx context.Context, req *notesv1.DeleteRequest) (*emptypb.Empty, error) {
	const op = "handlers.grpc.Delete"
	log := s.log.With(zap.String("op", op))

	noteIdStr := req.GetId()
	noteId, err := uuid.Parse(noteIdStr)
	if err != nil {
		log.Error("failed to parse note ID", zap.String("note_id", noteIdStr), zap.Error(err))
		return nil, status.Error(codes.InvalidArgument, "invalid note ID format")
	}

	if err := s.service.Delete(ctx, noteId); err != nil {
		if errors.Is(err, service.ErrNotFound) {
			log.Error("note not found", zap.String("id", noteId.String()))
			return nil, status.Error(codes.NotFound, "note not found")
		}

		log.Error("failed to delete note", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to delete note")
	}

	return &emptypb.Empty{}, nil
}
