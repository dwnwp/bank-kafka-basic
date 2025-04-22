package models

type OpenAccount struct {
	AccountHolder  string
	AccountType    int
	OpeningBalance float64
}

type DepositFund struct {
	ID     string
	Amount float64
}

type WithdrawFund struct {
	ID     string
	Amount float64
}

type CloseAccount struct {
	ID string
}