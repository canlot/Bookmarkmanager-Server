package Models

type Bookmark struct {
	ID         uint   `json:"id"`
	CategoryID uint   `json:"category_id"`
	Url        string `json:"url"`
}
