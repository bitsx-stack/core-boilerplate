package repos

import (
	"database/sql"
	"fmt"

	"bitpos.com/m/models"
)

type StoreRepo struct {
	db *sql.DB
}

func NewStoreRepo(db *sql.DB) *StoreRepo {
	return &StoreRepo{db: db}
}

func (r *StoreRepo) CreateStore(store models.Store) error {
	query := `
	INSERT INTO stores (name, type, category, location, owner, tags, created_by)
	VALUES ($1,$2,$3,$4,$5,$6,$7)
	`
	_, err := r.db.Exec(query,
		store.Name,
		store.Type,
		store.Category,
		store.Location,
		store.Owner,
		store.Tags,
		store.CreatedBy,
	)

	if err != nil {
		return fmt.Errorf("failed to create store: %w", err)
	}
	return nil
}

func (r *StoreRepo) GetSingleStore(id int64) (*models.Store, error) {
	store := &models.Store{}

	query := `
	SELECT id,name,type,category,location,owner,status,tags,created_at,created_by,deleted_at,deleted_by
	FROM stores WHERE id=$1 AND deleted_at IS NULL
	`

	err := r.db.QueryRow(query, id).Scan(
		&store.ID,
		&store.Name,
		&store.Type,
		&store.Category,
		&store.Location,
		&store.Owner,
		&store.Status,
		&store.Tags,
		&store.CreatedAt,
		&store.CreatedBy,
		&store.DeletedAt,
		&store.DeletedBy,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return store, nil
}

func (r *StoreRepo) UpdateStore(id int64, store models.Store) error {
	query := `UPDATE stores 
		SET name=$1, type=$2, category=$3, location=$4, owner=$5, tags=$6
		WHERE id=$7 AND deleted_at IS NULL`

	result, err := r.db.Exec(query,
		store.Name, store.Type, store.Category, store.Location,
		store.Owner, store.Tags, id)

	if err != nil {
		return fmt.Errorf("failed to update store: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("store not found")
	}

	return nil
}

func (r *StoreRepo) DeleteStore(id int64, deletedBy int64) error {
	query := `UPDATE stores SET deleted_at=NOW(), deleted_by=$1 WHERE id=$2 AND deleted_at IS NULL`

	result, err := r.db.Exec(query, deletedBy, id)
	if err != nil {
		return fmt.Errorf("failed to delete store: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("store not found")
	}

	return nil
}

func (r *StoreRepo) GetAllStores(search, store_type, category, location string, limit, offset int) ([]models.Store, int, error) {
	query := `
	SELECT s.id, s.name, s.type, s.category, s.location, s.owner, s.status, s.tags, s.created_at, s.created_by, s.deleted_at, s.deleted_by, u.name as created_username
	FROM stores s INNER JOIN users u ON s.created_by = u.id WHERE ($1 = '' OR s.name ILIKE $1)
		AND ($2 = '' OR s.type = $2)
		AND ($3 = '' OR s.category = $3)
		AND ($4 = '' OR s.location = $4) AND s.deleted_at IS NULL ORDER BY s.id DESC LIMIT $5 OFFSET $6`

	searchPattern := "%"+search+"%"
	rows, err := r.db.Query(query, searchPattern, store_type, category, location, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to fetch stores: %w", err)
	}
	defer rows.Close()

	var stores []models.Store
	for rows.Next() {
		var s models.Store
		err := rows.Scan(
			&s.ID, &s.Name, &s.Type, &s.Category, &s.Location,
			&s.Owner, &s.Status, &s.Tags, &s.CreatedAt, &s.CreatedBy,
			&s.DeletedAt, &s.DeletedBy,&s.Created_Username,
		)
		if err != nil {
			return nil, 0, err
		}
		stores = append(stores, s)
	}

	var total int

	return stores, total, nil
}