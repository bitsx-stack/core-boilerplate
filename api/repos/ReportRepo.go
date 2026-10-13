package repos

import (
	"database/sql"
	"fmt"
)

type ReportRepo struct {
	db *sql.DB
}

func NewReportRepo(db *sql.DB) *ReportRepo {
	return &ReportRepo{db: db}
}

type StockReportRow struct {
	ItemID        int64   `json:"item_id"`
	ItemName      string  `json:"item_name"`
	StoreName     string  `json:"store_name"`
	OpenStock     float64 `json:"open_stock"`
	PurchaseStock float64 `json:"purchase_stock"`
	TransferStock float64 `json:"transfer_stock"`
	AdjustStock   float64 `json:"adjust_stock"`
	SaleStock     float64 `json:"sale_stock"`
	CloseStock    float64 `json:"close_stock"`
	PurchaseAmt   float64 `json:"purchase_amt"`
	SaleAmt       float64 `json:"sale_amt"`
	BankAmt       float64 `json:"bank_amt"`
	CashAmt       float64 `json:"cash_amt"`
	MobileAmt     float64 `json:"mobile_amt"`
	Balance       float64 `json:"balance"`
	CostPrice     float64 `json:"cost_price"`
	SellPrice     float64 `json:"sell_price"`
}

func (r *ReportRepo) GetStockReport(itemID, storeID int64, startDate, endDate string) ([]StockReportRow, error) {

	// Base: every item+store combo that has stock or movements
	query := `
		WITH params AS (
			SELECT
				$1::bigint AS p_item,
				$2::bigint AS p_store,
				CASE WHEN $3 = '' THEN NULL ELSE $3::date END AS p_start,
        		CASE WHEN $4 = '' THEN NULL ELSE $4::date END AS p_end
		),
		combos AS (
			SELECT DISTINCT sm.item_id, sm.store_id
			FROM stocks_movements sm, params p
			WHERE (p.p_item = 0 OR sm.item_id = p.p_item)
			  AND (p.p_store = 0 OR sm.store_id = p.p_store)
			  AND sm.status = 'completed'
		),
		movements AS (
			SELECT sm.item_id, sm.store_id, sm.category, sm.type, sm.quantity, sm.created_at
			FROM stocks_movements sm, params p
			WHERE sm.status = 'completed'
			  AND (p.p_item = 0 OR sm.item_id = p.p_item)
			  AND (p.p_store = 0 OR sm.store_id = p.p_store)
		),
		open_stock AS (
				SELECT m.item_id, m.store_id,
					COALESCE(SUM(CASE WHEN m.type='positive' THEN m.quantity ELSE -m.quantity END), 0) AS qty
				FROM movements m, params p
				WHERE p.p_start IS NOT NULL AND m.created_at < p.p_start::timestamp
				GROUP BY m.item_id, m.store_id
			),
		period_moves AS (
					SELECT m.item_id, m.store_id,
						COALESCE(SUM(CASE WHEN m.category='purchase' AND m.type='positive' THEN m.quantity END), 0) AS purchase_in,
						COALESCE(SUM(CASE WHEN m.category IN ('transfer','request') THEN
							CASE WHEN m.type='positive' THEN m.quantity ELSE -m.quantity END END), 0) AS transfer_net,
						COALESCE(SUM(CASE WHEN m.category='adjustment' THEN
							CASE WHEN m.type='positive' THEN m.quantity ELSE -m.quantity END END), 0) AS adjust_net,
						COALESCE(SUM(CASE WHEN m.category='sale' AND m.type='negative' THEN m.quantity END), 0) AS sale_out
					FROM movements m, params p
					WHERE (p.p_start IS NULL OR m.created_at >= p.p_start::timestamp)
					AND (p.p_end IS NULL OR m.created_at < (p.p_end + INTERVAL '1 day')::timestamp)
					GROUP BY m.item_id, m.store_id
				),
		bill_data AS (
    SELECT bi.item_id, b.store_id,
        COALESCE(SUM(CASE WHEN b.type='purchase' THEN bi.item_price * bi.item_quantity - bi.discount_amount END), 0) AS purchase_amt,
        COALESCE(SUM(CASE WHEN b.type='sale' THEN bi.item_price * bi.item_quantity - bi.discount_amount END), 0) AS sale_amt
    FROM bill_items bi
    INNER JOIN bills b ON b.id = bi.bill_id
    CROSS JOIN params p
    WHERE b.deleted_at IS NULL
      AND (p.p_item = 0 OR bi.item_id = p.p_item)
      AND (p.p_store = 0 OR b.store_id = p.p_store)
      AND (p.p_start IS NULL OR b.created_at >= p.p_start::timestamp)
      AND (p.p_end IS NULL OR b.created_at < (p.p_end + INTERVAL '1 day')::timestamp)
    GROUP BY bi.item_id, b.store_id
),
pay_data AS (
    SELECT bi.item_id, b.store_id,
        COALESCE(SUM(CASE WHEN py.type='BANK' THEN py.payment_amount * (bi.item_price * bi.item_quantity - bi.discount_amount) / NULLIF(b.amount_billed, 0) END), 0) AS bank_amt,
        COALESCE(SUM(CASE WHEN py.type='CASH' THEN py.payment_amount * (bi.item_price * bi.item_quantity - bi.discount_amount) / NULLIF(b.amount_billed, 0) END), 0) AS cash_amt,
        COALESCE(SUM(CASE WHEN py.type='MOBILE' THEN py.payment_amount * (bi.item_price * bi.item_quantity - bi.discount_amount) / NULLIF(b.amount_billed, 0) END), 0) AS mobile_amt
    FROM payments py
    INNER JOIN bills b ON b.id = py.bill_id
    INNER JOIN bill_items bi ON bi.bill_id = b.id
    CROSS JOIN params p
    WHERE b.deleted_at IS NULL AND b.type = 'sale'
      AND (p.p_item = 0 OR bi.item_id = p.p_item)
      AND (p.p_store = 0 OR b.store_id = p.p_store)
      AND (p.p_start IS NULL OR py.created_at >= p.p_start::timestamp)
      AND (p.p_end IS NULL OR py.created_at < (p.p_end + INTERVAL '1 day')::timestamp)
    GROUP BY bi.item_id, b.store_id
)
		SELECT
			c.item_id, i.name, COALESCE(s.name, ''),
			COALESCE(os.qty, 0),
			COALESCE(pm.purchase_in, 0),
			COALESCE(pm.transfer_net, 0),
			COALESCE(pm.adjust_net, 0),
			COALESCE(pm.sale_out, 0),
			COALESCE(os.qty, 0) + COALESCE(pm.purchase_in, 0) + COALESCE(pm.transfer_net, 0) + COALESCE(pm.adjust_net, 0) - COALESCE(pm.sale_out, 0),
			COALESCE(bd.purchase_amt, 0),
			COALESCE(bd.sale_amt, 0),
			COALESCE(pd.bank_amt, 0),
			COALESCE(pd.cash_amt, 0),
			COALESCE(pd.mobile_amt, 0),
			COALESCE(bd.sale_amt, 0) - COALESCE(pd.bank_amt, 0) - COALESCE(pd.cash_amt, 0) - COALESCE(pd.mobile_amt, 0),
			COALESCE(i.unit_cost, 0),
			COALESCE(ist.selling_price, 0)
		FROM combos c
		INNER JOIN items i ON i.id = c.item_id
		LEFT JOIN stores s ON s.id = c.store_id
		LEFT JOIN open_stock os ON os.item_id = c.item_id AND os.store_id = c.store_id
		LEFT JOIN period_moves pm ON pm.item_id = c.item_id AND pm.store_id = c.store_id
		LEFT JOIN bill_data bd ON bd.item_id = c.item_id AND bd.store_id = c.store_id
		LEFT JOIN pay_data pd ON pd.item_id = c.item_id AND pd.store_id = c.store_id
		LEFT JOIN items_stocks ist ON ist.item_id = c.item_id AND ist.store_id = c.store_id
		ORDER BY i.name, s.name
	`

	rows, err := r.db.Query(query, itemID, storeID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("report query failed: %w", err)
	}
	defer rows.Close()

	var results []StockReportRow
	for rows.Next() {
		var row StockReportRow
		err := rows.Scan(
			&row.ItemID, &row.ItemName, &row.StoreName,
			&row.OpenStock, &row.PurchaseStock, &row.TransferStock, &row.AdjustStock, &row.SaleStock, &row.CloseStock,
			&row.PurchaseAmt, &row.SaleAmt,
			&row.BankAmt, &row.CashAmt, &row.MobileAmt, &row.Balance,
			&row.CostPrice, &row.SellPrice,
		)
		if err != nil {
			return nil, fmt.Errorf("scan failed: %w", err)
		}
		results = append(results, row)
	}

	if results == nil {
		results = []StockReportRow{}
	}

	return results, nil
}