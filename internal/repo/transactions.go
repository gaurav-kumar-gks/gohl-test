package repo


import (
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"

	"gks.com/gohl-test/internal/models"
	"context"
)

type TransactionsRepository struct {
	DB *pgx.Conn
	Logger *zap.Logger
}


func NewTransactionsRepository(db *pgx.Conn, logger *zap.Logger) *TransactionsRepository {
	return &TransactionsRepository{DB: db, Logger: logger}
}

func(u *TransactionsRepository) ListTransactions(ctx context.Context) ([]models.Transactions, error) {
	query := `SELECT id, type, amount, description, created_at, updated_at FROM transactions ORDER BY id`
	rows, err := u.DB.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	transactions := []models.Transactions{}
	for rows.Next() {
		t := models.Transactions{}
		rows.Scan(&t.Id, &t.Type, &t.Amount, &t.Description, &t.CreatedAt, &t.UpdatedAt)
		transactions = append(transactions, t)
	}
	return transactions, nil
}



func (u *TransactionsRepository) CreateTransactionsTx(
	ctx context.Context,
	tx pgx.Tx,
	model *models.CreateTransactions,
) (*models.Transactions, error) {

	query := `
		INSERT INTO transactions (user_id, type, amount, description)
		VALUES ($1, $2, $3, $4)
		RETURNING id, user_id, type, amount, description, created_at, updated_at
	`

	var t models.Transactions
	err := tx.QueryRow(ctx, query,
		model.UserId,
		model.Type,
		model.Amount,
		model.Description,
	).Scan(
		&t.Id,
		&t.UserId,
		&t.Type,
		&t.Amount,
		&t.Description,
		&t.CreatedAt,
		&t.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &t, nil
}


func(u *TransactionsRepository) CreateTransactions(ctx context.Context, model *models.CreateTransactions) (*models.Transactions, error) {
	query := `INSERT INTO transactions (user_id, type, amount, description) VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at`
	rows, err := u.DB.Query(ctx, query, model.UserId, model.Type, model.Amount, model.Description)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	transaction := models.Transactions{}
	rows.Scan(&transaction.Id, &transaction.UserId, &transaction.Type, &transaction.Amount, &transaction.Description, &transaction.CreatedAt, &transaction.UpdatedAt)
	return &transaction, nil
}