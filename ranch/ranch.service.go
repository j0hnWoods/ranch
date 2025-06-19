// ranch/service.go
package ranch

import (
	"fmt"
)

type RanchService struct {
}

// NewMiddleware - создание нового middleware
func (s *RanchService) NewMiddleware(domain *Domain) *Middleware {
	fmt.Printf("Creating middleware for domain: %s\n", domain.Domain)

	middleware := &Middleware{
		Domain:  domain,
		Service: s,
	}

	// Загрузка проекта
	middleware.Project = s.findProject(domain.ProjectID)

	// Выполнение handleProject
	middleware.handleProject()

	return middleware
}

// Вспомогательные методы RanchService

func (s *RanchService) findProject(projectID int) *Project {
	// TODO: Реальный поиск проекта в БД
	fmt.Printf("Finding project: %d\n", projectID)
	return &Project{
		ID:          projectID,
		Name:        "SEO Farm Project",
		Language:    "en",
		EnableDebug: false,
		TableName:   "seo_content_table",
	}
}
