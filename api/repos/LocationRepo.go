package repos

import (
	"database/sql"
	"fmt"
	"strings"
)

type LocationRepo struct {
	db *sql.DB
}

func NewLocationRepo(db *sql.DB) *LocationRepo {
	return &LocationRepo{db: db}
}

func (r *LocationRepo) GetAllLocations(
	model, search, country string,
	regionID, districtID int,
) ([]map[string]interface{}, error) {

	var query string
	var args []interface{}
	argIndex := 1
	var conditions []string

	switch model {

	case "continent":
		query = `SELECT id, name, code, status FROM continents WHERE 1=1`

		if search != "" {
			conditions = append(conditions, fmt.Sprintf("name ILIKE $%d", argIndex))
			args = append(args, "%"+search+"%")
			argIndex++
		}

	case "country":
		query = `SELECT id, name, status FROM countries WHERE 1=1`

		if search != "" {
			conditions = append(conditions, fmt.Sprintf("name ILIKE $%d", argIndex))
			args = append(args, "%"+search+"%")
			argIndex++
		}

	case "region":
		query = `SELECT id, country, zone, name, code, status FROM regions WHERE 1=1`

		if country != "" {
			conditions = append(conditions, fmt.Sprintf("country = $%d", argIndex))
			args = append(args, country)
			argIndex++
		}

		if search != "" {
			conditions = append(conditions, fmt.Sprintf("name ILIKE $%d", argIndex))
			args = append(args, "%"+search+"%")
			argIndex++
		}

	case "district":
		query = `
			SELECT d.id, d.name, d.region_id, r.name as region_name, d.status 
			FROM districts d 
			INNER JOIN regions r ON r.id = d.region_id 
			WHERE 1=1
		`

		if regionID > 0 {
			conditions = append(conditions, fmt.Sprintf("d.region_id = $%d", argIndex))
			args = append(args, regionID)
			argIndex++
		}

		if search != "" {
			conditions = append(conditions, fmt.Sprintf("d.name ILIKE $%d", argIndex))
			args = append(args, "%"+search+"%")
			argIndex++
		}

	case "village":
		query = `
			SELECT v.id, v.name, v.district_id, d.name as district_name, v.status 
			FROM villages v 
			INNER JOIN districts d ON d.id = v.district_id 
			WHERE 1=1
		`

		if districtID > 0 {
			conditions = append(conditions, fmt.Sprintf("v.district_id = $%d", argIndex))
			args = append(args, districtID)
			argIndex++
		}

		if search != "" {
			conditions = append(conditions, fmt.Sprintf("v.name ILIKE $%d", argIndex))
			args = append(args, "%"+search+"%")
			argIndex++
		}

	default:
		return nil, fmt.Errorf("unknown model: %s", model)
	}

	if len(conditions) > 0 {
		query += " AND " + strings.Join(conditions, " AND ")
	}

	query += " ORDER BY name ASC"

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("columns failed: %w", err)
	}

	var results []map[string]interface{}

	for rows.Next() {
		values := make([]interface{}, len(columns))
		pointers := make([]interface{}, len(columns))
		for i := range values {
			pointers[i] = &values[i]
		}

		if err := rows.Scan(pointers...); err != nil {
			return nil, fmt.Errorf("scan failed: %w", err)
		}

		row := make(map[string]interface{})
		for i, col := range columns {
			val := values[i]
			// Convert []byte to string for readability
			if b, ok := val.([]byte); ok {
				row[col] = string(b)
			} else {
				row[col] = val
			}
		}
		results = append(results, row)
	}

	return results, nil
}