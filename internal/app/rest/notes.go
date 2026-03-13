package restapp

import (
	"context"
	"net/http"
	restmiddleware "notesservice/internal/app/rest/middleware"
	"notesservice/internal/handlers"
	"notesservice/internal/handlers/rest/notes"
	"strconv"

	"github.com/prometheus/client_golang/prometheus/promhttp"
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
	const op = "rest.App.Run"
	log := a.log.With(zap.String("op", op))

	notesHandler := notes.NewHandler(log, a.service)

	base := mux.NewRouter()
	base.Use(restmiddleware.CORS)
	base.Use(restmiddleware.RequestLoggerMiddleware(log))
	base.Handle("/metrics", promhttp.Handler())

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

	log.Info("starting rest server", zap.String("op", op), zap.Int("port", a.port))

	a.server = &http.Server{
		Addr:    ":" + strconv.Itoa(a.port),
		Handler: base,
	}

	if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Error("failed to start server", zap.String("op", op), zap.Error(err))
		return err
	}

	return nil
}

func (a *App) Shutdown() error {
	return a.server.Shutdown(context.Background())
}
