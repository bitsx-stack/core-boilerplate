package repos

import (
	"database/sql"
	"fmt"
	"strings"

	"bitpos.com/m/models"
)

// ══════════════════════════════════════════════
// BILL REPO
// ══════════════════════════════════════════════

type BillRepo struct {
	db *sql.DB
}

func NewBillRepo(db *sql.DB) *BillRepo {
	return &BillRepo{db: db}
}

func (r *BillRepo) CreateBill(b models.Bill, items []models.BillItem, stockRepo *StockRepo) (int64, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, fmt.Errorf("tx begin failed: %w", err)
	}
	defer tx.Rollback()

	// Calculate total from items
	var totalBilled float64
	for _, item := range items {
		lineTotal := (item.ItemPrice * item.ItemQuantity) - item.DiscountAmount
		totalBilled += lineTotal
	}

	// For SALES: check stock and deduct immediately
	if b.Type == "sale" {
		for _, item := range items {
			// Check if item is a product (not service)
			var itemType string
			err := tx.QueryRow(`SELECT COALESCE(type, 'product') FROM items WHERE id = $1`, item.ItemID).Scan(&itemType)
			if err != nil {
				return 0, fmt.Errorf("failed to check item type: %w", err)
			}

			// Skip stock check for services
			if strings.ToLower(itemType) == "service" {
				continue
			}

			// Check stock availability
			var currentQty float64
			err = tx.QueryRow(`
				SELECT COALESCE(quantity, 0) FROM items_stocks
				WHERE item_id = $1 AND store_id = $2
			`, item.ItemID, b.StoreID).Scan(&currentQty)

			if err != nil {
				return 0, fmt.Errorf("stock check failed for item %d: %w", item.ItemID, err)
			}

			if currentQty < item.ItemQuantity {
				// Get item name for better error
				var itemName string
				tx.QueryRow(`SELECT name FROM items WHERE id = $1`, item.ItemID).Scan(&itemName)
				return 0, fmt.Errorf("insufficient stock for '%s' (available: %.0f, requested: %.0f)",
					itemName, currentQty, item.ItemQuantity)
			}

			// Deduct stock
			_, err = tx.Exec(`
				UPDATE items_stocks SET quantity = quantity - $1
				WHERE item_id = $2 AND store_id = $3
			`, item.ItemQuantity, item.ItemID, b.StoreID)
			if err != nil {
				return 0, fmt.Errorf("stock deduct failed: %w", err)
			}


			// Record movement
			_, err = tx.Exec(`
				INSERT INTO stocks_movements (item_id, store_id, type, category, quantity, status, created_by)
				VALUES ($1, $2, 'negative', 'sale', $3, 'completed', $4)
			`, item.ItemID, b.StoreID, item.ItemQuantity, b.CreatedBy)
			if err != nil {
				return 0, fmt.Errorf("movement record failed: %w", err)
			}
		}
	}

	// Insert bill
	var billID int64
	err = tx.QueryRow(`
		INSERT INTO bills (code, people_id, store_id, type, status, amount_billed, amount_paid, created_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id
	`, b.Code, b.PeopleID, b.StoreID, b.Type, "pending", totalBilled, 0, b.CreatedBy).Scan(&billID)

	if err != nil {
		return 0, fmt.Errorf("bill insert failed: %w", err)
	}

	// Insert bill items
	for _, item := range items {
		_, err = tx.Exec(`
			INSERT INTO bill_items (bill_id, item_id, item_price, item_quantity, discount_amount, discount_rate)
			VALUES ($1, $2, $3, $4, $5, $6)
		`, billID, item.ItemID, item.ItemPrice, item.ItemQuantity, item.DiscountAmount, item.DiscountRate)

		if err != nil {
			return 0, fmt.Errorf("bill item insert failed: %w", err)
		}
	}

	if err = tx.Commit(); err != nil {
		return 0, fmt.Errorf("tx commit failed: %w", err)
	}

	return billID, nil
}


