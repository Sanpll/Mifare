package dto

type CreateCardInput struct {
	CardNumber string `json:"number"     binding:"required,hexadecimal,min=8,max=14"`
	Balance    string `json:"balance"    binding:"required"`
	IsBlocked  bool   `json:"is_blocked" binding:"omitempty"`
	KeyValue   string `json:"key_value"  binding:"required,len=12,hexadecimal"`
}

type CardResponse struct {
	ID         uint   `json:"id"`
	CardNumber string `json:"number"`
	Balance    string `json:"balance"`
	IsBlocked  bool   `json:"is_blocked"`
	OwnerName  string `json:"owner_name"`
	KeyValue   string `json:"key_value"`
}

type CardsResponse struct {
	Cards []CardResponse `json:"cards"`
}

type CardUpdate struct {
	Balance   *string `json:"balance"    binding:"omitempty"`
	IsBlocked *bool   `json:"is_blocked" binding:"omitempty"`
	OwnerName *string `json:"owner_name" binding:"omitempty,min=3,max=100"`
	KeyValue  *string `json:"key_value"  binding:"omitempty,len=12,hexadecimal"`
}
