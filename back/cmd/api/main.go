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
        "challenge_enube/internal/usecase"
    )

    func main() {
        _ = godotenv.Load()
    
        // Conexão com o DB
        databaseConn, err := database.Connect(os.Getenv("DATABASE_URL"))
        if err != nil {
            log.Fatal("Failed to connect to database:", err)
        }
        defer databaseConn.Close()
    
        initDatabase(databaseConn)
    
        // Repositórios e serviços
        userRepo := models.NewUserRepository(databaseConn)
        authService := service.NewAuthService(userRepo, os.Getenv("JWT_SECRET"), time.Hour*24)
        authHandler := handlers.NewAuthHandler(authService)
    
        importUC := usecase.NewImportUseCase(databaseConn)
        importHandler := handlers.NewImportHandler(importUC)
    
        partnersHandler := handlers.NewPartnersHandler(databaseConn)
        customersHandler := handlers.NewCustomersHandler(databaseConn)
        productsHandler := handlers.NewProductsHandler(databaseConn)
        billingHandler := handlers.NewBillingHandler(databaseConn)
    
        // Rotas
        router := mux.NewRouter()
    
        // Auth
        router.HandleFunc("/register", authHandler.Register).Methods("POST")
        router.HandleFunc("/login", authHandler.Login).Methods("POST")
    
        // Excel
        router.HandleFunc("/excel", authService.JWTMiddleware(importHandler.HandlerExcel)).Methods("POST")
    
        // Partners
        router.HandleFunc("/partners", authService.JWTMiddleware(partnersHandler.PartnersSummary)).Methods("GET")
    
        // Customers
        router.HandleFunc("/customers", authService.JWTMiddleware(customersHandler.CustomersSummary)).Methods("GET")
        router.HandleFunc("/customers/{id}/billing", authService.JWTMiddleware(customersHandler.BillingByCustomer)).Methods("GET")
    
        // Products
        router.HandleFunc("/products", authService.JWTMiddleware(productsHandler.ProductsSummary)).Methods("GET")
        router.HandleFunc("/skus", authService.JWTMiddleware(productsHandler.SKUsSummary)).Methods("GET")
    
        // Billing
        router.HandleFunc("/billing/monthly", authService.JWTMiddleware(billingHandler.BillingByMonth)).Methods("GET")
    
        // Start server
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

    // Tables Excel
    _, err = db.Exec(`
    CREATE TABLE IF NOT EXISTS partners (
        partner_id TEXT PRIMARY KEY,
        partner_name TEXT NOT NULL
    );
    `)
    if err != nil {
        log.Fatalf("Erro ao criar tabela partners: %v", err)
    }

    _, err = db.Exec(`
    CREATE TABLE IF NOT EXISTS customers (
        customer_id TEXT PRIMARY KEY,
        customer_name TEXT NOT NULL,
        customer_domain TEXT,
        customer_country TEXT,
        partner_id TEXT REFERENCES partners(partner_id)
    );
    `)
    if err != nil {
        log.Fatalf("Erro ao criar tabela customers: %v", err)
    }

    _, err = db.Exec(`
    CREATE TABLE IF NOT EXISTS products (
        product_id TEXT PRIMARY KEY,
        product_name TEXT NOT NULL,
        publisher_id TEXT,
        publisher_name TEXT
    );
    `)
    if err != nil {
        log.Fatalf("Erro ao criar tabela products: %v", err)
    }

    _, err = db.Exec(`
    CREATE TABLE IF NOT EXISTS skus (
        sku_id TEXT PRIMARY KEY,
        sku_name TEXT,
        product_id TEXT REFERENCES products(product_id),
        availability_id TEXT
    );
    `)
    if err != nil {
        log.Fatalf("Erro ao criar tabela skus: %v", err)
    }

    _, err = db.Exec(`
    CREATE TABLE IF NOT EXISTS subscriptions (
        subscription_id TEXT PRIMARY KEY,
        subscription_description TEXT,
        charge_start_date DATE,
        charge_end_date DATE,
        usage_date DATE
    );
    `)
    if err != nil {
        log.Fatalf("Erro ao criar tabela subscriptions: %v", err)
    }

    _, err = db.Exec(`
    CREATE TABLE IF NOT EXISTS meters (
        meter_id TEXT PRIMARY KEY,
        meter_name TEXT,
        meter_type TEXT,
        meter_category TEXT,
        meter_sub_category TEXT,
        meter_region TEXT,
        unit TEXT
    );
    `)
    if err != nil {
        log.Fatalf("Erro ao criar tabela meters: %v", err)
    }

    _, err = db.Exec(`
    CREATE TABLE IF NOT EXISTS billing_records (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
        invoice_number TEXT,
        partner_id TEXT REFERENCES partners(partner_id),
        customer_id TEXT REFERENCES customers(customer_id),
        sku_id TEXT REFERENCES skus(sku_id),
        subscription_id TEXT REFERENCES subscriptions(subscription_id),
        meter_id TEXT REFERENCES meters(meter_id),
        resource_location TEXT,
        consumed_service TEXT,
        resource_group TEXT,
        resource_uri TEXT,
        charge_type TEXT,
        unit_type TEXT,
        unit_price NUMERIC,
        quantity NUMERIC,
        billing_pre_tax_total NUMERIC,
        billing_currency TEXT,
        pricing_pre_tax_total NUMERIC,
        pricing_currency TEXT,
        additional_info JSONB,
        effective_unit_price NUMERIC,
        pc_to_bc_exchange_rate NUMERIC,
        pc_to_bc_exchange_rate_date DATE,
        entitlement_id TEXT,
        entitlement_description TEXT,
        partner_earned_credit_percentage NUMERIC,
        credit_percentage NUMERIC,
        credit_type TEXT,
        benefit_order_id TEXT,
        benefit_id TEXT,
        benefit_type TEXT
    );
    `)
    if err != nil {
        log.Fatalf("Erro ao criar tabela billing_records: %v", err)
    }

	log.Println("Banco de dados inicializado com sucesso!")
}


