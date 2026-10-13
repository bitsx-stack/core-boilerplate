package repos

import (
	"database/sql"
	"encoding/json"
	"fmt"

)

type UserAccessRepo struct {
	db *sql.DB
}

func NewUserAccessRepo(db *sql.DB) *UserAccessRepo {
	return &UserAccessRepo{db: db}
}



// GetUserPermissions — returns permissions from users_access
func (r *UserAccessRepo) GetUserPermissions(userID int64) ([]string, error) {
	var raw string
	err := r.db.QueryRow(`
		SELECT COALESCE(permissions::text, '[]') 
		FROM users_access 
		WHERE user_id = $1 AND deleted_at IS NULL
		ORDER BY created_at DESC LIMIT 1
	`, userID).Scan(&raw)

	if err == sql.ErrNoRows {
		return []string{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get permissions: %w", err)
	}

	var perms []string
	if err := json.Unmarshal([]byte(raw), &perms); err != nil {
		return []string{}, nil
	}
	return perms, nil
}

func (r *UserAccessRepo) AssignRoleToUser(userID int64, roleID int64) error {

	// ─────────────────────────────
	// 1. GET ROLE PERMISSIONS
	// ─────────────────────────────
	var roleRaw []byte

	err := r.db.QueryRow(`
		SELECT COALESCE(permissions, '[]'::jsonb)
		FROM roles
		WHERE id = $1 AND deleted_at IS NULL
	`, roleID).Scan(&roleRaw)

	if err == sql.ErrNoRows {
		return fmt.Errorf("role not found")
	}
	if err != nil {
		return fmt.Errorf("failed to get role: %w", err)
	}

	var rolePerms []string
	if len(roleRaw) > 0 {
		_ = json.Unmarshal(roleRaw, &rolePerms)
	}

	// ─────────────────────────────
	// 2. ENSURE USER EXISTS
	// ─────────────────────────────
	var exists bool
	err = r.db.QueryRow(`
		SELECT EXISTS(
			SELECT 1 FROM users
			WHERE id = $1 AND deleted_at IS NULL
		)
	`, userID).Scan(&exists)

	if err != nil {
		return fmt.Errorf("failed to check user: %w", err)
	}

	if !exists {
		return fmt.Errorf("user not found")
	}

	// ─────────────────────────────
	// 3. CONVERT TO JSON
	// ─────────────────────────────
	finalJSON, err := json.Marshal(rolePerms)
	if err != nil {
		return fmt.Errorf("failed to encode permissions: %w", err)
	}

	// ─────────────────────────────
	// 4. UPDATE (REPLACE EVERYTHING)
	// ─────────────────────────────
	res, err := r.db.Exec(`
		UPDATE users_access
		SET permissions = $1
		WHERE user_id = $2 AND deleted_at IS NULL
	`, finalJSON, userID)

	if err != nil {
		return fmt.Errorf("update failed: %w", err)
	}

	rows, _ := res.RowsAffected()

	// ─────────────────────────────
	// 5. INSERT IF NOT EXISTS
	// ─────────────────────────────
	if rows == 0 {
		_, err = r.db.Exec(`
			INSERT INTO users_access (user_id, permissions)
			VALUES ($1, $2)
		`, userID, finalJSON)

		if err != nil {
			return fmt.Errorf("insert failed: %w", err)
		}
	}

	return nil
}


func (r *UserAccessRepo) UpdateUserPermissions(userID int64, permissions []string) error {

	// ─────────────────────────────
	// 1. GET EXISTING PERMISSIONS
	// ─────────────────────────────
	var existingRaw []byte
	var existing []string

	err := r.db.QueryRow(`
		SELECT COALESCE(permissions, '[]'::jsonb)
		FROM users_access
		WHERE user_id = $1 AND deleted_at IS NULL
		LIMIT 1
	`, userID).Scan(&existingRaw)

	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("failed to get existing permissions: %w", err)
	}

	if len(existingRaw) > 0 {
		_ = json.Unmarshal(existingRaw, &existing)
	}

	// ─────────────────────────────
	// 2. MERGE (APPEND + REMOVE DUPLICATES)
	// ─────────────────────────────
	permMap := map[string]struct{}{}

	for _, p := range existing {
		permMap[p] = struct{}{}
	}
	for _, p := range permissions {
		permMap[p] = struct{}{}
	}

	final := make([]string, 0, len(permMap))
	for p := range permMap {
		final = append(final, p)
	}

	finalJSON, err := json.Marshal(final)
	if err != nil {
		return fmt.Errorf("failed to marshal: %w", err)
	}

	// ─────────────────────────────
	// 3. UPDATE OR INSERT
	// ─────────────────────────────
	res, err := r.db.Exec(`
		UPDATE users_access
		SET permissions = $1
		WHERE user_id = $2 AND deleted_at IS NULL
	`, finalJSON, userID)

	if err != nil {
		return fmt.Errorf("update failed: %w", err)
	}

	rows, _ := res.RowsAffected()

	if rows == 0 {
		_, err = r.db.Exec(`
			INSERT INTO users_access (user_id, permissions)
			VALUES ($1, $2)
		`, userID, finalJSON)

		if err != nil {
			return fmt.Errorf("insert failed: %w", err)
		}
	}

	return nil
}




// HasPermission — check if user has a specific permission
func (r *UserAccessRepo) HasPermission(userID int64, permission string) (bool, error) {
	perms, err := r.GetUserPermissions(userID)
	if err != nil {
		return false, err
	}
	for _, p := range perms {
		if p == permission {
			return true, nil
		}
	}
	return false, nil
}