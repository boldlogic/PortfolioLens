package models

import "errors"

var (
	ErrAssetTypeCreating      = errors.New("не удалось создать класс актива")
	ErrInstrumentTypeCreating = errors.New("не удалось создать тип инструмента")
	ErrInstrumentTypeNotFound = errors.New("тип инструмента не существует")
	ErrInstrumentTypesMerging = errors.New("не заполнить таблицу типов инструментов")

	ErrInstrumentSubTypeCreating = errors.New("не удалось создать подтип инструмента")
	ErrInstrumentSubTypeNotFound = errors.New("подтип инструмента не существует")
	ErrInstrumentCreating        = errors.New("не удалось создать инструмент")
	ErrInstrumentNotFound        = errors.New("инструмент не найден")

	ErrBoardsMerging = errors.New("не заполнить таблицу бордов")
	ErrBoardCreating = errors.New("не удалось создать тип инструмента")

	ErrMLNotFound   = errors.New("текущие позиции по деньгам не найдены")
	ErrMLRetrieving = errors.New("ошибка при получении позиций по деньгам")

	ErrSLNotFound   = errors.New("позиции по бумагам не найдены")
	ErrSLRetrieving = errors.New("ошибка при получении позиций по бумагам")

	ErrLimitsNotFound = errors.New("позиции не найдены")
	//ErrLimitsRetrieving = errors.New("ошибка при получении позиций")

	ErrPortfolioNotFound   = errors.New("портфель не найден")
	ErrPortfolioRetrieving = errors.New("ошибка при получении портфеля")
)
