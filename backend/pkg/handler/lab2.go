package handler

import (
	"backend/pkg/handler/errorResponse"
	"backend/pkg/model"
	"backend/pkg/service"
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (h *Handler) UpdateUserVarianceLab2(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, handlerTimeout)
	defer cancel()

	userId, err := getUserId(c)
	if err != nil {
		return
	}

	minutesDuration, err := strconv.Atoi(os.Getenv("FIRST_LAB_DURATION_MINUTES"))
	if err != nil {
		err = fmt.Errorf("ошибка получения продолжительности работы")
		errorResponse.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	var lab2 model.Variance2
	if err := c.BindJSON(&lab2); err != nil {
		err = fmt.Errorf("ошибка получения информации о лабораторной работе")
		errorResponse.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.Service.UpdateUserVariance(userId, service.Lab2Id, lab2); err != nil {
		err = fmt.Errorf("ошибка сохранения варианта")
		errorResponse.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	userInfo, err := h.Service.GetUserInfo(userId, service.Lab2Id)
	if err != nil {
		err = fmt.Errorf("ошибка получения информации о лабораторной работе")
		errorResponse.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"user_id": userId,
		"variant": lab2,
	})

	go func() {
		logrus.Println(fmt.Sprintf("START user:%d lab:%d", userId, service.Lab2Id))

		time.Sleep(time.Duration(minutesDuration) * time.Minute)

		if h.Service.IsEmptyToken(userId, service.Lab2Id) {
			return
		}

		if userInfo.IsDone {
			return
		}

		userMark, err := h.Service.GetLabResult(ctx, userId, service.Lab2Id)
		if err != nil {
			logrus.Errorf("ERROR get result user:%d lab:%d", userId, service.Lab2Id)
			return
		}

		if err := h.Service.SendLabMark(ctx, userId, userInfo.ExternalLabId, userMark); err != nil {
			logrus.Errorf("ERROR LAB3A send result user:%d lab:%d", userId, userInfo.ExternalLabId)
			return
		}

		if err := h.Service.ClearToken(userId, service.Lab2Id); err != nil {
			logrus.Errorf("ERROR clear token user:%d lab:%d", userId, service.Lab2Id)
			return
		}

		logrus.Println(fmt.Sprintf("SEND user:%d lab:%d percentage:%d", userId, service.Lab2Id, userMark))
	}()
}

func (h *Handler) GetCurrentStepLab2(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, handlerTimeout)
	defer cancel()

	userId, err := getUserId(c)
	if err != nil {
		return
	}

	step, err := h.Service.GetLabCurrentStep(ctx, userId, service.Lab2Id)
	if err != nil {
		err = fmt.Errorf("необходимо открыть лабораторную работу")
		errorResponse.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	mark, err := h.Service.GetCurrentMark(userId, service.Lab2Id)
	if err != nil {
		err = fmt.Errorf("ошибка получения текущей оценки")
		errorResponse.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	userDone, err := h.Service.GetUserVariance(ctx, userId, service.Lab2Id)
	if err != nil {
		err = fmt.Errorf("ошибка получения варианта")
		errorResponse.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"user_id":    userId,
		"step":       step,
		"variance":   userDone,
		"percentage": mark,
	})
}

func (h *Handler) UpdateUserInfoLab2(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, handlerTimeout)
	defer cancel()

	var data model.UserStepPercentage
	if err := c.BindJSON(&data); err != nil {
		err = fmt.Errorf("ошибка получения данных")
		errorResponse.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	userId, err := getUserId(c)
	if err != nil {
		return
	}

	if h.Service.CheckIsEmptyVariant(userId, service.Lab2Id) {
		err = fmt.Errorf("ошибка получения варианта лабораторной работы")
		errorResponse.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.Service.UpdateLabStep(ctx, userId, service.Lab2Id, data.Step); err != nil {
		err = fmt.Errorf("ошибка получения шага")
		errorResponse.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	if err := h.Service.IncrementPercentageDone(ctx, userId, service.Lab2Id, data.Percentage); err != nil {
		err = fmt.Errorf("ошибка получения шага")
		errorResponse.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	userDone, err := h.Service.GetUserVariance(ctx, userId, service.Lab2Id)
	if err != nil {
		err = fmt.Errorf("ошибка получения варианта")
		errorResponse.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"user_id":  userId,
		"variance": userDone,
	})
}

func (h *Handler) OpenLab2ForStudent(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, handlerTimeout)
	defer cancel()

	// Define a struct to parse the JSON body
	var body struct {
		UserId int  `json:"user_id"`
		IsOpen bool `json:"is_open"`
		LabId  int  `json:"lab_id"`
	}

	// Parse the JSON body
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if body.IsOpen {
		if _, err := h.Service.OpenLabForStudent(ctx, body.UserId, service.Lab2Id, body.LabId); err != nil {
			err = fmt.Errorf("ошибка открытия лабораторной работы")
			errorResponse.NewErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}
	} else {
		if err := h.Service.CloseLabForStudent(ctx, body.UserId, service.Lab2Id); err != nil {
			err = fmt.Errorf("ошибка закрытия лабораторной работы")
			errorResponse.NewErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{})
}

func (h *Handler) SendUserResultLab2(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, handlerTimeout)
	defer cancel()

	var data model.SendUserResult
	if err := c.BindJSON(&data); err != nil {
		err = fmt.Errorf("ошибка получения данных")
		errorResponse.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	userId, err := getUserId(c)
	if err != nil {
		return
	}

	if h.Service.CheckIsEmptyVariant(userId, service.Lab2Id) {
		err = fmt.Errorf("ошибка получения варианта лабораторной работы")
		errorResponse.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	userInfo, err := h.Service.GetUserInfo(userId, service.Lab2Id)
	if err != nil {
		err = fmt.Errorf("ошибка получения информации о лабораторной работе")
		errorResponse.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	// if h.Service.IsEmptyToken(userId, service.Lab2Id) {
	// 	err = fmt.Errorf("ошибка получения информации о лабораторной работе: пустой токен")
	// 	errorResponse.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
	// 	return
	// }

	userMark, err := h.Service.GetLabResult(ctx, userId, service.Lab2Id)
	if err != nil {
		err = fmt.Errorf("ошибка получения результатов")
		errorResponse.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	if err := h.Service.SendLabMark(ctx, userId, userInfo.ExternalLabId, userMark); err != nil {
		err = fmt.Errorf("ошибка получения оценки")
		errorResponse.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		logrus.Errorf("ERROR LAB1A send result user:%d lab:%d", userId, userInfo.ExternalLabId)
		return
	}

	if err := h.Service.ClearToken(userId, service.Lab2Id); err != nil {
		logrus.Errorf("ERROR clear token user:%d lab:%d", userId, service.Lab2Id)
		err = fmt.Errorf("внутренняя ошибка")
		errorResponse.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	logrus.Println(fmt.Sprintf("SEND user:%d lab:%d percentage:%d", userId, service.Lab2Id, userMark))

	c.JSON(http.StatusOK, gin.H{})
}
