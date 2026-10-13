package repos

import (
	"database/sql"
	"fmt"
	"time"

	"bitpos.com/m/models"
)

type CategoryTypeRepo struct {
	db *sql.DB
}

func NewCategoryTypeRepo(db *sql.DB) *CategoryTypeRepo {
	return &CategoryTypeRepo{db: db}
}

// Get all category types
func (r *CategoryTypeRepo) GetCategoryTypes(search string, limit, offset int) ([]models.CategoryType, int, error) {

	searchPattern := "%" + search + "%"

	query := `
	SELECT ct.id, ct.name, ct.description, ct.created_at, u.name
	FROM category_types ct
	INNER JOIN users u ON ct.created_by = u.id
	WHERE ct.deleted_at IS NULL
	  AND ($1 = '' OR ct.name ILIKE $1)
	ORDER BY ct.created_at DESC
	LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(query, searchPattern, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var types []models.CategoryType

	for rows.Next() {
		var t models.CategoryType
		if err := rows.Scan(&t.ID, &t.Name, &t.Description, &t.CreatedAt, &t.Created_Username); err != nil {
			return nil, 0, err
		}
		types = append(types, t)
	}

	countQuery := `
	SELECT COUNT(*)
	FROM category_types
	WHERE deleted_at IS NULL
	  AND ($1 = '' OR name ILIKE $1)
	`

	var total int
	err = r.db.QueryRow(countQuery, searchPattern).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	return types, total, nil
}

// PataCategoryType - Get category type by ID
func (r *CategoryTypeRepo) GetCategoryTypeByID(id int64) (*models.CategoryType, error) {
	categoryType := &models.CategoryType{}
	query := `SELECT 
					ct.id, 
					ct.name, 
					ct.description, 
					ct.created_at, 
					u.name AS created_by
					FROM category_types ct
					INNER JOIN users u 
					ON ct.created_by = u.id
					WHERE ct.id = $1 AND ct.deleted_at IS NULL;`

	err := r.db.QueryRow(query, id).Scan(
		&categoryType.ID, &categoryType.Name, &categoryType.Description, &categoryType.CreatedAt, &categoryType.Created_Username,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("Failed to retrieve category type: %w", err)
	}

	return categoryType, nil
}




// CreateCategoryType - Create a new category type
func (r *CategoryTypeRepo) CreateCategoryType(categoryType models.CategoryType) error {
	query := `INSERT INTO category_types (name, description, created_by, created_at) VALUES ($1, $2, $3, NOW())`
	result, err := r.db.Exec(query, categoryType.Name, categoryType.Description, categoryType.CreatedBy)
	if err != nil {
		return fmt.Errorf("Failed to create category type: %w", err)
	}

	_, err = result.RowsAffected()
	if err != nil {
		return fmt.Errorf("Failed to get last inserted category type ID: %w", err)
	}

	return nil

}

// UpdateCategoryType - Update an existing category type
func (r *CategoryTypeRepo) UpdateCategoryType(id int64, categoryType *models.CategoryType) (*models.CategoryType, error) {
	query := `UPDATE category_types SET name = $1, description = $2 WHERE id = $3`
	_, err := r.db.Exec(query, categoryType.Name, categoryType.Description, id)
	if err != nil {
		return nil, fmt.Errorf("Failed to update category type: %w", err)
	}

	return r.GetCategoryTypeByID(id)
}

// DeleteCategoryType - Delete a category type
func (r *CategoryTypeRepo) DeleteCategoryType(id int64, userID int64) error {
	query := `UPDATE category_types SET deleted_at = $1, deleted_by = $2 WHERE id = $3`
	_, err := r.db.Exec(query, time.Now(), userID, id)
	if err != nil {
		return fmt.Errorf("Failed to delete category type: %w", err)
	}
	return nil
}
