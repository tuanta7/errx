package parsers

type Parser interface {
	Unmarshal([]byte) (map[string]map[string]string, error)
	Marshal(map[string]any) ([]byte, error)
}
