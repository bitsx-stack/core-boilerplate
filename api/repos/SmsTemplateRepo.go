package repos

import (
	"database/sql"

	"bitpos.com/m/models"
)

type SMSTemplateRepo struct {
	db *sql.DB
}

func NewSMSTemplateRepo(db *sql.DB) *SMSTemplateRepo {
	return &SMSTemplateRepo{db: db}
}

func (r *SMSTemplateRepo) CreateTemplate(t models.SMSTemplate) error {
	_, err := r.db.Exec(`
	INSERT INTO sms_templates (name,content,created_by)
	VALUES ($1,$2,$3)
	`,
		t.Name,
		t.Content,
		t.CreatedBy,
	)
	return err
}