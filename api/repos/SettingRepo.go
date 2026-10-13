package repos

import (
	"database/sql"
	"fmt"

	"bitpos.com/m/models"
)

type SettingsFilter struct {
	Search  string
	SortBy  string
	SortDir string
	Page    int
	PerPage int
}

type SettingsResult struct {
	Data []models.Setting `json:"data"`
	Meta PaginationMeta   `json:"meta"`
}

type PaginationMeta struct {
	Page       int `json:"page"`
	PerPage    int `json:"perPage"`
	Total      int `json:"total"`
	TotalPages int `json:"totalPages"`
}

type SettingRepo struct {
	db *sql.DB
}

func NewSettingRepo(db *sql.DB) *SettingRepo {
	return &SettingRepo{db: db}
}

// Get settings
func (r *SettingRepo) GetSettings(search, category string, limit, offset int) ([]models.Setting, int, error) {

	var searchPattern *string
	if search != "" {
		s := "%" + search + "%"
		searchPattern = &s
	}

	var categoryFilter *string
	if category != "" {
		categoryFilter = &category
	}

	var total int
	countQuery := `
		SELECT COUNT(*) 
		FROM settings 
		WHERE deleted_at IS NULL 
		  AND ($1::text IS NULL OR name ILIKE $1) 
		  AND ($2::text IS NULL OR category = $2)
	`

	if err := r.db.QueryRow(countQuery, searchPattern, categoryFilter).Scan(&total); err != nil {
		return nil, 0, err
	}

	query := `
		SELECT id, code, name, value, category
		FROM settings
		WHERE deleted_at IS NULL
		  AND ($1::text IS NULL OR name ILIKE $1) 
		  AND ($2::text IS NULL OR category = $2)
		ORDER BY created_at DESC, id DESC
		LIMIT $3 OFFSET $4
	`

	rows, err := r.db.Query(query, searchPattern, categoryFilter, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	settings := make([]models.Setting, 0)

	for rows.Next() {
		var model models.Setting
		if err := rows.Scan(
			&model.ID, &model.Code, &model.Name, &model.Value, &model.Category,
		); err != nil {
			return nil, 0, err
		}
		settings = append(settings, model)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return settings, total, nil
}

// Get setting by ID
func (r *SettingRepo) GetSetting(id int64) (*models.Setting, error) {
	setting := &models.Setting{}
	query := `SELECT * FROM settings WHERE id = $1 AND deleted_at IS NULL`

	err := r.db.QueryRow(query, id).Scan(
		&setting.ID, &setting.Code, &setting.Name, &setting.Value, &setting.Category, &setting.CreatedAt, &setting.DeletedAt, &setting.DeletedBy, &setting.CreatedBy)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("Failed to retrieve a setting: %w", err)
	}

	return setting, nil
}

// Create a new setting
func (r *SettingRepo) CreateSetting(setting *models.Setting) error {
	query := `INSERT INTO settings (code, name, value, category, created_by) VALUES ($1, $2, $3, $4, $5)`
	_, err := r.db.Exec(query, setting.Code, setting.Name, setting.Value, setting.Category, setting.CreatedBy)
	if err != nil {
		return fmt.Errorf("Failed to create a setting: %w", err)
	}
	return nil
}

// Update an existing setting
func (r *SettingRepo) UpdateSetting(setting *models.Setting) error {
	query := `UPDATE settings SET name = $1, value=$2, category = $3 WHERE id = $4`
	_, err := r.db.Exec(query, setting.Name, setting.Value, setting.Category, setting.ID)
	if err != nil {
		return fmt.Errorf("Failed to update a setting: %w", err)
	}
	return nil
}

// Soft delete a setting
func (r *SettingRepo) DeleteSetting(id int64, deletedBy int64) error {
	query := `UPDATE settings SET deleted_at = NOW(), deleted_by = $1 WHERE id = $2`
	_, err := r.db.Exec(query, deletedBy, id)
	if err != nil {
		return fmt.Errorf("Failed to delete a setting: %w", err)
	}
	return nil
}
