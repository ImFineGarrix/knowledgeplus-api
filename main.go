package main

import (
	"knowledgeplus/go-api/controllers"
	"net/http"

	"github.com/gin-contrib/cors"
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
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{
		"http://localhost:3000",
		"https://cp23sj2.sit.kmutt.ac.th:3000",
		"http://cp23sj2.sit.kmutt.ac.th:3001",
		"http://localhost:3001",
		"https://cp23sj2.sit.kmutt.ac.th:3001",
		"http://cp23sj2.sit.kmutt.ac.th:3001",
	}

	// config := cors.DefaultConfig()
	// config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE"}
	// config.AllowOrigins = []string{
	// 	"http://localhost:3000",
	// 	"https://cp23sj2.sit.kmutt.ac.th:3000",
	// 	"http://cp23sj2.sit.kmutt.ac.th:3001",
	// 	"http://localhost:3001",
	// 	"https://cp23sj2.sit.kmutt.ac.th:3001",
	// 	"http://cp23sj2.sit.kmutt.ac.th:3001",
	// }
	// config.AllowHeaders = []string{
	// 	"Origin",
	// 	"Content-Length",
	// 	"Content-Type",
	// 	"Authorization", // Add any other headers your API might use
	// }
	r.Use(cors.New(config))

	// r.Use(CORSMiddleware())

	defaultPath := r.Group("/api")

	defaultPath.GET("ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, "pong")
	})

	CareerRepo := controllers.NewCareerRepo()
	CategoriesRepo := controllers.NewCategoriesRepo()
	SkillRepo := controllers.NewSkillRepo()
	LevelsRepo := controllers.NewLevelsRepo()
	CourseRepo := controllers.NewCourseRepo()

	defaultPath.POST("/careers", CareerRepo.CreateCareer)
	defaultPath.GET("/careers", CareerRepo.GetCareers)
	defaultPath.GET("/careers_have_categories", CareerRepo.GetCareersWithHaveCategories)
	defaultPath.GET("/careers/:id", CareerRepo.GetCareer)
	defaultPath.PUT("/careers/:id", CareerRepo.UpdateCareer)
	defaultPath.DELETE("/careers/:id", CareerRepo.DeleteCareer)

	defaultPath.GET("/categories", CategoriesRepo.GetCategories)
	defaultPath.GET("/categories/:id", CategoriesRepo.GetCategoryById)
	defaultPath.POST("/categories", CategoriesRepo.CreateCategory)
	defaultPath.PUT("/categories/:id", CategoriesRepo.UpdateCategory)
	defaultPath.DELETE("/categories/:id", CategoriesRepo.DeleteCategoryById)

	defaultPath.GET("/skills", SkillRepo.GetSkills)
	defaultPath.GET("/skills/:id", SkillRepo.GetSkillById)
	defaultPath.POST("/skills", SkillRepo.CreateSkill)
	defaultPath.PUT("/skills/:id", SkillRepo.UpdateSkill)
	defaultPath.DELETE("/skills/:id", SkillRepo.DeleteSkillById)

	defaultPath.GET("/levels", LevelsRepo.GetLevels)
	defaultPath.GET("/levels/:id", LevelsRepo.GetLevelById)

	defaultPath.GET("/courses", CourseRepo.GetCourses)
	defaultPath.GET("/courses/:id", CourseRepo.GetCourse)

	return r
}

// func CORSMiddleware() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		// Get the origin from the request header
// 		origin := c.Request.Header.Get("Origin")

// 		// Check if the origin is allowed
// 		allowedOrigins := []string{
// 			"http://localhost:3000",
// 			"https://cp23sj2.sit.kmutt.ac.th:3000",
// 			"http://cp23sj2.sit.kmutt.ac.th:3001",
// 			"http://localhost:3001",
// 			"https://cp23sj2.sit.kmutt.ac.th:3001",
// 			"http://cp23sj2.sit.kmutt.ac.th:3001",
// 		}

// 		// Check if the request origin is in the allowed list
// 		allowed := false
// 		for _, allowedOrigin := range allowedOrigins {
// 			if origin == allowedOrigin {
// 				allowed = true
// 				break
// 			}
// 		}

// 		if allowed {
// 			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
// 			c.Writer.Header().Set("Access-Control-Max-Age", "86400")
// 			c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
// 			c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
// 			c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")
// 			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

// 			if c.Request.Method == "OPTIONS" {
// 				c.AbortWithStatus(200)
// 			} else {
// 				c.Next()
// 			}
// 		} else {
// 			c.AbortWithStatus(403) // Forbidden if the origin is not allowed
// 		}
// 	}
// }
