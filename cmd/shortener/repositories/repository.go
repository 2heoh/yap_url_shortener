package repositories

type Repository interface {
	Get(id string) (string, error)
	Add(id string, url string, userID string) error

	GetAllFor(userID string) []LinkItem
	Ping() error
}

type LinkItem struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}
