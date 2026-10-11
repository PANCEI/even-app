package controllers

import (
	"net/http"
	"strings"

	"even.com/even-app/config"
	"even.com/even-app/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// pembuatan struct untuk register
type RegisterInput struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// pembuatan struct untuk login
type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

func Register(context *gin.Context) {
	var input RegisterInput
	if err := context.ShouldBindJSON(&input); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Hash password
	hashedPassword, errHash := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if errHash != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to hash password",
		})
		return
	}

	// Simpan ke database
	user := models.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: string(hashedPassword),
	}

	userCreated := config.DB.Create(&user).Error
	if userCreated != nil {
		// Cek apakah error disebabkan oleh email yang sudah terdaftar (duplikat)
		// Pesan error ini umum untuk PostgreSQL/MySQL jika kolom email diset UNIQUE
		if strings.Contains(userCreated.Error(), "duplicate key") || strings.Contains(userCreated.Error(), "UNIQUE constraint failed") {
			context.JSON(http.StatusBadRequest, gin.H{
				"error": "Email sudah terdaftar",
			})
			return
		}

		context.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create user",
		})
		return
	}

	// Amankan data password agar tidak ikut terkirim ke client di dalam response JSON
	user.Password = ""

	context.JSON(http.StatusCreated, gin.H{
		"message": "User berhasil dibuat",
		"user":    user,
	})
}
