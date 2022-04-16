package entities

type ShortenResultURL struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

type URLItem struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

type LinkItem struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
	IsDeleted   bool
}

type ShortenURL struct {
	Key       string
	IsDeleted bool
}

type DeleteCandidate struct {
	Key        string
	UserID     string
	RetryCount int8
}
