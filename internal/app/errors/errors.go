package errors

import "errors"

var (
	// ErrInvalidToken indicates that the provided token is invalid
	ErrInvalidToken = errors.New("invalid token")

	// ErrInvalidSigningMethod indicates that an unsupported signing method was used
	ErrInvalidSigningMethod = errors.New("invalid signing method")

	// ErrEmptyName indicates that the name is empty or invalid
	ErrEmptyName = errors.New("empty name")

	// ErrEmptyDescription indicates that the description is empty or invalid
	ErrEmptyDescription = errors.New("empty description")

	// ErrEmptyIdentityNumber indicates that the identity number is empty or invalid
	ErrEmptyIdentityNumber = errors.New("empty identity number")

	// ErrEmptyPersonalCode indicates that the personal code is empty or invalid
	ErrEmptyPersonalCode = errors.New("empty personal code")

	// ErrEmptyFirstName indicates that the first name is empty or invalid
	ErrEmptyFirstName = errors.New("empty first name")

	// ErrEmptyLastName indicates that the last name is empty or invalid
	ErrEmptyLastName = errors.New("empty last name")

	// ErrInvalidArguments indicates that the provided request arguments are invalid
	ErrInvalidArguments = errors.New("invalid arguments")

	// ErrFailedToFetchResults indicates that failed to fetch results
	ErrFailedToFetchResults = errors.New("failed to fetch results")

	// ErrRecordNotFound indicates that the requested record could not be found
	ErrRecordNotFound = errors.New("record not found")

	// ErrFailedToCreateRecord indicates that failed to create record
	ErrFailedToCreateRecord = errors.New("failed to create record")

	// ErrFailedToUpdateRecord indicates that failed to update record
	ErrFailedToUpdateRecord = errors.New("failed to update record")

	// ErrFailedToDeleteRecord indicates that failed to delete record
	ErrFailedToDeleteRecord = errors.New("failed to delete record")

	// ErrPermissionNotFound indicates that the requested permission could not be found
	ErrPermissionNotFound = errors.New("permission not found")

	// ErrRoleNotFound indicates that the requested role could not be found
	ErrRoleNotFound = errors.New("role not found")

	// ErrScopeNotFound indicates that the requested scope could not be found
	ErrScopeNotFound = errors.New("scope not found")

	// ErrUserNotFound indicates that the requested user could not be found
	ErrUserNotFound = errors.New("user not found")

	// ErrForbidden indicates that the user is not allowed to perform the requested action
	ErrForbidden = errors.New("access forbidden")

	// ErrUnauthorized indicates that the user is not authorized to perform the requested action
	ErrUnauthorized = errors.New("unauthorized")
)

var (
	Is     = errors.Is
	As     = errors.As
	Unwrap = errors.Unwrap
)
