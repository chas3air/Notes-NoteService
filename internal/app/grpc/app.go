package grpcapp

import (
	"fmt"
	"net"
	"net/http"
	grpcmiddleware "notesservice/internal/app/grpc/middleware"
	"notesservice/internal/handlers"
	"notesservice/internal/handlers/grpc/notes"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
)

type App struct {
	log         *zap.Logger
	gRPCServer  *grpc.Server
	port        int
	metricsPort int
}

func New(log *zap.Logger, service handlers.Service, port int, metricsPort int) *App {
	kaParams := keepalive.ServerParameters{
		MaxConnectionIdle:     5 * time.Minute,  // разорвать idle соединение
		MaxConnectionAge:      30 * time.Minute, // максимальный возраст соединения
		MaxConnectionAgeGrace: 5 * time.Minute,  // дать время завершить запросы
		Time:                  2 * time.Minute,  // как часто сервер шлёт ping
		Timeout:               20 * time.Second, // сколько ждать ответа на ping
	}

	kaPolicy := keepalive.EnforcementPolicy{
		MinTime:             1 * time.Minute, // минимальный интервал между ping от клиента
		PermitWithoutStream: true,            // разрешить ping без активных RPC
	}

	gRPCServer := grpc.NewServer(
		grpc.KeepaliveParams(kaParams),
		grpc.KeepaliveEnforcementPolicy(kaPolicy),
		grpc.UnaryInterceptor(grpcmiddleware.UnaryMetricsInterceptor(log)),
	)

	notes.Register(gRPCServer, log, service)
	reflection.Register(gRPCServer)

	return &App{
		log:         log,
		gRPCServer:  gRPCServer,
		port:        port,
		metricsPort: metricsPort,
	}
}

func (a *App) Run() error {
	const op = "grpc.App.Run"
	a.log.Info("starting gRPC server", zap.String("op", op), zap.Int("port", a.port))

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		a.log.Error("failed to listen", zap.String("op", op), zap.Error(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	a.startMetricsServer()
	a.log.Info("starting gRPC metrics", zap.String("op", op), zap.Int("port", a.metricsPort))

	if err := a.gRPCServer.Serve(l); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (a *App) startMetricsServer() {
	http.Handle("/metrics", promhttp.Handler())
	go func() {
		if err := http.ListenAndServe(fmt.Sprintf(":%d", a.metricsPort), nil); err != nil {
			a.log.Error("failed to start metrics server", zap.Error(err))
		}
	}()
}

func (a *App) Stop() {
	const op = "grpcapp.Stop"

	a.log.With(zap.String("op", op)).Info("stopping gRPC server")
	a.gRPCServer.GracefulStop()
}
