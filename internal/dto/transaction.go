package dto

type CreateTransactionInput struct {
	CardNumber           string `json:"card_number"            binding:"required,hexadecimal,min=8,max=14"`
	Price                string `json:"price"                  binding:"required"`
	TerminalSerialNumber string `json:"terminal_serial_number" binding:"required,min=3,max=50"`
	Status               string `json:"status"                 binding:"omitempty,min=2"`
}

type TransactionResponse struct {
	ID                   uint   `json:"id"`
	CardNumber           string `json:"card_number"`
	Price                string `json:"price"`
	TerminalSerialNumber string `json:"terminal_serial_number"`
}

type TransactionsResponse struct {
	Transactions []TransactionResponse `json:"transactions"`
}

type TransactionUpdate struct {
	CardNumber           *string `json:"card_number"            binding:"omitempty,hexadecimal,min=8,max=14"`
	TerminalSerialNumber *string `json:"terminal_serial_number" binding:"omitempty,min=3,max=50"`
	Status               *string `json:"status"                 binding:"omitempty,min=2"`
}

type AuthorizeTransactionInput struct {
	CardNumber           string `json:"card_number"            binding:"required,hexadecimal,min=8,max=14"`
	Price                string `json:"price"                  binding:"required"`
	TerminalSerialNumber string `json:"terminal_serial_number" binding:"omitempty,min=3,max=50"`
}

type AuthorizeTransactionResponse struct {
	Authorized bool   `json:"authorized"`
	Message    string `json:"message"`
}
