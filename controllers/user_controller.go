package controllers

import (
	"net/http"
	"os"
	"strings"
	"time"

	"even.com/even-app/config"
	"even.com/even-app/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
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
func Login(context *gin.Context) {
	var input LoginInput
	if err := context.ShouldBindJSON(&input); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Ambil data user dari database berdasarkan email
	var user models.User
	errUser := config.DB.Where("email = ?", input.Email).First(&user).Error
	if errUser != nil {
		context.JSON(http.StatusUnauthorized, gin.H{
			"error": "email yang di masukan belum terdaftar",
		})
		return
	}

	// Bandingkan password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{
			"error": "password yang di masukan salah",
		})
		return
	}

	// 1. Buat klaim token JWT (tambahkan waktu kedaluwarsa agar lebih aman)
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // Token berlaku 24 jam
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 2. Tentukan secret key (sebaiknya simpan di file .env nantinya)

	// 3. Sign token menjadi string
	tokenString, errToken := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if errToken != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gagal menghasilkan token",
		})
		return
	}

	// 4. Kirim token kembali ke client
	context.JSON(http.StatusOK, gin.H{
		"message": "Login berhasil",
		"token":   tokenString,
		"user": gin.H{
			"id":     user.ID,
			"nama":   user.Name,
			"email":  user.Email,
			"events": user.Events,
		},
	})
}

func GetUserData(context *gin.Context) {
	// 1. Ambil user id dari context (pastikan key-nya sama persis dengan di middleware: "UserID")
	userId, exists := context.Get("userId")
	if !exists {
		context.JSON(http.StatusUnauthorized, gin.H{
			"error": "Tidak Terautentikasi",
		})
		return
	}

	var user models.User

	// 2. Query ke database menggunakan GORM dan tangkap error-nya dengan benar
	err := config.DB.Select("id", "name", "email").First(&user, userId).Error
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{
			"error": "User tidak ditemukan",
		})
		return
	}

	// 3. Kirim respons sukses jika user ditemukan
	context.JSON(http.StatusOK, gin.H{
		"message": "Berhasil mengambil data user",
		"user": gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
		},
	})
}
