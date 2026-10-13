package repos

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)


type RoleRepo struct {
	db *sql.DB
}

func NewRoleRepo(db *sql.DB) *RoleRepo {
	return &RoleRepo{db: db}
}

/* ================= CREATE ROLE ================= */

func (r *RoleRepo) CreateRole(name string, createdBy int64) error {
	var roleID int64
	err := r.db.QueryRow(`
		INSERT INTO roles (name, created_by)
		VALUES ($1,$2) RETURNING id
	`, name, createdBy).Scan(&roleID)

	if err != nil {
		return err
	}

	return nil
}


func (r *RoleRepo) UpdateRole(name string, permissions []string, roleID int64) error {

	name = strings.TrimSpace(name)

	query := "UPDATE roles SET "
	args := []interface{}{}
	i := 1

	updates := []string{}

	if name != "" {
		updates = append(updates, fmt.Sprintf("name=$%d", i))
		args = append(args, name)
		i++
	}

	if len(permissions) > 0 {
		updates = append(updates, fmt.Sprintf("permissions=$%d", i))
		args = append(args, permissions)
		i++
	}

	// nothing to update
	if len(updates) == 0 {
		return nil
	}

	query += strings.Join(updates, ", ")
	query += fmt.Sprintf(" WHERE id=$%d", i)
	args = append(args, roleID)

	_, err := r.db.Exec(query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (r *RoleRepo) GetAllRoles(search string, limit, offset int) ([]map[string]interface{}, int, error) {
	baseQuery := `
		FROM roles r 
		LEFT JOIN users u ON u.id = r.created_by 
		WHERE r.deleted_at IS NULL
	`
	var args []interface{}
	argIndex := 1

	if search != "" {
		baseQuery += fmt.Sprintf(" AND r.name ILIKE $%d", argIndex)
		args = append(args, "%"+search+"%")
		argIndex++
	}

	// Count
	var total int
	err := r.db.QueryRow("SELECT COUNT(*) "+baseQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("count failed: %w", err)
	}

	// Data
	query := fmt.Sprintf(`
    SELECT r.id, r.name, COALESCE(r.permissions::text, '[]') as permissions, 
           r.created_at, COALESCE(u.name, '') as created_username 
    %s ORDER BY r.created_at DESC LIMIT $%d OFFSET $%d
`, baseQuery, argIndex, argIndex+1)
	args = append(args, limit, offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	var results []map[string]interface{}
	for rows.Next() {
		var id int64
		var name, permissions, createdUsername string
		var createdAt time.Time

		err := rows.Scan(&id, &name, &permissions, &createdAt, &createdUsername)
		if err != nil {
			return nil, 0, fmt.Errorf("scan failed: %w", err)
		}

		results = append(results, map[string]interface{}{
			"id":               id,
			"name":             name,
			"permissions":      permissions,
			"created_at":       createdAt,
			"created_username": createdUsername,
		})
	}

	return results, total, nil
}


func (r *RoleRepo) GetSingleRole(id int64) (map[string]interface{}, error) {
	var roleID int64
	var name, permissions, createdUsername string
	var createdAt time.Time

	err := r.db.QueryRow(`
		SELECT r.id, r.name, COALESCE(r.permissions::text, '[]') as permissions,
		       r.created_at, COALESCE(u.name, '') as created_username
		FROM roles r
		LEFT JOIN users u ON u.id = r.created_by
		WHERE r.id = $1 AND r.deleted_at IS NULL
	`, id).Scan(&roleID, &name, &permissions, &createdAt, &createdUsername)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("role not found")
	}
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}

	return map[string]interface{}{
		"id":               roleID,
		"name":             name,
		"permissions":      permissions,
		"created_at":       createdAt,
		"created_username": createdUsername,
	}, nil
}

// ══════════════════════════════════════════════
// ROLE REPO — add these methods
// ══════════════════════════════════════════════

func (r *RoleRepo) UpdatePermissions(roleID int64, permissions []string) error {
	permJSON, err := json.Marshal(permissions)
	if err != nil {
		return fmt.Errorf("failed to marshal permissions: %w", err)
	}

	result, err := r.db.Exec(`
		UPDATE roles SET permissions = $1 WHERE id = $2 AND deleted_at IS NULL
	`, string(permJSON), roleID)

	if err != nil {
		return fmt.Errorf("failed to update permissions: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("role not found")
	}
	return nil
}

func (r *RoleRepo) GetRolePermissions(roleID int64) ([]string, error) {
	var raw string
	err := r.db.QueryRow(`
		SELECT COALESCE(permissions::text, '[]') FROM roles 
		WHERE id = $1 AND deleted_at IS NULL
	`, roleID).Scan(&raw)

	if err != nil {
		return nil, fmt.Errorf("role not found: %w", err)
	}

	var perms []string
	if err := json.Unmarshal([]byte(raw), &perms); err != nil {
		return []string{}, nil
	}
	return perms, nil
}