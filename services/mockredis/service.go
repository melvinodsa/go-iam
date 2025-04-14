package mockredis

// Service represents the Redis dummy service interface.
type Service interface {
	Set(key string, value string)
	Get(key string) (string, error)
	Delete(key string) error
}
