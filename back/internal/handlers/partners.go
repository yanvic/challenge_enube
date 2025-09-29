package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
)

type PartnersHandler struct {
	DB *sql.DB
}

func NewPartnersHandler(db *sql.DB) *PartnersHandler {
	return &PartnersHandler{DB: db}
}

func (h *PartnersHandler) PartnersSummary(w http.ResponseWriter, r *http.Request) {
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
		SELECT partner_id, SUM(billing_pre_tax_total) AS total
		FROM billing_records
		GROUP BY partner_id
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
			"partner_id": id,
			"total":      total,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
