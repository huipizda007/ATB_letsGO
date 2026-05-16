package handler

import (
	"ATB_letsGO/internal/models"
	"ATB_letsGO/internal/repository"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type FruitHandler struct {
	repo *repository.FruitRepository
}

func NewFruitHandler(repo *repository.FruitRepository) *FruitHandler {
	return &FruitHandler{repo: repo}
}

func (h *FruitHandler) CreateFruit(w http.ResponseWriter, r *http.Request) {
	var fruit models.Fruit
	if err := json.NewDecoder(r.Body).Decode(&fruit); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Некоректний формат JSON"})
		return
	}

	if fruit.Name == "" || fruit.PricePerKg <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Назва не може бути порожньою, а ціна має бути більшою за нуль"})
		return
	}

	id, err := h.repo.Create(r.Context(), fruit)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Помилка збереження в базу"})
		return
	}

	fruit.ID = id
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(fruit)
}

func (h *FruitHandler) GetAllFruits(w http.ResponseWriter, r *http.Request) {
	fruits, err := h.repo.GetAll(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Помилка отримання даних"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(fruits)
}

func (h *FruitHandler) GetFruitByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Некоректний ID"})
		return
	}

	fruit, err := h.repo.GetByID(r.Context(), id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Фрукт не знайдено"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(fruit)
}

func (h *FruitHandler) UpdateFruit(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Некоректний ID фрукта"})
		return
	}

	var fruit models.Fruit
	if err := json.NewDecoder(r.Body).Decode(&fruit); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Некоректний формат JSON"})
		return
	}

	if fruit.Name == "" || fruit.Brand == "" || fruit.PricePerKg <= 0 || fruit.StockKg < 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Некоректні дані: заповніть усі поля правильно"})
		return
	}

	if err := h.repo.Update(r.Context(), id, fruit); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Помилка оновлення в базі даних"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Дані фрукта успішно повністю оновлено"})
}

func (h *FruitHandler) UpdateFruitPrice(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Некоректний ID"})
		return
	}

	var req struct {
		PricePerKg float64 `json:"price_per_kg"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.PricePerKg <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Некоректна ціна"})
		return
	}

	if err := h.repo.UpdatePrice(r.Context(), id, req.PricePerKg); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Помилка оновлення ціни"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Ціну успішно оновлено"})
}

func (h *FruitHandler) UpdateFruitStock(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Некоректний ID"})
		return
	}

	var req struct {
		StockKg float64 `json:"stock_kg"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.StockKg < 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Некоректний залишок"})
		return
	}

	if err := h.repo.UpdateStock(r.Context(), id, req.StockKg); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Помилка оновлення залишку"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Залишок успішно оновлено"})
}

func (h *FruitHandler) DeleteFruit(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Некоректний ID"})
		return
	}

	if err := h.repo.Delete(r.Context(), id); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Помилка видалення"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Фрукт успішно видалено"})
}