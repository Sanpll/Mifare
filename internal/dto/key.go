package dto

type CreateKeyInput struct {
	KeyValue    string `json:"value"       binding:"required,len=12,hexadecimal"`
	KeyType     string `json:"type"        binding:"required,oneof=A B"`
	Description string `json:"description" binding:"required,min=5,max=200"`
}

type KeyResponse struct {
	ID          uint   `json:"id"`
	KeyValue    string `json:"value"`
	KeyType     string `json:"type"`
	Description string `json:"description"`
}

type KeysResponse struct {
	Keys []KeyResponse `json:"keys"`
}

type KeyUpdate struct {
	KeyValue    *string `json:"value"       binding:"omitempty,len=12,hexadecimal"`
	KeyType     *string `json:"type"        binding:"omitempty,oneof=A B"`
	Description *string `json:"description" binding:"omitempty,min=5,max=200"`
}