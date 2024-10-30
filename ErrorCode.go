package main

// Error codes
const (
	ErrMissingParameter       = 1001
	ErrUserCreationFailed     = 1002
	ErrUsingSameUsername      = 1003
	ErrFetchingUsersFailed    = 2001
	ErrLoginMissingParameter  = 3001
	ErrInvalidCredentials     = 3002
	ErrLoginDatabaseFailed    = 3003
	ErrUserNotFound           = 4001
	ErrCheckInFailed          = 4002
	ErrFetchingCheckInsFailed = 5001
)

// Error messages map
var errorMessages = map[int]string{
	ErrMissingParameter:       "Missing required parameters",
	ErrUsingSameUsername:      "Username already exists",
	ErrUserCreationFailed:     "Unable to create user",
	ErrFetchingUsersFailed:    "Unable to fetch users",
	ErrLoginMissingParameter:  "Missing login parameters",
	ErrInvalidCredentials:     "Invalid username or password",
	ErrLoginDatabaseFailed:    "Unable to check user credentials",
	ErrUserNotFound:           "User not found",
	ErrCheckInFailed:          "Unable to record check-in",
	ErrFetchingCheckInsFailed: "Unable to fetch check-ins",
}

// GetErrorMessage returns the error message for a given error code
func GetErrorMessage(code int) string {
	return errorMessages[code]
}
