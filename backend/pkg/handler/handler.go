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

	lab1a := router.Group("/lab3a")
	{
		lecturerLab3a := lab1a.Group("/open", h.CheckFirstAHeaderLecturer)
		{
			lecturerLab3a.POST("", h.OpenFirstALabForStudent)
		}

		studentLab1a := lab1a.Group("/variant")
		{
			securityLab3a := studentLab1a.Group("", h.CheckFirstAHeaderStudentForStart)
			{
				securityLab3a.GET("", h.OpenFirstALab)
			}

			notsecurityLab1a := studentLab1a.Group("", h.CheckFirstAHeaderStudent)
			{
			}
		}
	}

	lab1b := router.Group("/lab3b")
	{
		lecturerLab1b := lab1b.Group("/open", h.CheckFirstBHeaderLecturer)
		{
			lecturerLab1b.POST("", h.OpenFirstBLabForStudent)
		}

		studentLab1b := lab1b.Group("/variant")
		{
			securityLab1b := studentLab1b.Group("", h.CheckFirstBHeaderStudentForStart)
			{
				securityLab1b.GET("", h.OpenFirstBLab)
			}

			notsecurityLab1b := studentLab1b.Group("", h.CheckFirstBHeaderStudent)
			{
			}
		}
	}

	lab2 := router.Group("/lab3c")
	{
		lecturer := lab2.Group("/open", h.CheckSecondHeaderLecturer)
		{
			lecturer.POST("", h.OpenSecondLabForStudent)
		}

		student2 := lab2.Group("/variant")
		{
			security := student2.Group("", h.CheckSecondHeaderStudentForStart)
			{
				security.GET("", h.OpenSecondLab)
			}

			notsecurityLab2 := student2.Group("", h.CheckSecondHeaderStudent)
			{
			}
		}
	}

	return router
}
