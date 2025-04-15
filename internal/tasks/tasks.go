package tasks

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type TaskStatus int

const (
	TaskOngoing TaskStatus = iota
	TaskCompleted
	TaskClosed
)

var taskStatusName = map[TaskStatus]string{
	TaskOngoing:   "ongoing",
	TaskCompleted: "completed",
	TaskClosed:    "closed",
}

func IsValidStatus(status int) bool {
	switch TaskStatus(status) {
	case TaskOngoing, TaskCompleted, TaskClosed:
		return true
	default:
		return false
	}
}

func (ts TaskStatus) String() string {
	return taskStatusName[ts]
}

type Task struct {
	Id     uint32
	Title  string
	Status TaskStatus
}

type InvalidIdError struct {
	InvalidId uint32
}

func (e *InvalidIdError) Error() string {
	return fmt.Sprintf("No task found with id: %d\n", e.InvalidId)
}

type TaskStorage interface {
	Get(id uint32) (*Task, error)
	Add(task *Task) error
	Delete(id uint32) error
	Update(id uint32, title string, status TaskStatus) error
	GetAll() ([]*Task, error)
}

type LocalStorage struct {
	tasks map[uint32]*Task
	mu    sync.RWMutex
}

func (storage *LocalStorage) Get(id uint32) (*Task, error) {
	storage.mu.RLock()
	defer storage.mu.RUnlock()

	var pTask *Task
	pTask, ok := storage.tasks[id]
	if !ok {
		return pTask, &InvalidIdError{InvalidId: id}
	}

	return pTask, nil
}

func (storage *LocalStorage) GetAll() ([]*Task, error) {
	storage.mu.RLock()
	defer storage.mu.RUnlock()

	tasks := make([]*Task, 0, len(storage.tasks))

	for _, task := range storage.tasks {
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (storage *LocalStorage) Add(task *Task) error {
	storage.mu.Lock()
	defer storage.mu.Unlock()

	_, ok := storage.tasks[task.Id]
	if ok {
		return fmt.Errorf("Duplicate id\n")
	}

	storage.tasks[task.Id] = task
	return nil
}

func (storage *LocalStorage) Delete(id uint32) error {
	storage.mu.Lock()
	defer storage.mu.Unlock()

	_, ok := storage.tasks[id]
	if !ok {
		return &InvalidIdError{InvalidId: id}
	}

	delete(storage.tasks, id)
	return nil
}

func (storage *LocalStorage) Update(id uint32, title string, status TaskStatus) error {
	storage.mu.Lock()
	defer storage.mu.Unlock()

	task, ok := storage.tasks[id]
	if !ok {
		return &InvalidIdError{InvalidId: id}
	}

	task.Title = title
	task.Status = status
	return nil
}

// by default create a localstorage type
func newStorage() TaskStorage {
	return &LocalStorage{
		tasks: make(map[uint32]*Task),
	}
}

var (
	instance TaskStorage
	once     sync.Once
)

func GetStorage() TaskStorage {
	once.Do(func() {
		instance = newStorage()
	})

	return instance
}

type Manager struct {
	Storage TaskStorage
}

func NewManager() *Manager {
	return &Manager{
		Storage: GetStorage(),
	}
}

func (m *Manager) NewTask(title string) (uint32, error) {
	if title == "" {
		return 0, fmt.Errorf("Empty title\n")
	}

	rand.Seed(time.Now().UnixNano())
	id := rand.Uint32()
	task := &Task{
		Id:     id,
		Title:  title,
		Status: TaskOngoing,
	}

	err := m.Storage.Add(task)
	if err != nil {
		return 0, fmt.Errorf("Error creating new task: %v\n", err)
	}

	return id, nil
}

func (m *Manager) GetTask(id uint32) (*Task, error) {
	task, err := m.Storage.Get(id)
	if err != nil {
		return task, fmt.Errorf("Error retrieving task: %v\n", err)
	}

	return task, nil
}

func (m *Manager) DeleteTask(id uint32) error {
	err := m.Storage.Delete(id)
	if err != nil {
		return fmt.Errorf("Error deleting task: %v\n", err)
	}

	return nil
}

func (m *Manager) UpdateTask(id uint32, title string, status TaskStatus) error {
	if title == "" {
		return fmt.Errorf("Empty title\n")
	}

	err := m.Storage.Update(id, title, status)
	if err != nil {
		return fmt.Errorf("Error updating task: %v\n", err)
	}

	return nil
}

func (m *Manager) GetAllTasks() ([]*Task, error) {
	tasks, err := m.Storage.GetAll()
	if err != nil {
		return tasks, fmt.Errorf("Error getting all tasks: %v\n", err)
	}

	return tasks, nil
}
