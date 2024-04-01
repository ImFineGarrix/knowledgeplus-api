package controllers

import (
	"errors"
	"fmt"
	"knowledgeplus/go-api/models"
	"knowledgeplus/go-api/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type UserRepo struct {
	Db      *gorm.DB
	UserDb  *gorm.DB
	AdminDb *gorm.DB
}

func NewUserRepo(db *gorm.DB, userDb *gorm.DB, adminDb *gorm.DB) *UserRepo {
	db.AutoMigrate(&models.User{})
	userDb.AutoMigrate(&models.User{})
	adminDb.AutoMigrate(&models.User{})
	return &UserRepo{Db: db, AdminDb: adminDb, UserDb: userDb}
}

// GetUsers retrieves all User records from the database.
func (repository *UserRepo) GetUsers(c *gin.Context) {
	var users []models.UserGetResponse
	// Check if the key "userRole" exists in the context
	var repoByRole = repository.Db
	if userRole, ok := c.Get("userRole"); ok {
		switch userRole {
		case "user":
			repoByRole = repository.UserDb
		case "admin":
			repoByRole = repository.AdminDb
			fmt.Println("Admin role")
		case "owner":
			repoByRole = repository.Db
		default:
		}
	} else {
		repoByRole = repository.UserDb
	}

	err := models.GetAllUsers(repoByRole, &users)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, users)
}

// GetUserById retrieves a User by its ID from the database.
func (repository *UserRepo) GetUserById(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var user models.UserGetResponse

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

	err := models.GetUserById(repoByRole, &user, uint(id))
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
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			out := make([]response.ErrorMsg, len(ve))
			for i, fe := range ve {
				out[i] = response.ErrorMsg{
					Code:    http.StatusBadRequest,
					Field:   fe.Field(),
					Message: response.GetErrorMsg(fe),
				}
			}
			c.JSON(http.StatusBadRequest, out)
		}
		return
	}

	// Check if the name already exists in the database
	var existingUser models.User
	if err := repoByRole.Where("name = ?", user.Name).First(&existingUser).Error; err == nil {
		out := response.ErrorMsg{
			Code:    http.StatusBadRequest,
			Field:   "Name",
			Message: "Name already used.",
		}
		c.JSON(http.StatusBadRequest, out)
		return
	}

	// Hash the user's password before saving it to the database
	hashedPassword, err := hashPassword(user.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash the password"})
		return
	}

	user.Password = hashedPassword

	err = models.CreateUser(repoByRole, &user)
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

	err := repoByRole.First(&updatedUser, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Record not found!"})
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	if err := c.ShouldBindJSON(&updatedUser); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			out := make([]response.ErrorMsg, len(ve))
			for i, fe := range ve {
				out[i] = response.ErrorMsg{
					Code:    http.StatusBadRequest,
					Field:   fe.Field(),
					Message: response.GetErrorMsg(fe),
				}
			}
			c.JSON(http.StatusBadRequest, out)
		}
		return
	}

	// Check if the name already exists in the database
	var existingUser models.Organizations
	if err := repoByRole.Where("name = ? AND user_id != ?", updatedUser.Name, id).First(&existingUser).Error; err == nil {
		out := response.ErrorMsg{
			Code:    http.StatusBadRequest,
			Field:   "Name",
			Message: "Name already used.",
		}
		c.JSON(http.StatusBadRequest, out)
		return
	}

	// Hash the user's password before saving it to the database
	hashedPassword, err := hashPassword(updatedUser.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash the password"})
		return
	}

	updatedUser.Password = hashedPassword

	err = models.UpdateUser(repoByRole, &updatedUser)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	// Prepare a user response without the password field
	var userResponse models.UserResponse
	userResponse.ID = updatedUser.ID
	userResponse.Name = updatedUser.Name
	userResponse.Email = updatedUser.Email
	userResponse.Role = updatedUser.Role

	c.JSON(http.StatusOK, userResponse)
}

// DeleteUserById deletes a User record by ID.
func (repository *UserRepo) DeleteUserById(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
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

	err := models.DeleteUserById(repoByRole, uint(id))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
