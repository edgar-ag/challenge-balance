package service

import (
	"challenge/balance/models"
	"challenge/balance/repository"
	"context"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

const (
	directory   = "./data/"
	RegexMonths = `\b([1-9]|1[0-2])\/([0-9]|[012][0-9]|3[01])\b`
	RegexTxn    = `^[+-]?\d+(\.\d+)?`
)

type Balance struct {
	TxnsFile     string
	CustomerInfo *models.CustomerInfo
}

func NewBalance(txnsFile string, customerInfo *models.CustomerInfo) *Balance {
	return &Balance{
		TxnsFile:     txnsFile,
		CustomerInfo: customerInfo,
	}
}

func (b *Balance) ProccessFile() ([]*models.Transaction, error) {
	var txns []*models.Transaction
	file, err := os.Open(directory + b.TxnsFile)
	if err != nil {
		return txns, err
	}

	defer file.Close()

	reader := csv.NewReader(file)
	for {
		record, err := reader.Read()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			return txns, err
		}
		isDate := hasFormat(record[1], RegexMonths)
		isTxn := hasFormat(record[2], RegexTxn)
		if isDate && isTxn {
			txn := &models.Transaction{
				Date:   record[1],
				Amount: record[2],
			}
			txns = append(txns, txn)
		}
	}

	log.Println("the file was proccessed.")
	return txns, err
}

func (b *Balance) GetBalanceInfo(txns []*models.Transaction) (*models.BalanceInfo, error) {
	var debitTxns []float32
	var creditTxns []float32
	txnsPerMonth := []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

	for _, txn := range txns {
		month := getMothsTxn(txn.Date)
		if month > 0 {
			txnsPerMonth[month-1] += 1
		}

		txn := getTxn(txn.Amount)
		if txn > 0.0 {
			creditTxns = append(creditTxns, txn)
		} else if txn < 0.0 {
			debitTxns = append(debitTxns, txn)
		} else {
			continue
		}
	}

	totalDebit := calcTotalOf(debitTxns)
	totalCredit := calcTotalOf(creditTxns)
	totalBalance := totalDebit + totalCredit
	averageDebit := totalDebit / float32(len(debitTxns))
	averageCredit := totalCredit / float32(len(creditTxns))

	log.Println("the balance was calculated.")
	return &models.BalanceInfo{
		TotalBalance:  totalBalance,
		AverageDebit:  averageDebit,
		AverageCredit: averageCredit,
		TxnsJanuary:   txnsPerMonth[0],
		TxnsFebruary:  txnsPerMonth[1],
		TxnsMarch:     txnsPerMonth[2],
		TxnsApril:     txnsPerMonth[3],
		TxnsMay:       txnsPerMonth[4],
		TxnsJune:      txnsPerMonth[5],
		TxnsJuly:      txnsPerMonth[6],
		TxnsAgust:     txnsPerMonth[7],
		TxnsSeptember: txnsPerMonth[8],
		TxnsOctober:   txnsPerMonth[9],
		TxnsNovember:  txnsPerMonth[10],
		TxnsDecember:  txnsPerMonth[11],
		CustomerName:  b.CustomerInfo.Name,
	}, nil
}

func (b *Balance) InsertDataIntoDB(ctx context.Context, txns []*models.Transaction) {
	defer repository.Close()

	customerId, err := repository.InsertCustomerInfo(ctx, b.CustomerInfo)
	if err != nil {
		log.Printf("error trying to insert data into DB. %v\n", err)
	}
	log.Println("the customer info was inserted into DB.")
	for _, txn := range txns {
		err := repository.InsertTransaction(ctx, customerId, txn)
		if err != nil {
			log.Printf("error trying to insert data into DB. %v\n", err)
		}
	}
	log.Println("the transactions were inserted into DB.")
}

func getMothsTxn(date string) int {
	str := strings.Split(date, "/")
	month, err := strconv.ParseInt(str[0], 10, 16)
	if err != nil {
		fmt.Println(err.Error())
		return 0
	}
	return int(month)
}

func getTxn(record string) float32 {
	txn, err := strconv.ParseFloat(record, 64)
	if err != nil {
		fmt.Println(err.Error())
		return 0.0
	}
	return float32(txn)
}

func calcTotalOf(txns []float32) float32 {
	total := 0.0
	for _, txn := range txns {
		total += float64(txn)
	}
	return float32(total)
}

func hasFormat(str, expression string) bool {
	str = strings.TrimSpace(str)
	regex := regexp.MustCompile(expression)
	return regex.MatchString(str)
}
