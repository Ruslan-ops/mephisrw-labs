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

func (h *Handler) OpenFirstALab(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, handlerTimeout)
	defer cancel()

	minutesDuration, err := strconv.Atoi(os.Getenv("FIRST_LAB_DURATION_MINUTES"))
	if err != nil {
		err = fmt.Errorf("ошибка получения продолжительности работы")
		errorResponse.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	userId, err := getUserId(c)
	if err != nil {
		return
	}

	userDone, err := h.Service.GetVariance1A(ctx, userId)
	if err != nil {
		err = fmt.Errorf("ошибка получения варианта")
		errorResponse.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	userInfo, err := h.Service.GetUserInfo(userId, service.Lab1AId)
	if err != nil {
		err = fmt.Errorf("ошибка получения информации о лаблораторной работе")
		errorResponse.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"user_id": userId,
		"variant": userDone,
	})

	go func() {
		logrus.Println(fmt.Sprintf("START user:%d lab:%d", userId, service.Lab1AId))

		time.Sleep(time.Duration(minutesDuration) * time.Minute)

		if h.Service.IsEmptyToken(userId, service.Lab1AId) {
			return
		}

		userMark, err := h.Service.GetLabResult(ctx, userId, service.Lab1AId)
		if err != nil {
			logrus.Errorf("ERROR get result user:%d lab:%d", userId, service.Lab1AId)
			return
		}

		if err := h.Service.SendLabMark(ctx, userId, userInfo.ExternalLabId, userMark); err != nil {
			logrus.Errorf("ERROR LAB3A send result user:%d lab:%d", userId, userInfo.ExternalLabId)
			return
		}

		if err := h.Service.ClearToken(userId, service.Lab1AId); err != nil {
			logrus.Errorf("ERROR clear token user:%d lab:%d", userId, service.Lab1AId)
			return
		}

		logrus.Println(fmt.Sprintf("SEND user:%d lab:%d percentage:%d", userId, service.Lab1AId, userMark))
	}()
}

func (h *Handler) OpenFirstALabForStudent(c *gin.Context) {
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
		if _, err := h.Service.OpenLabForStudent(ctx, userId, service.Lab1AId, externalLabId); err != nil {
			err = fmt.Errorf("ошибка открытия лабораторной работы")
			errorResponse.NewErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}
		variance, data := h.Service.GenerateUserVariance1A(ctx)
		if err := h.Service.UpdateUserVariance1A(ctx, userId, model.Variance1A{
			Number: variance,
			Data:   data,
		}); err != nil {
			return
		}
	} else {
		if err := h.Service.CloseLabForStudent(ctx, userId, service.Lab1AId); err != nil {
			err = fmt.Errorf("ошибка закрытия лабораторной работы")
			errorResponse.NewErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{})
}

