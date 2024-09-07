package service

import (
	"challenge/balance/models"
	"testing"
)

func TestProccessFile(t *testing.T) {
	expectedTxns := []*models.Transaction{
		{Amount: "7/15", Date: "+60.5"},
		{Amount: "7/28", Date: "-10.3"},
		{Amount: "8/2", Date: "-20.46"},
		{Amount: "8/13", Date: "+10"},
	}
	txns, _ := ProccessFile("txns_test.csv")
	for i, txn := range txns {
		if txn.Amount != expectedTxns[i].Amount && txn.Date != expectedTxns[i].Date {
			t.Errorf("Transaction was incorrect, got %v expected %v", txn, expectedTxns[i])
		}
	}
}

func TestCetTotalOf(t *testing.T) {
	cases := []struct {
		txns           []float32
		expectedResult float32
	}{
		{[]float32{50.1, 10.2, 11.0, 0.2}, 71.5},
		{[]float32{60.5, -10.3, -20.46, 10}, 39.74},
	}
	for _, c := range cases {
		result := calcTotalOf(c.txns)
		if result != c.expectedResult {
			t.Errorf("Total was incorrect, got %v expected %v", result, c.expectedResult)
		}
	}
}

func TestHasFormat(t *testing.T) {
	cases := []struct {
		input          string
		regex          string
		expectedResult bool
	}{
		{"8/21", RegexMonths, true},
		{"13/33", RegexMonths, false},
		{"-12.5", RegexTxn, true},
		{"+0.5", RegexTxn, true},
		{"-.1", RegexTxn, false},
	}
	for _, c := range cases {
		result := hasFormat(c.input, c.regex)
		if result != c.expectedResult {
			t.Errorf("Validation was incorrect, got %v expected %v", result, c.expectedResult)
		}
	}
}
