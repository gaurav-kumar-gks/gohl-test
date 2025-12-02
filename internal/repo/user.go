package repo


import (
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"

	"gks.com/gohl-test/internal/models"
	"context"
)

type UserRepository struct {

	DB *pgx.Conn
	Logger *zap.Logger
}

func NewUserRepository(db *pgx.Conn, logger *zap.Logger) *UserRepository {
	return &UserRepository{DB: db, Logger: logger}
}

func(u *UserRepository) ListUsers(ctx context.Context) ([]models.User, error) {
	query := `SELECT id, name, email, balance, created_at, updated_at FROM users ORDER BY id`
	rows, err := u.DB.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []models.User{}
	for rows.Next() {
		u := models.User{}
		rows.Scan(&u.Id, &u.Name, &u.Email, &u.Balance, &u.CreatedAt, &u.UpdatedAt)
		users = append(users, u)
	}
	return users, nil
}

func(u *UserRepository) CreateUser(ctx context.Context, model *models.CreateUser) (*models.User, error) {
	query := `INSERT INTO users (name, email)
			  VALUES ($1, $2)
			  RETURNING id, name, email, balance, created_at, updated_at`

	var user models.User
	err := u.DB.QueryRow(ctx, query, model.Name, model.Email).Scan(
		&user.Id,
		&user.Name,
		&user.Email,
		&user.Balance,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *UserRepository) UpdateUser(ctx context.Context, model *models.User, tx *pgx.Tx) (*models.User, error) {
	query := `UPDATE users SET balance = $1, updated_at = NOW() WHERE id = $2 RETURNING id, name, email, balance, created_at, updated_at`
	rows, err := u.DB.Query(ctx, query, model.Balance, model.Id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	user := models.User{}
	rows.Scan(&user.Id, &user.Name, &user.Email, &user.Balance, &user.CreatedAt, &user.UpdatedAt)
	return &user, nil
}

func (u *UserRepository) UpdateUserTx(
	ctx context.Context,
	tx pgx.Tx,
	model *models.User,
) (*models.User, error) {

	query := `
		UPDATE users
		SET balance = $1, updated_at = NOW()
		WHERE id = $2
		RETURNING id, name, email, balance, created_at, updated_at
	`

	var user models.User
	err := tx.QueryRow(ctx, query, model.Balance, model.Id).Scan(
		&user.Id,
		&user.Name,
		&user.Email,
		&user.Balance,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}


func(u *UserRepository) GetUser(ctx context.Context, id string) (*models.User, error) {
	query := `SELECT id, name, email, balance, created_at, updated_at FROM users WHERE id = $1`

	var user models.User
	err := u.DB.QueryRow(ctx, query, id).Scan(
		&user.Id,
		&user.Name,
		&user.Email,
		&user.Balance,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}