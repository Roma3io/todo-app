package postgresql

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

type Storage struct {
	db *sql.DB
}

type Task struct {
	ID          int
	Title       string
	Description string
	Status      string
}

func NewStorage(storagePath string) (*Storage, error) {
	const op = "storage.postgresql.NewStorage"
	db, err := sql.Open("postgres", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	stmt, err := db.Prepare(`CREATE TABLE IF NOT EXISTS tasks(
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    status VARCHAR(50) NOT NULL);
`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) CreateTask(title, description, status string) (int, error) {
	const op = "storage.postgresql.CreateTask"

	stmt, err := s.db.Prepare("INSERT INTO tasks (title, description, status) VALUES ($1, $2, $3) RETURNING id")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	var id int
	err = stmt.QueryRow(title, description, status).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (s *Storage) GetTask(id int) (Task, error) {
	const op = "storage.postgresql.GetTask"

	stmt, err := s.db.Prepare("SELECT id, title, description, status FROM tasks WHERE id = $1")
	if err != nil {
		return Task{}, fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	var task Task
	err = stmt.QueryRow(id).Scan(&task.ID, &task.Title, &task.Description, &task.Status)
	if err != nil {
		return Task{}, fmt.Errorf("%s: %w", op, err)
	}
	return task, nil
}

func (s *Storage) UpdateTask(id int, title, description, status string) error {
	const op = "storage.postgresql.UpdateTask"

	stmt, err := s.db.Prepare("UPDATE tasks SET title = $1, description = $2, status = $3 WHERE id = $4")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(title, description, status, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) DeleteTask(id int) error {
	const op = "storage.postgresql.DeleteTask"

	stmt, err := s.db.Prepare("DELETE FROM tasks WHERE id = $1")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	result, err := stmt.Exec(id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (s *Storage) GetAllTasks() ([]Task, error) {
	const op = "storage.postgresql.GetAllTasks"

	rows, err := s.db.Query("SELECT id, title, description, status FROM tasks")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		if err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.Status); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return tasks, nil
}
