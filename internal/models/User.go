package models


import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)


type User struct {
	Id uuid.UUID `json:"id" validate:"required"`
	Name string `json:"name" validate:"required,max=255"`
	Email string `json:"email" validate:"required,email"`
	Balance decimal.Decimal `json:"balance" validate:"required,numeric"`
	CreatedAt time.Time `json:"created_at" validate:"required"`
	UpdatedAt time.Time `json:"updated_at" validate:"required"`
}

type CreateUser struct {
	Name string `json:"name" validate:"required,max=255"`
	Email string `json:"email" validate:"required,email"`
}


var userValidator = validator.New()

func ValidateUser(User *CreateUser) error {
	return userValidator.Struct(User)
}