package controllers

import (
	"errors"
	"knowledgeplus/go-api/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type LevelsRepo struct {
	Db      *gorm.DB
	UserDb  *gorm.DB
	AdminDb *gorm.DB
}

func NewLevelsRepo(db *gorm.DB, userDb *gorm.DB, adminDb *gorm.DB) *LevelsRepo {
	db.AutoMigrate(&models.Levels{}, &models.Levels{})
	userDb.AutoMigrate(&models.Levels{}, &models.Levels{})
	adminDb.AutoMigrate(&models.Levels{}, &models.Levels{})
	return &LevelsRepo{Db: db, UserDb: userDb, AdminDb: adminDb}
}

// get Levelss
func (repository *LevelsRepo) GetLevels(c *gin.Context) {
	var Levels []models.Level
	print(Levels)
	err := models.GetLevels(repository.Db, &Levels)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, Levels)
}

// get Levels by id
func (repository *LevelsRepo) GetLevelById(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var Levels models.Level
	err := models.GetLevelById(repository.Db, &Levels, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, Levels)
}
