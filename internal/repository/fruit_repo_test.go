package repository

import (
  "ATB_letsGO/internal/models"
  "context"
  "testing"
  "time"
  "github.com/DATA-DOG/go-sqlmock"
)

func TestCreateFruit(t *testing.T) {
  db, mock, err := sqlmock.New()
  if err != nil {
    t.Fatalf("Помилка створення мока: %s", err)
  }
  defer db.Close()

  repo := NewFruitRepository(db)
  fruit := models.Fruit{
    Name:       "Яблуко",
    Brand:      "Гала",
    PricePerKg: 45.50,
    StockKg:    120.0,
  }

  mock.ExpectQuery(`INSERT INTO fruits`).
    WithArgs(fruit.Name, fruit.Brand, fruit.PricePerKg, fruit.StockKg).
    WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

  id, err := repo.Create(context.Background(), fruit)

  if err != nil {
    t.Errorf("Очікувався err = nil, отримано %s", err)
  }
  if id != 1 {
    t.Errorf("Очікувався id = 1, отримано %d", id)
  }

  if err := mock.ExpectationsWereMet(); err != nil {
    t.Errorf("Не всі очікування бази виконались: %s", err)
  }
}

func TestGetFruitByID(t *testing.T) {
  db, mock, err := sqlmock.New()
  if err != nil {
    t.Fatalf("Помилка створення мока: %s", err)
  }
  defer db.Close()

  repo := NewFruitRepository(db)

  rows := sqlmock.NewRows([]string{"id", "name", "brand", "price_per_kg", "stock_kg", "created_at"}).
    AddRow(1, "Банан", "Еквадор", 60.00, 50.0, time.Now())

  mock.ExpectQuery(`SELECT id, name, brand, price_per_kg, stock_kg, created_at FROM fruits WHERE id = \$1`).
    WithArgs(1).
    WillReturnRows(rows)

  fruit, err := repo.GetByID(context.Background(), 1)

  if err != nil {
    t.Errorf("Очікувався err = nil, отримано %s", err)
  }
  if fruit.Name != "Банан" {
    t.Errorf("Очікувалося ім'я Банан, отримано %s", fruit.Name)
  }

  if err := mock.ExpectationsWereMet(); err != nil {
    t.Errorf("Не всі очікування бази виконались: %s", err)
  }
}
