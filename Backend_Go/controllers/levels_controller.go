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

type LevelsRepo struct {
	Db *gorm.DB
}

func NewLevelsRepo() *LevelsRepo {
	db := database.InitDb()
	db.AutoMigrate(&models.Levels{},&models.Levels{})
	return &LevelsRepo{Db: db}
}

//get Levelss
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

//get Levels by id
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
