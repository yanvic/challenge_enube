package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"

	"challenge_enube/internal/database"
	"challenge_enube/internal/handlers"
	"challenge_enube/internal/models"
	"challenge_enube/internal/service"
	"challenge_enube/internal/repository"
	"challenge_enube/internal/usecase"
    "challenge_enube/internal/jobs"
)

func main() {
	_ = godotenv.Load()

	databaseConn, err := database.Connect(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer databaseConn.Close()

	initDatabase(databaseConn)

    jobs.StartWorker()

	userRepo := models.NewUserRepository(databaseConn)
    clientRepo := repository.NewClientRepository(databaseConn)
    orderRepo := repository.NewOrderRepository(databaseConn)

	authService := service.NewAuthService(userRepo, os.Getenv("JWT_SECRET"), time.Hour*24)
	authHandler := handlers.NewAuthHandler(authService)
    importUseCase := usecase.NewImportUseCase(clientRepo, orderRepo)
    importHandler := handlers.NewImportHandler(importUseCase)

    // importHandler := handlers.NewImportHandler(importUseCase)

	router := mux.NewRouter()

	router.HandleFunc("/register", authHandler.Register).Methods("POST")
	router.HandleFunc("/login", authHandler.Login).Methods("POST")

	router.HandleFunc("/home", authService.JWTMiddleware(handlers.HomeHandler)).Methods("GET")
	router.HandleFunc("/import", authService.JWTMiddleware(handlers.TestHandler)).Methods("GET")
    router.HandleFunc("/excel", authService.JWTMiddleware(importHandler.HandlerExcel)).Methods("POST")
    


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

	_, err = db.Exec(`CREATE INDEX IF NOT EXISTS idx_refresh_tokens_token ON refresh_tokens(token);`)
	if err != nil {
		log.Fatalf("Erro ao criar índice refresh_tokens: %v", err)
	}

    _, err = db.Exec(`
       CREATE TABLE IF NOT EXISTS clients (
            id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
            algo_id VARCHAR(50) UNIQUE NOT NULL,
            nome VARCHAR(255) NOT NULL
        );
    `)
	if err != nil {
		log.Fatalf("Erro ao criar tabela clients: %v", err)
	}
    
    _, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS orders (
            id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
            order_id VARCHAR(50) UNIQUE NOT NULL,
            product VARCHAR(255) NOT NULL,
            amount INT NOT NULL,
            client_id UUID NOT NULL REFERENCES clients(id) ON DELETE CASCADE
        );
    `)
	if err != nil {
		log.Fatalf("Erro ao criar tabela orders: %v", err)
	}

    // Tables Excel
    _, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS partners (
            id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
            partner_id VARCHAR(50) UNIQUE NOT NULL,
            name VARCHAR(255) NOT NULL
        );
    `)
    if err != nil {
        log.Fatalf("Erro ao criar tabelas partners: %v", err)
    }

    _, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS customers (
            id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
            customer_id VARCHAR(50) UNIQUE NOT NULL,
            name VARCHAR(255),
            domain_name VARCHAR(255),
            country VARCHAR(50),
            partner_id UUID NOT NULL REFERENCES partners(id) ON DELETE CASCADE
        );
     `)
    if err != nil {
        log.Fatalf("Erro ao criar tabela customer: %v", err)
    }

    _, err = db.Exec(`
       CREATE TABLE IF NOT EXISTS products (
            id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
            product_id VARCHAR(50) NOT NULL,
            sku_id VARCHAR(50),
            sku_name VARCHAR(255),
            product_name VARCHAR(255),
            publisher_id VARCHAR(50),
            publisher_name VARCHAR(255)
        );
     `)
    if err != nil {
        log.Fatalf("Erro ao criar tabela products: %v", err)
    }

    _, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS subscriptions (
            id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
            subscription_id VARCHAR(50) NOT NULL,
            description TEXT,
            customer_id UUID REFERENCES customers(id),
            product_id UUID REFERENCES products(id)
        );
    `)
    if err != nil {
        log.Fatalf("Erro ao criar tabela subs: %v", err)
    }

    _, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS invoices (
            id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
            invoice_number VARCHAR(50) NOT NULL,
            customer_id UUID REFERENCES customers(id)
        );
    `)
    if err != nil {
    log.Fatalf("Erro ao criar tabela invoice: %v", err)
    }

    _, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS meters (
            id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
            meter_id VARCHAR(50) NOT NULL,
            meter_type VARCHAR(50),
            meter_category VARCHAR(50),
            meter_sub_category VARCHAR(50),
            meter_name VARCHAR(255),
            meter_region VARCHAR(50),
            invoice_id UUID REFERENCES invoices(id),
            product_id UUID REFERENCES products(id)
        );
    `)
    if err != nil {
    log.Fatalf("Erro ao criar tabela meters: %v", err)
    }

    _, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS billing (
            id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
            invoice_id UUID REFERENCES invoices(id),
            meter_id UUID REFERENCES meters(id),
            unit_price NUMERIC,
            quantity NUMERIC,
            unit_type VARCHAR(50),
            pre_tax_total NUMERIC,
            currency VARCHAR(10),
            effective_unit_price NUMERIC,
            pc_to_bc_exchange_rate NUMERIC,
            pc_to_bc_exchange_rate_date DATE,
            entitlement_id VARCHAR(50),
            entitlement_description TEXT,
            partner_earned_credit_percentage NUMERIC,
            credit_percentage NUMERIC,
            credit_type VARCHAR(50),
            benefit_order_id VARCHAR(50),
            benefit_id VARCHAR(50),
            benefit_type VARCHAR(50)
        );
    `)
    if err != nil {
    log.Fatalf("Erro ao criar tabela billings: %v", err)
    }

	log.Println("Banco de dados inicializado com sucesso!")
}


