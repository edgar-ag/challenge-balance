package repository

import (
	"challenge/balance/models"
	"context"
)

type Repository interface {
	GetCustomerId(ctx context.Context, customerInfo *models.CustomerInfo) (int64, error)
	InsertCustomerInfo(ctx context.Context, customerInfo *models.CustomerInfo) (int64, error)
	InsertTransaction(ctx context.Context, customertId int64, txn *models.Transaction) error
	Close() error
}

var implementation Repository

func SetRepository(repository Repository) {
	implementation = repository
}

func GetCustomerId(ctx context.Context, customerInfo *models.CustomerInfo) (int64, error) {
	return implementation.GetCustomerId(ctx, customerInfo)
}

func InsertCustomerInfo(ctx context.Context, customerInfo *models.CustomerInfo) (int64, error) {
	return implementation.InsertCustomerInfo(ctx, customerInfo)
}

func InsertTransaction(ctx context.Context, customertId int64, txn *models.Transaction) error {
	return implementation.InsertTransaction(ctx, customertId, txn)
}

func Close() error {
	return implementation.Close()
}
