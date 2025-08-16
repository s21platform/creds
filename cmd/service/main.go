package main

import (
	"fmt"
	"github.com/s21platform/creds/internal/infra"
	"github.com/s21platform/creds/internal/repository"
	logger_lib "github.com/s21platform/logger-lib"
	"github.com/s21platform/metrics-lib/pkg"
	"log"
	"net"

	"github.com/s21platform/creds/internal/config"
	"github.com/s21platform/creds/internal/service"
	pb "github.com/s21platform/creds/pkg/creds"
	"google.golang.org/grpc"
)

func main() {
	cfg := config.MustLoad()
	fmt.Println("config", cfg)
	logger := logger_lib.New(cfg.Logger.Host, cfg.Logger.Port, cfg.Service.Name, cfg.Platform.Env)

	fmt.Println("logger: loaded")
	dbRepo := repository.New(cfg)

	fmt.Println("db: loaded")

	thisService := service.New(dbRepo)

	metrics, err := pkg.NewMetrics(cfg.Metrics.Host, cfg.Metrics.Port, cfg.Service.Name, cfg.Platform.Env)
	if err != nil {
		log.Fatalf("cannot init metrics, err: %v", err)
	}
	defer metrics.Disconnect()

	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			infra.Logger(logger),
			infra.MetricsInterceptor(metrics),
		),
	)

	pb.RegisterCredentialsServiceServer(s, thisService)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.Service.Port))
	if err != nil {
		log.Fatalf("Cannnot listen port. Error: %s", err)
	}

	log.Println("Service is listening")

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Cannnot start server. Error: %s", err)
	}
}
