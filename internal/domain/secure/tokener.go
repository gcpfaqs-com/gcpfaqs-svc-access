package secure

type Tokener interface {
	Generate(subject string) (string, error)
}
