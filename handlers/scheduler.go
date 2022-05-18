package handlers

import (
	"fmt"
	"net/http"

	"github.com/ecoprohcm/DMS_BackendServer/models"
	"github.com/ecoprohcm/DMS_BackendServer/mqttSvc"
	"github.com/ecoprohcm/DMS_BackendServer/utils"
	"github.com/gin-gonic/gin"
)

type SchedulerHandler struct {
	deps *HandlerDependencies
}

func NewSchedulerHandler(deps *HandlerDependencies) *SchedulerHandler {
	return &SchedulerHandler{
		deps,
	}
}

// Find all scheduler info
// @Summary Find All Scheduler
// @Schemes
// @Description find all scheduler info
// @Produce json
// @Success 200 {array} []models.Scheduler
// @Failure 400 {object} utils.ErrorResponse
// @Router /v1/schedulers [get]
func (h *SchedulerHandler) FindAllScheduler(c *gin.Context) {
	sList, err := h.deps.SvcOpts.SchedulerSvc.FindAllScheduler(c)
	if err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Get all scheduler failed",
			ErrorMsg:   err.Error(),
		})
		return
	}
	utils.ResponseJson(c, http.StatusOK, sList)
}

// Find scheduler info by id
// @Summary Find Scheduler By ID
// @Schemes
// @Description find scheduler info by scheduler id
// @Produce json
// @Param        id	path	string	true	"Scheduler ID"
// @Success 200 {object} models.Scheduler
// @Failure 400 {object} utils.ErrorResponse
// @Router /v1/scheduler/{id} [get]
func (h *SchedulerHandler) FindSchedulerByID(c *gin.Context) {
	id := c.Param("id")

	s, err := h.deps.SvcOpts.SchedulerSvc.FindSchedulerByID(c, id)
	if err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Get scheduler failed",
			ErrorMsg:   err.Error(),
		})
		return
	}
	utils.ResponseJson(c, http.StatusOK, s)
}

// Create scheduler
// @Summary Create Scheduler
// @Schemes
// @Description Create scheduler
// @Accept  json
// @Produce json
// @Param	data	body	models.SwagCreateScheduler	true	"Fields need to create a scheduler"
// @Success 200 {object} models.Scheduler
// @Failure 400 {object} utils.ErrorResponse
// @Router /v1/scheduler [post]
func (h *SchedulerHandler) CreateScheduler(c *gin.Context) {
	s := &models.Scheduler{}
	err := c.ShouldBind(s)
	if err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Invalid req body",
			ErrorMsg:   err.Error(),
		})
		return
	}

	_, err = h.deps.SvcOpts.SchedulerSvc.CreateScheduler(c.Request.Context(), s)
	if err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Create scheduler failed",
			ErrorMsg:   err.Error(),
		})
		return
	}
	utils.ResponseJson(c, http.StatusOK, s)
}

// Update scheduler
// @Summary Update Scheduler By ID
// @Schemes
// @Description Update scheduler, must have "id" field. Send updated info to MQTT broker
// @Accept  json
// @Produce json
// @Param	data	body	models.UpdateScheduler	true	"Fields need to update a scheduler"
// @Success 200 {boolean} true
// @Failure 400 {object} utils.ErrorResponse
// @Router /v1/scheduler [patch]
func (h *SchedulerHandler) UpdateScheduler(c *gin.Context) {
	s := &models.UpdateScheduler{}
	err := c.ShouldBind(s)
	if err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Invalid req body",
			ErrorMsg:   err.Error(),
		})
		return
	}

	isSuccess, err := h.deps.SvcOpts.SchedulerSvc.UpdateScheduler(c.Request.Context(), &s.Scheduler)
	if err != nil || !isSuccess {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Update scheduler failed",
			ErrorMsg:   err.Error(),
		})
		return
	}

	t := h.deps.MqttClient.Publish(mqttSvc.TOPIC_SV_SCHEDULER_U, 1, false,
		mqttSvc.ServerUpdateRegisterPayload("0", s))
	if err := mqttSvc.HandleMqttErr(t); err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Update scheduler mqtt failed",
			ErrorMsg:   err.Error(),
		})
		return
	}

	utils.ResponseJson(c, http.StatusOK, isSuccess)
}

