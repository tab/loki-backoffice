package models

import "github.com/google/uuid"

type User struct {
	ID             uuid.UUID
	IdentityNumber string
	PersonalCode   string
	FirstName      string
	LastName       string

	RoleIDs  []uuid.UUID
	ScopeIDs []uuid.UUID
}
