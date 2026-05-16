package handler

import (
	"ATB_letsGO/internal/models"
	"ATB_letsGO/internal/repository"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
)

type FruitHandler struct {
	repo *repository.FruitRepository
}

func NewFruitHandler(repo *repository.FruitRepository) *FruitHandler {
	return &FruitHandler{repo: repo}
}

func (h *FruitHandler) CreateFruit(w http.ResponseWriter, r *http.Request) {
	var f models.Fruit
	if err := json.NewDecoder(r.Body).Decode(&f); err != nil {
		http.Error(w, "Некоректний JSON", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(f.Name) == "" {
		http.Error(w, "Назва фрукта (name) не може бути порожньою", http.StatusBadRequest)
		return
	}
	if strings.TrimSpace(f.Brand) == "" {
		http.Error(w, "Бренд (brand) не може бути порожнім", http.StatusBadRequest)
		return
	}
	if f.PricePerKg <= 0 {
		http.Error(w, "Ціна (price_per_kg) має бути більшою за нуль", http.StatusBadRequest)
		return
	}
	if f.StockKg < 0 {
		http.Error(w, "Кількість на складі (stock_kg) не може бути від'ємною", http.StatusBadRequest)
		return
	}

	id, err := h.repo.Create(r.Context(), f)
	if err != nil {
		http.Error(w, "Помилка збереження в базу: "+err.Error(), http.StatusInternalServerError)
		return
	}

	f.ID = id
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(f)
}

func (h *FruitHandler) GetAllFruits(w http.ResponseWriter, r *http.Request) {
	fruits, err := h.repo.GetAll(r.Context())
	if err != nil {
		http.Error(w, "Помилка отримання даних: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(fruits)
}

func (h *FruitHandler) GetFruitByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Некоректний ID", http.StatusBadRequest)
		return
	}

	f, err := h.repo.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, "Фрукт не знайдено", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(f)
}

func (h *FruitHandler) UpdateFruitPrice(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Некоректний ID", http.StatusBadRequest)
		return
	}

	var req struct {
		PricePerKg float64 `json:"price_per_kg"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Некоректний JSON", http.StatusBadRequest)
		return
	}

	if req.PricePerKg <= 0 {
		http.Error(w, "Нова ціна має бути більшою за нуль", http.StatusBadRequest)
		return
	}

	if err := h.repo.UpdatePrice(r.Context(), id, req.PricePerKg); err != nil {
		http.Error(w, "Помилка оновлення бази", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Ціну успішно оновлено"}`))
}

func (h *FruitHandler) UpdateFruitStock(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Некоректний ID", http.StatusBadRequest)
		return
	}

	var req struct {
		StockKg float64 `json:"stock_kg"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Некоректний JSON", http.StatusBadRequest)
		return
	}

	if req.StockKg < 0 {
		http.Error(w, "Кількість на складі не може бути від'ємною", http.StatusBadRequest)
		return
	}

	if err := h.repo.UpdateStock(r.Context(), id, req.StockKg); err != nil {
		http.Error(w, "Помилка оновлення бази", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Залишок на складі успішно оновлено"}`))
}

func (h *FruitHandler) DeleteFruit(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Некоректний ID", http.StatusBadRequest)
		return
	}

	if err := h.repo.Delete(r.Context(), id); err != nil {
		http.Error(w, "Помилка видалення", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Фрукт видалено"}`))
}