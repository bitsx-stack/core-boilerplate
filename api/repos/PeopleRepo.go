package repos

import (
	"database/sql"
	"fmt"

	"bitpos.com/m/models"
)

type PeopleRepo struct {
	db *sql.DB
}

func NewPeopleRepo(db *sql.DB) *PeopleRepo {
	return &PeopleRepo{db: db}
}

func (r *PeopleRepo) CreatePeople(people models.People) error {

	query := "INSERT INTO people (type, category, tin, vrn, nin, email, location, name, phone, logo_url, po_box, tags, created_by) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10, $11, $12,$13)"
	_, err := r.db.Exec(query, &people.Type, &people.Category, &people.Tin, &people.Vrn, &people.Nin, &people.Email, &people.Location, &people.Name, &people.Phone, &people.LogoURL, &people.PoBox, &people.Tags, &people.CreatedBy)

	if err != nil {
		return fmt.Errorf("Failed to create people: %w", err)
	}

	return nil
}

func (r *PeopleRepo) GetSinglePeople(id int64) (*models.People, error) {

	people := &models.People{}
	query := "SELECT id, type, category, tin, vrn, nin, email, location, name, phone, status, logo_url, po_box, tags, created_at, created_by, deleted_at, deleted_by FROM people WHERE id = $1 AND deleted_at IS NULL"
	err := r.db.QueryRow(query, id).Scan(&people.ID,
		&people.Type,
		&people.Category,
		&people.Tin,
		&people.Vrn,
		&people.Nin,
		&people.Email,
		&people.Location,
		&people.Name,
		&people.Phone,
		&people.Status,
		&people.LogoURL,
		&people.PoBox,
		&people.Tags,
		&people.CreatedAt,
		&people.CreatedBy,
		&people.DeletedAt,
		&people.DeletedBy)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("Failed to retrieve a setting: %w", err)
	}

	return people, nil
}

func (r *PeopleRepo) GetAllPeople(search, itemType, category string, limit, offset int) ([]models.People, int, error) {

	var people []models.People

	searchPattern := "%" + search + "%"

	query := `
	SELECT 
		p.id,
		p.type,
		p.category,
		p.tin,
		p.vrn,
		p.nin,
		p.email,
		p.location,
		p.name,
		p.phone,
		p.status,
		p.logo_url,
		p.po_box,
		p.tags,
		p.created_at,
		p.created_by,
		p.deleted_at,
		p.deleted_by,
		u.name as created_username
	FROM people p INNER JOIN users u ON p.created_by = u.id 
	WHERE
		($1 = '' OR p.name ILIKE $1)
		AND ($2 = '' OR p.type = $2)
		AND ($3 = '' OR p.category = $3)
		AND p.deleted_at IS NULL
	ORDER BY p.created_at DESC
	LIMIT $4 OFFSET $5
	`

	rows, err := r.db.Query(query, searchPattern, itemType, category, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("query error: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var p models.People

		err := rows.Scan(
			&p.ID,
			&p.Type,
			&p.Category,
			&p.Tin,
			&p.Vrn,
			&p.Nin,
			&p.Email,
			&p.Location,
			&p.Name,
			&p.Phone,
			&p.Status,
			&p.LogoURL,
			&p.PoBox,
			&p.Tags,
			&p.CreatedAt,
			&p.CreatedBy,
			&p.DeletedAt,
			&p.DeletedBy,
			&p.Created_Username,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("scan error: %w", err)
		}

		people = append(people, p)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	countQuery := `
	SELECT COUNT(*)
	FROM people p
	WHERE
		p.deleted_at IS NULL
		AND ($1 = '' OR p.name ILIKE $1)
		AND ($2 = '' OR p.type = $2)
		AND ($3 = '' OR p.category = $3)
	`

	var total int

	err = r.db.QueryRow(countQuery, searchPattern, itemType, category).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("count error: %w", err)
	}

	return people, total, nil
}


func (r *PeopleRepo) UpdatePeople(id int64, people models.People) error {
	query := `UPDATE people 
		SET type=$1, category=$2, tin=$3, vrn=$4, nin=$5, email=$6, 
		    location=$7, name=$8, phone=$9, logo_url=$10, po_box=$11, 
		    tags=$12
		WHERE id=$13 AND deleted_at IS NULL`

	result, err := r.db.Exec(query,
		&people.Type, &people.Category, &people.Tin, &people.Vrn,
		&people.Nin, &people.Email, &people.Location, &people.Name,
		&people.Phone, &people.LogoURL, &people.PoBox, &people.Tags, &id)

	if err != nil {
		return fmt.Errorf("failed to update people: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("people not found")
	}

	return nil
}

func (r *PeopleRepo) DeletePeople(id int64, deletedBy int64) error {
	query := `UPDATE people SET deleted_at=NOW(), deleted_by=$1 WHERE id=$2 AND deleted_at IS NULL`

	result, err := r.db.Exec(query, &deletedBy, &id)
	if err != nil {
		return fmt.Errorf("failed to delete people: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("people not found")
	}

	return nil
}


