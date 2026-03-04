package apperrors

import "errors"

var (
	ErrNotFound = errors.New("Данные по запросу не найдены")

	ErrRetrievingData = errors.New("ошибка при получении данных")
	ErrSavingData     = errors.New("ошибка при изменении данных")

	ErrConflict = errors.New("запись с таким ключом уже существует")
)
