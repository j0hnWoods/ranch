package ranch

import (
	"context"
	"fmt"
	"time"
)

// ===== Statistics Service =====

type statisticsService struct {
	publisher MessagePublisher
}

func NewStatisticsService(publisher MessagePublisher) StatisticsService {
	return &statisticsService{
		publisher: publisher,
	}
}

func (s *statisticsService) SaveRequestEvent(ctx context.Context, event *RequestEvent) error {
	fmt.Printf("Stats: Saving request event - %s\n", event.RequestID)
	return s.publisher.PublishRequestEvent(ctx, event)
}

func (s *statisticsService) SaveRedirectEvent(ctx context.Context, event *RedirectEvent) error {
	fmt.Printf("Stats: Saving redirect event - %s\n", event.RequestID)
	return s.publisher.PublishRedirectEvent(ctx, event)
}

func (s *statisticsService) SaveContentEvent(ctx context.Context, event *ContentEvent) error {
	fmt.Printf("Stats: Saving content event - %s\n", event.RequestID)
	return s.publisher.PublishContentEvent(ctx, event)
}

// ===== Content Service =====

type contentService struct {
	contentRepo  ContentRepository
	redirectRepo RedirectRepository
	cache        CacheService
}

func NewContentService(contentRepo ContentRepository, redirectRepo RedirectRepository, cache CacheService) ContentService {
	return &contentService{
		contentRepo:  contentRepo,
		redirectRepo: redirectRepo,
		cache:        cache,
	}
}

func (s *contentService) GetContent(ctx context.Context, domain *Domain, project *Project) (*ContentResult, error) {
	requestID := ctx.Value(RequestIDKey).(string)
	fmt.Printf("[%s] Getting content for domain: %s, project: %s\n",
		requestID, domain.Domain, project.Name)

	// Пытаемся получить из кэша
	cacheKey := fmt.Sprintf("content:%s:%d", domain.Domain, project.ID)
	if cached, err := s.cache.Get(ctx, cacheKey); err == nil && cached != nil {
		fmt.Printf("[%s] Content found in cache\n", requestID)
		return cached.(*ContentResult), nil
	}

	// Получаем из БД
	content, err := s.contentRepo.GetRandomContent(ctx, project.ID)
	if err != nil {
		fmt.Printf("[%s] Error getting content: %v\n", requestID, err)
		return nil, err
	}

	// Сохраняем в кэш на 5 минут
	s.cache.Set(ctx, cacheKey, content, 5*time.Minute)

	// Обогащаем контент данными домена
	content.Meta["domain"] = domain.Domain
	content.Meta["project"] = project.Name
	content.Meta["title"] = fmt.Sprintf("SEO Page - %s", domain.Domain)
	content.Meta["description"] = fmt.Sprintf("SEO optimized content for %s", domain.Domain)

	return content, nil
}

func (s *contentService) GetRedirectURL(ctx context.Context, domain *Domain) (string, error) {
	requestID := ctx.Value(RequestIDKey).(string)
	fmt.Printf("[%s] Checking redirect for domain: %s\n", requestID, domain.Domain)

	// Проверяем кэш редиректов
	cacheKey := fmt.Sprintf("redirect:%d", domain.ID)
	if cached, err := s.cache.Get(ctx, cacheKey); err == nil && cached != nil {
		fmt.Printf("[%s] Redirect found in cache\n", requestID)
		return cached.(string), nil
	}

	// Получаем активный редирект из БД
	redirectURL, err := s.redirectRepo.GetActiveRedirect(ctx, domain.ID)
	if err != nil {
		fmt.Printf("[%s] Error getting redirect: %v\n", requestID, err)
		return "", err
	}

	// Если нет активного редиректа, получаем случайный
	if redirectURL == "" {
		redirectURL, err = s.redirectRepo.GetRandomRedirect(ctx, domain.ID)
		if err != nil {
			fmt.Printf("[%s] Error getting random redirect: %v\n", requestID, err)
			return "", err
		}
	}

	// Сохраняем в кэш на 1 минуту
	if redirectURL != "" {
		s.cache.Set(ctx, cacheKey, redirectURL, 1*time.Minute)
	}

	return redirectURL, nil
}

// ===== Render Service =====

type renderService struct{}

func NewRenderService() RenderService {
	return &renderService{}
}

func (r *renderService) RenderForBot(ctx context.Context, content *ContentResult) (string, error) {
	requestID := ctx.Value(RequestIDKey).(string)
	fmt.Printf("[%s] Rendering for bot\n", requestID)

	// Расширенный HTML для ботов с SEO оптимизацией
	html := fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>%s</title>
    <meta name="description" content="%s">
    <meta name="keywords" content="%s">
    <meta name="robots" content="index, follow">
    <link rel="canonical" href="https://%s/">
</head>
<body>
    <header>
        <h1>%s</h1>
        <nav>
            <a href="/">Home</a>
            <a href="/about">About</a>
            <a href="/contact">Contact</a>
        </nav>
    </header>
    <main>
        <article>
            <h2>%s</h2>
            <div class="content">
                %s
            </div>
            <div class="snippets">
                <h3>Related Information:</h3>
                <ul>%s</ul>
            </div>
        </article>
    </main>
    <footer>
        <p>&copy; 2025 %s. All rights reserved.</p>
    </footer>
</body>
</html>`,
		content.Meta["title"], content.Meta["description"], content.Keyword,
		content.Meta["domain"], content.Meta["title"],
		content.Keyword, content.Content,
		formatSnippets(content.Snippets), content.Meta["domain"])

	return html, nil
}

func (r *renderService) RenderForHuman(ctx context.Context, content *ContentResult) (string, error) {
	requestID := ctx.Value(RequestIDKey).(string)
	fmt.Printf("[%s] Rendering for human\n", requestID)

	// Простая страница для людей
	html := fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Welcome - %s</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 40px; }
        .container { max-width: 800px; margin: 0 auto; }
        .header { text-align: center; margin-bottom: 40px; }
        .content { line-height: 1.6; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Welcome to %s</h1>
            <p>Thank you for visiting our website!</p>
        </div>
        <div class="content">
            <p>%s</p>
            <p>We hope you find what you're looking for. Feel free to explore our site.</p>
        </div>
    </div>
</body>
</html>`, content.Meta["domain"], content.Meta["domain"], content.Content)

	return html, nil
}

func (r *renderService) RenderError(message string) string {
	html := fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Error</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 40px; text-align: center; }
        .error { color: #d32f2f; }
    </style>
</head>
<body>
    <h1 class="error">Error</h1>
    <p>%s</p>
    <a href="/">Return to Home</a>
</body>
</html>`, message)

	return html
}

func formatSnippets(snippets []string) string {
	result := ""
	for _, snippet := range snippets {
		result += fmt.Sprintf("<li>%s</li>", snippet)
	}
	return result
}
