package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type ProductsHandler struct {
	DB *sql.DB
}

func NewProductsHandler(db *sql.DB) *ProductsHandler {
	return &ProductsHandler{DB: db}
}

func (h *ProductsHandler) ProductsSummary(w http.ResponseWriter, r *http.Request) {
	rows, err := h.DB.Query(`
		SELECT p.product_name, SUM(b.billing_pre_tax_total) AS total
		FROM billing_records b
		LEFT JOIN skus s ON b.sku_id = s.sku_id
		LEFT JOIN products p ON s.product_id = p.product_id
		GROUP BY p.product_name
		ORDER BY total DESC
		LIMIT 10
	`)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	result := []map[string]interface{}{}
	for rows.Next() {
		var name string
		var total float64
		rows.Scan(&name, &total)
		result = append(result, map[string]interface{}{
			"product": name,
			"total":   total,
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

func (h *ProductsHandler) CustomersByProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	product := vars["id"]

	rows, err := h.DB.Query(`
		SELECT customer_id, SUM(billing_pre_tax_total) AS total
		FROM billing_records
		WHERE product_name=$1
		GROUP BY customer_id
		ORDER BY total DESC
	`, product)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	result := []map[string]interface{}{}
	for rows.Next() {
		var cust string
		var total float64
		rows.Scan(&cust, &total)
		result = append(result, map[string]interface{}{
			"customer_id": cust,
			"total":       total,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
