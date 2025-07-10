package http

// TweetRequestDTO representa el request para publicar un tweet via HTTP
// Usado solo por el handler de tweets
type TweetRequestDTO struct {
	AuthorID string `json:"author_id" binding:"required"`
	Content  string `json:"content" binding:"required"`
}

// TweetPublishResponseDTO representa la respuesta del endpoint de publicación de tweet
// Usado solo por el handler de tweets
type TweetPublishResponseDTO struct {
	Message string `json:"message"`
}
