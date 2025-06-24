package ranch

import (
	"context"
	"net/http"
	"time"
)

// Контекстные ключи
type ContextKey string

const (
	DomainKey    ContextKey = "domain"
	ProjectKey   ContextKey = "project"
	IsBotKey     ContextKey = "is_bot"
	UserAgentKey ContextKey = "user_agent"
	RequestIDKey ContextKey = "request_id"
)

// Основные структуры
type Domain struct {
	ID        int    `json:"id"`
	Domain    string `json:"domain"`
	Type      string `json:"type"`
	ProjectID int    `json:"project_id"`
}

type Project struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Language    string `json:"language"`
	EnableDebug bool   `json:"enable_debug"`
	TableName   string `json:"table_name"`
}

type ContentResult struct {
	Keyword  string                 `json:"keyword"`
	Snippets []string               `json:"snippets"`
	Content  string                 `json:"content"`
	Meta     map[string]interface{} `json:"meta"`
}

// HTTP Middleware тип
type HTTPMiddleware func(http.HandlerFunc) http.HandlerFunc

// Структуры для статистики и событий
type RequestEvent struct {
	RequestID    string    `json:"request_id"`
	Domain       string    `json:"domain"`
	ProjectID    int       `json:"project_id"`
	IsBot        bool      `json:"is_bot"`
	UserAgent    string    `json:"user_agent"`
	IPAddress    string    `json:"ip_address"`
	Timestamp    time.Time `json:"timestamp"`
	Path         string    `json:"path"`
	Method       string    `json:"method"`
	StatusCode   int       `json:"status_code"`
	ResponseTime int64     `json:"response_time_ms"`
}

type RedirectEvent struct {
	RequestID   string    `json:"request_id"`
	Domain      string    `json:"domain"`
	ProjectID   int       `json:"project_id"`
	RedirectURL string    `json:"redirect_url"`
	UserAgent   string    `json:"user_agent"`
	IPAddress   string    `json:"ip_address"`
	Timestamp   time.Time `json:"timestamp"`
}

type ContentEvent struct {
	RequestID string    `json:"request_id"`
	Domain    string    `json:"domain"`
	ProjectID int       `json:"project_id"`
	Keyword   string    `json:"keyword"`
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
}

// Интерфейсы репозиториев
type DomainRepository interface {
	FindByHostname(hostname string) (*Domain, error)
	GetAll() ([]*Domain, error)
	Create(domain *Domain) error
	Update(domain *Domain) error
	Delete(id int) error
}

type ProjectRepository interface {
	FindByID(id int) (*Project, error)
	GetAll() ([]*Project, error)
	Create(project *Project) error
	Update(project *Project) error
	Delete(id int) error
}

type ContentRepository interface {
	GetContentByProject(ctx context.Context, projectID int, keyword string) (*ContentResult, error)
	GetRandomContent(ctx context.Context, projectID int) (*ContentResult, error)
	SaveContent(ctx context.Context, content *ContentResult) error
}

type RedirectRepository interface {
	GetActiveRedirect(ctx context.Context, domainID int) (string, error)
	GetRandomRedirect(ctx context.Context, domainID int) (string, error)
	CreateRedirect(ctx context.Context, domainID int, url string) error
	DisableRedirect(ctx context.Context, domainID int, url string) error
}

// Интерфейсы для RabbitMQ
type MessagePublisher interface {
	PublishRequestEvent(ctx context.Context, event *RequestEvent) error
	PublishRedirectEvent(ctx context.Context, event *RedirectEvent) error
	PublishContentEvent(ctx context.Context, event *ContentEvent) error
	Close() error
}

type MessageConsumer interface {
	ConsumeRequestEvents(ctx context.Context, handler func(*RequestEvent) error) error
	ConsumeRedirectEvents(ctx context.Context, handler func(*RedirectEvent) error) error
	ConsumeContentEvents(ctx context.Context, handler func(*ContentEvent) error) error
	Close() error
}

// Интерфейсы сервисов
type ContentService interface {
	GetContent(ctx context.Context, domain *Domain, project *Project) (*ContentResult, error)
	GetRedirectURL(ctx context.Context, domain *Domain) (string, error)
}

type RenderService interface {
	RenderForBot(ctx context.Context, content *ContentResult) (string, error)
	RenderForHuman(ctx context.Context, content *ContentResult) (string, error)
	RenderError(message string) string
}

type StatisticsService interface {
	SaveRequestEvent(ctx context.Context, event *RequestEvent) error
	SaveRedirectEvent(ctx context.Context, event *RedirectEvent) error
	SaveContentEvent(ctx context.Context, event *ContentEvent) error
}

type CacheService interface {
	Get(ctx context.Context, key string) (interface{}, error)
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	Delete(ctx context.Context, key string) error
}
