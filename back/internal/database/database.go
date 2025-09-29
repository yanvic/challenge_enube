package database

import (
    "database/sql"
    "log"
    "time"
    _ "github.com/lib/pq"
)

func Connect(connectionString string) (*sql.DB, error) {
    db, err := sql.Open("postgres", connectionString)
    if err != nil {
        return nil, err
    }

    db.SetMaxOpenConns(10)
    db.SetMaxIdleConns(5)
    db.SetConnMaxLifetime(time.Minute * 5)


    if err := db.Ping(); err != nil {
        return nil, err
    }

    log.Println("Connected to the database successfully")
    return db, nil
}
