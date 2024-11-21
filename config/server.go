package config

type Server struct {
	Host string
	Port string
}

type Deployment struct {
	Environment string
	Name        string
}
