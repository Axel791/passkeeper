package main

import (
	"net"

	"github.com/Axel791/auth/internal/config"
	"github.com/Axel791/auth/internal/db"
	grpcV1Handler "github.com/Axel791/auth/internal/grpc/auth/v1"
	apiV1Handlers "github.com/Axel791/auth/internal/rest/user/v1"
	"github.com/Axel791/auth/internal/services"
	"github.com/Axel791/passkeeper_grpc/pb"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"net/http"
)

func main() {
	log := logrus.New()

	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})
	log.SetLevel(logrus.InfoLevel)

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("error loading config: %v", err)
	}

	dbConn, err := db.ConnectDB(cfg.DatabaseDSN, cfg)
	if err != nil {
		log.Fatalf("error connecting to database: %v", err)
	}
	defer func() {
		if dbConn != nil {
			_ = dbConn.Close()
		}
	}()

	router := chi.NewRouter()
	router.Use(middleware.StripSlashes)
	router.Use(middleware.Logger)

	tokenService := services.NewTokenService(cfg.SecretKey)
	hashService := services.NewHashPasswordService(cfg.SecretKey)

	providers := config.NewProviders(dbConn)
	useCases := config.NewUseCases(log, providers, hashService, tokenService)

	router.Route("/api/v1", func(r chi.Router) {
		r.Method(
			http.MethodPost,
			"/users/registration",
			apiV1Handlers.NewRegister(useCases.Registration),
		)
		r.Method(
			http.MethodPost,
			"/users/login",
			apiV1Handlers.NewLogin(useCases.Login),
		)
	})

	grpcServer := grpc.NewServer()
	authServer := grpcV1Handler.NewAuthServer(useCases.GroupUseCases, useCases.Validate)

	pb.RegisterAuthServiceServer(grpcServer, authServer)

	lis, err := net.Listen(cfg.GrpcNetwork, cfg.GrpcAddress)
	if err != nil {
		log.Fatalf("failed to listen on :50051: %v", err)
	}

	if err = grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve gRPC: %v", err)
	}

	if err != nil {
		log.Fatalf("failed to listen on :50051: %v", err)
	}
	log.Println("Starting Loyalty gRPC server on", cfg.GrpcAddress)
}
