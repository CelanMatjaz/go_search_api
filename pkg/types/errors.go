package types

import "errors"

var (
	RecordDoesNotExist            = errors.New("requested record does not exist")
	UserIdNotProvidedErr          = errors.New("token not provided")
	UserDoesNotExistErr           = errors.New("user does not exist")
	InvalidBodyErr                = errors.New("provided JSON body is not valid")
	PasswordsDoNotMatchErr        = errors.New("passwords do not match")
	WronglyFormattedAuthHeaderErr = errors.New("authentication error is not formatted correctly")
	MissingRequiredHeaderErr      = errors.New("request is missing required header")
)
