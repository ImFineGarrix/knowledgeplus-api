package controllers

import (
	"errors"
	"knowledgeplus/go-api/database"
	"knowledgeplus/go-api/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserRepo struct {
	Db *gorm.DB
}

func NewUserRepo() *UserRepo {
	db := database.InitDb()
	db.AutoMigrate(&models.User{})
	return &UserRepo{Db: db}
}

// GetUsers retrieves all User records from the database.
func (repository *UserRepo) GetUsers(c *gin.Context) {
	var users []models.UserResponse

	err := models.GetAllUsers(repository.Db, &users)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, users)
}

// GetUserById retrieves a User by its ID from the database.
func (repository *UserRepo) GetUserById(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var user models.UserResponse

	err := models.GetUserById(repository.Db, &user, uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, user)
}

// CreateUser creates a new User record.
func (repository *UserRepo) CreateUser(c *gin.Context) {
	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Hash the user's password before saving it to the database
	hashedPassword, err := hashPassword(user.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash the password"})
		return
	}

	user.Password = hashedPassword

	err = models.CreateUser(repository.Db, &user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// Prepare a user response without the password field
	var userResponse models.UserResponse
	userResponse.ID = user.ID
	userResponse.Name = user.Name
	userResponse.Email = user.Email
	userResponse.Role = user.Role

	c.JSON(http.StatusCreated, userResponse)
}

// UpdateUser updates a User record by ID.
func (repository *UserRepo) UpdateUser(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var updatedUser models.UserResponse

	err := repository.Db.First(&updatedUser, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Record not found!"})
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	if err := c.ShouldBindJSON(&updatedUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = models.UpdateUser(repository.Db, &updatedUser)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, updatedUser)
}

// DeleteUserById deletes a User record by ID.
func (repository *UserRepo) DeleteUserById(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	err := models.DeleteUserById(repository.Db, uint(id))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
