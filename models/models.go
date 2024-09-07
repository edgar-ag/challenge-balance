package models

type Transaction struct {
	Amount string
	Date   string
}

type TemplateData struct {
	TotalBalance        float32
	AverageDebitAmount  float32
	AverageCreditAmount float32
	TxnsJanuary         int
	TxnsFebruary        int
	TxnsMarch           int
	TxnsApril           int
	TxnsMay             int
	TxnsJune            int
	TxnsJuly            int
	TxnsAgust           int
	TxnsSeptember       int
	TxnsOctober         int
	TxnsNovember        int
	TxnsDecember        int
	ImageSrc            string
}
