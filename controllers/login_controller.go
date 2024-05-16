package controllers

import (
	"knowledgeplus/go-api/models"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthRepo struct {
	Db      *gorm.DB
	UserDb  *gorm.DB
	AdminDb *gorm.DB
}

func NewAuthRepo(db *gorm.DB, userDb *gorm.DB, admiinDb *gorm.DB) *AuthRepo {
	return &AuthRepo{Db: db, AdminDb: admiinDb, UserDb: userDb}
}

var secretKey = os.Getenv("SECRET_KEY")

// LoginHandler handles user authentication and generates a JWT token
func (repository *AuthRepo) LoginHandler(c *gin.Context) {
	var user models.UserLogin

	// Check if the key "userRole" exists in the context
	var repoByRole = repository.Db
	if userRole, ok := c.Get("userRole"); ok {
		switch userRole {
		case "user":
			repoByRole = repository.UserDb
		case "admin":
			repoByRole = repository.AdminDb
		case "owner":
			repoByRole = repository.Db
		default:
		}
	} else {
		repoByRole = repository.UserDb
	}

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Replace this with your actual authentication logic
	currentUser, role, err := isValidUser(repoByRole, user.Email, user.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Generate and return a JWT token with user role
	token := generateToken(currentUser.Email, role)
	c.JSON(http.StatusOK, gin.H{"access-token": token})
}

func generateToken(email, role string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": email,
		"role":  role,
		"exp":   time.Now().Add(time.Hour * 12).Unix(), // Token expiration time
	})

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		panic(err)
	}

	return tokenString
}

func isValidUser(db *gorm.DB, email, password string) (models.User, string, error) {
	var user models.User
	result := db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return models.User{}, "", result.Error
	}

	// Compare the hashed password with the input password
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return models.User{}, "", err
	}

	// Return the user and role
	return user, user.Role, nil
}

func hashPassword(password string) (string, error) {
	// Generate a hashed version of the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}
