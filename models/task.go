package models

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func ConnectDatabase() error {

	db, err := sql.Open("sqlite3", "tasks.db")
	if err != nil {
		return err
	}
	DB = db
	return nil
}

type Task struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	DueDate     string `json:"due_date"`
	Status      string `json:"status"`
}

func GetTasks() ([]Task, error) {
	rows, err := DB.Query("SELECT id, title, description, due_date, status from tasks ")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	taks := make([]Task, 0)

	for rows.Next() {
		singleTask := Task{}
		err = rows.Scan(&singleTask.ID, &singleTask.Title, &singleTask.Description, &singleTask.DueDate, &singleTask.Status)
		if err != nil {
			return nil, err
		}
		taks = append(taks, singleTask)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return taks, err
}

func GetTaskById(id string) (Task, error) {
	stmt, err := DB.Prepare("SELECT id, title, description, due_date, status from tasks WHERE id = ?")
	if err != nil {
		return Task{}, err
	}
	task := Task{}
	sqlErr := stmt.QueryRow(id).Scan(&task.ID, &task.Title, &task.Description, &task.DueDate, &task.Status)
	if sqlErr != nil {
		if sqlErr == sql.ErrNoRows {
			return Task{}, nil
		}
		return Task{}, sqlErr
	}
	return task, nil
}

func AddTask(newtask Task) (bool, error) {
	tx, err := DB.Begin()
	if err != nil {
		return false, err
	}
	stmt, err := tx.Prepare("INSERT INTO tasks (title, description, due_date, status) VALUES (?,?,?,?)")

	if err != nil {
		return false, err
	}

	defer stmt.Close()

	_, err = stmt.Exec(newtask.Title, newtask.Description, newtask.DueDate, newtask.Status)
	if err != nil {
		return false, err
	}
	tx.Commit()
	return true, nil
}

func UpdateTask(ourTask Task, id int) (bool, error) {
	tx, err := DB.Begin()
	if err != nil {
		return false, err
	}
	stmt, err := tx.Prepare("UPDATE tasks SET title = ?, description = ?, due_date = ?, status = ? WHERE id = ? ")
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	_, err = stmt.Exec(ourTask.Title, ourTask.Description, ourTask.DueDate, ourTask.Status, ourTask.ID)
	if err != nil {
		return false, err
	}
	tx.Commit()
	return true, nil
}

func DeleteTask(taskId int) (bool, error) {
	tx, err := DB.Begin()
	if err != nil {
		return false, err
	}
	stmt, err := DB.Prepare("DELETE from tasks where id = ?")

	if err != nil {
		return false, err
	}

	defer stmt.Close()

	_, err = stmt.Exec(taskId)

	if err != nil {
		return false, err
	}
	tx.Commit()
	return true, nil
}