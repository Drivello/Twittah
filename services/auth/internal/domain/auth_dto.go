package domain

// UserCreatedEventDTO is the DTO for user created events in Kafka
type UserCreatedEventDTO struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}
