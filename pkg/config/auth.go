package config

// Auth - данные для авторизации в сервисах.
type Auth struct {
	PublicKeyBase64 string `env:"RSA_PUBLIC_KEY_BASE64"`
	ApiKey          string `env:"X_API_KEY"`
}
