package news

type News struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"descriptions"`
	Link        string `json:"link"`
}

type Repository interface {
}
