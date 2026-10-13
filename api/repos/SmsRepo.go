package repos

import (
	"database/sql"

	"bitpos.com/m/models"
)

type SMSRepo struct {
	db *sql.DB
}

func NewSMSRepo(db *sql.DB) *SMSRepo {
	return &SMSRepo{db: db}
}

func (r *SMSRepo) CreateSMS(s models.SMS) error {
	_, err := r.db.Exec(`
	INSERT INTO sms (content,sender,receiver,status,created_by)
	VALUES ($1,$2,$3,$4,$5)
	`,
		s.Content,
		s.Sender,
		s.Receiver,
		s.Status,
		s.CreatedBy,
	)
	return err
}