// Delete scheduler
// @Summary Delete Scheduler By ID
// @Schemes
// @Description Delete scheduler using "id" field. Send deleted info to MQTT broker
// @Accept  json
// @Produce json
// @Param	data	body	object{id=int}	true	"Scheduler ID"
// @Success 200 {boolean} true
// @Failure 400 {object} utils.ErrorResponse
// @Router /v1/scheduler [delete]
func (h *SchedulerHandler) DeleteScheduler(c *gin.Context) {
	dId := &models.DeleteID{}
	err := c.ShouldBind(dId)
	if err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Invalid req body",
			ErrorMsg:   err.Error(),
		})
		return
	}

	t := h.deps.MqttClient.Publish(mqttSvc.TOPIC_SV_SCHEDULER_D, 1, false,
		mqttSvc.ServerDeleteRegisterPayload("0", dId.ID))
	if err := mqttSvc.HandleMqttErr(t); err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Delete scheduler mqtt failed",
			ErrorMsg:   err.Error(),
		})
		return
	}

	isSuccess, err := h.deps.SvcOpts.SchedulerSvc.DeleteScheduler(c.Request.Context(), dId.ID)
	if err != nil || !isSuccess {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Delete scheduler failed",
			ErrorMsg:   err.Error(),
		})
		return
	}
	utils.ResponseJson(c, http.StatusOK, isSuccess)

}

// Append Scheduler base on Excel
// @Summary Append Scheduler Base On Excel
// @Schemes
// @Description Append Scheduler base on Excel
// @Accept  json
// @Produce json
// @Param	data	body	models.SwagCreateScheduler	true	"Fields need to append a scheduler"
// @Success 200 {object} models.Scheduler
// @Failure 400 {object} utils.ErrorResponse
// @Router /v1//scheduler/excel [post]
func (h *SchedulerHandler) AppendSchedulerOnExcel(c *gin.Context) {
	userScheduler := &models.UserScheduler{}
	err := c.ShouldBind(&userScheduler.ScheInfo)
	if err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Invalid req body",
			ErrorMsg:   err.Error(),
		})
		return
	}

	roomId := userScheduler.ScheInfo.RoomID
	dlList, err := h.deps.SvcOpts.DoorlockSvc.FindAllDoorlocksByRoomID(c, roomId)

	if err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Get doorlocks by room ID failed",
			ErrorMsg:   err.Error(),
		})
		return
	}
	userRole := userScheduler.ScheInfo.Role
	userScheduler, err = getUserInformation(c, h.deps.SvcOpts, userRole, userScheduler)
	if err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Get user failed",
			ErrorMsg:   err.Error(),
		})
		return
	}
	for i := 0; i < len(dlList); i++ {
		newScheduler := userScheduler.ScheInfo
		newScheduler.DoorID = uint(dlList[i].ID)
		_, err = h.deps.SvcOpts.SchedulerSvc.CreateScheduler(c.Request.Context(), &newScheduler)
		if err != nil {
			utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
				StatusCode: http.StatusBadRequest,
				Msg:        "Create scheduler failed",
				ErrorMsg:   err.Error(),
			})
			return
		}

		t := h.deps.MqttClient.Publish(mqttSvc.TOPIC_SV_SCHEDULER_C, 1, false, mqttSvc.ServerCreateRegisterPayload(
			dlList[i].GatewayID,
			dlList[i].DoorlockAddress,
			&newScheduler,
			&mqttSvc.UserIDPassword{
				UserId:     userScheduler.UserID,
				RfidPass:   userScheduler.RfidPass,
				KeypadPass: userScheduler.KeypadPass,
			}))

		if err := mqttSvc.HandleMqttErr(t); err != nil {
			utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
				StatusCode: http.StatusBadRequest,
				Msg:        "Create scheduler mqtt failed",
				ErrorMsg:   err.Error(),
			})
			return
		}

		if userRole == "employee" {
			_, err = h.deps.SvcOpts.EmployeeSvc.AppendEmployeeSchedulerExcel(c.Request.Context(), &newScheduler)
			if err != nil {
				utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
					StatusCode: http.StatusBadRequest,
					Msg:        "Update employee failed",
					ErrorMsg:   err.Error(),
				})
				return
			}
		} else if userRole == "student" {
			_, err = h.deps.SvcOpts.StudentSvc.AppendStudentSchedulerExcel(c.Request.Context(), &newScheduler)
			if err != nil {
				utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
					StatusCode: http.StatusBadRequest,
					Msg:        "Update student failed",
					ErrorMsg:   err.Error(),
				})
				return
			}
		} else if userRole == "customer" {
			_, err = h.deps.SvcOpts.CustomerSvc.AppendCustomerSchedulerExcel(c.Request.Context(), &newScheduler)
			if err != nil {
				utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
					StatusCode: http.StatusBadRequest,
					Msg:        "Update customer failed",
					ErrorMsg:   err.Error(),
				})
				return
			}
		} else {
			utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
				StatusCode: http.StatusBadRequest,
				Msg:        "Get role of user failed",
				ErrorMsg:   err.Error(),
			})
			return
		}
	}

	utils.ResponseJson(c, http.StatusOK, true)
}

