package repos

import (
	"database/sql"

	"bitpos.com/m/models"
)

type RegionRepo struct {
	db *sql.DB
}

func NewRegionRepo(db *sql.DB) *RegionRepo {
	return &RegionRepo{
		db: db,
	}
}

func (r *RegionRepo) GetAllRegions(country string) ([]*models.Region, error) {
	var regions []*models.Region
	query := "SELECT id, name FROM regions WHERE country = $1"
	rows, err := r.db.Query(query, country)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var region models.Region

		if err := rows.Scan(&region.ID, &region.Name); err != nil {
			return nil, err
		}

		regions = append(regions, &region)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return regions, nil
}
