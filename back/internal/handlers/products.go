package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
)

type ProductsHandler struct {
	DB *sql.DB
}

func NewProductsHandler(db *sql.DB) *ProductsHandler {
	return &ProductsHandler{DB: db}
}

func (h *ProductsHandler) ProductsSummary(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("pageSize"))
	if pageSize < 1 {
		pageSize = 50
	}
	offset := (page - 1) * pageSize

	rows, err := h.DB.Query(`
		SELECT s.product_id, SUM(b.billing_pre_tax_total) AS total
		FROM billing_records b
		JOIN skus s ON b.sku_id = s.sku_id
		GROUP BY s.product_id
		ORDER BY total DESC
		LIMIT $1 OFFSET $2
	`, pageSize, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	result := []map[string]interface{}{}
	for rows.Next() {
		var productID string
		var total float64
		rows.Scan(&productID, &total)
		result = append(result, map[string]interface{}{
			"product_id": productID,
			"total":      total,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func (h *ProductsHandler) SKUsSummary(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("pageSize"))
	if pageSize < 1 {
		pageSize = 50
	}
	offset := (page - 1) * pageSize

	rows, err := h.DB.Query(`
		SELECT sku_id, SUM(billing_pre_tax_total) AS total
		FROM billing_records
		GROUP BY sku_id
		ORDER BY total DESC
		LIMIT $1 OFFSET $2
	`, pageSize, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	result := []map[string]interface{}{}
	for rows.Next() {
		var id string
		var total float64
		rows.Scan(&id, &total)
		result = append(result, map[string]interface{}{
			"sku_id": id,
			"total":  total,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
