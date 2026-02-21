package storage

import (
	"errors"
	"fmt"
	"time"

	"github.com/boldlogic/PortfolioLens/pkg/models"
	"gorm.io/gorm"
)

func (st *Storage) GetActionId(code string) (models.Action, error) {
	var result models.Action

	err := st.db.Where("code=?", code).First(&result).Error

	if err != nil {

		return models.Action{}, err
	}
	return result, nil
}

func (st *Storage) GetAction(id int) (models.Action, error) {
	var result models.Action

	err := st.db.Where("id=?", id).First(&result).Error

	if err != nil {

		return models.Action{}, err
	}
	return result, nil
}

func (st *Storage) SaveAction(row *models.Action) error {

	result := st.db.Where(models.Action{Code: row.Code}).Assign(models.Action{
		Id:   row.Id,
		Code: row.Code,
		Name: row.Name,
	}).FirstOrCreate(row)
	if result.Error != nil {
		return fmt.Errorf("Произошла ошибка при сохранении Action. %w", result.Error)
	}

	return nil
}

// func (st *Storage) SaveTask(row *models.Task) error {

// 	result := st.db.Where(models.Task{: row.Code}).Assign(models.Task{
// 		Id:   row.Id,
// 		Code: row.Code,
// 		Name: row.Name,
// 	}).FirstOrCreate(row)
// 	if result.Error != nil {
// 		return fmt.Errorf("Произошла ошибка при сохранении Action. %w", result.Error)
// 	}

// 	return nil
// }

func (st *Storage) GetTask(uuid string) (models.Task, error) {
	var result models.Task

	err := st.db.Where("uuid=?", uuid).First(&result).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.Task{}, models.ErrTaskNotFound
		}
		return models.Task{}, err
	}
	return result, nil
}

const (
	createTaskOut = `
	insert into tasks (uuid, actionId, char_code, date_from, date_to)
	output inserted.id, inserted.uuid, inserted.created_at, inserted.scheduled_at 
	values (?, ?,?,?,?)
	`
	fetchTask = `
	UPDATE TOP (1) tasks
	SET started_at = GETDATE()
		,status_id = ?
		,updated_at = GETDATE()
	output inserted.*

	WHERE 1=1 
		and status_id = ?
		and scheduled_at <= GETDATE()
	`
	setTaskStatusCompleted = `
	UPDATE tasks
	SET status_id = ?
		,updated_at = GETDATE()
		,completed_at = GETDATE()
	WHERE id=?
	`
	setTaskStatusError = `
	UPDATE tasks
	SET status_id = ?
		,updated_at = GETDATE()
		,error = ?
	WHERE id=?
	`
)

func (st *Storage) CreateTask(row *models.Task) (models.Task, error) {

	type taskShort struct {
		Id          int
		Uuid        string
		CreatedAt   time.Time
		ScheduledAt time.Time
	}
	var created taskShort
	result := st.db.Raw(createTaskOut, row.Uuid, row.ActionId, row.CharCode, row.DateFrom, row.DateTo).Scan(&created)

	if result.Error != nil || result.RowsAffected == 0 {
		return models.Task{}, models.ErrTaskCreating
	}

	return models.Task{
		Id:          created.Id,
		Uuid:        created.Uuid,
		CreatedAt:   created.CreatedAt,
		ScheduledAt: created.ScheduledAt,
	}, nil
}

func (st *Storage) FetchTask(status int, newStatus int) (models.Task, error) {

	var out models.Task
	result := st.db.Raw(fetchTask, newStatus, status).Scan(&out)

	if result.Error != nil {
		return models.Task{}, result.Error
	}
	if result.RowsAffected == 0 || out.Id == 0 {
		return models.Task{}, models.ErrNoNewTasks
	}

	return out, nil
	//{
	// Id:          out.Id,
	// Uuid:        out.Uuid,
	// ActionId:    out.ActionId,
	// CreatedAt:   out.CreatedAt,
	// StartedAt:   out.StartedAt,
	// CharCode:    out.CharCode,
	// DateFrom:    out.DateFrom,
	// DateTo:      out.DateTo,
	// StatusId:    out.Id,
	// ScheduledAt: out.ScheduledAt,
	// CompletedAt: out.CompletedAt,
	//}

}

func (st *Storage) SetTaskStatusCompleted(id int) error {
	result := st.db.Exec(setTaskStatusCompleted, models.TaskStatusCompleted, id)

	if result.Error != nil {
		return models.ErrTaskUpdating
	}

	return nil
}

func (st *Storage) SetTaskStatusError(id int, err string) error {
	result := st.db.Exec(setTaskStatusError, models.TaskStatusError, err, id)

	if result.Error != nil {
		return models.ErrTaskUpdating
	}

	return nil
}