func (r *BillRepo) GetAllBills(search, billType, status, store string, peopleID int64, limit, offset int) ([]models.Bill, int, error) {
	baseQuery := `
		FROM bills b
		LEFT JOIN people p ON p.id = b.people_id
		LEFT JOIN stores s ON s.id = b.store_id
		LEFT JOIN users u ON u.id = b.created_by
		WHERE b.deleted_at IS NULL
	`
	var conditions []string
	var args []interface{}
	argIndex := 1

	if search != "" {
		conditions = append(conditions, fmt.Sprintf(
			"(b.code ILIKE $%d OR p.name ILIKE $%d OR s.name ILIKE $%d)", argIndex, argIndex, argIndex))
		args = append(args, "%"+search+"%")
		argIndex++
	}
	if billType != "" {
		conditions = append(conditions, fmt.Sprintf("b.type = $%d", argIndex))
		args = append(args, billType)
		argIndex++
	}
	if status != "" {
		conditions = append(conditions, fmt.Sprintf("b.status = $%d", argIndex))
		args = append(args, status)
		argIndex++
	}
	if store != "" {
		conditions = append(conditions, fmt.Sprintf("s.name = $%d", argIndex))
		args = append(args, store)
		argIndex++
	}
	if peopleID > 0 {
		conditions = append(conditions, fmt.Sprintf("b.people_id = $%d", argIndex))
		args = append(args, peopleID)
		argIndex++
	}

	if len(conditions) > 0 {
		baseQuery += " AND " + strings.Join(conditions, " AND ")
	}

	var total int
	err := r.db.QueryRow("SELECT COUNT(*) "+baseQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("count failed: %w", err)
	}

	selectQuery := fmt.Sprintf(`
		SELECT b.id, b.code, b.people_id, COALESCE(p.name,''),
		       COALESCE(b.store_id, 0), COALESCE(s.name,''),
		       b.type, b.status, b.amount_billed, b.amount_paid,
		       b.created_at, b.created_by, COALESCE(u.name,'')
		%s ORDER BY b.created_at DESC LIMIT $%d OFFSET $%d
	`, baseQuery, argIndex, argIndex+1)
	args = append(args, limit, offset)

	rows, err := r.db.Query(selectQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	var bills []models.Bill
	for rows.Next() {
		var b models.Bill
		err := rows.Scan(&b.ID, &b.Code, &b.PeopleID, &b.PeopleName,
			&b.StoreID, &b.StoreName,
			&b.Type, &b.Status, &b.AmountBilled, &b.AmountPaid,
			&b.CreatedAt, &b.CreatedBy, &b.CreatedName)
		if err != nil {
			return nil, 0, fmt.Errorf("scan failed: %w", err)
		}
		bills = append(bills, b)
	}

	return bills, total, nil
}



func (r *BillRepo) GetSingleBill(id int64) (*models.Bill, error) {
	var b models.Bill
	err := r.db.QueryRow(`
		SELECT b.id, b.code, b.people_id, COALESCE(p.name,''),
		       COALESCE(b.store_id, 0), COALESCE(s.name,''),
		       b.type, b.status, b.amount_billed, b.amount_paid,
		       b.created_at, b.created_by, COALESCE(u.name,'')
		FROM bills b
		LEFT JOIN people p ON p.id = b.people_id
		LEFT JOIN stores s ON s.id = b.store_id
		LEFT JOIN users u ON u.id = b.created_by
		WHERE b.id = $1 AND b.deleted_at IS NULL
	`, id).Scan(&b.ID, &b.Code, &b.PeopleID, &b.PeopleName,
		&b.StoreID, &b.StoreName,
		&b.Type, &b.Status, &b.AmountBilled, &b.AmountPaid,
		&b.CreatedAt, &b.CreatedBy, &b.CreatedName)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("bill not found")
	}
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}

	// Get items
	itemRows, err := r.db.Query(`
		SELECT bi.id, bi.bill_id, bi.item_id, COALESCE(i.name,''), bi.item_price,
		       bi.item_quantity, bi.discount_amount, bi.discount_rate
		FROM bill_items bi
		LEFT JOIN items i ON i.id = bi.item_id
		WHERE bi.bill_id = $1
	`, id)
	if err == nil {
		defer itemRows.Close()
		for itemRows.Next() {
			var item models.BillItem
			itemRows.Scan(&item.ID, &item.BillID, &item.ItemID, &item.ItemName,
				&item.ItemPrice, &item.ItemQuantity, &item.DiscountAmount, &item.DiscountRate)
			item.LineTotal = (item.ItemPrice * item.ItemQuantity) - item.DiscountAmount
			b.Items = append(b.Items, item)
		}
	}

	// Get payments
	payRows, err := r.db.Query(`
		SELECT py.id, py.bill_id, py.people_id, COALESCE(p.name,''), py.type, COALESCE(py.category,''),
		       py.payment_amount, py.status, py.created_at, py.created_by, COALESCE(u.name,'')
		FROM payments py
		LEFT JOIN people p ON p.id = py.people_id
		LEFT JOIN users u ON u.id = py.created_by
		WHERE py.bill_id = $1
		ORDER BY py.created_at DESC
	`, id)
	if err == nil {
		defer payRows.Close()
		for payRows.Next() {
			var pay models.Payment
			payRows.Scan(&pay.ID, &pay.BillID, &pay.PeopleID, &pay.PeopleName, &pay.Type,
				&pay.Category, &pay.PaymentAmount, &pay.Status, &pay.CreatedAt, &pay.CreatedBy, &pay.CreatedName)
			b.Payments = append(b.Payments, pay)
		}
	}

	return &b, nil
}

