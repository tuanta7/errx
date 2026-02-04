package yaml

import "errors"

type YAML struct{}

func Parser() *YAML {
	return &YAML{}
}

func (y *YAML) Unmarshal(bytes []byte) (map[string]string, error) {
	return nil, errors.New("not implemented yet")
}

func (y *YAML) Marshal(m map[string]string) ([]byte, error) {
	return nil, errors.New("not implemented yet")
}
