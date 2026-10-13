package repos

import (
	"database/sql"
	"fmt"
	"time"
	"bitpos.com/m/models"
)

type CategoryRepo struct {
	db *sql.DB
}

func NewCategoryRepo(db *sql.DB) *CategoryRepo {
	return &CategoryRepo{db: db}
}

// PataCategories - Get all categories
func (r *CategoryRepo) GetCategories(search, category_type string, limit, offset int) ([]models.Category, int, error) {

	searchPattern := "%" + search + "%"

	query := `
	SELECT 
		c.id,
		c.name,
		c.description,
		ct.name AS type,
		c.created_at,
		u.name AS created_by
	FROM categories c
	INNER JOIN users u ON c.created_by = u.id
	INNER JOIN category_types ct ON c.categorytype_id = ct.id
	WHERE c.deleted_at IS NULL
	  AND ($1 = '' OR c.name ILIKE $1)
	  AND ($2 = '' OR ct.name=$2)
	ORDER BY c.created_at DESC
	LIMIT $3 OFFSET $4
	`

	rows, err := r.db.Query(query, searchPattern, category_type, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var categories []models.Category

	for rows.Next() {
		var category models.Category

		err := rows.Scan(
			&category.ID,
			&category.Name,
			&category.Description,
			&category.Type,               // ct.name
			&category.CreatedAt,
			&category.Created_Username,  // u.name
		)
		if err != nil {
			return nil, 0, err
		}

		categories = append(categories, category)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	// ✅ COUNT QUERY
	countQuery := `
	SELECT COUNT(*)
	FROM categories c
	INNER JOIN category_types ct ON c.categorytype_id = ct.id
	WHERE c.deleted_at IS NULL
	  AND ($1 = '' OR c.name ILIKE $1)
	  AND ($2 = '' OR ct.name=$2)
	`

	var total int
	err = r.db.QueryRow(countQuery, searchPattern, category_type).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	return categories, total, nil
}

// PataCategory - Get category by ID
func (r *CategoryRepo) GetCategory(id int64) (*models.Category, error) {
	category := &models.Category{}
	query := `SELECT id, categorytype_id, name, description, created_at FROM categories WHERE id = $1 AND deleted_at IS NULL`

	err := r.db.QueryRow(query, id).Scan(
		&category.ID, &category.CategoryTypeID, &category.Name, &category.Description, &category.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("Failed to retrieve category: %w", err)
	}

	return category, nil
}



// CreateCategory - Create a new category
func (r *CategoryRepo) CreateCategory(category models.Category) error {
	query := `INSERT INTO categories (categorytype_id, name, description, created_at, created_by) VALUES ($1, $2, $3, NOW(), $4)`
	result, err := r.db.Exec(query, category.CategoryTypeID, category.Name, category.Description, category.CreatedBy)
	if err != nil {
		return fmt.Errorf("Failed to create category: %w", err)
	}

	_, err = result.RowsAffected()
	if err != nil {
		return fmt.Errorf("Failed to get last inserted category ID: %w", err)
	}

	return nil
}

// UpdateCategory - Update an existing category
func (r *CategoryRepo) UpdateCategory(id int64, category *models.Category) (*models.Category, error) {
	query := `UPDATE categories SET categorytype_id = $1, name = $2, description = $3 WHERE id = $4`
	_, err := r.db.Exec(query, category.CategoryTypeID, category.Name, category.Description, id)
	if err != nil {
		return nil, fmt.Errorf("Failed to update category: %w", err)
	}

	return r.GetCategory(id)
}

// DeleteCategory - Delete a category
func (r *CategoryRepo) DeleteCategory(id, userID int64) error {
	query := `UPDATE categories SET deleted_at = $1, deleted_by = $2 WHERE id = $3`
	_, err := r.db.Exec(query, time.Now(), userID, id)
	if err != nil {
		return fmt.Errorf("Failed to delete category: %w", err)
	}
	return nil
}
