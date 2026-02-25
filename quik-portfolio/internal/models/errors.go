package models

import "errors"

var (
	ErrInstrumentSubTypeCreating = errors.New("не удалось создать подтип инструмента")
	ErrInstrumentSubTypeNotFound = errors.New("подтип инструмента не существует")
	ErrInstrumentCreating        = errors.New("не удалось создать инструмент")
	ErrInstrumentNotFound        = errors.New("инструмент не найден")

	ErrBoardsMerging = errors.New("не заполнить таблицу бордов")
	ErrBoardCreating = errors.New("не удалось создать класс инструмента")

	ErrLimitsNotFound = errors.New("позиции не найдены")

	ErrPortfolioNotFound   = errors.New("портфель не найден")
	ErrPortfolioRetrieving = errors.New("ошибка при получении портфеля")
)
