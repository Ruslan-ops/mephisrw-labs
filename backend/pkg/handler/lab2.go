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

func (h *Handler) OpenSecondLab(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, handlerTimeout)
	defer cancel()

	minutesDuration, err := strconv.Atoi(os.Getenv("THIRD_LAB_DURATION_MINUTES"))
	if err != nil {
		err = fmt.Errorf("ошибка получения продолжительности работы")
		errorResponse.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	userId, err := getUserId(c)
	if err != nil {
		return
	}

	userInfo, err := h.Service.GetUserInfo(userId, service.Lab2Id)
	if err != nil {
		err = fmt.Errorf("ошибка получения информации о лаблораторной работе")
		errorResponse.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	userDone, err := h.Service.GetVariance2(ctx, userId)
	if err != nil {
		err = fmt.Errorf("ошибка получения варианта")
		errorResponse.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"user_id": userId,
		"variant": userDone,
	})

	go func() {
		logrus.Println(fmt.Sprintf("START user:%d lab:%d", userId, service.Lab2Id))

		time.Sleep(time.Duration(minutesDuration) * time.Minute)

		if h.Service.IsEmptyToken(userId, service.Lab2Id) {
			return
		}

		userMark, err := h.Service.GetLabResult(ctx, userId, service.Lab2Id)
		if err != nil {
			logrus.Errorf("ERROR get result user:%d lab:%d", userId, service.Lab2Id)
			return
		}

		if err := h.Service.SendLabMark(c, userId, userInfo.ExternalLabId, userMark); err != nil {
			logrus.Errorf("ERROR LAB3C send result user:%d lab:%d", userId, userInfo.ExternalLabId)
			return
		}

		if err := h.Service.ClearToken(userId, service.Lab2Id); err != nil {
			logrus.Errorf("ERROR clear token user:%d lab:%d", userId, service.Lab2Id)
			return
		}
		logrus.Println(fmt.Sprintf("SEND user:%d lab:%d percentage:%d", userId, service.Lab2Id, userMark))
	}()
}

func (h *Handler) OpenSecondLabForStudent(c *gin.Context) {
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
		if _, err := h.Service.OpenLabForStudent(ctx, userId, service.Lab2Id, externalLabId); err != nil {
			err = fmt.Errorf("ошибка открытия лабораторной работы")
			errorResponse.NewErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}

		variance, data := h.Service.GenerateUserVariance2(ctx)
		if err := h.Service.UpdateUserVariance2(ctx, userId, model.Variance2{
			Number: variance,
			Data:   data,
		}); err != nil {
			return
		}
	} else {
		if err := h.Service.CloseLabForStudent(ctx, userId, service.Lab2Id); err != nil {
			err = fmt.Errorf("ошибка открытия лабораторной работы")
			errorResponse.NewErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{})
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

	c.JSON(http.StatusOK, map[string]interface{}{
		"user_id":    userId,
		"step":       step,
		"percentage": mark,
	})
}
