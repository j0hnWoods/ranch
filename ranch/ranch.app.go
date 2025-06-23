package ranch

import (
	"net/http"
)

// App - основное приложение
type App struct {
	middlewares []HandlerFunc
	service     *RanchService
	render      *RanchRender
}

// NewApp - создание нового приложения
func NewApp() *App {
	return &App{
		service: &RanchService{},
		render:  &RanchRender{},
	}
}

// Use - добавление middleware
func (a *App) Use(mw HandlerFunc) {
	a.middlewares = append(a.middlewares, mw)
}

// Handle - обработка HTTP запроса
func (a *App) Handle(w http.ResponseWriter, r *http.Request) {
	// Проверка метода и пути
	if r.Method != "GET" || r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	ctx := &Context{
		Writer:   w,
		Request:  r,
		handlers: a.middlewares,
		Content:  make(map[string]interface{}),
	}

	// Инициализация сервисов в контексте
	ctx.Run()

	// Отправка ответа
	if ctx.Response != "" {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(ctx.Response))
	}
}
