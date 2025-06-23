package main

import (
	"fmt"
	"net/http"
	"ranch/ranch"
)

func main() {
	app := ranch.NewApp()

	// Регистрация middleware в нужном порядке
	app.Use(ranch.DomainMiddleware)
	app.Use(ranch.BotDetectionMiddleware)
	app.Use(ranch.StatisticsMiddleware)
	app.Use(ranch.RedirectMiddleware)
	app.Use(ranch.ContentMiddleware)
	app.Use(ranch.RenderMiddleware)

	// Запуск HTTP сервера
	http.HandleFunc("/", app.Handle)
	port := ":8080"
	fmt.Printf("SEO Farm server started on %s\n", port)

	if err := http.ListenAndServe(port, nil); err != nil {
		panic(err)
	}
}
