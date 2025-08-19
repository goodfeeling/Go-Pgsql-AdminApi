// scheduler/scheduler.go
package scheduler

import (
	"fmt"
	"time"

	"sync"

	"github.com/gbrayhan/microservices-go/src/domain"
	domainScheduledTask "github.com/gbrayhan/microservices-go/src/domain/sys/scheduled_task"
	"github.com/gbrayhan/microservices-go/src/infrastructure/executor"
	logger "github.com/gbrayhan/microservices-go/src/infrastructure/logger"
	"github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/sys/scheduled_task"
	"github.com/go-co-op/gocron"
	"go.uber.org/zap"
)

type TaskScheduler struct {
	scheduler *gocron.Scheduler
	repo      scheduled_task.IScheduledTaskRepository
	logger    *logger.Logger
	tasks     map[int]*gocron.Job
	executor  *executor.TaskExecutorManager
	mutex     sync.RWMutex
}

func NewTaskScheduler(
	repo scheduled_task.IScheduledTaskRepository,
	logger *logger.Logger,
	executor *executor.TaskExecutorManager,
) *TaskScheduler {
	return &TaskScheduler{
		scheduler: gocron.NewScheduler(time.UTC),
		repo:      repo,
		logger:    logger,
		tasks:     make(map[int]*gocron.Job),
		executor:  executor,
	}
}

func (s *TaskScheduler) Start() {
	s.loadTasks()
	s.scheduler.StartAsync()
	s.logger.Info("Task scheduler started")
}

func (s *TaskScheduler) Stop() {
	s.scheduler.Stop()
	s.logger.Info("Task scheduler stopped")
}
func (s *TaskScheduler) loadTasks() {
	filters := domain.DataFilters{
		Matches: map[string][]string{
			"status": {"1"}, // 只加载启用的任务
		},
	}

	result, err := s.repo.SearchPaginated(filters)
	if err != nil {
		s.logger.Error("Failed to load tasks", zap.Error(err))
		return
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	for _, task := range *result.Data {
		// 为每个任务创建本地副本以避免闭包问题
		taskCopy := task
		s.addTaskToScheduleInternal(&taskCopy)
	}
}

// 内部方法，不加锁
func (s *TaskScheduler) addTaskToScheduleInternal(task *domainScheduledTask.ScheduledTask) {
	// 创建一个闭包来捕获当前任务
	taskFunc := func() {
		s.executeTask(task)
	}

	// 使用gocron解析cron表达式并调度任务
	job, err := s.scheduler.Cron(task.CronExpression).Do(taskFunc)
	if err != nil {
		s.logger.Error("Failed to schedule task",
			zap.Int("task_id", task.ID),
			zap.String("task_name", task.TaskName),
			zap.Error(err))
		return
	}

	s.tasks[task.ID] = job
	s.logger.Info("Task scheduled",
		zap.Int("task_id", task.ID),
		zap.String("task_name", task.TaskName),
		zap.String("cron", task.CronExpression))
}

// 公共方法，加锁
func (s *TaskScheduler) addTaskToSchedule(task *domainScheduledTask.ScheduledTask) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.addTaskToScheduleInternal(task)
	return nil
}

func (s *TaskScheduler) executeTask(task *domainScheduledTask.ScheduledTask) {
	s.logger.Info("Executing task",
		zap.Int("task_id", task.ID),
		zap.String("task_name", task.TaskName))

	// 执行任务
	err := s.executor.Execute(task)
	if err != nil {
		s.logger.Error("Task execution failed",
			zap.Int("task_id", task.ID),
			zap.String("task_name", task.TaskName),
			zap.Error(err))
	} else {
		s.logger.Info("Task executed successfully",
			zap.Int("task_id", task.ID),
			zap.String("task_name", task.TaskName))
	}

	// 更新执行时间
	now := time.Now()
	updateData := map[string]interface{}{
		"last_execute_time": &now,
	}

	_, err = s.repo.Update(task.ID, updateData)
	if err != nil {
		s.logger.Error("Failed to update task execution time",
			zap.Int("task_id", task.ID),
			zap.Error(err))
	}
}

