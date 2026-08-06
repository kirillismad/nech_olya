package notesdb

import (
	"context"
	"database/sql"
	"errors"
	"testing"
)

// Тесты (in-memory + MigrateAccounts; helper заводит счета (1, "alice", 100) и (2, "bob", 50)):
// - TestTransfer_Success — после Transfer(ctx, db, 1, 2, 30) балансы стали 70 и 80
// - TestTransfer_InsufficientFunds — Transfer(ctx, db, 1, 2, 9999) возвращает ErrInsufficientFunds, балансы не изменились
// - TestTransfer_RollbackPreservesTotal — после неуспешного перевода SELECT sum(balance) FROM accounts равно исходной сумме (150)
// - TestTransfer_AccountNotFound — перевод с fromID = 999 возвращает ErrAccountNotFound

func helper2(t *testing.T) *sql.DB {
	db, err := OpenInMemory()
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	err = MigrateAccounts(ctx, db)
	if err != nil {
		t.Fatal(err)
	}

	const query = `INSERT INTO accounts VALUES(?,?,?)`

	_, err = db.ExecContext(ctx, query, 1, "alice", 100)
	if err != nil {
		t.Fatal(err)
	}
	_, err = db.ExecContext(ctx, query, 2, "bob", 50)
	if err != nil {
		t.Fatal(err)
	}
	return db
}

func TestTransfer_Success(t *testing.T) {
	db := helper2(t)
	ctx := context.Background()
	err := Transfer(ctx, db, 1, 2, 30)
	if err != nil {
		t.Fatal(err)
	}

	var aliceBalance, bobBalance int64
	_ = db.QueryRowContext(ctx, `SELECT balance FROM accounts WHERE id = 1`).Scan(&aliceBalance)
	_ = db.QueryRowContext(ctx, `SELECT balance FROM accounts WHERE id = 2`).Scan(&bobBalance)
	if aliceBalance != 70 {
		t.Fatal("alice balance expected 70")
	}
	if bobBalance != 80 {
		t.Fatal("bob balance expected 80")
	}

}

func TestTransfer_InsufficientFunds(t *testing.T) {
	db := helper2(t)
	ctx := context.Background()
	var startBalanceAlice int64
	_ = db.QueryRowContext(ctx, `SELECT balance FROM accounts WHERE id = 1`).Scan(&startBalanceAlice)
	err := Transfer(ctx, db, 1, 2, 9999)
	if !errors.Is(err, ErrInsufficientFunds) {
		t.Fatal(err)
	}
	var finalBalanceAlice int64
	_ = db.QueryRowContext(ctx, `SELECT balance FROM accounts WHERE id = 1`).Scan(&finalBalanceAlice)
	if startBalanceAlice != finalBalanceAlice {
		t.Fatal("balance written off")
	}
}

func TestTransfer_RollbackPreservesTotal(t *testing.T) {
	db := helper2(t)
	ctx := context.Background()
	err := Transfer(ctx, db, 1, 2, 9999)
	if !errors.Is(err, ErrInsufficientFunds) {
		t.Fatal(err)
	}
	var aliceBalance, bobBalance int64
	_ = db.QueryRowContext(ctx, "SELECT balance FROM accounts WHERE id = 1").Scan(&aliceBalance)
	_ = db.QueryRowContext(ctx, "SELECT balance FROM accounts WHERE id = 2").Scan(&bobBalance)
	sum := aliceBalance + bobBalance
	if sum != 150 {
		t.Fatal(err)
	}
}

func TestTransfer_AccountNotFound(t *testing.T) {
	db := helper2(t)
	ctx := context.Background()
	err := Transfer(ctx, db, 999, 2, 100)
	if !errors.Is(err, ErrAccountNotFound) {
		t.Fatal(err)
	}
}
