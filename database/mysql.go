package database

import (
	"challenge/balance/models"
	"context"
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

type MysqlRepository struct {
	db *sql.DB
}

func NewMysqlRepository(url string) (*MysqlRepository, error) {
	db, err := sql.Open("mysql", url)
	if err != nil {
		return nil, err
	}
	return &MysqlRepository{db}, nil
}

// Find if the customer info already exists.
func (r *MysqlRepository) GetCustomerId(ctx context.Context, customerInfo *models.CustomerInfo) (int64, error) {
	var customerId int64
	query := "SELECT id FROM customer WHERE account_number = ? AND name = ? AND email = ?"
	err := r.db.QueryRowContext(ctx, query, customerInfo.AccountNumber, customerInfo.Name, customerInfo.Email).Scan(&customerId)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		} else {
			return 0, err
		}
	}
	return customerId, nil
}

func (r *MysqlRepository) InsertCustomerInfo(ctx context.Context, customerInfo *models.CustomerInfo) (int64, error) {
	customerId, err := r.GetCustomerId(ctx, customerInfo)
	if err != nil {
		return 0, err
	}
	if customerId > 0 {
		//If the customer exists, return the id
		return customerId, nil
	}

	query := "INSERT INTO customer (account_number, name, email) VALUES (?, ?, ?)"
	result, err := r.db.ExecContext(ctx, query, customerInfo.AccountNumber, customerInfo.Name, customerInfo.Email)
	if err != nil {
		return 0, err
	}
	customerId, err = result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return customerId, nil
}

func (r *MysqlRepository) InsertTransaction(ctx context.Context, accountId int64, txn *models.Transaction) error {
	query := "INSERT INTO transaction (account_id, transaction, date) VALUES (?, ?, ?)"
	_, err := r.db.ExecContext(ctx, query, accountId, txn.Amount, txn.Date)
	if err != nil {
		return err
	}
	return nil
}

func (r *MysqlRepository) Close() error {
	return r.db.Close()
}
