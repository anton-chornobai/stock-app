package model

import "time"

type Stock struct {
	Name      string    `json:"name"`
	Symbol    string    `json:"symbol"`
	ID        int       `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
}
