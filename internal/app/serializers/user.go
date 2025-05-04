package serializers

import "github.com/google/uuid"

type UserSerializer struct {
	ID             uuid.UUID `json:"id"`
	IdentityNumber string    `json:"identity_number"`
	PersonalCode   string    `json:"personal_code"`
	FirstName      string    `json:"first_name"`
	LastName       string    `json:"last_name"`

	RoleIDs  []uuid.UUID `json:"role_ids,omitempty"`
	ScopeIDs []uuid.UUID `json:"scope_ids,omitempty"`
}
