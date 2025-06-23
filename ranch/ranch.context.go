package ranch

import (
	"net/http"
)

// Context - контекст запроса со всеми данными
type Context struct {
	Writer   http.ResponseWriter
	Request  *http.Request
	index    int
	handlers []HandlerFunc

	// SEO Farm специфичные поля
	Domain      *Domain
	Project     *Project
	IsBot       bool
	IsDebug     bool
	RedirectURL string
	Content     map[string]interface{}
	Response    string
	Aborted     bool
}

type HandlerFunc func(*Context)

// Next - переход к следующему middleware
func (c *Context) Next() {
	c.index++
	if c.index < len(c.handlers) && !c.Aborted {
		c.handlers[c.index](c)
	}
}

// Abort - прерывание цепочки middleware
func (c *Context) Abort() {
	c.Aborted = true
}

// AbortWithResponse - прерывание с ответом
func (c *Context) AbortWithResponse(response string) {
	c.Response = response
	c.Abort()
}

// Run - запуск цепочки middleware
func (c *Context) Run() {
	c.index = -1
	c.Next()
}

// GetHostname - получение hostname из запроса
func (c *Context) GetHostname() string {
	hostname := c.Request.Host
	if hostname == "" {
		hostname = "localhost:8080"
	}
	return hostname
}
