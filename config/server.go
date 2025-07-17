package config

type Server struct {
	Host                                 string
	Port                                 string
	EnableRedis                          bool
	TokenCacheTTLInMinutes               int64
	AuthProviderRefetchIntervalInMinutes int64
}

type Deployment struct {
	Environment string
	Name        string
}
