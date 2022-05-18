package handlers

import (
	"net/http"

	"github.com/ecoprohcm/DMS_BackendServer/models"
	"github.com/ecoprohcm/DMS_BackendServer/mqttSvc"
	"github.com/ecoprohcm/DMS_BackendServer/utils"
	"github.com/gin-gonic/gin"
)

type StudentHandler struct {
	deps *HandlerDependencies
}

func NewStudentHandler(deps *HandlerDependencies) *StudentHandler {
	return &StudentHandler{
		deps,
	}
}

// Find all students info
// @Summary Find All Student
// @Schemes
// @Description find all students info
// @Produce json
// @Success 200 {array} []models.Student
// @Failure 400 {object} utils.ErrorResponse
// @Router /v1/students [get]
func (h *StudentHandler) FindAllStudent(c *gin.Context) {
	sList, err := h.deps.SvcOpts.StudentSvc.FindAllStudent(c)
	if err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Get all students failed",
			ErrorMsg:   err.Error(),
		})
		return
	}
	utils.ResponseJson(c, http.StatusOK, sList)
}

// Find student info by mssv
// @Summary Find Student By MSSV
// @Schemes
// @Description find student info by student mssv
// @Produce json
// @Param        mssv	path	string	true	"Student MSSV"
// @Success 200 {object} models.Student
// @Failure 400 {object} utils.ErrorResponse
// @Router /v1/student/{mssv} [get]
func (h *StudentHandler) FindStudentByMSSV(c *gin.Context) {
	mssv := c.Param("mssv")

	s, err := h.deps.SvcOpts.StudentSvc.FindStudentByMSSV(c, mssv)
	if err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Get student failed",
			ErrorMsg:   err.Error(),
		})
		return
	}
	utils.ResponseJson(c, http.StatusOK, s)
}

// Create student
// @Summary Create Student
// @Schemes
// @Description Create student
// @Accept  json
// @Produce json
// @Param	data	body	models.SwagCreateStudent	true	"Fields need to create a student"
// @Success 200 {object} models.Student
// @Failure 400 {object} utils.ErrorResponse
// @Router /v1/student [post]
func (h *StudentHandler) CreateStudent(c *gin.Context) {
	s := &models.Student{}
	err := c.ShouldBind(s)
	if err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Invalid req body",
			ErrorMsg:   err.Error(),
		})
		return
	}

	_, err = h.deps.SvcOpts.StudentSvc.CreateStudent(c.Request.Context(), s)
	if err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Create student failed",
			ErrorMsg:   err.Error(),
		})
		return
	}
	utils.ResponseJson(c, http.StatusOK, s)
}

// Update student
// @Summary Update Student By ID
// @Schemes
// @Description Update student, must have correct "id" and "mssv" field. Send updated info to MQTT broker
// @Accept  json
// @Produce json
// @Param	data	body	models.Student	true	"Fields need to update a student"
// @Success 200 {boolean} true
// @Failure 400 {object} utils.ErrorResponse
// @Router /v1/student [patch]
func (h *StudentHandler) UpdateStudent(c *gin.Context) {
	s := &models.Student{}
	err := c.ShouldBind(s)
	if err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Invalid req body",
			ErrorMsg:   err.Error(),
		})
		return
	}

	isSuccess, err := h.deps.SvcOpts.StudentSvc.UpdateStudent(c.Request.Context(), s)
	if err != nil || !isSuccess {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Update student failed",
			ErrorMsg:   err.Error(),
		})
		return
	}

	t := h.deps.MqttClient.Publish(mqttSvc.TOPIC_SV_USER_U, 1, false,
		mqttSvc.ServerUpdateUserPayload("0", s.MSSV, s.RfidPass, s.KeypadPass))
	if err := mqttSvc.HandleMqttErr(t); err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Update student mqtt failed",
			ErrorMsg:   err.Error(),
		})
		return
	}

	utils.ResponseJson(c, http.StatusOK, isSuccess)
}

// Delete student
// @Summary Delete Student By MSSV
// @Schemes
// @Description Delete student using "mssv" field. Send deleted info to MQTT broker
// @Accept  json
// @Produce json
// @Param	data	body	models.DeleteStudent	true	"Student MSSV"
// @Success 200 {boolean} true
// @Failure 400 {object} utils.ErrorResponse
// @Router /v1/student [delete]
func (h *StudentHandler) DeleteStudent(c *gin.Context) {
	ds := &models.DeleteStudent{}
	err := c.ShouldBind(ds)
	if err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Invalid req body",
			ErrorMsg:   err.Error(),
		})
		return
	}

	isSuccess, err := h.deps.SvcOpts.StudentSvc.DeleteStudent(c.Request.Context(), ds.MSSV)
	if err != nil || !isSuccess {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Delete student failed",
			ErrorMsg:   err.Error(),
		})
		return
	}

	t := h.deps.MqttClient.Publish(mqttSvc.TOPIC_SV_USER_D, 1, false,
		mqttSvc.ServerDeleteUserPayload("0", ds.MSSV))
	if err := mqttSvc.HandleMqttErr(t); err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Delete student mqtt failed",
			ErrorMsg:   err.Error(),
		})
		return
	}

	utils.ResponseJson(c, http.StatusOK, isSuccess)
}

// Add student scheduler
// @Summary Add Door Open Scheduler For Student
// @Schemes
// @Description Add scheduler that allows student open specific door. Send updated info to MQTT broker
// @Accept  json
// @Produce json
// @Param	data	body	models.UserSchedulerReq	true	"Request with Scheduler, GatewayID, DoorlockAdress"
// @Success 200 {boolean} true
// @Failure 400 {object} utils.ErrorResponse
// @Router /v1/student/{msnv}/scheduler [post]
func (h *StudentHandler) AppendStudentScheduler(c *gin.Context) {
	usu := &models.UserSchedulerReq{}
	mssv := c.Param("mssv")
	err := c.ShouldBind(usu)
	if err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Invalid req body",
			ErrorMsg:   err.Error(),
		})
		return
	}

	s, err := h.deps.SvcOpts.StudentSvc.FindStudentByMSSV(c, mssv)
	if err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Get student failed",
			ErrorMsg:   err.Error(),
		})
		return
	}

	sche := &usu.Scheduler
	sche.StudentID = &s.MSSV
	sche.Role = "student"
	sche.UserID = s.MSSV
	_, err = h.deps.SvcOpts.SchedulerSvc.CreateScheduler(c, sche)

	if err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Create scheduler failed",
			ErrorMsg:   err.Error(),
		})
		return
	}

	t := h.deps.MqttClient.Publish(mqttSvc.TOPIC_SV_SCHEDULER_C, 1, false, mqttSvc.ServerCreateRegisterPayload(
		usu.GatewayID,
		usu.DoorlockAddress,
		sche,
		&mqttSvc.UserIDPassword{
			UserId:     s.MSSV,
			RfidPass:   s.RfidPass,
			KeypadPass: s.KeypadPass,
		},
	))
	if err := mqttSvc.HandleMqttErr(t); err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Create scheduler mqtt failed",
			ErrorMsg:   err.Error(),
		})
		return
	}

	_, err = h.deps.SvcOpts.StudentSvc.AppendStudentScheduler(c.Request.Context(), s, usu, sche)
	if err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Update student failed",
			ErrorMsg:   err.Error(),
		})
		return
	}

	utils.ResponseJson(c, http.StatusOK, true)
}
