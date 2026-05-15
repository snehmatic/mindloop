package task

import (
	"github.com/snehmatic/mindloop/models"
	"gorm.io/gorm"
)

type taskRepository struct {
	db *gorm.DB
}

func NewSQLTaskRepository(db *gorm.DB) TaskRepository {
	return &taskRepository{db: db}
}

func (r *taskRepository) CreateTask(title string, intentID, focusID *uint) (*models.Task, error) {
	t := &models.Task{
		Title:          title,
		IntentID:       intentID,
		FocusSessionID: focusID,
	}
	if err := r.db.Create(t).Error; err != nil {
		return nil, err
	}
	return t, nil
}

func (r *taskRepository) GetTask(id uint) (*models.Task, error) {
	var task models.Task
	if err := r.db.First(&task, id).Error; err != nil {
		return nil, err
	}
	return &task, nil
}

func (r *taskRepository) CompleteTask(id uint) (*models.Task, error) {
	var task models.Task
	if err := r.db.Preload("SubTasks").First(&task, id).Error; err != nil {
		return nil, err
	}

	task.Status = models.TaskStatusCompleted
	if err := r.db.Save(&task).Error; err != nil {
		return nil, err
	}
	return &task, nil
}

func (r *taskRepository) ListTasks() ([]models.Task, error) {
	var tasks []models.Task
	if err := r.db.Preload("SubTasks", func(db *gorm.DB) *gorm.DB {
		return db.Order("Position ASC, CreatedAt ASC")
	}).Order("Position ASC, CreatedAt DESC").Find(&tasks).Error; err != nil {
		return nil, err
	}
	return tasks, nil
}

func (r *taskRepository) AddSubTask(taskID uint, title string) (*models.SubTask, error) {
	st := &models.SubTask{
		TaskID: taskID,
		Title:  title,
	}
	if err := r.db.Create(st).Error; err != nil {
		return nil, err
	}
	return st, nil
}

func (r *taskRepository) GetSubTask(id uint) (*models.SubTask, error) {
	var st models.SubTask
	if err := r.db.First(&st, id).Error; err != nil {
		return nil, err
	}
	return &st, nil
}

func (r *taskRepository) CompleteSubTask(id uint) (*models.SubTask, error) {
	var st models.SubTask
	if err := r.db.First(&st, id).Error; err != nil {
		return nil, err
	}

	st.Status = models.TaskStatusCompleted
	if err := r.db.Save(&st).Error; err != nil {
		return nil, err
	}
	return &st, nil
}

func (r *taskRepository) GetTasksByIntent(intentID uint) ([]models.Task, error) {
	var tasks []models.Task
	if err := r.db.Where("IntentID = ?", intentID).Preload("SubTasks", func(db *gorm.DB) *gorm.DB {
		return db.Order("Position ASC, CreatedAt ASC")
	}).Order("Position ASC, CreatedAt DESC").Find(&tasks).Error; err != nil {
		return nil, err
	}
	return tasks, nil
}

func (r *taskRepository) GetTasksByFocusSession(focusID uint) ([]models.Task, error) {
	var tasks []models.Task
	if err := r.db.Where("FocusSessionID = ?", focusID).Preload("SubTasks", func(db *gorm.DB) *gorm.DB {
		return db.Order("Position ASC, CreatedAt ASC")
	}).Order("Position ASC, CreatedAt DESC").Find(&tasks).Error; err != nil {
		return nil, err
	}
	return tasks, nil
}

func (r *taskRepository) DeleteTask(id uint) error {
	if err := r.db.Where("TaskID = ?", id).Delete(&models.SubTask{}).Error; err != nil {
		return err
	}
	return r.db.Delete(&models.Task{}, id).Error
}

func (r *taskRepository) DeleteSubTask(id uint) error {
	return r.db.Delete(&models.SubTask{}, id).Error
}

func (r *taskRepository) ReorderTasks(ids []uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		for i, id := range ids {
			if err := tx.Model(&models.Task{}).Where("id = ?", id).Update("position", i).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *taskRepository) ReorderSubTasks(ids []uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		for i, id := range ids {
			if err := tx.Model(&models.SubTask{}).Where("id = ?", id).Update("position", i).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *taskRepository) UpdateTask(task *models.Task) error {
	return r.db.Save(task).Error
}

func (r *taskRepository) UpdateSubTask(subTask *models.SubTask) error {
	return r.db.Save(subTask).Error
}

func (r *taskRepository) GetDB() *gorm.DB {
	return r.db
}
