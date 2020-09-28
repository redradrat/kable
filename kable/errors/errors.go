package errors

import "errors"

var (
	StaleRepoCacheIndexError      = errors.New("attempted to write stale cache")
	RepositoryInvalidError        = errors.New("given repository is not a valid kable repository")
	RepositoryAlreadyExistsError  = errors.New("repository is already configured")
	RepositoryNotInitializedError = errors.New("repository is not yet initialized")
	ConfigNotInitializedError     = errors.New("currentConfig is not yet initialized")
	ConfigAlreadyInitializedError = errors.New("currentConfig is already initialized")
	ConceptTypeUnsupportedError   = errors.New("given concept type is not supported")
	RenderTargetUnsupportedError  = errors.New("desired render target is not supported")
	InvalidConceptIdentifierError = errors.New("given concept identifier is invalid")
	ConceptDirInvalidError        = errors.New("directory is not a concept directory")
	InvalidRenderNameError        = errors.New("given app name is invalid (only allowed: 'a-z', '-', '_')")
	ValueTypeNotSupported         = errors.New("given value type is not supported")
)
