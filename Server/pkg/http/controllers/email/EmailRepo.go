package email

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const ErrInvalidAPIKey = "ERROR INVALID API KEY"

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
	Createdt    time.Time
	UpdatedAt   time.Time
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
	client := &ClientSaved{}

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	query := `
		WITH ap AS (
			SELECT id FROM apikeys WHERE apikey = $4
		)
		INSERT INTO clients (id, name, email, api_key_id)
		SELECT $1, $2, $3, ap.id
		FROM ap
		RETURNING id, name, email, api_key_id, created_at, updated_at
	`

	id := uuid.New()

	err := r.db.QueryRow(
		ctx,
		query,
		id,
		name,
		email,
		apikey,
	).Scan(
		&client.Id,
		&client.Name,
		&client.Email,
		&client.ApikeyId,
		&client.CreatedAt,
		&client.UpdatedAt,
	)

	// if no matching API key, ap CTE returns no row → ErrNoRows
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New(ErrInvalidAPIKey)
		}
		return nil, err
	}

	return client, nil
}
func (r *PgPool) SaveEmailSent(ctx context.Context, client *ClientSaved, subject string, body string) (*EmailSentSave, error) {
	var emailSent EmailSentSave
	id := uuid.New()
	query := `INSERT INTO emails_sent (id , client_id, body , subject, client_email) VALUES ($1, $2, $3, $4, $5) RETURNING id, client_id, body, subject, client_email, created_at, updated_at`

	err := r.db.QueryRow(ctx, query, id, client.Id, body, subject, client.Email).Scan(&emailSent.Id, &emailSent.ClientId, &emailSent.Body, &emailSent.Subject, &emailSent.ClientEmail, &emailSent.Createdt, &emailSent.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return &emailSent, nil
}
