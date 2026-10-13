package repos

import (
	"database/sql"

	"bitpos.com/m/models"
)

type SalesRepo struct {
	db *sql.DB
}

func NewSalesRepo(db *sql.DB) *SalesRepo {
	return &SalesRepo{db: db}
}

func (r *SalesRepo) CreateSale(
	bill models.Bill,
	items []models.BillItem,
	payments []models.Payment,
) error {

	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	var billID int64

	err = tx.QueryRow(`
	INSERT INTO bills (code,type,amount_billed,amount_paid,created_by)
	VALUES ($1,$2,$3,$4,$5) RETURNING id
	`,
		bill.Code,
		bill.Type,
		bill.AmountBilled,
		bill.AmountPaid,
		bill.CreatedBy,
	).Scan(&billID)

	if err != nil {
		tx.Rollback()
		return err
	}

	for _, item := range items {

		_, err = tx.Exec(`
		INSERT INTO bill_items 
		(bill_id,item_id,item_price,item_quantity)
		VALUES ($1,$2,$3,$4)
		`,
			billID,
			item.ItemID,
			item.ItemPrice,
			item.ItemQuantity,
		)

		if err != nil {
			tx.Rollback()
			return err
		}

		// reduce stock
		_, err = tx.Exec(`
		UPDATE item_stocks 
		SET quantity = quantity - $1
		WHERE item_id=$2
		`,
			item.ItemQuantity,
			item.ItemID,
		)

		if err != nil {
			tx.Rollback()
			return err
		}
	}

	for _, p := range payments {
		_, err = tx.Exec(`
		INSERT INTO payments (bill_id,type,payment_amount,status,created_by)
		VALUES ($1,$2,$3,'completed',$4)
		`,
			billID,
			p.Type,
			p.PaymentAmount,
			p.CreatedBy,
		)

		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}