package repos

import (
	"database/sql"
	"fmt"
)

type StockRepo struct {
	db *sql.DB
}

func NewStockRepo(db *sql.DB) *StockRepo {
	return &StockRepo{db: db}
}

// ADD STOCK (IN)
func (r *StockRepo) AddStock(itemID, storeID int64, qty float64) error {

	_, err := r.db.Exec(`
	UPDATE items_stocks 
	SET quantity = quantity + $1
	WHERE item_id=$2 AND store_id=$3
	`, qty, itemID, storeID)

	if err != nil {
		return err
	}

	return nil
}

// REDUCE STOCK (OUT)
func (r *StockRepo) ReduceStock(itemID, storeID int64, qty float64) error {

	var current float64

	err := r.db.QueryRow(`
	SELECT quantity FROM items_stocks
	WHERE item_id=$1 AND store_id=$2
	`, itemID, storeID).Scan(&current)

	if err != nil {
		return err
	}

	if current < qty {
		return fmt.Errorf("insufficient stock")
	}

	_, err = r.db.Exec(`
	UPDATE items_stocks 
	SET quantity = quantity - $1
	WHERE item_id=$2 AND store_id=$3
	`, qty, itemID, storeID)

	return err
}