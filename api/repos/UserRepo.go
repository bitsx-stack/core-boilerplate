package repos

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"bitpos.com/m/models"
)

type UserRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) CreateUser(user *models.User) error {
	query := `INSERT INTO users (email, password, name, phone) VALUES ($1, $2, $3, $4)`
	result, err := r.db.Exec(query, &user.Email, &user.Password, &user.Name, &user.Phone)
	if err != nil {
		return fmt.Errorf("Failed to create user: %w", err)
	}

	_, err = result.RowsAffected()
	if err != nil {
		return fmt.Errorf("Failed to get last inserted user ID: %w", err)
	}

	return nil
}

func (r *UserRepo) GetUserByEmail(email string) (*models.User, error) {
	user := &models.User{}
	query := `SELECT * FROM users WHERE email = $1`

	err := r.db.QueryRow(query, email).Scan(
		&user.ID, &user.Email, &user.Password, &user.Name, &user.Phone, &user.Status, &user.CreatedAt, &user.DeletedAt, &user.DeletedBy,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("Failed to retrieve user information by email: %w", err)
	}

	return user, nil
}

func (r *UserRepo) GetUserByID(id int64) (*models.User, error) {
	user := &models.User{}
	query := `SELECT * FROM users WHERE id = $1 AND deleted_at IS NULL`

	err := r.db.QueryRow(query, id).Scan(
		&user.ID, &user.Email, &user.Password, &user.Name, &user.Phone, &user.Status, &user.CreatedAt, &user.DeletedAt, &user.DeletedBy,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("Failed to retrieve user information by ID: %w", err)
	}

	return user, nil
}

func (r *UserRepo) SahihishaTaarifaZaUser(user *models.User) error {
	query := `UPDATE users SET name = ?, phone = ?, updated_at = NOW() WHERE id = ?`
	_, err := r.db.Exec(query, user.Name, user.Phone, user.ID)
	if err != nil {
		return fmt.Errorf("Failed to update user information: %w", err)
	}
	return nil
}

// PataWatumiaji - Get all users
func (r *UserRepo) PataWatumiaji() ([]models.User, error) {
	query := `SELECT id, email, name, phone, created_at, updated_at FROM users ORDER BY created_at DESC`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("Failed to retrieve users: %w", err)
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		user := models.User{}
		err := rows.Scan(&user.ID, &user.Email, &user.Name, &user.Phone, &user.CreatedAt, &user.DeletedAt)
		if err != nil {
			return nil, fmt.Errorf("Failed to scan user: %w", err)
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("Failed to iterate users: %w", err)
	}

	return users, nil
}

// FutaUser - Delete a user
func (r *UserRepo) FutaUser(id int64) error {
	query := `DELETE FROM users WHERE id = ?`
	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("Failed to delete user: %w", err)
	}
	return nil
}





func (r *UserRepo) GetAllUsers(search string, limit, offset int) ([]models.User, int, error) {
	baseQuery := ` FROM users WHERE deleted_at IS NULL`
	var conditions []string
	var args []interface{}
	argIndex := 1

	if search != "" {
		conditions = append(conditions, fmt.Sprintf(
			"(name ILIKE $%d OR email ILIKE $%d OR phone ILIKE $%d)",
			argIndex, argIndex, argIndex,
		))
		args = append(args, "%"+search+"%")
		argIndex++
	}

	if len(conditions) > 0 {
		baseQuery += " AND " + strings.Join(conditions, " AND ")
	}

	// Count
	var total int
	err := r.db.QueryRow("SELECT COUNT(*)"+baseQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("count failed: %w", err)
	}

	// Data
	selectQuery := "SELECT id, email, name, phone, status, created_at" + baseQuery +
		fmt.Sprintf(" ORDER BY id DESC LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, limit, offset)

	rows, err := r.db.Query(selectQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		err := rows.Scan(&u.ID, &u.Email, &u.Name, &u.Phone, &u.Status, &u.CreatedAt)
		if err != nil {
			return nil, 0, fmt.Errorf("scan failed: %w", err)
		}
		users = append(users, u)
	}

	return users, total, nil
}

// func (r *UserRepo) GetUserByID(id int64) (*models.User, error) {
// 	var u models.User
// 	err := r.db.QueryRow(`
// 		SELECT id, email, name, phone, status, created_at
// 		FROM users WHERE id = $1 AND deleted_at IS NULL
// 	`, id).Scan(&u.ID, &u.Email, &u.Name, &u.Phone, &u.Status, &u.CreatedAt)

// 	if err == sql.ErrNoRows {
// 		return nil, nil
// 	}
// 	if err != nil {
// 		return nil, fmt.Errorf("query failed: %w", err)
// 	}
// 	return &u, nil
// }

// func (r *UserRepo) CreateUser(u models.User) error {
// 	// Check duplicate email
// 	var exists bool
// 	r.db.QueryRow(`SELECT EXISTS(SELECT 1 FROM users WHERE email = $1 AND deleted_at IS NULL)`, u.Email).Scan(&exists)
// 	if exists {
// 		return fmt.Errorf("email already exists")
// 	}

// 	// Check duplicate phone
// 	r.db.QueryRow(`SELECT EXISTS(SELECT 1 FROM users WHERE phone = $1 AND deleted_at IS NULL)`, u.Phone).Scan(&exists)
// 	if exists {
// 		return fmt.Errorf("phone already exists")
// 	}

// 	_, err := r.db.Exec(`
// 		INSERT INTO users (email, password, name, phone, status)
// 		VALUES ($1, $2, $3, $4, $5)
// 	`, u.Email, u.Password, u.Name, u.Phone, u.Status)

// 	if err != nil {
// 		return fmt.Errorf("failed to create user: %w", err)
// 	}
// 	return nil
// }

func (r *UserRepo) UpdateUser(id int64, u models.User) error {
	// Check duplicate email (exclude self)
	var exists bool
	r.db.QueryRow(`SELECT EXISTS(SELECT 1 FROM users WHERE email = $1 AND id != $2 AND deleted_at IS NULL)`, u.Email, id).Scan(&exists)
	if exists {
		return fmt.Errorf("email already exists")
	}

	// Check duplicate phone (exclude self)
	r.db.QueryRow(`SELECT EXISTS(SELECT 1 FROM users WHERE phone = $1 AND id != $2 AND deleted_at IS NULL)`, u.Phone, id).Scan(&exists)
	if exists {
		return fmt.Errorf("phone already exists")
	}

	result, err := r.db.Exec(`
		UPDATE users SET email = $1, name = $2, phone = $3, status = $4
		WHERE id = $5 AND deleted_at IS NULL
	`, u.Email, u.Name, u.Phone, u.Status, id)

	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("user not found")
	}
	return nil
}

