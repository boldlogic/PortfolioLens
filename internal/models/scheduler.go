package models

import "time"

const (
	DateFormat = "2006-01-02"
)

const (
	TaskStatusCreated int = iota
	TaskStatusInProgress
	TaskStatusCompleted
	TaskStatusError
)

type Task struct {
	Id        int        `gorm:"column:id;primaryKey;autoincrement"`
	Uuid      string     `gorm:"column:uuid;unique;not null"`
	ActionId  int        `gorm:"column:actionId;type:int"`
	Action    Action     `gorm:"foreignKey:ActionId;references:Id"`
	CreatedAt time.Time  `gorm:"column:created_at;type:datetime2(3);default:SYSDATETIME()"`
	StartedAt *time.Time `gorm:"column:started_at;type:datetime2(3)"`

	CharCode *string    `gorm:"column:char_code;size:13"`
	DateFrom *time.Time `gorm:"column:date_from;type:date"` //2006-01-11
	DateTo   *time.Time `gorm:"column:date_to;type:date"`   //2006-01-11

	StatusId    int        `gorm:"column:status_id;default:0"`
	ScheduledAt time.Time  `gorm:"column:scheduled_at;type:datetime2(3);default:SYSDATETIME()"`
	Error       string     `gorm:"column:error;type:nvarchar(100)"`
	CompletedAt *time.Time `gorm:"column:completed_at;type:datetime2(3)"`

	UpdatedAt time.Time `gorm:"column:updated_at;type:datetime2(3);default:SYSDATETIME()"`
}

// /Содержит коды действий
type Action struct {
	Id   int    `gorm:"column:id;primaryKey;autoincrement"`
	Code string `gorm:"column:code;unique;type:nvarchar(50);not null"` //Пример: currency.cb.fetch.currency_list
	Name string `gorm:"column:name;type:nvarchar(150)"`                //Пример: Получение справочника валют из ЦБ по www.cbr.ru/scripts/XML_valFull.asp
}
