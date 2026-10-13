package repos

import (
	"database/sql"
	"strconv"
	"strings"

	"bitpos.com/m/models"
)

type BulkRepo struct {
	db *sql.DB
}

func NewBulkRepo(db *sql.DB) *BulkRepo {
	return &BulkRepo{db: db}
}

// ITEMS
func (r *BulkRepo) BulkInsertItems(rows [][]string, userID int64) models.BulkResult {

	result := models.BulkResult{}

	for i, row := range rows {
		if i == 0 {
			continue
		}

		// VALIDATION
		if len(row) < 5 {
			result.Failed++
			result.Errors = append(result.Errors, models.BulkError{i, "invalid column count"})
			continue
		}

		name := strings.TrimSpace(row[0])
		itemType := strings.TrimSpace(row[1])
		category := strings.TrimSpace(row[2])
		uom := strings.TrimSpace(row[3])
		costStr := strings.TrimSpace(row[4])

		if name == "" || itemType == "" || category == "" || uom == "" {
			result.Failed++
			result.Errors = append(result.Errors, models.BulkError{i, "missing required fields"})
			continue
		}

		cost, err := strconv.ParseFloat(costStr, 64)
		if err != nil {
			result.Failed++
			result.Errors = append(result.Errors, models.BulkError{i, "invalid unit_cost"})
			continue
		}

		_, err = r.db.Exec(`
		INSERT INTO items (name,type,category,uom,unit_cost,created_by)
		VALUES ($1,$2,$3,$4,$5,$6)
		ON CONFLICT (name) DO NOTHING
		`,
			strings.ToUpper(name),
			itemType,
			category,
			uom,
			cost,
			userID,
		)

		if err != nil {
			result.Failed++
			result.Errors = append(result.Errors, models.BulkError{Row: i, Error: err.Error()})
			continue
		}

		result.Success++
	}

	return result
}

// STOCK
func (r *BulkRepo) BulkInsertStock(rows [][]string) models.BulkResult {

	result := models.BulkResult{}

	for i, row := range rows {
		if i == 0 {
			continue
		}

		if len(row) < 4 {
			result.Failed++
			result.Errors = append(result.Errors, models.BulkError{i, "invalid column count"})
			continue
		}

		itemName := strings.TrimSpace(row[0])
		storeStr := strings.TrimSpace(row[1])
		qtyStr := strings.TrimSpace(row[2])
		priceStr := strings.TrimSpace(row[3])

		if itemName == "" {
			result.Failed++
			result.Errors = append(result.Errors, models.BulkError{i, "item_name required"})
			continue
		}

		storeID, err := strconv.ParseInt(storeStr, 10, 64)
		if err != nil {
			result.Failed++
			result.Errors = append(result.Errors, models.BulkError{i, "invalid store_id"})
			continue
		}

		qty, err := strconv.ParseFloat(qtyStr, 64)
		if err != nil {
			result.Failed++
			result.Errors = append(result.Errors, models.BulkError{i, "invalid quantity"})
			continue
		}

		price, err := strconv.ParseFloat(priceStr, 64)
		if err != nil {
			result.Failed++
			result.Errors = append(result.Errors, models.BulkError{i, "invalid selling_price"})
			continue
		}

		// find item
		var itemID int64
		err = r.db.QueryRow(`SELECT id FROM items WHERE name=$1`, itemName).Scan(&itemID)
		if err != nil {
			result.Failed++
			result.Errors = append(result.Errors, models.BulkError{i, "item not found"})
			continue
		}

		_, err = r.db.Exec(`
		INSERT INTO items_stocks (item_id,store_id,quantity,selling_price)
		VALUES ($1,$2,$3,$4)
		ON CONFLICT (item_id,store_id)
		DO UPDATE SET quantity = EXCLUDED.quantity
		`,
			itemID,
			storeID,
			qty,
			price,
		)

		if err != nil {
			result.Failed++
			result.Errors = append(result.Errors, models.BulkError{i, err.Error()})
			continue
		}

		result.Success++
	}

	return result
}