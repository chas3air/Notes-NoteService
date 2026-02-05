package restapp

import (
	"context"
	"net/http"
	"notesservice/internal/app/rest/middleware"
	"notesservice/internal/handlers"
	"notesservice/internal/handlers/rest/notes"
	"strconv"

	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/gorilla/mux"
	"go.uber.org/zap"

	_ "notesservice/docs"
)

type App struct {
	log     *zap.Logger
	service handlers.Service
	port    int
	server  *http.Server
}

func New(log *zap.Logger, service handlers.Service, port int) *App {
	return &App{
		log:     log,
		service: service,
		port:    port,
	}
}

func (a *App) Run() error {
	const op = "restapp.Run"
	log := a.log.With(zap.String("op", op))

	notesHandler := notes.NewHandler(log, a.service)

	base := mux.NewRouter()

	base.Use(middleware.CORS)
	router := base.PathPrefix("/api/v1").Subrouter()

	router.HandleFunc("/notes", notesHandler.GetNotes).Methods(http.MethodGet)
	router.HandleFunc("/notes/by-user", notesHandler.GetNotesByUser).Methods(http.MethodGet)
	router.HandleFunc("/notes/detail", notesHandler.GetNoteById).Methods(http.MethodGet)
	router.HandleFunc("/notes", notesHandler.Insert).Methods(http.MethodPost)
	router.HandleFunc("/notes", notesHandler.Update).Methods(http.MethodPut)
	router.HandleFunc("/notes", notesHandler.Delete).Methods(http.MethodDelete)

	router.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("http://localhost:" + strconv.Itoa(a.port) + "/api/v1/swagger/doc.json"),
	))

	log.Info("router was assembled")

	a.server = &http.Server{
		Addr:    ":" + strconv.Itoa(a.port),
		Handler: router,
	}

	if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Error("failed to start server", zap.Error(err))
		return err
	}

	return nil
}

func (a *App) Shutdown() error {
	return a.server.Shutdown(context.Background())
}
