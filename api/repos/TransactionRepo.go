package repos

import (
	"database/sql"
	"fmt"
	"strings"

	"bitpos.com/m/models"
)




type TransactionRepo struct {
	db *sql.DB
}

func NewTransactionRepo(db *sql.DB) *TransactionRepo {
	return &TransactionRepo{db: db}
}

func (r *TransactionRepo) CreateTransaction(t models.Transaction) error {
	// Check if item exists in store's stock
	var exists bool
	err := r.db.QueryRow(`
		SELECT EXISTS(
			SELECT 1 FROM items_stocks 
			WHERE item_id = $1 AND store_id = $2
		)
	`, t.ItemID, t.StoreID).Scan(&exists)

	if err != nil {
		return fmt.Errorf("stock check failed: %w", err)
	}

	// If not exists, create with 0 quantity
	if !exists {
		_, err = r.db.Exec(`
			INSERT INTO items_stocks (item_id, store_id, quantity, selling_price)
			VALUES ($1, $2, 0, 0)
		`, t.ItemID, t.StoreID)
		if err != nil {
			return fmt.Errorf("failed to create stock record: %w", err)
		}
	}

	// Insert transaction
	_, err = r.db.Exec(`
		INSERT INTO stocks_movements 
		(item_id, store_id, type, category, quantity, status, created_by)
		VALUES ($1,$2,$3,$4,$5,$6,$7)
	`,
		t.ItemID, t.StoreID, t.Type, t.Category,
		t.Quantity, t.Status, t.CreatedBy,
	)

	if err != nil {
		return fmt.Errorf("transaction failed: %w", err)
	}
	return nil
}


// func (r *TransactionRepo) CompleteTransaction(txID int64, billRepo *BillRepo, billItemRepo *BillItemRepo, stockRepo *StockRepo) error {

// 	// 1. Fetch all transactions in this batch (same batch = same txID group)
// 	// If you batch by a parent ID, change this query. Otherwise single tx:
// 	rows, err := r.db.Query(`
// 		SELECT id, item_id, store_id, type, category, quantity, status, created_by
// 		FROM stocks_movements 
// 		WHERE id = $1
// 	`, txID)
// 	if err != nil {
// 		return fmt.Errorf("transaction not found: %w", err)
// 	}
// 	defer rows.Close()

// 	var transactions []models.Transaction
// 	for rows.Next() {
// 		var t models.Transaction
// 		err := rows.Scan(&t.ID, &t.ItemID, &t.StoreID, &t.Type, &t.Category, &t.Quantity, &t.Status, &t.CreatedBy)
// 		if err != nil {
// 			return fmt.Errorf("scan failed: %w", err)
// 		}
// 		if t.Status == "completed" {
// 			return fmt.Errorf("transaction %d already completed", t.ID)
// 		}
// 		transactions = append(transactions, t)
// 	}

// 	if len(transactions) == 0 {
// 		return fmt.Errorf("no transactions found")
// 	}

// 	first := transactions[0]

// 	// 2. Update stock for each item
// 	for _, t := range transactions {
// 		switch t.Type {
// 			case "positive":
// 				err = stockRepo.AddStock(t.ItemID, t.StoreID, t.Quantity)
// 			case "negative":
// 				err = stockRepo.ReduceStock(t.ItemID, t.StoreID, t.Quantity)
// 			default:
// 				return fmt.Errorf("unknown type: %s", t.Type)
// 			}
// 			if err != nil {
// 				return fmt.Errorf("stock update failed for item %d: %w", t.ItemID, err)
// 			}
// 	}

// 	// 3. Create bill if purchase or sale
// 	if first.Category == "purchase" || first.Category == "sale" {

// 		var totalAmount float64
// 		var billItems []models.BillItem

// 		for _, t := range transactions {
// 			var price float64

// 			if first.Category == "sale" {
// 				// Sale: use selling_price from item_stocks
// 				err = r.db.QueryRow(`
// 					SELECT selling_price FROM items_stocks 
// 					WHERE item_id = $1 AND store_id = $2
// 				`, t.ItemID, t.StoreID).Scan(&price)
// 			} else {
// 				// Purchase: use unit_cost from items
// 				err = r.db.QueryRow(`
// 					SELECT unit_cost FROM items WHERE id = $1
// 				`, t.ItemID).Scan(&price)
// 			}

