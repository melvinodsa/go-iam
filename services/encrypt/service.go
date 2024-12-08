package encrypt

type Service interface {
	Encrypt(rawMessage string) (string, error)
	Decrypt(encryptedMessage string) (string, error)
}
