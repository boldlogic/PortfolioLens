package converters

import (
	"errors"
	"fmt"

	"gopkg.in/yaml.v3"
)

var (
	ErrWrongYAML = errors.New("некорректный формат YAML")
)

func DecodeYAML[T any](filebody []byte) (T, error) {
	var v T
	err := yaml.Unmarshal(filebody, &v)
	if err != nil {
		return v, fmt.Errorf("%w: %w", ErrWrongYAML, err)
	}
	return v, nil

}
