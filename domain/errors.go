package domain

import "errors"

var (
	COULDNT_INIT_CONTAINER_LAUNCHER      = errors.New("couldn't init container launcher")
	CONNECTION_COULDNT_BE_ESTABLISHED    = errors.New("connection couldn't be established")
	CONNECTION_WAS_NOT_ESTABLISHED       = errors.New("connection wasn't established")
	UNKNOWN_COMPONENT_FOR_TESTING        = errors.New("unknown component for testing")
	NO_REQUIRED_ENV_VAR_KEY              = errors.New("couldn't find required env var for container")
	COULDNT_CLOSE_CONTAINER_STATS_READER = errors.New("couldn't close containers stats reader")
)
