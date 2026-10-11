package controllers

import (
	"net/http"

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
	// hash password
	/*
			seperti itu karena fungsi bcrypt.GenerateFromPassword yang ada di pustaka standar Go (golang.org/x/crypto/bcrypt) dirancang secara khusus untuk menerima parameter bertipe slice of bytes ([]byte), bukan string biasa.

		Berikut adalah alasan teknis mengapa fungsi tersebut menggunakan tipe []byte:

		1. Standar Parameter Fungsi
		Pencipta pustaka bcrypt mendesain fungsi tersebut untuk menerima []byte agar lebih fleksibel dalam menangani berbagai macam format data sandi. Di Go, konversi dari string ke byte ([]byte(input.Password)) adalah hal yang sangat lumrah dan sangat ringan dilakukan.

		2. Keamanan Memori (Memory Safety)
		Dalam pemrograman tingkat rendah atau kriptografi, tipe data byte slice ([]byte) lebih disukai daripada string karena byte slice dapat dimodifikasi langsung di memori (mutable). Meskipun bcrypt tidak selalu langsung menghapusnya, penanganan data sensitif seperti password umumnya dibiasakan menggunakan []byte agar isinya bisa dikosongkan (wipe) dari memori RAM setelah selesai digunakan demi alasan keamanan. (Sebaliknya, tipe data string di Go bersifat immutable atau tidak bisa diubah-ubah setelah dibuat).
	*/
	hashedPassword, errHash := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if errHash != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to hash password",
		})
		return
	}
	// simpan ke database
	user := models.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: string(hashedPassword),
	}
	userCreated := config.DB.Create(&user).Error
	if userCreated != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create user",
		})
		return
	}
	context.JSON(http.StatusCreated, gin.H{
		"message": "User berhasil di buat",
		"user":    user,
	})
}
