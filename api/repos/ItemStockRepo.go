package repos

import (
	"database/sql"
	"fmt"

	"bitpos.com/m/models"
)

type ItemStockRepo struct {
	db *sql.DB
}

func NewItemStockRepo(db *sql.DB) *ItemStockRepo {
	return &ItemStockRepo{db: db}
}

func (r *ItemStockRepo) CreateStock(s models.ItemStock) error {
	query := `
	INSERT INTO item_stocks 
	(item_id, store_id, batch, quantity, selling_price, discount_amount, discount_rate, expiry_date)
	VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
	`
	_, err := r.db.Exec(query,
		s.ItemID,
		s.StoreID,
		s.Batch,
		s.Quantity,
		s.SellingPrice,
		s.DiscountAmount,
		s.DiscountRate,
		s.ExpiryDate,
	)

	if err != nil {
		return fmt.Errorf("create stock failed: %w", err)
	}
	return nil
}

func (r *ItemStockRepo) UpdateQuantity(id int64, qty float64) error {
	_, err := r.db.Exec(`UPDATE item_stocks SET quantity=$1 WHERE id=$2`, qty, id)
	return err
}

