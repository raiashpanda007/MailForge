package types

type LoginCredentials struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type SignUpCredentials struct {
	Email    string `json:"email" validate:"required,email"`
	Name     string `json:"name" validate:"required"`
	Password string `json:"password" validate:"required,min=6"`
}

type GenerateApiKeysCredentials struct {
	Organization     string `json:"organization" validate:"required,organization"`
	EmailAppPassword string `json:"emailAppPassword" validate:"omitempty,emailAppPassword"`
}