func (r *UserRepo) DeleteUser(id int64, deletedBy int64) error {
	result, err := r.db.Exec(`
		UPDATE users SET deleted_at = NOW(), deleted_by = $1
		WHERE id = $2 AND deleted_at IS NULL
	`, deletedBy, id)

	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("user not found")
	}
	return nil
}



func (r *UserRepo) GetUserPermissionsMap(userID int64) (map[string]bool, error) {

	rows, err := r.db.Query(`
		SELECT ua.permissions
		FROM users_access ua INNER JOIN users u ON u.id = ua.user_id
		WHERE ua.user_id = $1
	`, userID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	perms := make(map[string]bool)

	for rows.Next() {
		var name string
		rows.Scan(&name)
		perms[name] = true
	}

	return perms, nil
}


// ══════════════════════════════════════════════
// USER REPO — add these methods
// ══════════════════════════════════════════════

func (r *UserRepo) AssignRole(userID int64, roleID int64) error {
	result, err := r.db.Exec(`
		UPDATE users SET role_id = $1 WHERE id = $2 AND deleted_at IS NULL
	`, roleID, userID)

	if err != nil {
		return fmt.Errorf("failed to assign role: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("user not found")
	}
	return nil
}

func (r *UserRepo) GetUserPermissions(userID int64) ([]string, error) {
	// Get permissions from user's role
	var raw string
	err := r.db.QueryRow(`
		SELECT COALESCE(r.permissions::text, '[]')
		FROM users u
		LEFT JOIN roles r ON r.id = u.role_id
		WHERE u.id = $1 AND u.deleted_at IS NULL
	`, userID).Scan(&raw)

	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	var rolePerms []string
	json.Unmarshal([]byte(raw), &rolePerms)

	// Get user's own direct permissions
	var userRaw sql.NullString
	r.db.QueryRow(`
		SELECT COALESCE(permissions::text, '[]') FROM users WHERE id = $1
	`, userID).Scan(&userRaw)

	var userPerms []string
	if userRaw.Valid {
		json.Unmarshal([]byte(userRaw.String), &userPerms)
	}

	// Merge: role permissions + user direct permissions (deduplicated)
	merged := make(map[string]bool)
	for _, p := range rolePerms {
		merged[p] = true
	}
	for _, p := range userPerms {
		merged[p] = true
	}

	result := make([]string, 0, len(merged))
	for p := range merged {
		result = append(result, p)
	}
	return result, nil
}

func (r *UserRepo) UpdateUserPermissions(userID int64, permissions []string) error {
	permJSON, err := json.Marshal(permissions)
	if err != nil {
		return fmt.Errorf("failed to marshal: %w", err)
	}

	result, err := r.db.Exec(`
		UPDATE users SET permissions = $1 WHERE id = $2 AND deleted_at IS NULL
	`, string(permJSON), userID)

	if err != nil {
		return fmt.Errorf("failed to update: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("user not found")
	}
	return nil
}


