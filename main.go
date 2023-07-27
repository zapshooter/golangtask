package main

import (
	"database/sql"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

type Task struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Due_Date    string `json:"due_date"`
	Status      string `json:"status"`
}

var db *sql.DB

func main() {
	// Connecting to sqlite database
	var err error
	db, err = sql.Open("sqlite3", "./taskk.db")
	if err != nil {
		log.Fatal("Error opening database:", err)
	}
	defer db.Close()

	// Creating the "taskk" table if it doesn't exist
	createTable()

	// Initialize the Gin router
	r := gin.Default()

	// Define routes for CRUD operations
	r.POST("/taskk", createTask)
	r.GET("/taskk/:id", getTaskbyid)
	r.GET("/taskk", gettaskk)
	r.PUT("/taskk/:id", updateTask)
	r.DELETE("/taskk/:id", deleteTask)

	// Start the server
	r.Run(":8000")
}
func createTable() {
	query := `
	CREATE TABLE IF NOT EXISTS taskk (
		id INTEGER PRIMARY KEY,
		title TEXT NOT NULL,
		description TEXT NOT NULL,
		due_date DATE NOT NULL,
		status TEXT NOT NULL
	);`

	_, err := db.Exec(query)
	if err != nil {
		log.Fatal("Error in creating a table:", err)
	}
}

func createTask(c *gin.Context) {
	var task Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(400, gin.H{"error": "Invalid data"})
		return
	}
	if task.Title == "" || task.Description == "" || task.Due_Date == "" || task.Status == "" {
		c.JSON(400, gin.H{"error": "Missing required fields"})
		return
	}

	task.ID = time.Now().UnixNano() / int64(time.Millisecond)

	stmt, err := db.Prepare("INSERT INTO taskk (id, title, description, due_date, status) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		c.JSON(500, gin.H{"error": "Internal Error"})
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(task.ID, task.Title, task.Description, task.Due_Date, task.Status)
	if err != nil {
		c.JSON(500, gin.H{"error": "something went wrong"})
		return
	}

	c.JSON(201, task)
}

func getTaskbyid(c *gin.Context) {
	var task Task
	id := c.Param("id")

	row := db.QueryRow("SELECT id, title, description, due_date, status FROM taskk WHERE id = ?", id)
	err := row.Scan(&task.ID, &task.Title, &task.Description, &task.Due_Date, &task.Status)
	if err != nil {
		c.JSON(404, gin.H{"error": "Task not found"})
		return
	}

	c.JSON(200, task)
}

func gettaskk(c *gin.Context) {
	var taskk []Task

	rows, err := db.Query("SELECT id, title, description, due_date, status FROM taskk")
	if err != nil {
		c.JSON(500, gin.H{"error": "something went wrong"})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var task Task
		err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.Due_Date, &task.Status)
		if err != nil {
			c.JSON(500, gin.H{"error": "something went wrong"})
			return
		}
		taskk = append(taskk, task)
	}

	c.JSON(200, taskk)
}

func updateTask(c *gin.Context) {
	id := c.Param("id")
	var task Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(400, gin.H{"error": "Invalid data"})
		return
	}

	stmt, err := db.Prepare("UPDATE taskk SET title = ?, description = ?, due_date = ?, status = ? WHERE id = ?")
	if err != nil {
		c.JSON(500, gin.H{"error": "Something went wrong"})
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(task.Title, task.Description, task.Due_Date, task.Status, id)
	if err != nil {
		c.JSON(500, gin.H{"error": "Something went wrong"})
		return
	}

	c.JSON(200, task)
}

func deleteTask(c *gin.Context) {
	id := c.Param("id")

	stmt, err := db.Prepare("DELETE FROM taskk WHERE id = ?")
	if err != nil {
		c.JSON(500, gin.H{"error": "Something went wrong"})
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	if err != nil {
		c.JSON(500, gin.H{"error": "Something Went wrong"})
		return
	}

	c.JSON(200, gin.H{"message": "Task deleted successfully"})
}
