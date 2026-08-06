package notesdb

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// В пакет добавить:
// - функцию MigrateAccounts(ctx, db) error — создаёт accounts(id INTEGER PRIMARY KEY, owner TEXT NOT NULL, balance INTEGER NOT NULL CHECK(balance >= 0))
// - ошибки ErrAccountNotFound и ErrInsufficientFunds
// - функцию Transfer(ctx context.Context, db *sql.DB, fromID, toID int64, amount int64) error:
// - открывает транзакцию через db.BeginTx
// - сразу ставит defer tx.Rollback() (после Commit он безопасен)
// - читает текущий баланс отправителя через tx.QueryRowContext
// - если баланса не хватает — возвращает ErrInsufficientFunds (Rollback сработает автоматически)
// - выполняет два UPDATE и Commit

var ErrAccountNotFound = errors.New("account not found")
var ErrInsufficientFunds = errors.New("insufficient funds")

func MigrateAccounts(ctx context.Context, db *sql.DB) error {
	const query = `CREATE TABLE IF NOT EXISTS accounts (id INTEGER PRIMARY KEY, owner TEXT NOT NULL, balance INTEGER NOT NULL CHECK(balance >= 0))`
	_, err := db.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to create the table: %w", err)
	}
	return nil
}

func Transfer(ctx context.Context, db *sql.DB, fromID, toID int64, amount int64) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback()
	var balance int64
	err = tx.QueryRowContext(ctx, "SELECT balance FROM accounts WHERE id = ?", fromID).Scan(&balance)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrAccountNotFound
		}
		return fmt.Errorf("failed to get sender balance: %w", err)
	}
	if balance < amount {
		return ErrInsufficientFunds
	}
	query := `UPDATE accounts SET balance = balance - ? WHERE id = ?`
	_, err = tx.ExecContext(ctx, query, amount, fromID)
	if err != nil {
		return fmt.Errorf("unable to debit the funds: %w", err)
	}

	query = `UPDATE accounts SET balance = balance + ? WHERE id=?`
	result, err := tx.ExecContext(ctx, query, amount, toID)
	if err != nil {
		return fmt.Errorf("failed to credit balance: %w", err)
	}

	rowAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed row affected: %w", err)
	}

	if rowAffected == 0 {
		return ErrAccountNotFound
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}
