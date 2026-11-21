package secure

type Tokener interface {
	Generate(subject string) (string, error)
	DecodeToken(tokenString string) (string, error)
}
