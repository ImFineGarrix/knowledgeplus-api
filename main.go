package main

import (
	"knowledgeplus/go-api/controllers"
	"net/http"

	"github.com/gin-gonic/gin"
)

// func init() {
// 	initializers.LoadEnvVariables()
// }

func main() {
	r := setupRouter()
	_ = r.Run(":8081")
}

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.Use(corsMiddleware())

	defaultPath := r.Group("/api")

	defaultPath.GET("ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, "pong")
	})

	CareerRepo := controllers.NewCareerRepo()
	CategoriesRepo := controllers.NewCategoriesRepo()
	SkillsRepo := controllers.NewSkillsRepo()
	LevelsRepo := controllers.NewLevelsRepo()
	CourseRepo := controllers.NewCourseRepo()

	defaultPath.POST("/careers", CareerRepo.CreateCareer)
	defaultPath.GET("/careers", CareerRepo.GetCareers)
	defaultPath.GET("/careers/:id", CareerRepo.GetCareer)
	// defaultPath.PUT("/Careers/:id", Repo.UpdateCareer)
	// defaultPath.DELETE("/Careers/:id", Repo.DeleteCareer)

	defaultPath.POST("/categories", CategoriesRepo.CreateCategory)
	defaultPath.GET("/categories", CategoriesRepo.GetCategories)
	defaultPath.GET("/categories/:id", CategoriesRepo.GetCategoryById)

	defaultPath.POST("/skills", SkillsRepo.CreateSkill)
	defaultPath.GET("/skills", SkillsRepo.GetSkills)
	defaultPath.GET("/skills/:id", SkillsRepo.GetSkillById)

	defaultPath.GET("/levels", LevelsRepo.GetLevels)
	defaultPath.GET("/levels/:id", LevelsRepo.GetLevelById)

	defaultPath.GET("/courses", CourseRepo.GetCourses)
	defaultPath.GET("/courses/:id", CourseRepo.GetCourse)

	return r
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
