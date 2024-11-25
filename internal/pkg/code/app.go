package code

const (
	// ErrAppMapMode - 500: Error occurred while attempt to index from appTemplate using mode.
	ErrAppMapMode int = iota + 110101
	// ErrAppNotFoundRelatedConfig - 500: App config is not found.
	ErrAppNotFoundRelatedConfig
	// ErrAppStatusNotNormal - 500: App is not active.
	ErrAppStatusNotNormal
	// ErrAppCodeNotFound - 500: App code not found.
	ErrAppCodeNotFound
)
