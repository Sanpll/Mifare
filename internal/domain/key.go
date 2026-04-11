package domain

type Key struct {
	ID          uint   `db:"id" json:"id"`
	KeyValue    string `db:"key_value" json:"keyValue"`
	KeyType     string `db:"key_type" json:"keyType"`
	Description string `db:"description" json:"description"`
}
