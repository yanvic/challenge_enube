package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

type BillingHandler struct {
	DB *sql.DB
}

func NewBillingHandler(db *sql.DB) *BillingHandler {
	return &BillingHandler{DB: db}
}

func (h *BillingHandler) BillingByMonth(w http.ResponseWriter, r *http.Request) {
	startDateStr := r.URL.Query().Get("startDate")
	endDateStr := r.URL.Query().Get("endDate")

	now := time.Now()
	var startDate, endDate time.Time
	var err error

	if startDateStr == "" {
		startDate = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	} else {
		startDate, err = time.Parse("2006-01-02", startDateStr)
		if err != nil {
			http.Error(w, "startDate inválido", http.StatusBadRequest)
			return
		}
	}

	if endDateStr == "" {
		endDate = now
	} else {
		endDate, err = time.Parse("2006-01-02", endDateStr)
		if err != nil {
			http.Error(w, "endDate inválido", http.StatusBadRequest)
			return
		}
	}

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
		SELECT DATE_TRUNC('month', s.charge_start_date) AS month,
			SUM(b.billing_pre_tax_total) AS total
		FROM billing_records b
		JOIN subscriptions s ON b.subscription_id = s.subscription_id
		WHERE s.charge_start_date BETWEEN $1 AND $2
		GROUP BY month
		ORDER BY month
		LIMIT $3 OFFSET $4
	`, startDate, endDate, pageSize, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	result := []map[string]interface{}{}
	for rows.Next() {
		var month time.Time
		var total float64
		rows.Scan(&month, &total)
		result = append(result, map[string]interface{}{
			"month": month.Format("2006-01"),
			"total": total,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
