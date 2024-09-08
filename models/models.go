package models

type CustomerInfo struct {
	AccountNumber string
	Name          string
	Email         string
}

type Transaction struct {
	Amount string
	Date   string
}

type BalanceInfo struct {
	TotalBalance  float32
	AverageDebit  float32
	AverageCredit float32
	TxnsJanuary   int
	TxnsFebruary  int
	TxnsMarch     int
	TxnsApril     int
	TxnsMay       int
	TxnsJune      int
	TxnsJuly      int
	TxnsAgust     int
	TxnsSeptember int
	TxnsOctober   int
	TxnsNovember  int
	TxnsDecember  int
	ImageSrc      string
	CustomerName  string
}
