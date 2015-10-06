package route

type URLMapping struct {
	ShortURL string `json:"shorturl"`
	LongURL  string `json:"longurl"`
}

type APIResponse struct {
	StatusMsg string `json:"message"`
}
