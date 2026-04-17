package domain

type Key struct {
	ID          uint   `db:"id"          json:"id"`
	KeyValue    string `db:"key_value"   json:"key_value"`
	KeyType     string `db:"key_type"    json:"key_type"`
	Description string `db:"description" json:"description"`
}
