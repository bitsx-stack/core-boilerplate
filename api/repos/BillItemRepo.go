package repos

import (
	"database/sql"
	"fmt"

	"bitpos.com/m/models"
)

type BillItemRepo struct {
	db *sql.DB
}

func NewBillItemRepo(db *sql.DB) *BillItemRepo {
	return &BillItemRepo{db: db}
}

func (r *BillItemRepo) CreateBillItem(billID int64, item models.BillItem) error {
	_, err := r.db.Exec(`
		INSERT INTO bill_items (bill_id, item_id, item_price, item_quantity, discount_amount, discount_rate)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, billID, item.ItemID, item.ItemPrice, item.ItemQuantity, item.DiscountAmount, item.DiscountRate)

	if err != nil {
		return fmt.Errorf("bill item insert failed: %w", err)
	}
	return nil
}

func (r *BillItemRepo) GetBillItems(billID int64) ([]models.BillItem, error) {
	rows, err := r.db.Query(`
		SELECT bi.id, bi.bill_id, bi.item_id, COALESCE(i.name, ''),
		       bi.item_price, bi.item_quantity, bi.discount_amount, bi.discount_rate
		FROM bill_items bi
		LEFT JOIN items i ON i.id = bi.item_id
		WHERE bi.bill_id = $1
	`, billID)

	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	var items []models.BillItem
	for rows.Next() {
		var item models.BillItem
		err := rows.Scan(&item.ID, &item.BillID, &item.ItemID, &item.ItemName,
			&item.ItemPrice, &item.ItemQuantity, &item.DiscountAmount, &item.DiscountRate)
		if err != nil {
			return nil, fmt.Errorf("scan failed: %w", err)
		}
		item.LineTotal = (item.ItemPrice * item.ItemQuantity) - item.DiscountAmount
		items = append(items, item)
	}

	return items, nil
}

func (r *BillItemRepo) UpdateBillItem(id int64, item models.BillItem) error {
	result, err := r.db.Exec(`
		UPDATE bill_items 
		SET item_price = $1, item_quantity = $2, discount_amount = $3, discount_rate = $4
		WHERE id = $5
	`, item.ItemPrice, item.ItemQuantity, item.DiscountAmount, item.DiscountRate, id)

	if err != nil {
		return fmt.Errorf("update failed: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("bill item not found")
	}

	// Recalculate bill total
	err = r.recalcBillTotal(item.BillID)
	if err != nil {
		return fmt.Errorf("recalc failed: %w", err)
	}

	return nil
}

func (r *BillItemRepo) DeleteBillItem(id int64) error {
	// Get bill_id before deleting
	var billID int64
	err := r.db.QueryRow(`SELECT bill_id FROM bill_items WHERE id = $1`, id).Scan(&billID)
	if err != nil {
		return fmt.Errorf("bill item not found: %w", err)
	}

	_, err = r.db.Exec(`DELETE FROM bill_items WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete failed: %w", err)
	}

	return r.recalcBillTotal(billID)
}

func (r *BillItemRepo) recalcBillTotal(billID int64) error {
	var total float64
	err := r.db.QueryRow(`
		SELECT COALESCE(SUM((item_price * item_quantity) - discount_amount), 0)
		FROM bill_items WHERE bill_id = $1
	`, billID).Scan(&total)

	if err != nil {
		return fmt.Errorf("sum failed: %w", err)
	}

	_, err = r.db.Exec(`UPDATE bills SET amount_billed = $1 WHERE id = $2`, total, billID)
	return err
}