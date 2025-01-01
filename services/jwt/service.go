package jwt

type Service interface {
	GenerateToken(claims map[string]interface{}, expiryTimeInSeconds int64) (string, error)
	ValidateToken(token string) (map[string]interface{}, error)
}
