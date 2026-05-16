package repository

import (
	"ATB_letsGO/internal/models"
	"context"
	"database/sql"
)

type FruitRepository struct {
	db *sql.DB
}

func NewFruitRepository(db *sql.DB) *FruitRepository {
	return &FruitRepository{db: db}
}

func (r *FruitRepository) Create(ctx context.Context, fruit models.Fruit) (int, error) {
	var id int
	query := `INSERT INTO fruits (name, brand, price_per_kg, stock_kg) VALUES ($1, $2, $3, $4) RETURNING id`
	err := r.db.QueryRowContext(ctx, query, fruit.Name, fruit.Brand, fruit.PricePerKg, fruit.StockKg).Scan(&id)
	return id, err
}

func (r *FruitRepository) GetAll(ctx context.Context) ([]models.Fruit, error) {
	query := `SELECT id, name, brand, price_per_kg, stock_kg, created_at FROM fruits`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var fruits []models.Fruit
	for rows.Next() {
		var f models.Fruit
		if err := rows.Scan(&f.ID, &f.Name, &f.Brand, &f.PricePerKg, &f.StockKg, &f.CreatedAt); err != nil {
			return nil, err
		}
		fruits = append(fruits, f)
	}
	return fruits, nil
}

func (r *FruitRepository) GetByID(ctx context.Context, id int) (models.Fruit, error) {
	query := `SELECT id, name, brand, price_per_kg, stock_kg, created_at FROM fruits WHERE id = $1`
	var f models.Fruit
	err := r.db.QueryRowContext(ctx, query, id).Scan(&f.ID, &f.Name, &f.Brand, &f.PricePerKg, &f.StockKg, &f.CreatedAt)
	return f, err
}

func (r *FruitRepository) Update(ctx context.Context, id int, fruit models.Fruit) error {
	query := `UPDATE fruits SET name = $1, brand = $2, price_per_kg = $3, stock_kg = $4 WHERE id = $5`
	_, err := r.db.ExecContext(ctx, query, fruit.Name, fruit.Brand, fruit.PricePerKg, fruit.StockKg, id)
	return err
}

func (r *FruitRepository) UpdatePrice(ctx context.Context, id int, price float64) error {
	query := `UPDATE fruits SET price_per_kg = $1 WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, price, id)
	return err
}

func (r *FruitRepository) UpdateStock(ctx context.Context, id int, stock float64) error {
	query := `UPDATE fruits SET stock_kg = $1 WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, stock, id)
	return err
}

func (r *FruitRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM fruits WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}