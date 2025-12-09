package email

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ClientSaved struct {
	Id        uuid.UUID
	Name      string
	Email     string
	ApikeyId  uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
}

type EmailSentSave struct {
	Id          uuid.UUID
	ClientId    uuid.UUID
	Body        string
	Subject     string
	ClientEmail string
}

type EmailRepo interface {
	SaveClient(ctx context.Context, name string, email string, apikey string) (*ClientSaved, error)
	SaveEmailSent(ctx context.Context, client *ClientSaved, subject string, body string) (*EmailSentSave, error)
}

type PgPool struct {
	db *pgxpool.Pool
}

func NewEmailRepo(db *pgxpool.Pool) EmailRepo {
	return &PgPool{db: db}
}

func (r *PgPool) SaveClient(ctx context.Context, name string, email string, apikey string) (*ClientSaved, error) {
	var client ClientSaved
	query := `
		INSERT INTO clients (id, name, email, api_key_id)
		VALUES ($1, $2, $3, (SELECT id FROM apikeys WHERE apikey = $4))
		RETURNING id, name, email, api_key_id, created_at, updated_at
	`

	id := uuid.New()

	err := r.db.QueryRow(
		ctx,
		query,
		id,
		name,
		email,
		apikey, // text, not uuid
	).Scan(
		&client.Id,
		&client.Name,
		&client.Email,
		&client.ApikeyId,
		&client.CreatedAt,
		&client.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &client, nil
}
