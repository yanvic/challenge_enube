package usecase

import (
	"database/sql"
	// "encoding/csv"
	"encoding/json"
	// "golang.org/x/text/encoding/charmap"
    "github.com/google/uuid"
    "github.com/xuri/excelize/v2"
	"fmt"
	"log"
	// "os"
	"strconv"
	"strings"
	"time"
)

type BillingCSV struct {
	PartnerId                string
	PartnerName              string
	CustomerId               string
	CustomerName             string
	CustomerDomainName       string
	CustomerCountry          string
	MpnId                    string
	Tier2MpnId               string
	InvoiceNumber            string
	ProductId                string
	SkuId                    string
	AvailabilityId           string
	SkuName                  string
	ProductName              string
	PublisherName            string
	PublisherId              string
	SubscriptionDescription  string
	SubscriptionId           string
	ChargeStartDate          string
	ChargeEndDate            string
	UsageDate                string
	MeterType                string
	MeterCategory            string
	MeterId                  string
	MeterSubCategory         string
	MeterName                string
	MeterRegion              string
	Unit                     string
	ResourceLocation         string
	ConsumedService          string
	ResourceGroup            string
	ResourceURI              string
	ChargeType               string
	UnitPrice                string
	Quantity                 string
	UnitType                 string
	BillingPreTaxTotal       string
	BillingCurrency          string
	PricingPreTaxTotal       string
	PricingCurrency          string
	ServiceInfo1             string
	ServiceInfo2             string
	Tags                     string
	AdditionalInfo           string
	EffectiveUnitPrice       string
	PCToBCExchangeRate       string
	PCToBCExchangeRateDate   string
	EntitlementId            string
	EntitlementDescription   string
	PartnerEarnedCreditPct   string
	CreditPercentage         string
	CreditType               string
	BenefitOrderId           string
	BenefitId                string
	BenefitType              string
}

type ImportUseCase struct {
	DB *sql.DB
}

func NewImportUseCase(db *sql.DB) *ImportUseCase {
	return &ImportUseCase{DB: db}
}

func normalizeUUID(id string) (string, error) {
    id = strings.TrimSpace(id)
    id = strings.ReplaceAll(id, " ", "")
    id = strings.ToLower(id)

    if u, err := uuid.Parse(id); err == nil {
        return u.String(), nil
    }

    if len(id) == 32 {
        formatted := fmt.Sprintf("%s-%s-%s-%s-%s",
            id[0:8], id[8:12], id[12:16], id[16:20], id[20:])
        if u, err := uuid.Parse(formatted); err == nil {
            return u.String(), nil
        }
    }

    return "", fmt.Errorf("uuid inválido: %s", id)
}

func parseNumber(s string) float64 {
	if s == "" {
		return 0
	}
	s = strings.ReplaceAll(s, ",", ".")
	f, _ := strconv.ParseFloat(s, 64)
	return f
}

func parseDate(s string) string {
	if s == "" {
		return "NULL"
	}
	layouts := []string{"1/2/2006", "01/02/2006"}
	for _, layout := range layouts {
		if t, err := time.Parse(layout, s); err == nil {
			return fmt.Sprintf("'%s'", t.Format("2006-01-02"))
		}
	}
	return "NULL"
}

func (uc *ImportUseCase) ImportFromExcel(filePath string) error {
	log.Println("Iniciando importação...")

	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return fmt.Errorf("erro ao abrir excel: %w", err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Println("erro ao fechar excel:", err)
		}
	}()

	rows, err := f.GetRows("Planilha1")
	if err != nil {
		return fmt.Errorf("erro ao ler planilha: %w", err)
	}

	var data []BillingCSV

	for i, row := range rows {
		if i == 0 {
			// pula header
			continue
		}

		for len(row) < 55 {
			row = append(row, "")
		}

		rec := BillingCSV{
			PartnerId:               row[0],
			PartnerName:             row[1],
			CustomerId:              row[2],
			CustomerName:            row[3],
			CustomerDomainName:      row[4],
			CustomerCountry:         row[5],
			MpnId:                   row[6],
			Tier2MpnId:              row[7],
			InvoiceNumber:           row[8],
			ProductId:               row[9],
			SkuId:                   row[10],
			AvailabilityId:          row[11],
			SkuName:                 row[12],
			ProductName:             row[13],
			PublisherName:           row[14],
			PublisherId:             row[15],
			SubscriptionDescription: row[16],
			SubscriptionId:          row[17],
			ChargeStartDate:         parseDate(row[18]),
			ChargeEndDate:           parseDate(row[19]),
			UsageDate:               parseDate(row[20]),
			MeterType:               row[21],
			MeterCategory:           row[22],
			MeterId:                 row[23],
			MeterSubCategory:        row[24],
			MeterName:               row[25],
			MeterRegion:             row[26],
			Unit:                    row[27],
			ResourceLocation:        row[28],
			ConsumedService:         row[29],
			ResourceGroup:           row[30],
			ResourceURI:             row[31],
			ChargeType:              row[32],
			UnitPrice:               row[33],
			Quantity:                row[34],
			UnitType:                row[35],
			BillingPreTaxTotal:      row[36],
			BillingCurrency:         row[37],
			PricingPreTaxTotal:      row[38],
			PricingCurrency:         row[39],
			ServiceInfo1:            row[40],
			ServiceInfo2:            row[41],
			Tags:                    row[42],
			AdditionalInfo:          row[43],
			EffectiveUnitPrice:      row[44],
			PCToBCExchangeRate:      row[45],
			PCToBCExchangeRateDate:  row[46],
			EntitlementId:           row[47],
			EntitlementDescription:  row[48],
			PartnerEarnedCreditPct:  row[49],
			CreditPercentage:        row[50],
			CreditType:              row[51],
			BenefitOrderId:          row[52],
			BenefitId:               row[53],
			BenefitType:             row[54],
		}

		data = append(data, rec)
	}

	log.Printf("Total de registros lidos: %d", len(data))

	return uc.bulkInsert(data)
}

