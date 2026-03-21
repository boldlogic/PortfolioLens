package commonconfig

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/boldlogic/PortfolioLens/pkg/utils/converters"
)

var (
	ErrFileStat     = errors.New("не удалось получить информацию о файле")
	ErrFileRead     = errors.New("не удалось прочитать файл")
	ErrWrongFileExt = errors.New("некорректное расширение файла конфигурации")
)

func readFile(path string) ([]byte, error) {
	_, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("%w %s: %w", ErrFileStat, path, err)
	}

	bs, err := os.ReadFile(path)

	if err != nil {
		return nil, fmt.Errorf("%w %s: %w", ErrFileRead, path, err)
	}
	return bs, nil

}

func DecodeConfig[T any](path string) (T, error) {
	var v T

	fileBody, err := readFile(path)
	if err != nil {
		return v, err
	}

	ext := strings.ToLower(filepath.Ext(path))
	switch {
	case ext == ".yaml" || ext == ".yml":
		v, err = converters.DecodeYAML[T](fileBody)

	case ext == ".json":
		v, err = converters.DecodeJSON[T](fileBody)

	default:
		return v, fmt.Errorf("%w %s", ErrWrongFileExt, path)
	}

	if err != nil {
		return v, err
	}
	return v, nil
}
