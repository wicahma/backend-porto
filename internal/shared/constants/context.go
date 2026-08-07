package constants

// contextKey is an unexported type used to prevent collisions in context keys.
type contextKey string

const (
	// UserIDKey is the context key for the authenticated user's ID.
	UserIDKey contextKey = "userID"

	// UserEmailKey is the context key for the authenticated user's email.
	UserEmailKey contextKey = "userEmail"
)
