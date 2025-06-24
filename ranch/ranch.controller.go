package ranch

import (
	"context"
	"net/http"
	"time"
)

type RanchController struct {
	contentService    ContentService
	renderService     RenderService
	statisticsService StatisticsService
}

func NewRanchController(contentService ContentService, renderService RenderService, statisticsService StatisticsService) *RanchController {
	return &RanchController{
		contentService:    contentService,
		renderService:     renderService,
		statisticsService: statisticsService,
	}
}

// HandleHTTP - главный обработчик (теперь чистый от middleware логики)
func (c *RanchController) HandleHTTP(w http.ResponseWriter, r *http.Request) {
	// Проверяем только основной роут
	if r.Method != "GET" || r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	// Получаем данные из контекста (установленные middleware)
	domain := r.Context().Value(DomainKey).(*Domain)
	project := r.Context().Value(ProjectKey).(*Project)
	isBot := r.Context().Value(IsBotKey).(bool)

	// Основная бизнес-логика
	response, err := c.processRequest(r.Context(), domain, project, isBot)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Отправляем ответ
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(response))
}

// processRequest - основная бизнес-логика
func (c *RanchController) processRequest(ctx context.Context, domain *Domain, project *Project, isBot bool) (string, error) {
	requestID := ctx.Value(RequestIDKey).(string)

	// Для ботов проверяем редиректы
	if isBot {
		redirectURL, err := c.contentService.GetRedirectURL(ctx, domain)
		if err == nil && redirectURL != "" {
			// Сохраняем событие редиректа
			redirectEvent := &RedirectEvent{
				RequestID:   requestID,
				Domain:      domain.Domain,
				ProjectID:   domain.ProjectID,
				RedirectURL: redirectURL,
				UserAgent:   ctx.Value(UserAgentKey).(string),
				IPAddress:   getClientIPFromContext(ctx),
				Timestamp:   time.Now(),
			}

			go c.statisticsService.SaveRedirectEvent(ctx, redirectEvent)

			// TODO: настоящий HTTP редирект
			return c.renderService.RenderError("REDIRECT: " + redirectURL), nil
		}
	}

	// Получаем контент
	content, err := c.contentService.GetContent(ctx, domain, project)
	if err != nil {
		return c.renderService.RenderError("Content error: " + err.Error()), nil
	}

	// Сохраняем событие контента
	contentEvent := &ContentEvent{
		RequestID: requestID,
		Domain:    domain.Domain,
		ProjectID: domain.ProjectID,
		Keyword:   content.Keyword,
		Content:   content.Content,
		Timestamp: time.Now(),
	}

	go c.statisticsService.SaveContentEvent(ctx, contentEvent)

	// Рендерим в зависимости от типа пользователя
	if isBot || project.EnableDebug {
		return c.renderService.RenderForBot(ctx, content)
	}

	return c.renderService.RenderForHuman(ctx, content)
}

func getClientIPFromContext(ctx context.Context) string {
	if ip := ctx.Value("client_ip"); ip != nil {
		return ip.(string)
	}
	return "127.0.0.1" // fallback
}

// HandleHealth - health check endpoint
func (c *RanchController) HandleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok","service":"seo-farm"}`))
}

// HandleRobots - обработка robots.txt
func (c *RanchController) HandleRobots(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("User-agent: *\nDisallow:"))
}

// HandleSitemap - обработка sitemap.xml
func (c *RanchController) HandleSitemap(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/xml")
	w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
</urlset>`))
}
