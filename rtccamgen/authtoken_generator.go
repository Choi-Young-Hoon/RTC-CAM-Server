package rtccamgen

import "math/rand"

type AuthTokenGeneratorInterface interface {
	GenerateAuthToken() string
}

func NewAuthTokenGenerator() *AuthTokenGenerator {
	return &AuthTokenGenerator{}
}

type AuthTokenGenerator struct {
}

func (r *AuthTokenGenerator) GenerateAuthToken() string {
	const charset = "abcdefghijklmnopqrstuvwxyz" +
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	authToken := make([]byte, 60)
	for i := range authToken {
		authToken[i] = charset[rand.Intn(len(charset))]
	}

	return string(authToken)
}
