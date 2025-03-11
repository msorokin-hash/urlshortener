package entity

type InternalStorage struct {
	Alias string
	URL   string
}

type FileStorage struct {
	UUID        string `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}
