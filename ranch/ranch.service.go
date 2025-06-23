package ranch

import (
	"fmt"
	"strings"
)

// RanchService - сервисный слой (пока заглушка)
type RanchService struct {
}

// Вспомогательные функции

func shouldBlockBot(userAgent string) bool {
	// Список заблокированных ботов
	blockedBots := []string{"badbot", "spambot"}
	ua := strings.ToLower(userAgent)

	for _, bot := range blockedBots {
		if strings.Contains(ua, bot) {
			return true
		}
	}
	return false
}

func getRedirectURL(c *Context) string {
	// Логика получения редиректа (заглушка)
	// В реальности здесь будет запрос к БД
	if shouldRedirect() {
		return "https://target-site.com/seo-page"
	}
	return ""
}

func shouldRedirect() bool {
	// 30% вероятность редиректа (заглушка)
	return true // Для демонстрации
}

func saveStatistics(c *Context) {
	fmt.Printf("Saving stats: domain=%s, bot=%v\n",
		c.Domain.Domain, c.IsBot)
	// TODO: Отправка в RabbitMQ/БД
}

func setupDatabase(tableName string) {
	fmt.Printf("Setting up database table: %s\n", tableName)
	// TODO: Настройка подключения к БД
}

// Функции поиска в БД (заглушки)
func findDomain(hostname string) *Domain {
	fmt.Printf("Searching domain: %s\n", hostname)

	// TODO: Реальный поиск в БД
	// Заглушка - возвращаем домен только для известных хостов
	allowedDomains := []string{"localhost:8080", "test.com", "example.com"}

	for _, allowed := range allowedDomains {
		if hostname == allowed {
			return &Domain{
				ID:        1,
				Domain:    hostname,
				Type:      "PROJECT",
				ProjectID: 1,
			}
		}
	}

	return nil // Домен не найден
}

func findProject(projectID int) *Project {
	fmt.Printf("Searching project: %d\n", projectID)

	// TODO: Реальный поиск в БД
	// Заглушка
	if projectID == 1 {
		return &Project{
			ID:          projectID,
			Name:        "SEO Farm Project",
			Language:    "en",
			EnableDebug: false,
			TableName:   "seo_content_table",
		}
	}

	return nil // Проект не найден
}
