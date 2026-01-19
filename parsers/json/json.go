package json

import "encoding/json"

type JSON struct{}

func Parser() *JSON {
	return &JSON{}
}

func (j *JSON) Unmarshal(bytes []byte) (map[string]string, error) {
	var m map[string]string
	if err := json.Unmarshal(bytes, &m); err != nil {
		return nil, err
	}

	return m, nil
}

func (j *JSON) Marshal(m map[string]string) ([]byte, error) {
	return json.Marshal(m)
}