// ReloadTasks 重新加载所有任务（用于运行时刷新）
func (s *TaskScheduler) ReloadTasks() error {
	s.logger.Info("Reloading all tasks")

	// 停止所有当前任务
	s.StopAllTasks()

	// 清空任务映射
	s.mutex.Lock()
	s.tasks = make(map[int]*gocron.Job)
	s.mutex.Unlock()

	// 重新加载任务
	s.loadTasks()

	s.logger.Info("Tasks reloaded successfully")
	return nil
}

// StopAllTasks 停止所有任务
func (s *TaskScheduler) StopAllTasks() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for _, job := range s.tasks {
		s.scheduler.RemoveByReference(job)
	}
	s.logger.Info("All tasks stopped")
}

// StartTask 启动单个任务
func (s *TaskScheduler) StartTask(taskID int) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// 如果任务已经在运行，先停止它
	if job, exists := s.tasks[taskID]; exists {
		s.scheduler.RemoveByReference(job)
		delete(s.tasks, taskID)
	}

	// 从数据库获取任务
	task, err := s.repo.GetByID(taskID)
	if err != nil {
		s.logger.Error("Failed to get task by ID",
			zap.Int("task_id", taskID),
			zap.Error(err))
		return err
	}

	// 检查任务状态是否为启用
	if task.Status != 1 {
		s.logger.Warn("Task is not enabled, cannot start", zap.Int("task_id", taskID))
		return fmt.Errorf("task is not enabled")
	}

	// 调度任务
	return s.addTaskToSchedule(task)
}

// StopTask 停止单个任务
func (s *TaskScheduler) StopTask(taskID int) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	job, exists := s.tasks[taskID]
	if !exists {
		s.logger.Warn("Task not found", zap.Int("task_id", taskID))
		return fmt.Errorf("task not found")
	}

	s.scheduler.RemoveByReference(job)
	delete(s.tasks, taskID)

	s.logger.Info("Task stopped", zap.Int("task_id", taskID))
	return nil
}

// AddTask 添加新任务
func (s *TaskScheduler) AddTask(task *domainScheduledTask.ScheduledTask) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// 如果任务已存在，先移除
	if job, exists := s.tasks[task.ID]; exists {
		s.scheduler.RemoveByReference(job)
		delete(s.tasks, task.ID)
	}

	// 如果任务是启用状态，则调度它
	if task.Status == 1 {
		return s.addTaskToSchedule(task)
	}

	s.logger.Info("Task added but not scheduled (disabled)", zap.Int("task_id", task.ID))
	return nil
}

// RemoveTask 移除任务
func (s *TaskScheduler) RemoveTask(taskID int) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// 停止任务调度
	if job, exists := s.tasks[taskID]; exists {
		s.scheduler.RemoveByReference(job)
		delete(s.tasks, taskID)
	}

	s.logger.Info("Task removed", zap.Int("task_id", taskID))
	return nil
}

// UpdateTask 更新任务
func (s *TaskScheduler) UpdateTask(task *domainScheduledTask.ScheduledTask) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// 停止现有任务
	if job, exists := s.tasks[task.ID]; exists {
		s.scheduler.RemoveByReference(job)
		delete(s.tasks, task.ID)
	}

	// 如果任务启用，则重新调度
	if task.Status == 1 {
		return s.addTaskToSchedule(task)
	}

	s.logger.Info("Task updated", zap.Int("task_id", task.ID))
	return nil
}

// GetTaskStatus 获取任务状态
func (s *TaskScheduler) GetTaskStatus(taskID int) (bool, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	_, exists := s.tasks[taskID]
	return exists, nil
}

// ListAllTasks 列出所有任务及其状态
func (s *TaskScheduler) ListAllTasks() map[int]bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	status := make(map[int]bool)
	for id, job := range s.tasks {
		status[id] = job.IsRunning()
	}
	return status
}
