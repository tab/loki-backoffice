package serializers

import (
	"time"

	"github.com/google/uuid"
)

type TokenSerializer struct {
	ID        uuid.UUID `json:"id"`
	UserId    uuid.UUID `json:"user_id"`
	Type      string    `json:"type"`
	Value     string    `json:"value"`
	ExpiresAt time.Time `json:"expires_at"`
}
