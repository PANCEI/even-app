package main

import (
	"log"

	"even.com/even-app/config"
	"even.com/even-app/controllers"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error saat get env")
	}
	config.ConnectDB()
	server := gin.Default()
	// membuat route grup
	// membuat ruote grup untuk api kurung kurawal harus ada di baris baru
	api := server.Group("/api")
	{
		api.GET("/events", controllers.GetEvent)
		api.POST("/events", controllers.CreateEvent)
		api.GET("/events/:id", controllers.GetEventById)
		api.PUT("/events/:id", controllers.UpdateEvent)
		api.DELETE("/events/:id", controllers.DeleteEvent)

		api.POST("/register", controllers.Register)
		api.POST("/login", controllers.Login)
	}

	// route
	// server.GET("events", getEvents)
	// fmt.Println("Hello, World!")

	server.Run(":8089")
}

// function handler
// func getEvents(c *gin.Context) {
// 	events := models.GetAllEvents()
// 	if len(events) == 0 {
// 		c.JSON(http.StatusNotFound, gin.H{
// 			"message": "Data tidak ditemukan",
// 		})
// 		return
// 	}
// 	c.JSON(http.StatusOK, gin.H{
// 		"message": "Data Berhasil di dapatkan!",
// 		"data":    events,
// 	})
// }
// func createEvent(c *gin.Context) {
// 	var event models.Event
// 	err := c.ShouldBindJSON(&event)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"message": "request tidak valid",
// 			"error":   err.Error(),
// 		})
// 		return
// 	}

// 	// dummy data

// 	event.UserId = 1

// 	// save data event
// 	event.Save()

// 	c.JSON(http.StatusOK, gin.H{
// 		"message": "event berhasil dibuat",
// 		"data":    event,
// 	})
// }