// Update Scheduler base on Excel
// @Summary Update Scheduler Base On Excel
// @Schemes
// @Description Update Scheduler base on Excel
// @Accept  json
// @Produce json
// @Param	data	body	models.SwagUpdateScheduler	true	"Fields need to append a scheduler"
// @Success 200 {object} models.Scheduler
// @Failure 400 {object} utils.ErrorResponse
// @Router /v1//scheduler/excel [Patch]
func (h *SchedulerHandler) UpdateSchedulerOnExcel(c *gin.Context) {
	userScheduler := &models.UserScheduler{}
	err := c.ShouldBind(&userScheduler.ScheInfo)
	if err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Invalid req body",
			ErrorMsg:   err.Error(),
		})
		return
	}

	roomId := userScheduler.ScheInfo.RoomID
	dlList, err := h.deps.SvcOpts.DoorlockSvc.FindAllDoorlocksByRoomID(c, roomId)

	if err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Get doorlocks by room ID failed",
			ErrorMsg:   err.Error(),
		})
		return
	}
	userRole := userScheduler.ScheInfo.Role
	userScheduler, err = getUserInformation(c, h.deps.SvcOpts, userRole, userScheduler)
	if err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Get user failed",
			ErrorMsg:   err.Error(),
		})
		return
	}

	isSuccess, err := h.deps.SvcOpts.SchedulerSvc.UpdateScheduler(c.Request.Context(), &userScheduler.ScheInfo)
	if err != nil || !isSuccess {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Update scheduler failed",
			ErrorMsg:   err.Error(),
		})
		return
	}

	for i := 0; i < len(dlList); i++ {

		t := h.deps.MqttClient.Publish(mqttSvc.TOPIC_SV_SCHEDULER_U, 1, false, mqttSvc.ServerCreateRegisterPayload(
			dlList[i].GatewayID,
			dlList[i].DoorlockAddress,
			&userScheduler.ScheInfo,
			&mqttSvc.UserIDPassword{
				UserId:     userScheduler.UserID,
				RfidPass:   userScheduler.RfidPass,
				KeypadPass: userScheduler.KeypadPass,
			}))

		if err := mqttSvc.HandleMqttErr(t); err != nil {
			utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
				StatusCode: http.StatusBadRequest,
				Msg:        "Create scheduler mqtt failed",
				ErrorMsg:   err.Error(),
			})
			return
		}
	}

	utils.ResponseJson(c, http.StatusOK, true)
}

func getUserInformation(c *gin.Context, optSvc *models.ServiceOptions, userRole string, userScheduler *models.UserScheduler) (*models.UserScheduler, error) {
	if userRole == "employee" {
		userEmp, err := optSvc.EmployeeSvc.FindEmployeeByMSNV(c, userScheduler.ScheInfo.UserID)
		if err != nil {
			return nil, err
		}
		userScheduler.UserID = userEmp.MSNV
		userScheduler.RfidPass = userEmp.RfidPass
		userScheduler.KeypadPass = userEmp.KeypadPass
		userScheduler.ScheInfo.EmployeeID = &userEmp.MSNV
	} else if userRole == "student" {
		userStu, err := optSvc.StudentSvc.FindStudentByMSSV(c, userScheduler.ScheInfo.UserID)
		if err != nil {
			return nil, err
		}
		userScheduler.UserID = userStu.MSSV
		userScheduler.RfidPass = userStu.RfidPass
		userScheduler.KeypadPass = userStu.KeypadPass
		userScheduler.ScheInfo.StudentID = &userStu.MSSV
	} else if userRole == "customer" {
		userCus, err := optSvc.CustomerSvc.FindCustomerByCCCD(c, userScheduler.ScheInfo.UserID)
		if err != nil {
			return nil, err
		}
		userScheduler.UserID = userCus.CCCD
		userScheduler.RfidPass = userCus.RfidPass
		userScheduler.KeypadPass = userCus.KeypadPass
		userScheduler.ScheInfo.CustomerID = &userCus.CCCD

	} else {

		return nil, fmt.Errorf("no user record")
	}
	return userScheduler, nil
}
