package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

type MetricsHandler struct {
	DB *sql.DB
}

func NewMetricsHandler(db *sql.DB) *MetricsHandler {
	return &MetricsHandler{DB: db}
}

func (h *MetricsHandler) UsageSummary(w http.ResponseWriter, r *http.Request) {
	row := h.DB.QueryRow(`
		SELECT 
			COUNT(*) AS total_records,
			COUNT(DISTINCT b.customer_id) AS total_customers,
			COUNT(DISTINCT p.product_id) AS total_products,
			COALESCE(SUM(b.billing_pre_tax_total), 0) AS total_revenue
		FROM billing_records b
		LEFT JOIN skus s ON b.sku_id = s.sku_id
		LEFT JOIN products p ON s.product_id = p.product_id
	`)

	var records, customers, products int
	var revenue float64
	err := row.Scan(&records, &customers, &products, &revenue)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	result := map[string]interface{}{
		"total_records":   records,
		"total_customers": customers,
		"total_products":  products,
		"total_revenue":   revenue,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