// 			if err != nil {
// 				return fmt.Errorf("failed to get price for item %d: %w", t.ItemID, err)
// 			}

// 			lineTotal := price * t.Quantity
// 			totalAmount += lineTotal

// 			billItems = append(billItems, models.BillItem{
// 				ItemID:         t.ItemID,
// 				ItemPrice:      price,
// 				ItemQuantity:   t.Quantity,
// 				DiscountAmount: 0,
// 				DiscountRate:   0,
// 			})
// 		}

// 		// Create bill
// 		code := fmt.Sprintf("%s-%d-%d", strings.ToUpper(first.Category[:3]), first.StoreID, txID)
// 		amountPaid := 0.0
// 		if first.Category == "sale" {
// 			amountPaid = totalAmount // sale = paid upfront
// 		}

// 		billID, err := billRepo.CreateBill(models.Bill{
// 			Code:         code,
// 			Type:         first.Category,
// 			AmountBilled: totalAmount,
// 			AmountPaid:   amountPaid,
// 			CreatedBy:    first.CreatedBy,
// 		})
// 		if err != nil {
// 			return fmt.Errorf("bill creation failed: %w", err)
// 		}

// 		// Create bill items
// 		for _, bi := range billItems {
// 			bi.BillID = billID
// 			err = billItemRepo.CreateBillItem(bi)
// 			if err != nil {
// 				return fmt.Errorf("bill item failed: %w", err)
// 			}
// 		}
// 	}

// 	// 4. Mark all transactions as completed
// 	_, err = r.db.Exec(`
// 		UPDATE stocks_movements SET status = 'completed' WHERE id = $1
// 	`, txID)
// 	if err != nil {
// 		return fmt.Errorf("status update failed: %w", err)
// 	}

// 	return nil
// }


func (r *TransactionRepo) CompleteTransaction(txID int64, stockRepo *StockRepo) error {

	// 1. Fetch the transaction being completed
	var t models.Transaction
	err := r.db.QueryRow(`
		SELECT id, item_id, store_id, type, category, quantity, status
		FROM stocks_movements WHERE id = $1
	`, txID).Scan(&t.ID, &t.ItemID, &t.StoreID, &t.Type, &t.Category, &t.Quantity, &t.Status)

	if err != nil {
		return fmt.Errorf("transaction not found: %w", err)
	}

	if t.Status == "completed" {
		return fmt.Errorf("transaction already completed")
	}

	// 2. If category is "request", find the matching pair and complete both
	if t.Category == "request" {

		// Find the matching transaction:
		// Same item, same quantity, same category, opposite type, still pending
		oppositeType := "positive"
		if t.Type == "positive" {
			oppositeType = "negative"
		}

		var matchID int64
		var matchStoreID int64
		err = r.db.QueryRow(`
			SELECT id, store_id FROM stocks_movements
			WHERE item_id = $1
			  AND quantity = $2
			  AND category = 'request'
			  AND type = $3
			  AND status = 'pending'
			  AND id != $4
			ORDER BY id DESC
			LIMIT 1
		`, t.ItemID, t.Quantity, oppositeType, t.ID).Scan(&matchID, &matchStoreID)

		if err != nil {
			return fmt.Errorf("matching request not found: %w", err)
		}

		// 3. Check source has enough stock before proceeding
		// The negative one is the source — must have enough
		var negativeStoreID int64
		var positiveStoreID int64

		if t.Type == "negative" {
			negativeStoreID = t.StoreID
			positiveStoreID = matchStoreID
		} else {
			negativeStoreID = matchStoreID
			positiveStoreID = t.StoreID
		}

		// Verify source has enough stock
		var currentQty float64
		err = r.db.QueryRow(`
			SELECT COALESCE(quantity, 0) FROM items_stocks
			WHERE item_id = $1 AND store_id = $2
		`, t.ItemID, negativeStoreID).Scan(&currentQty)

		if err != nil {
			return fmt.Errorf("failed to check source stock: %w", err)
		}

		if currentQty < t.Quantity {
			return fmt.Errorf("insufficient stock in source store (available: %.0f, requested: %.0f)", currentQty, t.Quantity)
		}

		// 4. Deduct from source
		err = stockRepo.ReduceStock(t.ItemID, negativeStoreID, t.Quantity)
		if err != nil {
			return fmt.Errorf("failed to deduct from source: %w", err)
		}

		// 5. Add to destination
		err = stockRepo.AddStock(t.ItemID, positiveStoreID, t.Quantity)
		if err != nil {
			return fmt.Errorf("failed to add to destination: %w", err)
		}

		// 6. Mark both as completed
		_, err = r.db.Exec(`
			UPDATE stocks_movements SET status = 'completed' WHERE id IN ($1, $2)
		`, t.ID, matchID)
		if err != nil {
			return fmt.Errorf("failed to update status: %w", err)
		}

		return nil
	}

	// 7. Non-request transactions: handle normally
	switch t.Type {
		case "positive":
			err = stockRepo.AddStock(t.ItemID, t.StoreID, t.Quantity)
		case "negative":
			// Check stock before deducting
			var currentQty float64
			r.db.QueryRow(`
				SELECT COALESCE(quantity, 0) FROM items_stocks
				WHERE item_id = $1 AND store_id = $2
			`, t.ItemID, t.StoreID).Scan(&currentQty)

			if currentQty < t.Quantity {
				return fmt.Errorf("insufficient stock (available: %.0f, requested: %.0f)", currentQty, t.Quantity)
			}

			err = stockRepo.ReduceStock(t.ItemID, t.StoreID, t.Quantity)
		}

	if err != nil {
		return fmt.Errorf("stock update failed: %w", err)
	}

	_, err = r.db.Exec(`UPDATE stocks_movements SET status = 'completed' WHERE id = $1`, t.ID)
	if err != nil {
		return fmt.Errorf("status update failed: %w", err)
	}

	return nil
}



