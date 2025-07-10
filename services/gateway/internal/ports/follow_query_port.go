package ports

// FollowQueryUseCasePort define la interface para la lógica de consulta de followers
// desde cualquier fuente (HTTP, DB, etc)
type FollowQueryUseCasePort interface {
	GetFollowers(userID string) ([]string, error)
}
