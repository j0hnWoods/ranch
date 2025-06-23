package ranch

import (
	"fmt"
	"strings"
)

// DomainMiddleware - поиск и проверка домена
func DomainMiddleware(c *Context) {
	hostname := c.GetHostname()
	fmt.Printf("DomainMiddleware: Processing %s\n", hostname)

	// Поиск домена в БД
	domain := findDomain(hostname)
	if domain == nil {
		c.AbortWithResponse(renderError("Domain not found"))
		return
	}

	// Поиск проекта
	project := findProject(domain.ProjectID)
	if project == nil {
		c.AbortWithResponse(renderError("Project not found"))
		return
	}

	c.Domain = domain
	c.Project = project
	c.Next()
}

// BotDetectionMiddleware - определение типа пользователя
func BotDetectionMiddleware(c *Context) {
	fmt.Println("BotDetectionMiddleware: Detecting user type")

	userAgent := c.Request.UserAgent()

	// Простая проверка на бота
	botKeywords := []string{"bot", "crawler", "spider", "scraper", "googlebot", "bingbot"}
	isBot := false

	for _, keyword := range botKeywords {
		if strings.Contains(strings.ToLower(userAgent), keyword) {
			isBot = true
			break
		}
	}

	c.IsBot = isBot
	c.IsDebug = c.Project.EnableDebug

	fmt.Printf("BotDetectionMiddleware: IsBot=%v, IsDebug=%v\n", c.IsBot, c.IsDebug)
	c.Next()
}

// StatisticsMiddleware - сбор и сохранение статистики
func StatisticsMiddleware(c *Context) {
	fmt.Println("StatisticsMiddleware: Collecting statistics")

	// Логирование запроса
	fmt.Printf("Request: %s %s from %s (UA: %s)\n",
		c.Request.Method, c.Request.URL.Path,
		c.GetHostname(), c.Request.UserAgent())

	// Проверка и блокировка плохих ботов
	if c.IsBot && shouldBlockBot(c.Request.UserAgent()) {
		c.AbortWithResponse(renderError("Access denied"))
		return
	}

	// Сохранение статистики (заглушка)
	saveStatistics(c)

	c.Next()
}

// RedirectMiddleware - обработка редиректов для ботов
func RedirectMiddleware(c *Context) {
	fmt.Println("RedirectMiddleware: Processing redirects")

	if c.IsBot {
		redirectURL := getRedirectURL(c)
		if redirectURL != "" {
			c.RedirectURL = redirectURL
			response := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta http-equiv="refresh" content="0; url=%s">
    <title>Redirecting...</title>
</head>
<body>
    <p>Redirecting to <a href="%s">%s</a></p>
</body>
</html>`, redirectURL, redirectURL, redirectURL)

			c.AbortWithResponse(response)
			return
		}
	}

	c.Next()
}

// ContentMiddleware - получение контента из БД
func ContentMiddleware(c *Context) {
	fmt.Println("ContentMiddleware: Loading content")

	// Настройка БД таблицы
	setupDatabase(c.Project.TableName)

	// Получение контента (заглушка)
	content := map[string]interface{}{
		"keyword":     "sample SEO keyword",
		"title":       "SEO Optimized Title",
		"description": "SEO meta description",
		"snippets":    []string{"snippet1", "snippet2", "snippet3"},
		"content":     generateSEOContent(),
		"domain":      c.Domain.Domain,
	}

	c.Content = content
	c.Next()
}

// RenderMiddleware - рендеринг финального ответа
func RenderMiddleware(c *Context) {
	fmt.Println("RenderMiddleware: Rendering response")

	var response string

	if c.IsBot || c.IsDebug {
		// SEO контент для ботов
		response = renderFarmTemplate(c.Content)
	} else {
		// "Белая" страница для людей
		response = renderWhiteTemplate(c.Content)
	}

	c.AbortWithResponse(response)
}
