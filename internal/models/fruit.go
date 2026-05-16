package models

import "time"

type Fruit struct {
  ID         int       `json:"id"`
  Name       string    `json:"name"`
  Brand      string    `json:"brand"`
  PricePerKg float64   `json:"price_per_kg"`
  StockKg    float64   `json:"stock_kg"`
  CreatedAt  time.Time `json:"created_at"`
}