func (r *BillRepo) UpdateBillItems(billID int64, items []models.BillItem) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("tx begin failed: %w", err)
	}
	defer tx.Rollback()

	// Delete old items
	_, err = tx.Exec("DELETE FROM bill_items WHERE bill_id = $1", billID)
	if err != nil {
		return fmt.Errorf("delete old items failed: %w", err)
	}

	// Insert new items and calculate total
	var totalBilled float64
	for _, item := range items {
		lineTotal := (item.ItemPrice * item.ItemQuantity) - item.DiscountAmount
		totalBilled += lineTotal

		_, err = tx.Exec(`
			INSERT INTO bill_items (bill_id, item_id, item_price, item_quantity, discount_amount, discount_rate)
			VALUES ($1, $2, $3, $4, $5, $6)
		`, billID, item.ItemID, item.ItemPrice, item.ItemQuantity, item.DiscountAmount, item.DiscountRate)

		if err != nil {
			return fmt.Errorf("item insert failed: %w", err)
		}
	}

	// Update bill total
	_, err = tx.Exec("UPDATE bills SET amount_billed = $1 WHERE id = $2", totalBilled, billID)
	if err != nil {
		return fmt.Errorf("update total failed: %w", err)
	}

	return tx.Commit()
}

// ══════════════════════════════════════════════
// PAYMENT REPO
// ══════════════════════════════════════════════

type PaymentRepo struct {
	db *sql.DB
}

func NewPaymentRepo(db *sql.DB) *PaymentRepo {
	return &PaymentRepo{db: db}
}

