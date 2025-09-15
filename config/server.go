package config

// Server holds HTTP server configuration settings.
// All fields are public and can be accessed directly.
type Server struct {
	Host                                 string // HTTP server host address
	Port                                 string // HTTP server port number
	EnableRedis                          bool   // Whether Redis caching is enabled
	TokenCacheTTLInMinutes               int64  // Token cache time-to-live in minutes
	AuthProviderRefetchIntervalInMinutes int64  // Auth provider data refresh interval in minutes
}

// Deployment holds deployment environment configuration settings.
// All fields are public and can be accessed directly.
type Deployment struct {
	Environment string // Deployment environment name (e.g., development, production)
	Name        string // Application deployment name for identification
}
