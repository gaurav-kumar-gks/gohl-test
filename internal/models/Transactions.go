package models

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Transactions struct {
	Id uuid.UUID `json:"id" validate:"required"`
	UserId uuid.UUID `json:"user_id" validate:"required"`
	Type string `json:"type" validate:"required,oneof=debit credit"`
	Amount decimal.Decimal `json:"amount" validate:"required,numeric"`
	Description string `json:"description" validate:"required,max=255"`
	CreatedAt time.Time `json:"created_at" validate:"required"`
	UpdatedAt time.Time `json:"updated_at" validate:"required"`
}


type CreateTransactions struct {
	UserId uuid.UUID `json:"user_id" validate:"required"`
	Type string `json:"type" validate:"required,oneof=debit credit"`
	Amount decimal.Decimal `json:"amount" validate:"required,numeric"`
	Description string `json:"description" validate:"required,max=255"`
}


var txnValidator = validator.New()

func ValidateTransactions(Transactions *CreateTransactions) error {
	return txnValidator.Struct(Transactions)
}

/*

Models:

Accounts
- userId
- balance
- ? 

Transactions
- amount
- type ? (debit | credit)
- metadata

making payment:
- debit txn

deposit fund
- credit txn

atomicity of txn and account ?
- 

APIs:
- Create user (balance 0): User
- current balance 

- Deposit fund, make payment: User, Txn
- get transactions: User, Transactions
*/