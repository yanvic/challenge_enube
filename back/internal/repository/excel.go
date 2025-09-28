package repository

import (
    "database/sql"
    // "strconv"
    "challenge_enube/internal/models"
)

type ClientRepository struct {
    db *sql.DB
}

func NewClientRepository(db *sql.DB) *ClientRepository {
    return &ClientRepository{db: db}
}

func (r *ClientRepository) SaveBatch(clients []models.Client) error {
    if len(clients) == 0 {
        return nil
    }

    query := `INSERT INTO clients (algo_id, nome) VALUES ($1, $2) ON CONFLICT (algo_id) DO NOTHING`
    stmt, err := r.db.Prepare(query)
    if err != nil {
        return err
    }
    defer stmt.Close()

    for _, c := range clients {
        _, err := stmt.Exec(c.AlgoID, c.Nome)
        if err != nil {
            return err
        }
    }
    return nil
}

type OrderRepository struct {
    db *sql.DB
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
    return &OrderRepository{db: db}
}

func (r *OrderRepository) SaveBatch(orders []models.Order) error {
    if len(orders) == 0 {
        return nil
    }

    query := `INSERT INTO orders (order_id, product, amount, client_id) VALUES ($1, $2, $3, (SELECT id FROM clients WHERE algo_id = $4)) ON CONFLICT (order_id) DO NOTHING`
    stmt, err := r.db.Prepare(query)
    if err != nil {
        return err
    }
    defer stmt.Close()

    for _, o := range orders {
        _, err := stmt.Exec(o.OrderID, o.Product, o.Amount, o.ClientID)
        if err != nil {
            return err
        }
    }
    return nil
}
