package http

import "net/http"

// GenericResponse is a simple structure for API responses
// that should not leak internal details.
type GenericResponse struct {
	Message string `json:"message"`
}

// 4xx responses
var (
	ErrBadRequest   = GenericResponse{Message: "La solicitud es inválida o incompleta."}
	ErrUnauthorized = GenericResponse{Message: "No autorizado para realizar esta acción."}
	ErrForbidden    = GenericResponse{Message: "Acción no permitida."}
	ErrNotFound     = GenericResponse{Message: "El recurso solicitado no existe."}
	ErrConflict     = GenericResponse{Message: "Conflicto con el estado actual del recurso."}
)

// 5xx responses
var (
	ErrInternal           = GenericResponse{Message: "Ha ocurrido un error inesperado. Intenta nuevamente más tarde."}
	ErrServiceUnavailable = GenericResponse{Message: "El servicio no está disponible temporalmente."}
)

// Helper to write a generic error response
func WriteGenericError(c interface{ JSON(int, interface{}) }, status int) {
	switch status {
	case http.StatusBadRequest:
		c.JSON(status, ErrBadRequest)
	case http.StatusUnauthorized:
		c.JSON(status, ErrUnauthorized)
	case http.StatusForbidden:
		c.JSON(status, ErrForbidden)
	case http.StatusNotFound:
		c.JSON(status, ErrNotFound)
	case http.StatusConflict:
		c.JSON(status, ErrConflict)
	case http.StatusServiceUnavailable:
		c.JSON(status, ErrServiceUnavailable)
	default:
		c.JSON(status, ErrInternal)
	}
}