func (r *TransactionRepo) GetAllTransactions(
	search, store, typeFilter, statusFilter, categoryFilter string,
	limit, offset int,
) ([]models.Transaction, int, error) {

	baseQuery := `
		FROM stocks_movements sm
		INNER JOIN items i ON i.id = sm.item_id
		INNER JOIN users u ON u.id = sm.created_by
		INNER JOIN stores s ON s.id = sm.store_id
	`

	var conditions []string
	var args []interface{}
	argIndex := 1

	if search != "" {
		conditions = append(conditions, fmt.Sprintf(
			"(i.name ILIKE $%d OR u.name ILIKE $%d OR s.name ILIKE $%d)",
			argIndex, argIndex, argIndex,
		))
		args = append(args, "%"+search+"%")
		argIndex++
	}

	if store != "" {
		conditions = append(conditions, fmt.Sprintf("s.name = $%d", argIndex))
		args = append(args, store)
		argIndex++
	}

	if typeFilter != "" {
		conditions = append(conditions, fmt.Sprintf("sm.type = $%d", argIndex))
		args = append(args, typeFilter)
		argIndex++
	}

	if statusFilter != "" {
		conditions = append(conditions, fmt.Sprintf("sm.status = $%d", argIndex))
		args = append(args, statusFilter)
		argIndex++
	}

	if categoryFilter != "" {
		conditions = append(conditions, fmt.Sprintf("sm.category = $%d", argIndex))
		args = append(args, categoryFilter)
		argIndex++
	}

	if len(conditions) > 0 {
		baseQuery += " AND " + strings.Join(conditions, " AND ")
	}

	// Count
	var total int
	err := r.db.QueryRow("SELECT COUNT(*) "+baseQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("count failed: %w", err)
	}

	// Data
	selectQuery := `
		SELECT sm.id, sm.item_id, sm.store_id, sm.type, sm.category,
		       sm.quantity, sm.status, sm.created_by, sm.created_at,
		       i.name, u.name, s.name
	` + baseQuery + fmt.Sprintf(" ORDER BY sm.id DESC LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, limit, offset)

	rows, err := r.db.Query(selectQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	var transactions []models.Transaction
	for rows.Next() {
		var t models.Transaction
		err := rows.Scan(
			&t.ID, &t.ItemID, &t.StoreID, &t.Type, &t.Category,
			&t.Quantity, &t.Status, &t.CreatedBy,&t.CreatedAt,
			&t.ItemName, &t.Created_Username, &t.StoreName,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("scan failed: %w", err)
		}
		transactions = append(transactions, t)
	}

	return transactions, total, nil
}