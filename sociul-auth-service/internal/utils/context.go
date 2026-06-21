package utils

import (
	"context"

	"github.com/google/uuid"
)

// Define unexported key type to avoid collisions
type contextKey string

// Define key
const userIdKey contextKey = "userId"

// Setter function
func NewContext(ctx context.Context, id uuid.UUID) context.Context {
	return context.WithValue(ctx, userIdKey, id)
}

// Getter function
func FromContext(ctx context.Context) (uuid.UUID, bool) {
	id, ok := ctx.Value(userIdKey).(uuid.UUID)
	return id, ok
}
