package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/task_management_api/models"
)

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func getTasks(c *gin.Context) {

	tasks, err := models.GetTasks()
	checkErr(err)
	if tasks == nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "no record found"})
		return
	} else {
		c.IndentedJSON(http.StatusOK, gin.H{"data": tasks})
	}
}
func getTaskById(c *gin.Context) {
	id := c.Param("id")

	task, err := models.GetTaskById(id)
	checkErr(err)

	if task.Title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no recoeds found"})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"data": task})
	}
}

func addTask(c *gin.Context) {

	var json models.Task

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	success, err := models.AddTask(json)

	if success {
		c.JSON(http.StatusOK, gin.H{"message": "Success", "data": json})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
	}

}

func updateTask(c *gin.Context) {
	var json models.Task

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	taskId, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "invalid id"})
	}

	success, err := models.UpdateTask(json, taskId)
	if success {
		c.JSON(http.StatusOK, gin.H{"message": "success", "data": json})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
	}
}

func deleteTask(c *gin.Context) {
	taskId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid id"})
	}
	success, err := models.DeleteTask(taskId)

	if success {
		c.JSON(http.StatusOK, gin.H{"message": "Success"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"message": err})
	}
}

func main() {
	err := models.ConnectDatabase()
	checkErr(err)

	router := gin.Default()

	router.GET("task", getTasks)
	router.GET("task/:id", getTaskById)
	router.POST("task", addTask)
	router.PUT("task/:id", updateTask)
	router.DELETE("task/:id", deleteTask)

	router.Run("localhost:8080")

}
