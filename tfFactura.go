package producer

import "time"

type TfFactura struct {
	USER             int        `json:"[user]"`
	WH               int        `json:"wh"`
	FACTURA          string     `json:"factura"`
	CashN            int        `json:"cash_n"`
	CheckN           int        `json:"check_n"`
	Date             time.Time  `json:"[date]"`
	Customer         *string    `json:"customer"`
	Rnn              *string    `json:"rnn"`
	IinBin           *string    `json:"iin_bin"`
	Legal            bool       `json:"legal"`
	FirmId           *int       `json:"firm_id"`
	FacturaYear      int        `json:"factura_year"`
	Account          *string    `json:"account"`
	BankName         *string    `json:"bank_name"`
	Bik              *string    `json:"bik"`
	Address          *string    `json:"address"`
	DoverennostNomer *string    `json:"doverennost_nomer"`
	DoverennostData  *time.Time `json:"doverennost_data"`
	DoverennostKomu  *string    `json:"doverennost_komu"`
	NDog             *int       `json:"n_dog"`
	DataDog          *time.Time `json:"data_dog"`
	BaseCashN        *int       `json:"base_cash_n"`
	BaseCheckN       *int       `json:"base_check_n"`
	SaveDate         *time.Time `json:"save_date"`
	Contacts         *string    `json:"contacts"`
	Approved         *bool      `json:"approved"`
}
type TfFacturaLoad struct {
	Count int         `json:"CountItems"`
	Items []TfFactura `json:"Items"`
}
