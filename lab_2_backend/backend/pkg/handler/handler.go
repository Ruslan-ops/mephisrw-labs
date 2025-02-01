package handler

import (
	"backend/pkg/service"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	Service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{Service: service}
}

func (h *Handler) InitRoutes() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://mephi22.undersite.ru", "https://mephi71.ru", "http://127.0.0.1:9000", "http://localhost:9000", "http://127.0.0.1:9001", "http://localhost:9001", "http://localhost:9002", "http://127.0.0.1:9002", "http://localhost:9003", "http://127.0.0.1:9003"},
		AllowMethods:     []string{"PUT", "GET", "POST", "DELETE"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type", "lab-token"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	lab1a := router.Group("/lab1a")
	{
		lecturerLab1a := lab1a.Group("/open", h.CheckFirstAHeaderLecturer)
		{
			lecturerLab1a.POST("", h.OpenLab1AForStudent)
		}

		studentLab1a := lab1a.Group("/variant")
		{
			securityLab1a := studentLab1a.Group("", h.CheckFirstAHeaderStudentForStart)
			{
				securityLab1a.POST("", h.UpdateUserVarianceLab1A)
			}

			notsecurityLab1a := studentLab1a.Group("", h.CheckFirstAHeaderStudent)
			{
				notsecurityLab1a.POST("/info", h.UpdateUserInfoLab1A)
				notsecurityLab1a.GET("/info", h.GetCurrentStepLab1A)
				notsecurityLab1a.POST("/result", h.SendUserResultLab1A)
			}
		}
	}

	lab1b := router.Group("/lab1b")
	{
		lecturerLab1b := lab1b.Group("/open", h.CheckFirstBHeaderLecturer)
		{
			lecturerLab1b.POST("", h.OpenLab1BForStudent)
		}

		studentLab1b := lab1b.Group("/variant")
		{
			securityLab1b := studentLab1b.Group("", h.CheckFirstBHeaderStudentForStart)
			{
				securityLab1b.GET("/ideal", h.GetLab1BVariance)
				securityLab1b.POST("", h.UpdateUserVarianceLab1B)
			}

			notSecurityLab1b := studentLab1b.Group("", h.CheckFirstBHeaderStudent)
			{
				notSecurityLab1b.POST("/info", h.UpdateUserInfoLab1B)
				notSecurityLab1b.GET("/info", h.GetCurrentStepLab1B)
				notSecurityLab1b.POST("/result", h.SendUserResultLab1B)
			}
		}
	}

	lab2 := router.Group("/lab2")
	{
		lecturerLab2 := lab2.Group("/open", h.CheckSecondHeaderLecturer)
		{
			lecturerLab2.PATCH("", h.OpenLab2ForStudent)
		}

		studentLab2 := lab2.Group("/variant")
		{
			securityLab2 := studentLab2.Group("", h.CheckSecondHeaderStudentForStart)
			{
				securityLab2.POST("", h.UpdateUserVarianceLab2)
			}

			notSecurityLab2 := studentLab2.Group("", h.CheckSecondHeaderStudent)
			{
				notSecurityLab2.POST("/info", h.UpdateUserInfoLab2)
				notSecurityLab2.GET("/info", h.GetCurrentStepLab2)
				notSecurityLab2.POST("/result", h.SendUserResultLab2)
			}
		}
	}

	return router
}
