package apikeys

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/raiashpanda007/MailForge/pkg/http/controllers/auth"
)

type ApiKeyResult struct {
	Id           uuid.UUID `json:"id"`
	Organization string    `json:"organization"`
	Apikey       string    `json:"key"`
	EmailPass    string    `json:"emailpass"`
	Created_at   time.Time `json:"created_at"`
	Updated_at   time.Time `json:"updated_at"`
}

type ApiKeyService interface {
	GenerateKey(ctx context.Context, organization string, emaiAppPassword string) (*ApiKeyResult, error)
	DeleteKey(ctx context.Context, id string) (bool, error)
}

type apiKeyUtils struct{ db *pgxpool.Pool }

func NewApiKeyService(db *pgxpool.Pool) ApiKeyService {
	return &apiKeyUtils{db: db}
}

func (r *apiKeyUtils) GenerateKey(ctx context.Context, organization string, emailAppPassword string) (*ApiKeyResult, error) {
	var apiKeyresult ApiKeyResult
	user := ctx.Value("USER")
	if user == nil {
		return nil, errors.New("PLEASE LOGIN AGAIN")
	}
	verifiedUser, ok := user.(*auth.User)
	if !ok {
		return nil, errors.New("PLEASE PROVIDE A VALID LOGIN TOKEN , TRY AGAIN LOGIN ")
	}
	id := uuid.New()
	apiKey := uuid.New()
	err := r.db.QueryRow(ctx, "INSERT INTO apikeys (id, organization, apikey, email_app_password,user_id ) VALUES ($1, $2, $3, $4, $5) RETURNING id, organization, apikey, email_app_password, created_at, updated_at", id, organization, apiKey, emailAppPassword, verifiedUser.Id).Scan(&apiKeyresult.Id, &apiKeyresult.Organization, &apiKeyresult.Apikey, &apiKeyresult.EmailPass, &apiKeyresult.Created_at, &apiKeyresult.Updated_at)
	if err != nil {
		return nil, err
	}

	return &apiKeyresult, nil
}

func (r *apiKeyUtils) DeleteKey(ctx context.Context, id string) (bool, error) {
	// Parse key id
	keyId, err := uuid.Parse(id)
	if err != nil {
		return false, errors.New("invalid key id")
	}

	// Check user
	userValue := ctx.Value("USER")
	if userValue == nil {
		return false, errors.New("please login again")
	}

	verifiedUser, ok := userValue.(*auth.User)
	if !ok {
		return false, errors.New("invalid login token")
	}

	// Delete only if belongs to this user
	res, err := r.db.Exec(
		ctx,
		"DELETE FROM apikeys WHERE id = $1 AND user_id = $2",
		keyId,
		verifiedUser.Id,
	)
	if err != nil {
		return false, err
	}

	// Check if anything was deleted
	if res.RowsAffected() == 0 {
		return false, nil
	}

	return true, nil
}
