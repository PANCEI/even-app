package controllers

import (
	"net/http"

	"even.com/even-app/config"
	"even.com/even-app/models"
	"github.com/gin-gonic/gin"
)

func CreateEvent(context *gin.Context) {
	var event models.Event
	err := context.ShouldBindJSON(&event)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	event.UserId = 1
	config.DB.Create(&event)
	context.JSON(http.StatusCreated, gin.H{
		"message": "Data berhasil di buat",
		"event":   event,
	})
}
func GetEvent(context *gin.Context) {
	var events []models.Event
	config.DB.Find(&events)
	context.JSON(http.StatusOK, gin.H{
		"message": "Data berhasil di get",
		"event":   events,
	})
}
func GetEventById(context *gin.Context) {
	var events models.Event
	paramsId := context.Param("id")
	if paramsId == "" {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": "ID is required",
		})
		return
	}
	config.DB.First(&events, paramsId)
	if events.ID == 0 {
		context.JSON(http.StatusNotFound, gin.H{
			"error": "Data Tidak Ditemukan",
		})
		return
	}
	context.JSON(http.StatusOK, gin.H{
		"message": "Data berhasil di get",
		"event":   events,
	})
}
