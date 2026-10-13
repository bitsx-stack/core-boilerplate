package repos

import (
	"database/sql"
	"fmt"

	"bitpos.com/m/models"
)

type ItemRepo struct {
	db *sql.DB
}

func NewItemRepo(db *sql.DB) *ItemRepo {
	return &ItemRepo{db: db}
}

// CREATE
func (r *ItemRepo) CreateItem(item models.Item) error {

	query := `
			INSERT INTO items 
			(name, type, category, uom, unit_cost, tags, created_by) 
			VALUES ($1,$2,$3,$4,$5,$6,$7)
			`

	_, err := r.db.Exec(query,
		item.Name,
		item.Type,
		item.Category,
		item.UOM,
	    item.UnitCost,
		item.Tags,
		item.CreatedBy,
	)

	if err != nil {
		return fmt.Errorf("failed to create item: %w", err)
	}

	return nil
}

// GET SINGLE
func (r *ItemRepo) GetSingleItem(id int64) (*models.Item, error) {

	item := &models.Item{}

	query := `
	SELECT id, name, type, category, unit_cost, uom, status,
	       created_at, created_by, deleted_at, deleted_by
	FROM items
	WHERE id = $1 AND deleted_at IS NULL
	`

	err := r.db.QueryRow(query, id).Scan(
		&item.ID,
		&item.Name,
		&item.Type,
		&item.Category,
		&item.UnitCost,
		&item.UOM,
		&item.Status,
		&item.CreatedAt,
		&item.CreatedBy,
		&item.DeletedAt,
		&item.DeletedBy,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get item: %w", err)
	}

	return item, nil
}

// GET ALL
func (r *ItemRepo) GetAllItems(search, itemType, category string, storeID int, limit, offset int) ([]models.Item, int, error) {

	var items []models.Item

	searchPattern := "%" + search + "%"

	query := `
	SELECT 
		i.id,
		i.name,
		i.type,
		i.category,
		i.unit_cost,
		i.uom,
		i.status,
		i.created_at,
		i.created_by,
		i.deleted_at,
		u.name as created_username,
		COALESCE(
			json_agg(
				json_build_object(
					'store_id', s.store_id,
					'quantity', s.quantity,
					'selling_price', s.selling_price
				)
			) FILTER (WHERE s.store_id IS NOT NULL),
			'[]'
		) AS stock
	FROM items i
	LEFT JOIN items_stocks s ON s.item_id = i.id
	INNER JOIN users u ON u.id = i.created_by
	WHERE
		($1 = '' OR i.name ILIKE $1)
		AND ($2 = '' OR i.type = $2)
		AND ($3 = '' OR i.category = $3)
		AND ($4 = 0 OR s.store_id = $4)
		AND i.deleted_at IS NULL
	GROUP BY i.id, u.name
	ORDER BY i.created_at DESC
	LIMIT $5 OFFSET $6
	`

	rows, err := r.db.Query(query,
		searchPattern,
		itemType,
		category,
		storeID,
		limit,
		offset,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var i models.Item

		err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Type,
			&i.Category,
			&i.UnitCost,
			&i.UOM,
			&i.Status,
			&i.CreatedAt,
			&i.CreatedBy,
			&i.DeletedAt,
			&i.Created_Username,
			&i.Stock,
		)
		if err != nil {
			return nil, 0, err
		}

		items = append(items, i)
	}

	// COUNT FIXED
	countQuery := `
	SELECT COUNT(DISTINCT i.id)
	FROM items i
	LEFT JOIN items_stocks s ON s.item_id = i.id
	WHERE
		($1 = '' OR i.name ILIKE $1)
		AND ($2 = '' OR i.type = $2)
		AND ($3 = '' OR i.category = $3)
		AND ($4 = 0 OR s.store_id = $4)
		AND i.deleted_at IS NULL
	`

	var total int
	err = r.db.QueryRow(countQuery,
		searchPattern,
		itemType,
		category,
		storeID,
	).Scan(&total)

	if err != nil {
		return nil, 0, err
	}

	return items, total, nil
}


func (r *ItemRepo) UpdateItem(id int64, item models.Item) error {
	query := `UPDATE items 
		SET name=$1, type=$2, category=$3, uom=$4, unit_cost=$5
		WHERE id=$6 AND deleted_at IS NULL`

	result, err := r.db.Exec(query, item.Name, item.Type, item.Category, item.UOM, item.UnitCost, id)
	if err != nil {
		return fmt.Errorf("failed to update item: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("item not found")
	}

	return nil
}

func (r *ItemRepo) DeleteItem(id int64, deletedBy int64) error {
	query := `UPDATE items SET deleted_at=NOW(), deleted_by=$1 WHERE id=$2 AND deleted_at IS NULL`

	result, err := r.db.Exec(query, deletedBy, id)
	if err != nil {
		return fmt.Errorf("failed to delete item: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("item not found")
	}

	return nil
}