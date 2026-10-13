package repos

import (
	"database/sql"

	"bitpos.com/m/models"
)

type AuditRepo struct {
	db *sql.DB
}

func NewAuditRepo(db *sql.DB) *AuditRepo {
	return &AuditRepo{db: db}
}

func (r *AuditRepo) Log(a models.AuditLog) error {
	_, err := r.db.Exec(`
	INSERT INTO audit_logs (actor,action,entity,entity_id,description)
	VALUES ($1,$2,$3,$4,$5)
	`,
		a.Actor,
		a.Action,
		a.Entity,
		a.EntityID,
		a.Description,
	)
	return err
}