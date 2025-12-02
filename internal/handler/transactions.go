package handler


import (
	"encoding/json"
	"net/http"
	"fmt"

	"gks.com/gohl-test/internal/repo"
		"gks.com/gohl-test/internal/models"
	"go.uber.org/zap"
)

type TransactionsHandler struct {
	userRepo *repo.UserRepository
	repo *repo.TransactionsRepository
	logger *zap.Logger	
}

func NewTransactionsHandler(repo *repo.TransactionsRepository, userRepo *repo.UserRepository, logger *zap.Logger) *TransactionsHandler {
	return &TransactionsHandler{repo: repo, userRepo: userRepo, logger: logger}
}

func (t *TransactionsHandler) ListTransactions(w http.ResponseWriter, r *http.Request) {
	transactions, err := t.repo.ListTransactions(r.Context())
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, "Failed to fetch transactions")
		t.logger.Error("List transactions failed", zap.Error(err))
		return
	}
	jsonResponse(w, http.StatusOK, transactions)
}

func (t *TransactionsHandler) CreateTransactions(w http.ResponseWriter, r *http.Request) {
	var transaction models.CreateTransactions
	if err := json.NewDecoder(r.Body).Decode(&transaction); err != nil {
		errorResponse(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	if err := models.ValidateTransactions(&transaction); err != nil {
		message := fmt.Sprintf("invalid transaction: %v", err)
		errorResponse(w, http.StatusBadRequest, message)
		return
	}

	userId := transaction.UserId
	user, err := t.userRepo.GetUser(r.Context(), userId.String())
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, "Failed to get user")
		t.logger.Error("Get user failed", zap.Error(err))
		return
	}

	// Check balance for debit
	if transaction.Type == "debit" && user.Balance.LessThan(transaction.Amount) {
		errorResponse(w, http.StatusBadRequest, "Insufficient balance")
		return
	}

	ctx := r.Context()
	tx, err := t.repo.DB.Begin(ctx)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, "Failed to begin transaction")
		t.logger.Error("Begin transaction failed", zap.Error(err))
		return
	}
	defer tx.Rollback(ctx) // safe: no-op if already committed

	// update balance
	if transaction.Type == "debit" {
		user.Balance = user.Balance.Sub(transaction.Amount)
	} else {
		user.Balance = user.Balance.Add(transaction.Amount)
	}

	_, err = t.userRepo.UpdateUserTx(ctx, tx, user)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, "Failed to update user balance")
		t.logger.Error("Update user balance failed", zap.Error(err))
		return
	}

	transactionModel, err := t.repo.CreateTransactionsTx(ctx, tx, &transaction)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, "Failed to create transaction")
		t.logger.Error("Create transaction failed", zap.Error(err))
		return
	}

	if err := tx.Commit(ctx); err != nil {
		errorResponse(w, http.StatusInternalServerError, "Failed to commit transaction")
		t.logger.Error("Commit transaction failed", zap.Error(err))
		return
	}

	jsonResponse(w, http.StatusCreated, transactionModel)
}
