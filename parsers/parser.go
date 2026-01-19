package parsers

type Parser interface {
	Unmarshal([]byte) (map[string]string, error)
	Marshal(map[string]string) ([]byte, error)
}
