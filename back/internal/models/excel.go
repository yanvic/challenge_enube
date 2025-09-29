package models

type Order struct {
	ID      string `db:"id"`
	OrderID string `db:"order_id"`
	Product string `db:"product"`
	Amount  int    `db:"amount"`
    UnitPrice float64   `db:"unit_price"`
}

type Partner struct {
	ID        string `db:"id"`
	// PartnerID string `db:"partner_id"`
	Name      string `db:"name"`
}

type Customer struct {
	ID         string `db:"id"`
	CustomerID string `db:"customer_id"`
	Name       string `db:"name"`
	DomainName string `db:"domain_name"`
	Country    string `db:"country"`
	PartnerID  string `db:"partner_id"`
    MpnID      string `db:"mpn_id"`
    Tier2MpnID string `db:"tier2_mpn_id"`
}

type Product struct {
	ID            string `db:"id"`
	ProductID     string `db:"product_id"`
	SkuID         string `db:"sku_id"`
	SkuName       string `db:"sku_name"`
	ProductName   string `db:"product_name"`
	PublisherID   string `db:"publisher_id"`
	PublisherName string `db:"publisher_name"`
}

type Subscription struct {
	ID             string `db:"id"`
	SubscriptionID string `db:"subscription_id"`
	Description    string `db:"description"`
	CustomerID     string `db:"customer_id"`
	ProductID      string `db:"product_id"`
}

type Invoice struct {
	ID            string `db:"id"`
	InvoiceNumber string `db:"invoice_number"`
	CustomerID    string `db:"customer_id"`
}

type Meter struct {
	ID            string `db:"id"`
	MeterID       string `db:"meter_id"`
	MeterType     string `db:"meter_type"`
	MeterCategory string `db:"meter_category"`
	MeterSubCat   string `db:"meter_sub_category"`
	MeterName     string `db:"meter_name"`
	MeterRegion   string `db:"meter_region"`
	InvoiceID     string `db:"invoice_id"`
	ProductID     string `db:"product_id"`
}

type BillingFinancial struct {
	ID                     string  `db:"id"`
	InvoiceID              string  `db:"invoice_id"`
	MeterID                string  `db:"meter_id"`
	SubscriptionID         string  `db:"subscription_id"`
	UnitPrice              float64 `db:"unit_price"`
	Quantity               float64 `db:"quantity"`
	UnitType               string  `db:"unit_type"`
	BillingPreTaxTotal     float64 `db:"billing_pre_tax_total"`
	BillingCurrency        string  `db:"billing_currency"`
	PricingPreTaxTotal     float64 `db:"pricing_pre_tax_total"`
	PricingCurrency        string  `db:"pricing_currency"`
	ServiceInfo1           string  `db:"service_info1"`
	ServiceInfo2           string  `db:"service_info2"`
	Tags                   string  `db:"tags"`
	AdditionalInfo         string  `db:"additional_info"`
	EffectiveUnitPrice     float64 `db:"effective_unit_price"`
	PCToBCExchangeRate     float64 `db:"pc_to_bc_exchange_rate"`
	PCToBCExchangeRateDate string  `db:"pc_to_bc_exchange_rate_date"`
}

type Billing struct {
	ID                       string  `db:"id"`
	BillingFinancialID       string  `db:"billing_financial_id"`
	EntitlementID            string  `db:"entitlement_id"`
	EntitlementDescription   string  `db:"entitlement_description"`
	PartnerEarnedCreditPerc  float64 `db:"partner_earned_credit_percentage"`
	CreditPercentage         float64 `db:"credit_percentage"`
	CreditType               string  `db:"credit_type"`
	BenefitOrderID           string  `db:"benefit_order_id"`
	BenefitID                string  `db:"benefit_id"`
	BenefitType              string  `db:"benefit_type"`
}
