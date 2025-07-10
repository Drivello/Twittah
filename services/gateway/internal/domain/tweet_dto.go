package domain

type TweetPublishEventDTO struct {
	ID        string `json:"id"`
	AuthorID  string `json:"author_id"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
}
