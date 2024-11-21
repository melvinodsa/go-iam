package config

type DB struct {
	host string
}

func (d DB) Host() string {
	return d.host
}
