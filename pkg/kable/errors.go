package kable

import "errors"

var (
	StaleRepoCacheIndexError      = errors.New("attempted to write stale cache")
	RepositoryInvalidError        = errors.New("given repository is not a valid kable repository")
	RepositoryAlreadyExistsError  = errors.New("repository is already configured")
	RepositoryNotInitializedError = errors.New("repository is not yet initialized")
	ConfigNotInitializedError     = errors.New("currentConfig is not yet initialized")
	ConfigAlreadyInitializedError = errors.New("currentConfig is already initialized")
	ConceptTypeUnsupported        = errors.New("given concept type is not supported")
	RenderTargetUnsupported       = errors.New("desired render target is not supported")
)
