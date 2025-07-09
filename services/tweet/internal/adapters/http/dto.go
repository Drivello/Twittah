package http

// PostTweetRequestDTO is the DTO for creating a tweet via HTTP
type PostTweetRequestDTO struct {
	AuthorID  string   `json:"author_id" binding:"required"`
	Content   string   `json:"content" binding:"required"`
	Followers []string `json:"followers"`
}

// TweetResponseDTO is the DTO for returning a tweet via HTTP
type TweetResponseDTO struct {
	ID        string `json:"id"`
	AuthorID  string `json:"author_id"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
}
