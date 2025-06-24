package ranch

// ComponentFactory - фабрика для создания компонентов
type ComponentFactory struct {
	domainRepo    DomainRepository
	projectRepo   ProjectRepository
	contentRepo   ContentRepository
	redirectRepo  RedirectRepository
	contentSvc    ContentService
	renderSvc     RenderService
	statisticsSvc StatisticsService
	cacheService  CacheService
	msgPublisher  MessagePublisher
	msgConsumer   MessageConsumer
}

func NewComponentFactory() *ComponentFactory {
	// TODO: инициализация реальных зависимостей из конфигурации
	factory := &ComponentFactory{}

	// Инициализация репозиториев
	factory.domainRepo = NewDomainRepository()
	factory.projectRepo = NewProjectRepository()
	factory.contentRepo = NewContentRepository()
	factory.redirectRepo = NewRedirectRepository()

	// Инициализация RabbitMQ
	factory.msgPublisher = NewRabbitMQPublisher()
	factory.msgConsumer = NewRabbitMQConsumer()

	// Инициализация сервисов
	factory.cacheService = NewCacheService()
	factory.statisticsSvc = NewStatisticsService(factory.msgPublisher)
	factory.contentSvc = NewContentService(factory.contentRepo, factory.redirectRepo, factory.cacheService)
	factory.renderSvc = NewRenderService()

	return factory
}

func (f *ComponentFactory) CreateController() *RanchController {
	return NewRanchController(f.contentSvc, f.renderSvc, f.statisticsSvc)
}

func (f *ComponentFactory) GetDomainRepository() DomainRepository {
	return f.domainRepo
}

func (f *ComponentFactory) GetProjectRepository() ProjectRepository {
	return f.projectRepo
}

func (f *ComponentFactory) Close() error {
	// Закрываем все соединения
	if f.msgPublisher != nil {
		f.msgPublisher.Close()
	}
	if f.msgConsumer != nil {
		f.msgConsumer.Close()
	}
	return nil
}
