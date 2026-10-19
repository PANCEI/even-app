package controllers

import (
	"net/http"

	"even.com/even-app/config"
	"even.com/even-app/models"
	"github.com/gin-gonic/gin"
)

func CreateEvent(context *gin.Context) {
	userId := context.GetInt("userId")
	var event models.Event
	err := context.ShouldBindJSON(&event)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	event.UserId = userId
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
func UpdateEvent(context *gin.Context) {
	userId := context.GetInt("userId")
	var event models.Event
	paramsId := context.Param("id")
	if paramsId == "" {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": "ID is required",
		})
		return
	}
	var evenData = config.DB.First(&event, paramsId).Error
	if evenData != nil {
		context.JSON(http.StatusNotFound, gin.H{
			"error": "Data Tidak Ditemukan",
		})
		return
	}
	if event.UserId != userId {
		context.JSON(http.StatusForbidden, gin.H{
			"error": "Anda Tidak Memiliki Akses Untuk Mengubah Data Ini",
		})
		return
	}
	var updatedEvent models.Event
	err := context.ShouldBindJSON(&updatedEvent)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	config.DB.Model(&event).Updates(updatedEvent)
	context.JSON(http.StatusOK, gin.H{
		"message": "Data berhasil di update",
		"event":   event,
	})

}
func DeleteEvent(context *gin.Context) {
	userId := context.GetInt("userId")
	var event models.Event
	paramsId := context.Param("id")
	if paramsId == "" {
		context.JSON(http.StatusBadRequest, gin.H{
			"errorparameter": "ID is required",
		})
		return
	}
	var evenData = config.DB.First(&event, paramsId).Error
	if evenData != nil {
		context.JSON(http.StatusNotFound, gin.H{
			"error": "Data Tidak Ditemukan",
		})
		return
	}
	if event.UserId != userId {
		context.JSON(http.StatusForbidden, gin.H{
			"error": "Anda Tidak Memiliki Akses Untuk Mengubah Data Ini",
		})
		return
	}
	//config.DB.Delete(&event, paramsId)// soft delete karena menggunakan gorm.Model
	config.DB.Unscoped().Delete(&event, paramsId) // hard delete
	context.JSON(http.StatusOK, gin.H{
		"message": "Data berhasil di delete",
	})

}
