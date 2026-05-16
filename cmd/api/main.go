package main

import (
	"ATB_letsGO/internal/config"
	"fmt"
	"log"
)

func main() {
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("Неможливо завантажити конфіг:", err)
	}

	fmt.Println("Базова структура готова!")
	fmt.Println("Сервер буде на порту:", cfg.ServerAddress)
	fmt.Println("Зв'язок з базою налаштовано.")
}