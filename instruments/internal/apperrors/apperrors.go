package apperrors

import "errors"

var (
	ErrValidation = errors.New("некорректный запрос")

	ErrNotFound = errors.New("Данные по запросу не найдены")

	ErrRetrievingData     = errors.New("ошибка при получении данных")
	ErrSavingData         = errors.New("ошибка при изменении данных")
	ErrBusinessValidation = errors.New("некорректные данные в запросе")

	ErrSettleCode = errors.New("settleCode должен быть T0, T1, T2 или Tx")
	ErrConflict   = errors.New("запись с таким ключом уже существует")
)
