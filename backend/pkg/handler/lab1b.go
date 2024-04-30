package handler

import (
	"backend/pkg/handler/errorResponse"
	"backend/pkg/model"
	"backend/pkg/service"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"strconv"
	"time"
)

func (h *Handler) UpdateUserVarianceLab1B(c *gin.Context) {
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

	var lab1b model.Variance1B
	if err := c.BindJSON(&lab1b); err != nil {
		err = fmt.Errorf("ошибка получения информации о лабораторной работе")
		errorResponse.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.Service.UpdateUserVariance(userId, service.Lab1BId, lab1b); err != nil {
		err = fmt.Errorf("ошибка сохранения варианта")
		errorResponse.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	userInfo, err := h.Service.GetUserInfo(userId, service.Lab1BId)
	if err != nil {
		err = fmt.Errorf("ошибка получения информации о лабораторной работе")
		errorResponse.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"user_id": userId,
		"variant": lab1b,
	})

	go func() {
		logrus.Println(fmt.Sprintf("START user:%d lab:%d", userId, service.Lab1BId))

		time.Sleep(time.Duration(minutesDuration) * time.Minute)

		if h.Service.IsEmptyToken(userId, service.Lab1BId) {
			return
		}

		if userInfo.IsDone {
			return
		}

		userMark, err := h.Service.GetLabResult(ctx, userId, service.Lab1BId)
		if err != nil {
			logrus.Errorf("ERROR get result user:%d lab:%d", userId, service.Lab1BId)
			return
		}

		if err := h.Service.SendLabMark(ctx, userId, userInfo.ExternalLabId, userMark); err != nil {
			logrus.Errorf("ERROR LAB3A send result user:%d lab:%d", userId, userInfo.ExternalLabId)
			return
		}

		if err := h.Service.ClearToken(userId, service.Lab1BId); err != nil {
			logrus.Errorf("ERROR clear token user:%d lab:%d", userId, service.Lab1BId)
			return
		}

		logrus.Println(fmt.Sprintf("SEND user:%d lab:%d percentage:%d", userId, service.Lab1BId, userMark))
	}()
}

func (h *Handler) GetLab1BVariance(c *gin.Context) {
	_, cancel := context.WithTimeout(c, handlerTimeout)
	defer cancel()

	userId, err := getUserId(c)
	if err != nil {
		return
	}
	logrus.Println(fmt.Sprintf("GET VARIANCE user:%d", userId))

	variance, err := h.Service.GetIdealVariant1B()
	if err != nil {
		errorResponse.NewErrorResponse(c, http.StatusInternalServerError, fmt.Errorf("can't get variant").Error())
		return
	}

	c.JSON(http.StatusOK, variance)
}

func (h *Handler) GetCurrentStepLab1B(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, handlerTimeout)
	defer cancel()

	userId, err := getUserId(c)
	if err != nil {
		return
	}

	step, err := h.Service.GetLabCurrentStep(ctx, userId, service.Lab1BId)
	if err != nil {
		err = fmt.Errorf("необходимо открыть лабораторную работу")
		errorResponse.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	mark, err := h.Service.GetCurrentMark(userId, service.Lab1BId)
	if err != nil {
		err = fmt.Errorf("ошибка получения текущей оценки")
		errorResponse.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	userDone, err := h.Service.GetUserVariance(ctx, userId, service.Lab1BId)
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

func (h *Handler) UpdateUserInfoLab1B(c *gin.Context) {
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

	if err := h.Service.UpdateLabStep(ctx, userId, service.Lab1BId, data.Step); err != nil {
		err = fmt.Errorf("ошибка получения шага")
		errorResponse.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	if err := h.Service.IncrementPercentageDone(ctx, userId, service.Lab1BId, data.Percentage); err != nil {
		err = fmt.Errorf("ошибка получения шага")
		errorResponse.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	userDone, err := h.Service.GetUserVariance(ctx, userId, service.Lab1BId)
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

func (h *Handler) OpenLab1BForStudent(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, handlerTimeout)
	defer cancel()

	user := c.Query("user_id")
	isOpen := c.Query("is_open")
	externalLab := c.Query("lab_id")

	userId, err := strconv.Atoi(user)
	if err != nil {
		err = fmt.Errorf("ошибка получения студента")
		errorResponse.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	externalLabId, err := strconv.Atoi(externalLab)
	if err != nil {
		err = fmt.Errorf("ошибка получения лабораторной работы")
		errorResponse.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	isOpenBool, err := strconv.ParseBool(isOpen)
	if err != nil {
		err = fmt.Errorf("ошибка получения изменения проведения лабораторной работы")
		errorResponse.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if isOpenBool {
		if _, err := h.Service.OpenLabForStudent(ctx, userId, service.Lab1BId, externalLabId); err != nil {
			err = fmt.Errorf("ошибка открытия лабораторной работы")
			errorResponse.NewErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}
	} else {
		if err := h.Service.CloseLabForStudent(ctx, userId, service.Lab1BId); err != nil {
			err = fmt.Errorf("ошибка закрытия лабораторной работы")
			errorResponse.NewErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{})
}

func (h *Handler) SendUserResultLab1B(c *gin.Context) {
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

	userInfo, err := h.Service.GetUserInfo(userId, service.Lab1BId)
	if err != nil {
		err = fmt.Errorf("ошибка получения информации о лабораторной работе")
		errorResponse.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	if h.Service.IsEmptyToken(userId, service.Lab1BId) {
		err = fmt.Errorf("ошибка получения информации о лабораторной работе: пустой токен")
		errorResponse.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	userMark, err := h.Service.GetLabResult(ctx, userId, service.Lab1BId)
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

	if err := h.Service.ClearToken(userId, service.Lab1BId); err != nil {
		logrus.Errorf("ERROR clear token user:%d lab:%d", userId, service.Lab1BId)
		err = fmt.Errorf("внутренняя ошибка")
		errorResponse.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	logrus.Println(fmt.Sprintf("SEND user:%d lab:%d percentage:%d", userId, service.Lab1BId, userMark))

	c.JSON(http.StatusOK, gin.H{})
}
