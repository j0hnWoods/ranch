package ranch

// Domain структура
type Domain struct {
	ID        int
	Domain    string
	Type      string
	ProjectID int
}

// Project структура
type Project struct {
	ID          int
	Name        string
	Language    string
	EnableDebug bool
	TableName   string
}

// Под БД и ее использование - отдельная структура + функции
// Под Брокер и его использование - отдельная структура + функции
// ...
