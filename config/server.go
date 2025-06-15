package config

type Server struct {
	Host        string
	Port        string
	EnableRedis bool
}

type Deployment struct {
	Environment string
	Name        string
}
