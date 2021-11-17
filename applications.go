package main

import "time"

type Application struct {
	Uuid        string    `json:"id,omitempty" gorm:"primaryKey"`
	Company     string    `json:"company,omitempty"`
	Position    string    `json:"position,omitempty"`
	Link        string    `json:"link,omitempty"`
	Status      string    `json:"status,omitempty"`
	LastUpdated time.Time `json:"last_updated,omitempty"`
}