func (uc *ImportUseCase) bulkInsert(data []BillingCSV) error {
	tx, err := uc.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// --- Partners ---
	partners := map[string]string{}
	for i := range data {
		if data[i].PartnerId == "" {
			continue
		}
		pid, err := normalizeUUID(data[i].PartnerId)
		if err != nil {
			log.Println("Ignorando PartnerId inválido:", data[i].PartnerId)
			continue
		}
		data[i].PartnerId = pid
		partners[pid] = strings.TrimSpace(data[i].PartnerName)
	}

	if len(partners) > 0 {
		vals := []string{}
		for id, name := range partners {
			vals = append(vals, fmt.Sprintf("('%s','%s')", id, strings.ReplaceAll(name, "'", "''")))
		}
		query := "INSERT INTO partners(partner_id, partner_name) VALUES " + strings.Join(vals, ",") + " ON CONFLICT DO NOTHING"
		if _, err := tx.Exec(query); err != nil {
			return fmt.Errorf("erro ao inserir partners: %w", err)
		}
	}

	// --- Customers ---
	customers := map[string]BillingCSV{}
	for i := range data {
		if data[i].CustomerId == "" || data[i].PartnerId == "" {
			continue
		}
		if _, ok := partners[data[i].PartnerId]; !ok {
			log.Println("Ignorando customer com PartnerId inválido:", data[i].CustomerId, data[i].PartnerId)
			continue
		}
		cid := strings.TrimSpace(data[i].CustomerId)
		data[i].CustomerId = cid
		data[i].CustomerName = strings.TrimSpace(data[i].CustomerName)
		data[i].CustomerDomainName = strings.TrimSpace(data[i].CustomerDomainName)
		data[i].CustomerCountry = strings.TrimSpace(data[i].CustomerCountry)
		customers[cid] = data[i]
	}

	if len(customers) > 0 {
		vals := []string{}
		for _, c := range customers {
			vals = append(vals, fmt.Sprintf(
				"('%s','%s','%s','%s','%s')",
				c.CustomerId,
				strings.ReplaceAll(c.CustomerName, "'", "''"),
				strings.ReplaceAll(c.CustomerDomainName, "'", "''"),
				strings.ReplaceAll(c.CustomerCountry, "'", "''"),
				c.PartnerId,
			))
		}
		query := "INSERT INTO customers(customer_id, customer_name, customer_domain, customer_country, partner_id) VALUES " + strings.Join(vals, ",") + " ON CONFLICT DO NOTHING"
		if _, err := tx.Exec(query); err != nil {
			return fmt.Errorf("erro ao inserir customers: %w", err)
		}
	}

	// --- Products ---
	products := map[string]BillingCSV{}
	for _, d := range data {
		if d.ProductId == "" {
			continue
		}
		products[d.ProductId] = d
	}

	if len(products) > 0 {
		vals := []string{}
		for _, p := range products {
			vals = append(vals, fmt.Sprintf(
				"('%s','%s','%s','%s')",
				p.ProductId,
				strings.ReplaceAll(p.ProductName, "'", "''"),
				p.PublisherId,
				strings.ReplaceAll(p.PublisherName, "'", "''"),
			))
		}
		query := "INSERT INTO products(product_id, product_name, publisher_id, publisher_name) VALUES " + strings.Join(vals, ",") + " ON CONFLICT DO NOTHING"
		if _, err := tx.Exec(query); err != nil {
			return fmt.Errorf("erro ao inserir products: %w", err)
		}
	}

	// --- SKUs ---
	skus := map[string]BillingCSV{}
	for _, d := range data {
		if d.SkuId == "" {
			continue
		}
		skus[d.SkuId] = d
	}

	if len(skus) > 0 {
		vals := []string{}
		for _, s := range skus {
			vals = append(vals, fmt.Sprintf(
				"('%s','%s','%s','%s')",
				s.SkuId,
				strings.ReplaceAll(s.SkuName, "'", "''"),
				s.ProductId,
				s.AvailabilityId,
			))
		}
		query := "INSERT INTO skus(sku_id, sku_name, product_id, availability_id) VALUES " + strings.Join(vals, ",") + " ON CONFLICT DO NOTHING"
		if _, err := tx.Exec(query); err != nil {
			return fmt.Errorf("erro ao inserir skus: %w", err)
		}
	}

	// --- Subscriptions ---
	subscriptions := map[string]BillingCSV{}
	for _, d := range data {
		if d.SubscriptionId == "" {
			continue
		}
		subscriptions[d.SubscriptionId] = d
	}

	if len(subscriptions) > 0 {
		vals := []string{}
		for _, s := range subscriptions {
			vals = append(vals, fmt.Sprintf(
				"('%s','%s',%s,%s,%s)",
				s.SubscriptionId,
				strings.ReplaceAll(s.SubscriptionDescription, "'", "''"),
				s.ChargeStartDate,
				s.ChargeEndDate,
				s.UsageDate,
			))
		}
		query := "INSERT INTO subscriptions(subscription_id, subscription_description, charge_start_date, charge_end_date, usage_date) VALUES " + strings.Join(vals, ",") + " ON CONFLICT DO NOTHING"
		if _, err := tx.Exec(query); err != nil {
			return fmt.Errorf("erro ao inserir subscriptions: %w", err)
		}
	}

	// --- Meters ---
	meters := map[string]BillingCSV{}
	for _, d := range data {
		if d.MeterId == "" {
			continue
		}
		meters[d.MeterId] = d
	}

	if len(meters) > 0 {
		vals := []string{}
		for _, m := range meters {
			vals = append(vals, fmt.Sprintf(
				"('%s','%s','%s','%s','%s','%s','%s')",
				m.MeterId,
				strings.ReplaceAll(m.MeterName, "'", "''"),
				m.MeterType,
				m.MeterCategory,
				m.MeterSubCategory,
				m.MeterRegion,
				m.Unit,
			))
		}
		query := "INSERT INTO meters(meter_id, meter_name, meter_type, meter_category, meter_sub_category, meter_region, unit) VALUES " + strings.Join(vals, ",") + " ON CONFLICT DO NOTHING"
		if _, err := tx.Exec(query); err != nil {
			return fmt.Errorf("erro ao inserir meters: %w", err)
		}
	}

	// --- Billing Records ---
	vals := []string{}
	for _, r := range data {
		unitPrice := parseNumber(r.UnitPrice)
		qty := parseNumber(r.Quantity)
		billingTotal := parseNumber(r.BillingPreTaxTotal)
		pricingTotal := parseNumber(r.PricingPreTaxTotal)
		effPrice := parseNumber(r.EffectiveUnitPrice)
		pcToBc := parseNumber(r.PCToBCExchangeRate)
		partnerEarnedCredit := parseNumber(r.PartnerEarnedCreditPct)
		creditPercentage := parseNumber(r.CreditPercentage)

		additional := "{}"
		if json.Valid([]byte(r.AdditionalInfo)) {
			additional = strings.ReplaceAll(r.AdditionalInfo, "'", "''")
		}

		pcToBcDate := parseDate(r.PCToBCExchangeRateDate)
		vals = append(vals, fmt.Sprintf(
			"('%s','%s','%s','%s','%s','%s','%s','%s','%s','%s','%s','%s',%f,%f,%f,%f,%f,'%s',%f,%s,'%s','%s',%f,%f,'%s','%s','%s','%s')",
			r.InvoiceNumber,
			r.PartnerId,
			r.CustomerId,
			r.SkuId,
			r.SubscriptionId,
			r.MeterId,
			strings.ReplaceAll(r.ResourceLocation, "'", "''"),
			strings.ReplaceAll(r.ConsumedService, "'", "''"),
			strings.ReplaceAll(r.ResourceGroup, "'", "''"),
			strings.ReplaceAll(r.ResourceURI, "'", "''"),
			r.ChargeType,
			r.UnitType,
			unitPrice,
			qty,
			billingTotal,
			pricingTotal,
			effPrice,
			additional,
			pcToBc,
			pcToBcDate,
			r.EntitlementId,
			strings.ReplaceAll(r.EntitlementDescription, "'", "''"),
			partnerEarnedCredit,
			creditPercentage,
			r.CreditType,
			r.BenefitOrderId,
			r.BenefitId,
			r.BenefitType,
		))
		}

	if len(vals) > 0 {
		query := "INSERT INTO billing_records(invoice_number, partner_id, customer_id, sku_id, subscription_id, meter_id, resource_location, consumed_service, resource_group, resource_uri, charge_type, unit_type, unit_price, quantity, billing_pre_tax_total, pricing_pre_tax_total, effective_unit_price, additional_info, pc_to_bc_exchange_rate, pc_to_bc_exchange_rate_date, entitlement_id, entitlement_description, partner_earned_credit_percentage, credit_percentage, credit_type, benefit_order_id, benefit_id, benefit_type) VALUES " + strings.Join(vals, ",")
		if _, err := tx.Exec(query); err != nil {
			return fmt.Errorf("erro ao inserir billing_records: %w", err)
		}
	}

	return tx.Commit()
}
