package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type CustomersHandler struct {
	DB *sql.DB
}

func NewCustomersHandler(db *sql.DB) *CustomersHandler {
	return &CustomersHandler{DB: db}
}

func (h *CustomersHandler) CustomersSummary(w http.ResponseWriter, r *http.Request) {
	rows, err := h.DB.Query(`
		SELECT c.customer_name, SUM(b.billing_pre_tax_total) AS total
		FROM billing_records b
		LEFT JOIN customers c ON b.customer_id = c.customer_id
		GROUP BY c.customer_name
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
			"customer": name,
			"total":    total,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func (h *CustomersHandler) BillingByCustomer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	custID := vars["id"]

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
		SELECT invoice_number, sku_id, subscription_id, meter_id, billing_pre_tax_total, billing_currency, charge_type
		FROM billing_records
		WHERE customer_id=$1
		ORDER BY invoice_number DESC
		LIMIT $2 OFFSET $3
	`, custID, pageSize, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	result := []map[string]interface{}{}
	for rows.Next() {
		var invoice, sku, sub, meter, charge, currency string
		var total float64
		rows.Scan(&invoice, &sku, &sub, &meter, &total, &currency, &charge)
		result = append(result, map[string]interface{}{
			"invoice":               invoice,
			"sku_id":                sku,
			"subscription_id":       sub,
			"meter_id":              meter,
			"billing_pre_tax_total": total,
			"billing_currency":      currency,
			"charge_type":           charge,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
