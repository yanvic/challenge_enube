package main

import (
    "log"
    "net/http"
    "os"
    "time"
    "database/sql"

    "github.com/gorilla/mux"
    "github.com/joho/godotenv"


    "challenge_enube/internal/service"
    db "challenge_enube/internal/database"
    "challenge_enube/internal/handlers"
    "challenge_enube/internal/models"
)

func main() {
    err := godotenv.Load()
    if err != nil {
        log.Println("No .env file found, using environment variables")
    }

    database, err := db.Connect(os.Getenv("DATABASE_URL"))
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }
    defer database.Close()

	initDatabase(database)

    userRepo := models.NewUserRepository(database)
    authService := service.NewAuthService(userRepo, os.Getenv("JWT_SECRET"), time.Hour*24) // Token válido 24h
    authHandler := handlers.NewAuthHandler(authService)

    router := mux.NewRouter()

    router.HandleFunc("/register", authHandler.Register).Methods("POST")
    router.HandleFunc("/login", authHandler.Login).Methods("POST")

    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    log.Println("Server running on port", port)
    log.Fatal(http.ListenAndServe(":"+port, router))
}

func initDatabase(db *sql.DB) {
    _, err := db.Exec(`CREATE EXTENSION IF NOT EXISTS "pgcrypto";`)
    if err != nil {
        log.Fatalf("Erro ao criar extensão pgcrypto: %v", err)
    }

    // Table users
    _, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS users (
            id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
            email VARCHAR(255) UNIQUE NOT NULL,
            username VARCHAR(50) UNIQUE NOT NULL,
            password_hash VARCHAR(255) NOT NULL,
            created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
            last_login TIMESTAMP
        );
    `)
    if err != nil {
        log.Fatalf("Erro ao criar tabela users: %v", err)
    }

    // Table refresh_tokens
    _, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS refresh_tokens (
            id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
            user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
            token VARCHAR(255) UNIQUE NOT NULL,
            expires_at TIMESTAMP NOT NULL,
            created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
            revoked BOOLEAN NOT NULL DEFAULT FALSE
        );
    `)
    if err != nil {
        log.Fatalf("Erro ao criar tabela refresh_tokens: %v", err)
    }

    // Index refresh_tokens
    _, err = db.Exec(`
        CREATE INDEX IF NOT EXISTS idx_refresh_tokens_token ON refresh_tokens(token);
    `)
    if err != nil {
        log.Fatalf("Erro ao criar índice refresh_tokens: %v", err)
    }

    log.Println("Banco de dados inicializado com sucesso!")
}
