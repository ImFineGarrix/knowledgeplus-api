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
	r.Use()

	r.GET("ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, "pong")
	})

	CareerRepo := controllers.NewCareerRepo()
	CategoriesRepo := controllers.NewCategoriesRepo()
	SkillsRepo := controllers.NewSkillsRepo()
	LevelsRepo := controllers.NewLevelsRepo()
	CourseRepo := controllers.NewCourseRepo()

	//not test
	r.POST("/careers", CareerRepo.CreateCareer)
	r.GET("/careers", CareerRepo.GetCareers)
	r.GET("/careers/:id", CareerRepo.GetCareer)
	// r.PUT("/Careers/:id", Repo.UpdateCareer)
	// r.DELETE("/Careers/:id", Repo.DeleteCareer)

	r.POST("/categories", CategoriesRepo.CreateCategory)
	r.GET("/categories", CategoriesRepo.GetCategories)
	r.GET("/categories/:id", CategoriesRepo.GetCategoryById)

	r.GET("/skills", SkillsRepo.GetSkills)
	r.GET("/skills/:id", SkillsRepo.GetSkillById)

	r.GET("/levels", LevelsRepo.GetLevels)
	r.GET("/levels/:id", LevelsRepo.GetLevelById)

	r.GET("/courses", CourseRepo.GetCourses)
	r.GET("/courses/:id", CourseRepo.GetCourse)

	return r
}
