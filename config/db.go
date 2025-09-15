package config

// DB holds database configuration settings.
type DB struct {
	host string // MongoDB connection string (private field)
}

// Host returns the database connection string.
// This is the primary method to access the MongoDB connection URL.
//
// Returns the MongoDB connection string configured via DB_HOST environment variable.
func (d DB) Host() string {
	return d.host
}
