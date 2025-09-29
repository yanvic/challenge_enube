
package repository

import (
	"database/sql"
	"challenge_enube/internal/models"
)

// -------------------- OrderRepository --------------------

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

	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(`
		INSERT INTO orders (order_id, product, amount, unit_price)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (order_id) DO NOTHING
	`)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	for _, o := range orders {
		if _, err := stmt.Exec(o.OrderID, o.Product, o.Amount, o.UnitPrice); err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

// -------------------- PartnerRepository --------------------

type PartnerRepository struct {
	db *sql.DB
}

func NewPartnerRepository(db *sql.DB) *PartnerRepository {
	return &PartnerRepository{db: db}
}

func (r *PartnerRepository) SaveBatch(partners []models.Partner) error {
	if len(partners) == 0 {
		return nil
	}

	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(`
		INSERT INTO partners (partner_id, name)
		VALUES ($1, $2)
		ON CONFLICT (partner_id) DO NOTHING
	`)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	for _, p := range partners {
		if _, err := stmt.Exec(p.ID, p.Name); err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

// -------------------- CustomerRepository --------------------

type CustomerRepository struct {
	db *sql.DB
}

func NewCustomerRepository(db *sql.DB) *CustomerRepository {
	return &CustomerRepository{db: db}
}

func (r *CustomerRepository) SaveBatch(customers []models.Customer) error {
	if len(customers) == 0 {
		return nil
	}

	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(`
		INSERT INTO customers (customer_id, name, domain_name, country, partner_id, mpn_id, tier2_mpn_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (customer_id) DO NOTHING
	`)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	for _, c := range customers {
		if _, err := stmt.Exec(c.CustomerID, c.Name, c.DomainName, c.Country, c.PartnerID, c.MpnID, c.Tier2MpnID); err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

// -------------------- ProductRepository --------------------

type ProductRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) SaveBatch(products []models.Product) error {
	if len(products) == 0 {
		return nil
	}

	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(`
		INSERT INTO products (product_id, sku_id, sku_name, product_name, publisher_id, publisher_name)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (product_id) DO NOTHING
	`)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	for _, p := range products {
		if _, err := stmt.Exec(p.ProductID, p.SkuID, p.SkuName, p.ProductName, p.PublisherID, p.PublisherName); err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

// -------------------- SubscriptionRepository --------------------

type SubscriptionRepository struct {
	db *sql.DB
}

func NewSubscriptionRepository(db *sql.DB) *SubscriptionRepository {
	return &SubscriptionRepository{db: db}
}

func (r *SubscriptionRepository) SaveBatch(subs []models.Subscription) error {
	if len(subs) == 0 {
		return nil
	}

	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(`
		INSERT INTO subscriptions (subscription_id, description, customer_id, product_id)
		VALUES ($1, $2, (SELECT id FROM customers WHERE customer_id = $3), (SELECT id FROM products WHERE product_id = $4))
		ON CONFLICT (subscription_id) DO NOTHING
	`)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	for _, s := range subs {
		if _, err := stmt.Exec(s.SubscriptionID, s.Description, s.CustomerID, s.ProductID); err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

// -------------------- InvoiceRepository --------------------

type InvoiceRepository struct {
	db *sql.DB
}

func NewInvoiceRepository(db *sql.DB) *InvoiceRepository {
	return &InvoiceRepository{db: db}
}

func (r *InvoiceRepository) SaveBatch(invoices []models.Invoice) error {
	if len(invoices) == 0 {
		return nil
	}

	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(`
		INSERT INTO invoices (invoice_number, customer_id)
		VALUES ($1, (SELECT id FROM customers WHERE customer_id = $2))
		ON CONFLICT (invoice_number) DO NOTHING
	`)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	for _, i := range invoices {
		if _, err := stmt.Exec(i.InvoiceNumber, i.CustomerID); err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

// -------------------- MeterRepository --------------------

type MeterRepository struct {
	db *sql.DB
}

func NewMeterRepository(db *sql.DB) *MeterRepository {
	return &MeterRepository{db: db}
}

func (r *MeterRepository) SaveBatch(meters []models.Meter) error {
	if len(meters) == 0 {
		return nil
	}

	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(`
		INSERT INTO meters (meter_id, meter_type, meter_category, meter_sub_category, meter_name, meter_region, invoice_id, product_id)
		VALUES ($1, $2, $3, $4, $5, $6, (SELECT id FROM invoices WHERE invoice_number = $7), (SELECT id FROM products WHERE product_id = $8))
		ON CONFLICT (meter_id) DO NOTHING
	`)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	for _, m := range meters {
		if _, err := stmt.Exec(m.MeterID, m.MeterType, m.MeterCategory, m.MeterSubCat, m.MeterName, m.MeterRegion, m.InvoiceID, m.ProductID); err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

// -------------------- BillingFinancialRepository --------------------

type BillingFinancialRepository struct {
	db *sql.DB
}

func NewBillingFinancialRepository(db *sql.DB) *BillingFinancialRepository {
	return &BillingFinancialRepository{db: db}
}

func (r *BillingFinancialRepository) SaveBatch(bfs []models.BillingFinancial) error {
	if len(bfs) == 0 {
		return nil
	}

	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(`
		INSERT INTO billing_financial (
			invoice_id, meter_id, subscription_id,
			unit_price, quantity, unit_type,
			billing_pre_tax_total, billing_currency,
			pricing_pre_tax_total, pricing_currency,
			service_info1, service_info2, tags, additional_info,
			effective_unit_price, pc_to_bc_exchange_rate, pc_to_bc_exchange_rate_date
		)
		VALUES (
			(SELECT id FROM invoices WHERE invoice_number = $1),
			(SELECT id FROM meters WHERE meter_id = $2),
			(SELECT id FROM subscriptions WHERE subscription_id = $3),
			$4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17
		)
		ON CONFLICT DO NOTHING
	`)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	for _, b := range bfs {
		if _, err := stmt.Exec(
			b.InvoiceID, b.MeterID, b.SubscriptionID,
			b.UnitPrice, b.Quantity, b.UnitType,
			b.BillingPreTaxTotal, b.BillingCurrency,
			b.PricingPreTaxTotal, b.PricingCurrency,
			b.ServiceInfo1, b.ServiceInfo2, b.Tags, b.AdditionalInfo,
			b.EffectiveUnitPrice, b.PCToBCExchangeRate, b.PCToBCExchangeRateDate,
		); err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

// -------------------- BillingRepository --------------------

type BillingRepository struct {
	db *sql.DB
}

func NewBillingRepository(db *sql.DB) *BillingRepository {
	return &BillingRepository{db: db}
}

func (r *BillingRepository) SaveBatch(entitlements []models.Billing) error {
	if len(entitlements) == 0 {
		return nil
	}

	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(`
		INSERT INTO billing (
			billing_financial_id, entitlement_id, entitlement_description,
			partner_earned_credit_percentage, credit_percentage, credit_type,
			benefit_order_id, benefit_id, benefit_type
		)
		VALUES (
			(SELECT id FROM billing_financial WHERE id = $1),
			$2, $3, $4, $5, $6, $7, $8, $9
		)
		ON CONFLICT DO NOTHING
	`)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	for _, e := range entitlements {
		if _, err := stmt.Exec(
			e.BillingFinancialID, e.EntitlementID, e.EntitlementDescription,
			e.PartnerEarnedCreditPerc, e.CreditPercentage, e.CreditType,
			e.BenefitOrderID, e.BenefitID, e.BenefitType,
		); err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}