func (h *Handler) GetCurrentStepLab1A(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, handlerTimeout)
	defer cancel()

	userId, err := getUserId(c)
	if err != nil {
		return
	}

	step, err := h.Service.GetLabCurrentStep(ctx, userId, service.Lab1AId)
	if err != nil {
		err = fmt.Errorf("необходимо открыть лабораторную работу")
		errorResponse.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	mark, err := h.Service.GetCurrentMark(userId, service.Lab1AId)
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

func (h *Handler) Send1AImportanceMatrix(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, handlerTimeout)
	defer cancel()

	var userRes model.AnswerLab1AImportanceMatrix
	if err := c.BindJSON(&userRes); err != nil {
		err = fmt.Errorf("ошибка отправки ответа")
		errorResponse.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	userId, err := getUserId(c)
	if err != nil {
		return
	}

	if step, err := h.Service.GetLabCurrentStep(ctx, userId, service.Lab1AId); err != nil {
		err = fmt.Errorf("необходимо открыть лабораторную работу")
		errorResponse.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	} else if step != 0 {
		err = fmt.Errorf("необходимо проходить работу пошагово")
		errorResponse.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	mark, res, consistencyRatio, err := h.Service.CheckLab1AImportanceMatrix(ctx, userId, userRes.Matrix)
	if err != nil {
		err = fmt.Errorf("ошибка со стороны сервера")
		errorResponse.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	go func() {
		if err := h.Service.IncrementPercentageDone(context.Background(), userId, service.Lab1AId, mark); err != nil {
			logrus.Errorf("can't change percentage done user_id:%d labId:%d: %v", userId, service.Lab1AId, err)
			return
		}
		if err := h.Service.UpdateLabStep(ctx, userId, service.Lab1AId, 1); err != nil {
			logrus.Errorf("can't change lab step user_id:%d labId:%d: %v", userId, service.Lab1AId, err)
			return
		}
	}()

	c.JSON(http.StatusOK, map[string]interface{}{
		"percentage":  mark,
		"result":      res,
		"consistency": consistencyRatio,
	})
}

func (h *Handler) Send1AImportanceMatrixFirstCriteria(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, handlerTimeout)
	defer cancel()

	var userRes model.AnswerLab1AImportanceMatrixFirstCriteria
	if err := c.BindJSON(&userRes); err != nil {
		err = fmt.Errorf("ошибка отправки ответа")
		errorResponse.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	userId, err := getUserId(c)
	if err != nil {
		return
	}

	if step, err := h.Service.GetLabCurrentStep(ctx, userId, service.Lab1AId); err != nil {
		err = fmt.Errorf("необходимо открыть лабораторную работу")
		errorResponse.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	} else if step != 1 {
		err = fmt.Errorf("необходимо проходить работу пошагово")
		errorResponse.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	mark, res, consistencyRatio, err := h.Service.CheckLab1AImportanceMatrixFirstCriteria(ctx, userId, userRes.Matrix)
	if err != nil {
		err = fmt.Errorf("ошибка со стороны сервера")
		errorResponse.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	go func() {
		if err := h.Service.IncrementPercentageDone(context.Background(), userId, service.Lab1AId, mark); err != nil {
			logrus.Errorf("can't change percentage done user_id:%d labId:%d: %v", userId, service.Lab1AId, err)
			return
		}
		if err := h.Service.UpdateLabStep(ctx, userId, service.Lab1AId, 2); err != nil {
			logrus.Errorf("can't change lab step user_id:%d labId:%d: %v", userId, service.Lab1AId, err)
			return
		}
	}()

	c.JSON(http.StatusOK, map[string]interface{}{
		"percentage":  mark,
		"result":      res,
		"consistency": consistencyRatio,
	})
}

func (h *Handler) Send1AImportanceMatrixSecondCriteria(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, handlerTimeout)
	defer cancel()

	var userRes model.AnswerLab1AImportanceMatrixSecondCriteria
	if err := c.BindJSON(&userRes); err != nil {
		err = fmt.Errorf("ошибка отправки ответа")
		errorResponse.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	userId, err := getUserId(c)
	if err != nil {
		return
	}

	if step, err := h.Service.GetLabCurrentStep(ctx, userId, service.Lab1AId); err != nil {
		err = fmt.Errorf("необходимо открыть лабораторную работу")
		errorResponse.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	} else if step != 2 {
		err = fmt.Errorf("необходимо проходить работу пошагово")
		errorResponse.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	mark, res, consistencyRatio, err := h.Service.CheckLab1AImportanceMatrixSecondCriteria(ctx, userId, userRes.Matrix)
	if err != nil {
		err = fmt.Errorf("ошибка со стороны сервера")
		errorResponse.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	go func() {
		if err := h.Service.IncrementPercentageDone(context.Background(), userId, service.Lab1AId, mark); err != nil {
			logrus.Errorf("can't change percentage done user_id:%d labId:%d: %v", userId, service.Lab1AId, err)
			return
		}
		if err := h.Service.UpdateLabStep(ctx, userId, service.Lab1AId, 3); err != nil {
			logrus.Errorf("can't change lab step user_id:%d labId:%d: %v", userId, service.Lab1AId, err)
			return
		}
	}()

	c.JSON(http.StatusOK, map[string]interface{}{
		"percentage":  mark,
		"result":      res,
		"consistency": consistencyRatio,
	})
}

func (h *Handler) Send1AImportanceMatrixThirdCriteria(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, handlerTimeout)
	defer cancel()

	var userRes model.AnswerLab1AImportanceMatrixSecondCriteria
	if err := c.BindJSON(&userRes); err != nil {
		err = fmt.Errorf("ошибка отправки ответа")
		errorResponse.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	userId, err := getUserId(c)
	if err != nil {
		return
	}

	if step, err := h.Service.GetLabCurrentStep(ctx, userId, service.Lab1AId); err != nil {
		err = fmt.Errorf("необходимо открыть лабораторную работу")
		errorResponse.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	} else if step != 3 {
		err = fmt.Errorf("необходимо проходить работу пошагово")
		errorResponse.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	mark, res, consistencyRatio, err := h.Service.CheckLab1AImportanceMatrixThirdCriteria(ctx, userId, userRes.Matrix)
	if err != nil {
		err = fmt.Errorf("ошибка со стороны сервера")
		errorResponse.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	go func() {
		if err := h.Service.IncrementPercentageDone(context.Background(), userId, service.Lab1AId, mark); err != nil {
			logrus.Errorf("can't change percentage done user_id:%d labId:%d: %v", userId, service.Lab1AId, err)
			return
		}
		if err := h.Service.UpdateLabStep(ctx, userId, service.Lab1AId, 4); err != nil {
			logrus.Errorf("can't change lab step user_id:%d labId:%d: %v", userId, service.Lab1AId, err)
			return
		}
	}()

	c.JSON(http.StatusOK, map[string]interface{}{
		"percentage":  mark,
		"result":      res,
		"consistency": consistencyRatio,
	})
}

func (h *Handler) Send1AImportanceMatrixFourthCriteria(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, handlerTimeout)
	defer cancel()

	var userRes model.AnswerLab1AImportanceMatrixFourthCriteria
	if err := c.BindJSON(&userRes); err != nil {
		err = fmt.Errorf("ошибка отправки ответа")
		errorResponse.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	userId, err := getUserId(c)
	if err != nil {
		return
	}

	if step, err := h.Service.GetLabCurrentStep(ctx, userId, service.Lab1AId); err != nil {
		err = fmt.Errorf("необходимо открыть лабораторную работу")
		errorResponse.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	} else if step != 4 {
		err = fmt.Errorf("необходимо проходить работу пошагово")
		errorResponse.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	mark, res, consistencyRatio, err := h.Service.CheckLab1AImportanceMatrixFourthCriteria(ctx, userId, userRes.Matrix)
	if err != nil {
		err = fmt.Errorf("ошибка со стороны сервера")
		errorResponse.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	go func() {
		if err := h.Service.IncrementPercentageDone(context.Background(), userId, service.Lab1AId, mark); err != nil {
			logrus.Errorf("can't change percentage done user_id:%d labId:%d: %v", userId, service.Lab1AId, err)
			return
		}
		if err := h.Service.UpdateLabStep(ctx, userId, service.Lab1AId, 5); err != nil {
			logrus.Errorf("can't change lab step user_id:%d labId:%d: %v", userId, service.Lab1AId, err)
			return
		}
	}()

	c.JSON(http.StatusOK, map[string]interface{}{
		"percentage":  mark,
		"result":      res,
		"consistency": consistencyRatio,
	})
}

func (h *Handler) Send1AChosenAlternative(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, handlerTimeout)
	defer cancel()

	var userRes model.AnswerLab1AChosenAlternative
	if err := c.BindJSON(&userRes); err != nil {
		err = fmt.Errorf("ошибка отправки ответа")
		errorResponse.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	userId, err := getUserId(c)
	if err != nil {
		return
	}

	if step, err := h.Service.GetLabCurrentStep(ctx, userId, service.Lab1AId); err != nil {
		err = fmt.Errorf("необходимо открыть лабораторную работу")
		errorResponse.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	} else if step != 5 {
		err = fmt.Errorf("необходимо проходить работу пошагово")
		errorResponse.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	mark, res, consistencyRatio, err := h.Service.CheckLab1AChosenAlternative(ctx, userId, userRes.Matrix)
	if err != nil {
		err = fmt.Errorf("ошибка со стороны сервера")
		errorResponse.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	go func() {
		if err := h.Service.IncrementPercentageDone(context.Background(), userId, service.Lab1AId, mark); err != nil {
			logrus.Errorf("can't change percentage done user_id:%d labId:%d: %v", userId, service.Lab1AId, err)
			return
		}
		if err := h.Service.UpdateLabStep(ctx, userId, service.Lab1AId, 6); err != nil {
			logrus.Errorf("can't change lab step user_id:%d labId:%d: %v", userId, service.Lab1AId, err)
			return
		}
	}()

	c.JSON(http.StatusOK, map[string]interface{}{
		"percentage":  mark,
		"result":      res,
		"consistency": consistencyRatio,
	})
}
