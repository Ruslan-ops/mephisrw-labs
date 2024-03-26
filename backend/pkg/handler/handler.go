package handler

import (
	"backend/pkg/service"
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
				securityLab1b.POST("", h.UpdateUserVarianceLab1B)
			}

			notSecurityLab1b := studentLab1b.Group("", h.CheckFirstBHeaderStudent)
			{
				notSecurityLab1b.POST("/info", h.UpdateUserInfoLab1B)
				notSecurityLab1b.GET("/info", h.GetCurrentStepLab1B)
			}
		}
	}

	lab2 := router.Group("/lab2")
	{
		lecturerLab2 := lab2.Group("/open", h.CheckSecondHeaderLecturer)
		{
			lecturerLab2.POST("", h.OpenLab2ForStudent)
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
			}
		}
	}

	return router
}
