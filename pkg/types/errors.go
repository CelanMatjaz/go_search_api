package types

import "errors"

var (
	AccountAlreadyExistsErr = errors.New("account already exists")
	AccountDoesNotExistErr  = errors.New("account does not")
	RecordDoesNotExistErr   = errors.New("record does not exist")
	InvalidJsonBodyErr      = errors.New("invalid json body")
	InvalidTokenErr         = errors.New("invalid token")
	MissingCookieErr        = errors.New("missing cookie")
	UnknownOAuthProviderErr = errors.New("unknown oauth provider")
	OAuthProviderIssuesErr  = errors.New("authentication issues with oauth provider")
)

var (
	AccountAlreadyExistsErrors     = []string{"Account with provided email already exists"}
	InternalServerErrors           = []string{"Internal server error"}
	UnparseableJsonBodyErrors      = []string{"Provided JSON body is not parseable"}
	InvalidJsonBodyErrors          = []string{"Provided JSON body is not valid"}
	InvalidPasswordErrors          = []string{"Invalid password"}
	InvalidTokenErrors             = []string{"Invalid auth token"}
	UnauthenticatedErrors          = []string{"Unauthenticated"}
	LoggedOutDueToInactivityErrors = []string{"Logged out due to inactivity"}
	UnknownOAuthProviderErrors     = []string{"Unknown OAuth provider"}
	OAuthProviderIssueErrors       = []string{"Issue with authenticating with OAuth provider"}
	UnverifiedOAuthEmailErrors     = []string{"Email used to authenticate with OAuth is unverified"}
	InvalidPathParamErrors         = []string{"Path param is not valid"}
)