func (r *PaymentRepo) CreatePayment(p models.Payment, billRepo *BillRepo) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("tx begin failed: %w", err)
	}
	defer tx.Rollback()

	// Get bill info
	var amountBilled, amountPaid float64
	var billStatus, billType string
	var storeID int64
	err = tx.QueryRow(`
		SELECT amount_billed, amount_paid, status, type, COALESCE(store_id, 0)
		FROM bills WHERE id = $1 AND deleted_at IS NULL
	`, p.BillID).Scan(&amountBilled, &amountPaid, &billStatus, &billType, &storeID)

	if err == sql.ErrNoRows {
		return fmt.Errorf("bill not found")
	}
	if err != nil {
		return fmt.Errorf("bill query failed: %w", err)
	}

	if billStatus == "paid" {
		return fmt.Errorf("bill is already fully paid")
	}

	remaining := amountBilled - amountPaid
	if p.PaymentAmount > remaining {
		return fmt.Errorf("payment (%.2f) exceeds balance (%.2f)", p.PaymentAmount, remaining)
	}

	// Insert payment
	_, err = tx.Exec(`
		INSERT INTO payments (bill_id, people_id, type, category, payment_amount, status, created_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, p.BillID, p.PeopleID, p.Type, p.Category, p.PaymentAmount, "completed", p.CreatedBy)
	if err != nil {
		return fmt.Errorf("payment insert failed: %w", err)
	}

	// Update bill
	newPaid := amountPaid + p.PaymentAmount
	newStatus := "partial"
	if newPaid >= amountBilled {
		newStatus = "paid"
	}

	_, err = tx.Exec(`
		UPDATE bills SET amount_paid = $1, status = $2 WHERE id = $3
	`, newPaid, newStatus, p.BillID)
	if err != nil {
		return fmt.Errorf("bill update failed: %w", err)
	}

	// PURCHASE: add stock on first payment (check if stock already added)
	if billType == "purchase" {
		// Check if this is the first payment
		var paymentCount int
		tx.QueryRow(`SELECT COUNT(*) FROM payments WHERE bill_id = $1`, p.BillID).Scan(&paymentCount)

		// paymentCount == 1 means we just inserted the first one above
		// In CreatePayment, replace the purchase stock section with:

	if paymentCount <= 1 && storeID > 0 {
    // First: collect all items into a slice
    itemRows, err := tx.Query(`
        SELECT bi.item_id, bi.item_quantity, COALESCE(i.type, 'product') as item_type
        FROM bill_items bi
        INNER JOIN items i ON i.id = bi.item_id
        WHERE bi.bill_id = $1
    `, p.BillID)
    if err != nil {
        return fmt.Errorf("failed to get bill items: %w", err)
    }

    type stockItem struct {
        itemID   int64
        qty      float64
        itemType string
    }
    var stockItems []stockItem

    for itemRows.Next() {
        var si stockItem
        if err := itemRows.Scan(&si.itemID, &si.qty, &si.itemType); err != nil {
            itemRows.Close()
            return fmt.Errorf("scan failed: %w", err)
        }
        stockItems = append(stockItems, si)
    }
    itemRows.Close() // Close BEFORE running more queries

    // Second: update stock for each item
    for _, si := range stockItems {
        if strings.ToLower(si.itemType) == "service" {
            continue
        }

        // Ensure items_stocks record exists
        var exists bool
        err = tx.QueryRow(`
            SELECT EXISTS(SELECT 1 FROM items_stocks WHERE item_id = $1 AND store_id = $2)
        `, si.itemID, storeID).Scan(&exists)
        if err != nil {
            return fmt.Errorf("stock check failed: %w", err)
        }

        if !exists {
            _, err = tx.Exec(`
                INSERT INTO items_stocks (item_id, store_id, quantity, selling_price)
                VALUES ($1, $2, 0, 0)
            `, si.itemID, storeID)
            if err != nil {
                return fmt.Errorf("failed to create stock: %w", err)
            }
        }

        _, err = tx.Exec(`
            UPDATE items_stocks SET quantity = quantity + $1
            WHERE item_id = $2 AND store_id = $3
        `, si.qty, si.itemID, storeID)
        if err != nil {
            return fmt.Errorf("stock add failed: %w", err)
        }

		// Record movement
		_, err = tx.Exec(`
			INSERT INTO stocks_movements (item_id, store_id, type, category, quantity, status, created_by)
			VALUES ($1, $2, 'positive', 'purchase', $3, 'completed', $4)
		`, si.itemID, storeID, si.qty, p.CreatedBy)
		if err != nil {
			return fmt.Errorf("movement record failed: %w", err)
		}
    }
	}
	}

	return tx.Commit()
}

func (r *PaymentRepo) GetPaymentsByBill(billID int64) ([]models.Payment, error) {
	rows, err := r.db.Query(`
		SELECT py.id, py.bill_id, py.people_id, COALESCE(p.name,''), py.type,
		       COALESCE(py.category,''), py.payment_amount, py.status, py.created_at,
		       py.created_by, COALESCE(u.name,'')
		FROM payments py
		LEFT JOIN people p ON p.id = py.people_id
		LEFT JOIN users u ON u.id = py.created_by
		WHERE py.bill_id = $1
		ORDER BY py.created_at DESC
	`, billID)

	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	var payments []models.Payment
	for rows.Next() {
		var pay models.Payment
		err := rows.Scan(&pay.ID, &pay.BillID, &pay.PeopleID, &pay.PeopleName, &pay.Type,
			&pay.Category, &pay.PaymentAmount, &pay.Status, &pay.CreatedAt, &pay.CreatedBy, &pay.CreatedName)
		if err != nil {
			return nil, fmt.Errorf("scan failed: %w", err)
		}
		payments = append(payments, pay)
	}

	return payments, nil
}