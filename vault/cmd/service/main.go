package main

import (
	"net/http"

	"github.com/Axel791/passkeeper_grpc/pb"
	"github.com/Axel791/vault/internal/config"
	"github.com/Axel791/vault/internal/db"
	apiV1Handlers "github.com/Axel791/vault/internal/rest/vault_item/v1"
	"github.com/Axel791/vault/internal/services"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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

	conn, err := grpc.NewClient(
		cfg.GrpcAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("failed to create new client: %v", err)
	}

	defer conn.Close()

	authClient := pb.NewAuthServiceClient(conn)

	vaultProviders := config.NewProviders(dbConn, authClient)
	vaultUseCases := config.NewUseCases(log, vaultProviders)

	authService := services.NewValidateToken(vaultProviders.AuthProvider)

	router.Route("/api/v1", func(r chi.Router) {
		r.Method(
			http.MethodGet,
			"/vault/{vault_id}",
			apiV1Handlers.NewGetVaultItem(authService, vaultUseCases.GetVaultItems),
		)
		r.Method(
			http.MethodPost,
			"/vault",
			apiV1Handlers.NewCreateVaultItemV1(authService, vaultUseCases.CreateVaultItem),
		)
		r.Method(
			http.MethodPatch,
			"/vault",
			apiV1Handlers.NewUpdateVaultItem(authService, vaultUseCases.UpdateVaultItem),
		)
	})

	log.Infof("server started on %s", cfg.Address)
	err = http.ListenAndServe(cfg.Address, router)
	if err != nil {
		log.Fatalf("error starting server: %v", err)
	}
}
