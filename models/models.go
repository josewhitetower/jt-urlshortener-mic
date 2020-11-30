package models

// URL Schema for the urls table
type URL struct {
	OriginalURL string `json:"original_url"`
	ShortURL    int64  `json:"short_url"`
}
