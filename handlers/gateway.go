package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/ecoprohcm/DMS_BackendServer/models"
	"github.com/ecoprohcm/DMS_BackendServer/mqttSvc"
	"github.com/ecoprohcm/DMS_BackendServer/utils"
	"github.com/gin-gonic/gin"
)

type GatewayHandler struct {
	deps *HandlerDependencies
}

func NewGatewayHandler(deps *HandlerDependencies) *GatewayHandler {
	return &GatewayHandler{
		deps,
	}
}

// Find all gateways and doorlocks info
// @Summary Find All Gateway
// @Schemes
// @Description find all gateways info
// @Produce json
// @Success 200 {array} []models.Gateway
// @Failure 400 {object} utils.ErrorResponse
// @Router /v1/gateways [get]
func (h *GatewayHandler) FindAllGateway(c *gin.Context) {
	gwList, err := h.deps.SvcOpts.GatewaySvc.FindAllGateway(c)
	if err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Get all gateways failed",
			ErrorMsg:   err.Error(),
		})
		return
	}
	utils.ResponseJson(c, http.StatusOK, gwList)
}

// Find gateway and doorlock info by id
// @Summary Find Gateway By ID
// @Schemes
// @Description find gateway and doorlock info by gateway id
// @Produce json
// @Param        id	path	string	true	"Gateway ID"
// @Success 200 {object} models.Gateway
// @Failure 400 {object} utils.ErrorResponse
// @Router /v1/gateway/{id} [get]
func (h *GatewayHandler) FindGatewayByID(c *gin.Context) {
	id := c.Param("id")

	gw, err := h.deps.SvcOpts.GatewaySvc.FindGatewayByID(c, id)
	if err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Get gateway failed",
			ErrorMsg:   err.Error(),
		})
		return
	}
	utils.ResponseJson(c, http.StatusOK, gw)
}

// Create gateway
// @Summary Create Gateway
// @Schemes
// @Description Create gateway. Send created info to MQTT broker
// @Accept  json
// @Produce json
// @Param	data	body	models.SwagCreateGateway	true	"Fields need to create a gateway"
// @Success 200 {object} models.Gateway
// @Failure 400 {object} utils.ErrorResponse
// @Router /v1/gateway [post]
func (h *GatewayHandler) CreateGateway(c *gin.Context) {
	gw := &models.Gateway{}
	err := c.ShouldBind(gw)
	if err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Invalid req body",
			ErrorMsg:   err.Error(),
		})
		return
	}
	if len(gw.GatewayID) <= 0 || len(gw.Name) <= 0 {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Please fulfill these fields: name, gateway id",
			ErrorMsg:   "Missing on required fields",
		})
		return
	}

	gw, err = h.deps.SvcOpts.GatewaySvc.CreateGateway(c.Request.Context(), gw)
	if err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Create gateway failed",
			ErrorMsg:   err.Error(),
		})
		return
	}
	utils.ResponseJson(c, http.StatusOK, gw)
}

// Update gateway
// @Summary Update Gateway By Gateway ID
// @Schemes
// @Description Update gateway, must have "gateway_id" field. Send updated info to MQTT broker
// @Accept  json
// @Produce json
// @Param	data	body	models.SwagUpateGateway	true	"Fields need to update a gateway"
// @Success 200 {boolean} true
// @Failure 400 {object} utils.ErrorResponse
// @Router /v1/gateway [patch]
func (h *GatewayHandler) UpdateGateway(c *gin.Context) {
	gw := &models.Gateway{}
	err := c.ShouldBind(gw)
	if err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Invalid req body",
			ErrorMsg:   err.Error(),
		})
		return
	}

	isSuccess, err := h.deps.SvcOpts.GatewaySvc.UpdateGateway(c.Request.Context(), gw)
	if err != nil || !isSuccess {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Update gateway failed",
			ErrorMsg:   err.Error(),
		})
		return
	}

	t := h.deps.MqttClient.Publish(mqttSvc.TOPIC_SV_GATEWAY_U, 1, false, mqttSvc.ServerUpdateGatewayPayload(gw))
	if err := mqttSvc.HandleMqttErr(t); err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Update gateway mqtt failed",
			ErrorMsg:   err.Error(),
		})
		return
	}

	utils.ResponseJson(c, http.StatusOK, isSuccess)
}

// Delete gateway
// @Summary Delete Gateway By Gateway ID
// @Schemes
// @Description Delete gateway using "" field. Send deleted info to MQTT broker
// @Accept  json
// @Produce json
// @Param	data	body	object{gateway_id=string}	true	"Gateway ID"
// @Success 200 {boolean} true
// @Failure 400 {object} utils.ErrorResponse
// @Router /v1/gateway [delete]
func (h *GatewayHandler) DeleteGateway(c *gin.Context) {
	dgw := &models.DeleteGateway{}
	err := c.ShouldBind(dgw)
	if err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Invalid req body",
			ErrorMsg:   err.Error(),
		})
		return
	}

	dls, err := h.deps.SvcOpts.DoorlockSvc.FindAllDoorlockByGatewayID(c.Request.Context(), dgw.GatewayID)
	if err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Find doorlock failed",
			ErrorMsg:   err.Error(),
		})
		return
	}

	//delete gateway first
	isSuccess, err := h.deps.SvcOpts.GatewaySvc.DeleteGateway(c.Request.Context(), dgw.GatewayID)
	if err != nil || !isSuccess {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Delete gateway failed",
			ErrorMsg:   err.Error(),
		})
		return
	}

	t := h.deps.MqttClient.Publish(mqttSvc.TOPIC_SV_GATEWAY_D, 1, false, mqttSvc.ServerDeleteGatewayPayload(dgw.GatewayID))
	if err := mqttSvc.HandleMqttErr(t); err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Delete gateway mqtt failed",
			ErrorMsg:   err.Error(),
		})
		return
	}

	// delete doorlock belong to this gateway
	for i := 0; i < len(dls); i++ {
		isSuccess, err := h.deps.SvcOpts.DoorlockSvc.DeleteDoorlock(c.Request.Context(), strconv.FormatUint(uint64(dls[i].ID), 10))
		if err != nil || !isSuccess {
			utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
				StatusCode: http.StatusBadRequest,
				Msg:        "Delete doorlock failed",
				ErrorMsg:   err.Error(),
			})
			return
		}

		t := h.deps.MqttClient.Publish(mqttSvc.TOPIC_SV_DOORLOCK_D, 1, false,
			mqttSvc.ServerDeleteDoorlockPayload(&dls[i]))
		if err := mqttSvc.HandleMqttErr(t); err != nil {
			utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
				StatusCode: http.StatusBadRequest,
				Msg:        "Delete doorlock mqtt failed",
				ErrorMsg:   err.Error(),
			})
			return
		}

	}

	utils.ResponseJson(c, http.StatusOK, isSuccess)
}

