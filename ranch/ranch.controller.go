package ranch

import (
	"fmt"
	"net/http"
)

type RanchController struct {
	Service *RanchService
	Render  *RanchRender
}

// Domain структура
type Domain struct {
	ID        int
	Domain    string
	Type      string
	ProjectID int
}

// Project структура
type Project struct {
	ID          int
	Name        string
	Language    string
	EnableDebug bool
	TableName   string
}

// HandleHTTP - основной HTTP обработчик (обёртка для ActionIndex)
func (c *RanchController) HandleHTTP(w http.ResponseWriter, r *http.Request) {
	// Определяем домен из Host заголовка
	hostname := r.Host
	if hostname == "" {
		hostname = "localhost:8080" // Fallback для тестирования
	}

	// Логирование запроса
	fmt.Printf("HTTP Request: %s %s from %s (User-Agent: %s)\n",
		r.Method, r.URL.Path, hostname, r.UserAgent())

	// Обработка только GET запросов на главную
	if r.Method != "GET" || r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	// Вызов основной логики
	response := c.ActionIndex(hostname, r)

	// Отправка ответа
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(response))
}

// ActionIndex - главный метод контроллера
func (c *RanchController) ActionIndex(hostname string, r *http.Request) string {
	fmt.Printf("ActionIndex called for domain: %s\n", hostname)

	// 1. Поиск домена
	domain := c.findDomain(hostname)
	if domain == nil {
		return c.Render.RenderError("Domain not found")
	}

	// 2. Создание Middleware/Service
	middleware := c.Service.NewMiddleware(domain)

	// 3. Проверка типа пользователя
	isBot := c.isBot()
	isDebugEnabled := c.isDebugEnabled(domain)

	// 4. Логика редиректов для ботов
	if isBot {
		redirect := middleware.GetRedirect()
		if redirect != "" {
			return c.executeRedirect(redirect)
		}
	}

	// 5. Получение контента
	result := middleware.GetQueryResult()
	fmt.Printf("%v\n", result)

	// 6. Рендеринг в зависимости от типа пользователя
	if isBot || isDebugEnabled {
		return c.Render.RenderFarmTemplate(middleware, result)
	} else {
		return c.Render.RenderWhiteTemplate(result)
	}
}

// Вспомогательные методы контроллера

func (c *RanchController) findDomain(hostname string) *Domain {
	// TODO: Реальный поиск в БД
	fmt.Printf("Finding domain: %s\n", hostname)
	return &Domain{
		ID:        1,
		Domain:    hostname,
		Type:      "PROJECT",
		ProjectID: 1,
	}
}

func (c *RanchController) isBot() bool {
	// TODO: Реальная проверка User-Agent
	fmt.Println("Checking if request is from bot...")
	return true // Заглушка - считаем что это бот
}

func (c *RanchController) isDebugEnabled(domain *Domain) bool {
	// TODO: Проверка debug режима
	fmt.Println("Checking debug mode...")
	return false
}

func (c *RanchController) executeRedirect(redirectURL string) string {
	// TODO: Реальный HTTP редирект
	fmt.Printf("Executing redirect to: %s\n", redirectURL)
	return fmt.Sprintf("REDIRECT: %s", redirectURL)
}
