package log

import "errors"

var (
	ErrRemoteConfig   = errors.New("config is required in remote mode")
	ErrRemoteEndpoint = errors.New("endpoint is required")
)