func (h *GatewayHandler) DeleteGatewayDoorlock(c *gin.Context) {
	d := &models.Doorlock{}
	gwId := c.Param("id")
	err := c.ShouldBind(d)
	if err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Invalid req body",
			ErrorMsg:   err.Error(),
		})
		return
	}

	gw, err := h.deps.SvcOpts.GatewaySvc.FindGatewayByID(c, gwId)
	if err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Get gateway failed",
			ErrorMsg:   err.Error(),
		})
		return
	}

	_, err = h.deps.SvcOpts.GatewaySvc.DeleteGatewayDoorlock(c.Request.Context(), gw, d)
	if err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Delete gateway door failed",
			ErrorMsg:   err.Error(),
		})
		return
	}

	utils.ResponseJson(c, http.StatusOK, true)
}

// Unlock or Lock all gateway's doorlocks by BlockID
// @Summary Unlock or Lock all gateway's doorlocks by BlockID
// @Schemes
// @Description Unlock or Lock all doorlocks in a Block
// @Accept  json
// @Produce json
// @Param	data	body	models.GatewayBlockCmd	true	"Gateway Block command"
// @Success 200 {boolean} true
// @Failure 400 {object} utils.ErrorResponse
// @Router /v1/block/cmd [post]
func (h *GatewayHandler) UpdateGatewayCmdByBlockID(c *gin.Context) {
	cmd := &models.GatewayBlockCmd{}
	err := c.ShouldBind(cmd)
	if err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Invalid request body",
			ErrorMsg:   err.Error(),
		})
		return
	}
	// Find all gateways based on Block ID
	gwList, err := h.deps.SvcOpts.GatewaySvc.FindAllGatewaysByBlockID(c.Request.Context(), cmd.BlockId)
	if err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        fmt.Sprintf("Failed to find all gateways in Block ID %s\n", cmd.BlockId),
			ErrorMsg:   err.Error(),
		})
		return
	}
	// Send command action to all gateways
	for _, v := range gwList {
		t := h.deps.MqttClient.Publish(mqttSvc.TOPIC_SV_DOORLOCK_CMD, 1, false, mqttSvc.ServerUpdateGatewayCmd(v, cmd.Action))
		if err := mqttSvc.HandleMqttErr(t); err != nil {
			utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
				StatusCode: http.StatusBadRequest,
				Msg:        fmt.Sprintf("Failed to send command %s to all doorlocks at gatewayID %s\n", cmd.Action, v),
				ErrorMsg:   err.Error(),
			})
		}
	}
	// Update all doorlock status
	_, err = h.deps.SvcOpts.GatewaySvc.UpdateAllDoorlocksStateByBlockID(c.Request.Context(), cmd.BlockId, cmd.Action)
	if err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        fmt.Sprintf("Failed to update Doorlock state %s for blockID %s\n", cmd.Action, cmd.BlockId),
			ErrorMsg:   err.Error(),
		})
		return
	}
	utils.ResponseJson(c, http.StatusOK, true)
}

// Add gateway doorlock
// @Summary Add Doorlock for Gateway
// @Schemes
// @Description Add Doorlock for Gateway
// @Accept  json
// @Produce json
// @Param	data	body	models.Doorlock	true	"Request with Doorlock"
// @Success 200 {boolean} true
// @Failure 400 {object} utils.ErrorResponse
// @Router /v1/gateway/{id}/doorlock [post]
func (h *GatewayHandler) AppendGatewayDoorlock(c *gin.Context) {
	dl := &models.Doorlock{}
	gwID := c.Param("id")
	err := c.ShouldBind(dl)
	if err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Invalid req body",
			ErrorMsg:   err.Error(),
		})
		return
	}

	gw, err := h.deps.SvcOpts.GatewaySvc.FindGatewayByID(c, gwID)
	if err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Get gateway failed",
			ErrorMsg:   err.Error(),
		})
		return
	}
	dl.GatewayID = gw.GatewayID

	isSuccess, err := h.deps.SvcOpts.GatewaySvc.AppendGatewayDoorlock(c, gw, dl)
	if err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Update doorlock failed",
			ErrorMsg:   err.Error(),
		})
		return
	}

	t := h.deps.MqttClient.Publish(mqttSvc.TOPIC_SV_GATEWAY_U, 1, false, mqttSvc.ServerUpdateGatewayPayload(gw))
	if err := mqttSvc.HandleMqttErr(t); err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Update gateway mqtt failed",
			ErrorMsg:   err.Error(),
		})
		return
	}

	utils.ResponseJson(c, http.StatusOK, isSuccess)
}
