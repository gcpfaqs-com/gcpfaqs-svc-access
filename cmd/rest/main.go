package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"gcpaccesss-svc-accesss/internal/infraestructure/persistence/postgres"
	"gcpaccesss-svc-accesss/internal/infraestructure/rest"
	"gcpaccesss-svc-accesss/internal/infraestructure/secure"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	secretmanagerpb "cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	pgdriver "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const databaseDSN = "accesss-database-dsn"

func main() {
	db := connectDatabase()

	hasher := secure.NewBCryptHasher()
	tokener := secure.NewJWTGenerator([]byte(os.Getenv("JWT_SECRET")), os.Getenv("JWT_ISSUER"), 24*time.Hour)

	repo := postgres.NewPostgresAccessRepository(db)
	registerHandler := rest.NewRESTRegisterUseCase(repo, hasher, tokener)
	updateHandler := rest.NewRESTUpdateUseCase(repo)
	loginHandler := rest.NewRESTLoginUseCase(repo, hasher, tokener)
	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/accesss/register", func(w http.ResponseWriter, r *http.Request) {
		registerHandler.Register(w, r)
	})

	mux.HandleFunc("POST /api/accesss/login", func(w http.ResponseWriter, r *http.Request) {
		loginHandler.Login(w, r)
	})

	mux.HandleFunc("PUT /api/accesss/{id}", func(w http.ResponseWriter, r *http.Request) {
		updateHandler.Update(w, r)
	})

	mux.HandleFunc("GET /api/accesss/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	address := resolveAddress()
	log.Printf("Database connectivity OK. Starting REST service on port %s...", address)
	if err := http.ListenAndServe(address, corsMiddleware(mux)); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Password-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func connectDatabase() *gorm.DB {
	dsn, err := resolveDSN()
	if err != nil {
		log.Fatalf("database connectivity check failed: %v", err)
	}

	db, err := gorm.Open(pgdriver.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("failed to get database instance: %v", err)
	}

	if err := sqlDB.Ping(); err != nil {
		_ = sqlDB.Close()
		log.Fatalf("database connectivity check failed: %v", err)
	}

	return db
}

func resolveDSN() (string, error) {
	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
	if projectID == "" {
		projectID = os.Getenv("GCP_PROJECT")
	}

	var secretValue string
	if projectID != "" {
		value, err := readDatabaseDSNSecret(projectID)
		if err != nil {
			return "", err
		}
		secretValue = value
	}

	return resolveDatabaseDSN(secretValue, nil, os.Getenv("DATABASE_DSN"))
}

func readDatabaseDSNSecret(projectID string) (string, error) {
	ctx := context.Background()
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return "", fmt.Errorf("error with Secret Manager client: %w", err)
	}
	defer client.Close()

	result, err := client.AccessSecretVersion(ctx, &secretmanagerpb.AccessSecretVersionRequest{
		Name: fmt.Sprintf("projects/%s/secrets/%s/versions/latest", projectID, databaseDSN),
	})
	if status.Code(err) == codes.NotFound {
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("access Secret Manager secret %q: %w", databaseDSN, err)
	}
	if result == nil || result.GetPayload() == nil {
		return "", nil
	}

	return string(result.GetPayload().GetData()), nil
}

func resolveDatabaseDSN(secretValue string, secretErr error, environmentValue string) (string, error) {
	if secretErr != nil {
		return "", secretErr
	}
	if strings.TrimSpace(secretValue) != "" {
		return secretValue, nil
	}
	if environmentValue != "" {
		return environmentValue, nil
	}

	return "", errors.New("DATABASE_DSN must be set in Secret Manager secret accesss-database-dsn or in the DATABASE_DSN environment variable")
}

func resolveAddress() string {
	address := os.Getenv("SERVER_ADDRESS")
	if address == "" {
		log.Fatal("SERVER_ADDRESS must be set")
	}

	return address
}
