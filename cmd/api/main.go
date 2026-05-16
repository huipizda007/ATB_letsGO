package main

import (
	"ATB_letsGO/internal/config"
	"ATB_letsGO/internal/handler"
	"ATB_letsGO/internal/repository"
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/lib/pq" 
)

func main() {
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("Неможливо завантажити конфіг:", err)
	}

	db, err := sql.Open("postgres", cfg.DBSource)
	if err != nil {
		log.Fatal("Неможливо підключитися до бази даних:", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("База даних недоступна:", err)
	}

	fmt.Println("Успішно підключено до PostgreSQL!")

	query := `
	CREATE TABLE IF NOT EXISTS fruits (
		id SERIAL PRIMARY KEY,
		name VARCHAR(100) NOT NULL,
		brand VARCHAR(100) NOT NULL,
		price_per_kg NUMERIC(10, 2) NOT NULL,
		stock_kg NUMERIC(10, 2) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`
	
	if _, err := db.Exec(query); err != nil {
		log.Fatal("Помилка створення таблиці:", err)
	}
	fmt.Println("Таблиця fruits готова до роботи!")

	fruitRepo := repository.NewFruitRepository(db)
	fruitHandler := handler.NewFruitHandler(fruitRepo)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/fruits", func(r chi.Router) {
		r.Post("/", fruitHandler.CreateFruit)             
		r.Get("/", fruitHandler.GetAllFruits)               
		r.Get("/{id}", fruitHandler.GetFruitByID)            
		r.Put("/{id}", fruitHandler.UpdateFruit)              
		r.Patch("/{id}/price", fruitHandler.UpdateFruitPrice) 
		r.Patch("/{id}/stock", fruitHandler.UpdateFruitStock) 
		r.Delete("/{id}", fruitHandler.DeleteFruit)   
	})

	fmt.Printf("АТБ Сервер успішно запущено на порту %s 🚀\n", cfg.ServerAddress)
	if err := http.ListenAndServe(cfg.ServerAddress, r); err != nil {
		log.Fatal("Помилка запуску сервера:", err)
	}
}