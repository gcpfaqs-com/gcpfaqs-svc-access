package secure

type Hasher interface {
	Hash(data string) (string, error)
	Verify(hashedData, data string) error
}
