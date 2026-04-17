package domain

type Terminal struct {
	ID           uint   `db:"id"            json:"id"`
	SerialNumber string `db:"serial_number" json:"serial_number"`
	Address      string `db:"address"       json:"address"`
	Name         string `db:"name"          json:"name"`
}
