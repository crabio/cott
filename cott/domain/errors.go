package domain

import "errors"

var (
	CONNECTION_COULDNT_BE_ESTABLISHED = errors.New("connection couldn't be established")
	CONNECTION_WAS_NOT_ESTABLISHED    = errors.New("connection wasn't established")
	UNKNOWN_COMPONENT_FOR_TESTING     = errors.New("unknown component for testing")
)
