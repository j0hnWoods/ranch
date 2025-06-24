package ranch

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// WithMiddleware - композиция middleware функций
func WithMiddleware(handler http.HandlerFunc, middlewares ...HTTPMiddleware) http.HandlerFunc {
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}
	return handler
}

// LoggingMiddleware - логирование запросов
func LoggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Генерируем request ID
		requestID := generateRequestID()
		ctx := context.WithValue(r.Context(), RequestIDKey, requestID)
		r = r.WithContext(ctx)

		fmt.Printf("[%s] %s %s %s - START\n",
			requestID, r.Method, r.URL.Path, r.Host)

		next(w, r)

		fmt.Printf("[%s] %s %s %s - END (%v)\n",
			requestID, r.Method, r.URL.Path, r.Host, time.Since(start))
	}
}

// DomainMiddleware - обработка домена
func DomainMiddleware(domainRepo DomainRepository, projectRepo ProjectRepository) HTTPMiddleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			hostname := r.Host
			if hostname == "" {
				hostname = "localhost:8080"
			}

			// Поиск домена через репозиторий
			domain, err := domainRepo.FindByHostname(hostname)
			if err != nil || domain == nil {
				http.Error(w, "Domain not found", http.StatusNotFound)
				return
			}

			// Поиск проекта через репозиторий
			project, err := projectRepo.FindByID(domain.ProjectID)
			if err != nil || project == nil {
				http.Error(w, "Project not found", http.StatusNotFound)
				return
			}

			// Добавляем IP адрес в контекст
			clientIP := getClientIP(r)

			// Добавляем все данные в контекст
			ctx := context.WithValue(r.Context(), DomainKey, domain)
			ctx = context.WithValue(ctx, ProjectKey, project)
			ctx = context.WithValue(ctx, "client_ip", clientIP)
			r = r.WithContext(ctx)

			next(w, r)
		}
	}
}

// BotDetectionMiddleware - определение ботов
func BotDetectionMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userAgent := r.UserAgent()
		isBot := detectBot(userAgent)

		ctx := context.WithValue(r.Context(), IsBotKey, isBot)
		ctx = context.WithValue(ctx, UserAgentKey, userAgent)
		r = r.WithContext(ctx)

		next(w, r)
	}
}

// StatisticsMiddleware - сбор статистики
func StatisticsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Создаем wrapper для ResponseWriter чтобы захватить status code
		wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		// Выполняем основной обработчик
		next(wrapped, r)

		// Асинхронно сохраняем статистику
		go func() {
			domain := r.Context().Value(DomainKey).(*Domain)
			project := r.Context().Value(ProjectKey).(*Project)
			isBot := r.Context().Value(IsBotKey).(bool)
			userAgent := r.Context().Value(UserAgentKey).(string)
			requestID := r.Context().Value(RequestIDKey).(string)

			// Создаем событие запроса
			event := &RequestEvent{
				RequestID:    requestID,
				Domain:       domain.Domain,
				ProjectID:    project.ID,
				IsBot:        isBot,
				UserAgent:    userAgent,
				IPAddress:    getClientIP(r),
				Timestamp:    time.Now(),
				Path:         r.URL.Path,
				Method:       r.Method,
				StatusCode:   wrapped.statusCode,
				ResponseTime: time.Since(start).Milliseconds(),
			}

			// Здесь можно инжектить StatisticsService
			saveRequestEvent(r.Context(), event)
		}()
	}
}

// responseWriter wrapper для захвата status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *responseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

// SecurityMiddleware - безопасность и блокировки
func SecurityMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userAgent := r.UserAgent()

		// Блокировка нежелательных ботов
		if isBlockedBot(userAgent) {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		// Проверка rate limiting (можно добавить)
		// if isRateLimited(r.RemoteAddr) {
		//     http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
		//     return
		// }

		next(w, r)
	}
}

// Вспомогательные функции
func generateRequestID() string {
	bytes := make([]byte, 4)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

func detectBot(userAgent string) bool {
	botKeywords := []string{"bot", "crawler", "spider", "googlebot", "bingbot"}
	userAgentLower := strings.ToLower(userAgent)

	for _, keyword := range botKeywords {
		if strings.Contains(userAgentLower, keyword) {
			return true
		}
	}
	return false
}

func isBlockedBot(userAgent string) bool {
	blockedBots := []string{"badbot", "spam", "hack"}
	userAgentLower := strings.ToLower(userAgent)

	for _, blocked := range blockedBots {
		if strings.Contains(userAgentLower, blocked) {
			return true
		}
	}
	return false
}

func saveRequestEvent(ctx context.Context, event *RequestEvent) {
	// TODO: реальная отправка в RabbitMQ через MessagePublisher
	fmt.Printf("[%s] Saving request event: domain=%s, bot=%v, status=%d, time=%dms\n",
		event.RequestID, event.Domain, event.IsBot, event.StatusCode, event.ResponseTime)
}

func getClientIP(r *http.Request) string {
	// Проверяем заголовки прокси
	if ip := r.Header.Get("X-Forwarded-For"); ip != "" {
		return strings.Split(ip, ",")[0]
	}
	if ip := r.Header.Get("X-Real-IP"); ip != "" {
		return ip
	}
	return r.RemoteAddr
}
