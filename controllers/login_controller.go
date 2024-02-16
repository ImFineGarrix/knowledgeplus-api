package controllers

import (
	"knowledgeplus/go-api/database"
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
	Db *gorm.DB
}

func NewAuthRepo() *AuthRepo {
	db := database.InitDb02()
	return &AuthRepo{Db: db}
}

var secretKey = os.Getenv("SECRET_KEY")

// LoginHandler handles user authentication and generates a JWT token
func (repository *AuthRepo) LoginHandler(c *gin.Context) {
	var user models.UserLogin

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Replace this with your actual authentication logic
	currentUser, role, err := isValidUser(repository.Db, user.Email, user.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Generate and return a JWT token with user role
	token := generateToken(currentUser.Email, role)
	c.JSON(http.StatusOK, gin.H{"access-token": token})
}

// // CreateUserHandler creates a new user with hashed password
// func (repository *AuthRepo) CreateUserHandler(c *gin.Context) {
// 	var user models.User

// 	if err := c.ShouldBindJSON(&user); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	// Hash the user's password before saving it to the database
// 	hashedPassword, err := hashPassword(user.Password)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash the password"})
// 		return
// 	}

// 	user.Password = hashedPassword

// 	err = repository.Db.Create(&user).Error
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
// 		return
// 	}

// 	var userResponse models.UserResponse
// 	userResponse.ID = user.ID
// 	userResponse.Name = user.Name
// 	userResponse.Email = user.Email
// 	userResponse.Role = user.Role

// 	c.JSON(http.StatusCreated, userResponse)
// }

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
