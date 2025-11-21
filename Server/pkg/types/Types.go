package types

type LoginCredentials struct {
	Email    string `json:"email" validator:"required"`
	Password string `json:"password" validator:"required"`
}

type SignUpCredentials struct {
	Email            string `validate:"required,email"`
	Name             string `validate:"required"`
	Password         string `validate:"required,min=6"`
	Institution      string `validate:"required"`
	EmailAppPassword string `validate:"omitempty"`
}
