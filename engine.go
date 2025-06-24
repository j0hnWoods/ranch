package main

import (
	"fmt"
	"net/http"
	"ranch/ranch"
)

func main() {
	// Инициализация компонентов через Factory
	factory := ranch.NewComponentFactory()

	// Graceful shutdown
	defer func() {
		if err := factory.Close(); err != nil {
			fmt.Printf("Error closing factory: %v\n", err)
		}
	}()

	// Создание основного контроллера
	controller := factory.CreateController()

	// Настройка middleware цепочки
	handler := ranch.WithMiddleware(
		controller.HandleHTTP,
		ranch.LoggingMiddleware,
		ranch.DomainMiddleware(factory.GetDomainRepository(), factory.GetProjectRepository()),
		ranch.BotDetectionMiddleware,
		ranch.StatisticsMiddleware,
		ranch.SecurityMiddleware,
	)

	// Роутинг
	http.HandleFunc("/", handler)
	http.HandleFunc("/robots.txt", ranch.WithMiddleware(controller.HandleRobots, ranch.LoggingMiddleware))
	http.HandleFunc("/sitemap.xml", ranch.WithMiddleware(controller.HandleSitemap, ranch.LoggingMiddleware))
	http.HandleFunc("/health", ranch.WithMiddleware(controller.HandleHealth, ranch.LoggingMiddleware))

	fmt.Println("SEO Farm server starting on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
