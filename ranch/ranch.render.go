// ranch/render.go
package ranch

import (
	"fmt"
)

type RanchRender struct {
}

// RenderFarmTemplate - рендеринг
func (r *RanchRender) RenderFarmTemplate(middleware *Middleware, result map[string]interface{}) string {
	// TODO: Реальный рендеринг HTML шаблона
	fmt.Println("Rendering farm template for bot...")
	return fmt.Sprintf("SEO_CONTENT: domain=%s, keyword=%s, snippets=%v",
		middleware.Domain.Domain, result["keyword"], result["snippets"])
}

// RenderWhiteTemplate - рендеринг
func (r *RanchRender) RenderWhiteTemplate(result map[string]interface{}) string {
	// TODO: Реальный рендеринг "белого" HTML шаблона
	fmt.Println("Rendering white template for human...")
	return fmt.Sprintf("WHITE_PAGE: content=%s", result["content"])
}

// RenderError - рендеринг страницы ошибки
func (r *RanchRender) RenderError(message string) string {
	fmt.Printf("Rendering error page: %s\n", message)
	return fmt.Sprintf("ERROR_PAGE: %s", message)
}

// TODO: Дополнительные методы рендеринга
// RenderDoorTemplate - для door страниц
// RenderRedirectPage - для редиректов
// RenderRobotsTxt - для robots.txt
// RenderSitemap - для sitemap.xml
