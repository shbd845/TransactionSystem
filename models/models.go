package models

import (
	"database/sql"
	"fmt"
	"strconv"
)

type Account struct {
	AccountID int    `json:"account_id"`
	Balance   string `json:"balance"`
}

func CreateAccount(db *sql.DB, accountID int, initialBalance string) error {
	query := `INSERT INTO accounts (account_id, balance) VALUES ($1, $2)`
	_, err := db.Exec(query, accountID, initialBalance)
	return err
}

func GetAccount(db *sql.DB, accountID int) (Account, error) {
	var acc Account
	query := `SELECT account_id, balance FROM accounts WHERE account_id = $1`
	row := db.QueryRow(query, accountID)
	err := row.Scan(&acc.AccountID, &acc.Balance)
	if err != nil {
		if err == sql.ErrNoRows {
			return acc, fmt.Errorf("no account has been created with this account_id")
		}
		return acc, err
	}
	return acc, nil
}

func ProcessTransaction(db *sql.DB, srcID int, dstID int, amount string) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var balanceStr string
	checkBalance := `SELECT balance FROM accounts WHERE account_id = $1`
	err = db.QueryRow(checkBalance, srcID).Scan(&balanceStr)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("source account %v not found", srcID)
		}
		return fmt.Errorf("failed to fetch source account balance: %v", err)
	}
	// log.Println("Current Balance:", balanceStr)
	balance, err := strconv.ParseFloat(balanceStr, 32)
	// log.Println("Current Balance:", balance)
	if err != nil {
		return fmt.Errorf("balance conversion failed")
	}

	transactionAmount, err := strconv.ParseFloat(amount, 32)
	// log.Println("Transaction Amount:", transactionAmount)
	if err != nil {
		return fmt.Errorf("amount conversion failed")
	}

	if balance < transactionAmount {
		return fmt.Errorf("user %v doesn't have the required balance for transaction, current balance: %v", srcID, balance)
	}
	updateSrc := `UPDATE accounts SET balance = balance - $1 WHERE account_id = $2`
	_, err = tx.Exec(updateSrc, amount, srcID)
	if err != nil {
		return fmt.Errorf("failed to update source account: %v", err)
	}

	updateDst := `UPDATE accounts SET balance = balance + $1 WHERE account_id = $2`
	_, err = tx.Exec(updateDst, amount, dstID)
	if err != nil {
		return fmt.Errorf("failed to update destination account: %v", err)
	}

	logQuery := `INSERT INTO transactions (source_account_id, destination_account_id, amount) VALUES ($1, $2, $3)`
	_, err = tx.Exec(logQuery, srcID, dstID, amount)
	if err != nil {
		return fmt.Errorf("failed to log transaction: %v", err)
	}

	return tx.Commit()
}
