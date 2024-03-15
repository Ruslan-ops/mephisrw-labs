package handler

import (
	"backend/pkg/handler/errorResponse"
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

func (h *Handler) OpenFirstBLab(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, handlerTimeout)
	defer cancel()

	minutesDuration, err := strconv.Atoi(os.Getenv("SECOND_LAB_DURATION_MINUTES"))
	if err != nil {
		err = fmt.Errorf("ошибка получения продолжительности работы")
		errorResponse.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	userId, err := getUserId(c)
	if err != nil {
		return
	}

	userInfo, err := h.Service.GetUserInfo(userId, service.Lab1BId)
	if err != nil {
		err = fmt.Errorf("ошибка получения информации о лаблораторной работе")
		errorResponse.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	userDone, matrix, err := h.Service.GetVariance(ctx, userId, service.Lab1BId)
	if err != nil {
		err = fmt.Errorf("ошибка получения варианта")
		errorResponse.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"user_id": userId,
		"variant": userDone,
		"matrix":  matrix,
	})

	go func() {
		logrus.Println(fmt.Sprintf("START user:%d lab:%d", userId, service.Lab1BId))

		time.Sleep(time.Duration(minutesDuration) * time.Minute)

		if h.Service.IsEmptyToken(userId, service.Lab1BId) {
			return
		}

		userMark, err := h.Service.GetLabResult(ctx, userId, service.Lab1BId)
		if err != nil {
			logrus.Errorf("ERROR get result user:%d lab:%d", userId, service.Lab1BId)
			return
		}

		if err := h.Service.SendLabMark(ctx, userId, userInfo.ExternalLabId, userMark); err != nil {
			logrus.Errorf("ERROR LAB3B send result user:%d lab:%d", userId, userInfo.ExternalLabId)
			return
		}

		if err := h.Service.ClearToken(userId, service.Lab1BId); err != nil {
			logrus.Errorf("ERROR clear token user:%d lab:%d", userId, service.Lab1BId)
			return
		}
		logrus.Println(fmt.Sprintf("SEND user:%d lab:%d percentage:%d", userId, service.Lab1BId, userMark))
	}()
}

func (h *Handler) OpenFirstBLabForStudent(c *gin.Context) {
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
		isDone, err := h.Service.OpenLabForStudent(ctx, userId, service.Lab1BId, externalLabId)
		if err != nil {
			err = fmt.Errorf("ошибка открытия лабораторной работы")
			errorResponse.NewErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}
		if isDone {
			go func() {
				routineCtx, routineCancel := context.WithTimeout(context.Background(), 30*time.Second)
				defer routineCancel()
				flag := true
				for flag {
					select {
					case <-routineCtx.Done():
						return
					default:
						variance, err := h.Service.GenerateTask(ctx, userId)
						if err != nil {
							continue
						}
						setLabA, err := h.Service.ValidateLab3AResult(ctx, variance)
						if err != nil {
							continue
						}
						setLabB, err := h.Service.ValidateLab3BResult(ctx, variance)
						if err != nil {
							continue
						}
						setLabC, err := h.Service.ValidateLab3CResult(ctx, variance)
						if err != nil {
							continue
						}

						var firstMaxA, secondMaxA float64 = 0, 0
						for _, v := range setLabA {
							if firstMaxA <= v {
								secondMaxA = firstMaxA
								firstMaxA = v
							} else if secondMaxA <= v {
								secondMaxA = v
							}
						}

						var firstMaxB, secondMaxB float64 = 0, 0
						for _, v := range setLabB {
							if firstMaxB <= v {
								secondMaxB = firstMaxB
								firstMaxB = v
							} else if secondMaxB <= v {
								secondMaxB = v
							}
						}

						var firstMaxC, secondMaxC float64 = 0, 0
						for _, v := range setLabC {
							if firstMaxC <= v {
								secondMaxC = firstMaxC
								firstMaxC = v
							} else if secondMaxC <= v {
								secondMaxC = v
							}
						}

						if secondMaxA <= 0 || secondMaxB <= 0 || secondMaxC <= 0 {
							continue
						}

						if firstMaxA > secondMaxA+0.03 && firstMaxB > secondMaxB+0.03 && firstMaxC > secondMaxC+0.03 {
							if err := h.Service.UpdateUserVariance(ctx, userId, service.Lab1BId, variance); err != nil {
								return
							}
							return
						}
					}
				}
			}()
		}
	} else {
		if err := h.Service.CloseLabForStudent(ctx, userId, service.Lab1BId); err != nil {
			err = fmt.Errorf("ошибка открытия лабораторной работы")
			errorResponse.NewErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{})
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

	c.JSON(http.StatusOK, map[string]interface{}{
		"user_id":    userId,
		"step":       step,
		"percentage": mark,
	})
}
