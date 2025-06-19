package main

import (
	"fmt"
	"net/http"
	"ranch/ranch"
)

func main() {
	controller := ranch.RanchController{}
	middleware := ranch.Middleware{}
	service := ranch.RanchService{}
	render := ranch.RanchRender{}

	// Инициализация зависимостей
	controller.Service = &service
	controller.Render = &render

	// Настройка HTTP роутинга
	http.HandleFunc("/", controller.HandleHTTP)

	// Дополнительные эндпоинты
	// http.HandleFunc("/robots.txt", controller.HandleRobots)
	// http.HandleFunc("/sitemap.xml", controller.HandleSitemap)
	// http.HandleFunc("/ping", controller.HandlePing)

	// Запуск HTTP сервера
	port := ":8080"

	if err := http.ListenAndServe(port, nil); err != nil {
		panic(err)
	}

	fmt.Printf("%v\n%v\n%v\n%v\n", controller, middleware, service, render)
}
