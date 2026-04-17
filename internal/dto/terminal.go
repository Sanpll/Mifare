package dto

type CreateTerminalInput struct {
	SerialNumber string `json:"serial_number" binding:"required,min=3,max=50"`
	Address      string `json:"address"       binding:"required,min=5,max=200"`
	Name         string `json:"name"          binding:"required,min=3,max=100"`
}

type TerminalResponse struct {
	ID           uint   `json:"id"`
	SerialNumber string `json:"serial_number"`
	Address      string `json:"address"`
	Name         string `json:"name"`
}

type TerminalsResponse struct {
	Terminals []TerminalResponse `json:"terminals"`
}

type TerminalUpdate struct {
	SerialNumber *string `json:"serial_number" binding:"omitempty,min=3,max=50"`
	Address      *string `json:"address"       binding:"omitempty,min=5,max=200"`
	Name         *string `json:"name"          binding:"omitempty,min=3,max=100"`
}
