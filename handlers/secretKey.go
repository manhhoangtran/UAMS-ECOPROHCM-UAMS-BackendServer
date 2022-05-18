package handlers

import (
	"net/http"

	"github.com/ecoprohcm/DMS_BackendServer/models"
	"github.com/ecoprohcm/DMS_BackendServer/mqttSvc"
	"github.com/ecoprohcm/DMS_BackendServer/utils"
	"github.com/gin-gonic/gin"
)

type SecretKeyHandler struct {
	deps *HandlerDependencies
}

func NewSecretKeyHandler(deps *HandlerDependencies) *SecretKeyHandler {
	return &SecretKeyHandler{
		deps,
	}
}

// Find Mifare Card Secret Key
// @Summary Find Mifare Card Secret Key
// @Schemes
// @Description Find Mifare Card Secret Key
// @Produce json
// @Success 200 {object} models.SecretKey
// @Failure 400 {object} utils.ErrorResponse
// @Router /v1/secretkeys [get]
func (h *SecretKeyHandler) FindSecretKey(c *gin.Context) {
	sk, err := h.deps.SvcOpts.SecretKeySvc.FindSecretKey(c)
	if err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Get secret key failed",
			ErrorMsg:   err.Error(),
		})
		return
	}
	utils.ResponseJson(c, http.StatusOK, sk)
}

// Create Mifare Card secret key
// @Summary Create Mifare Card Secret Key
// @Schemes
// @Description Create Mifare Card secret key
// @Accept  json
// @Produce json
// @Param	data	body	models.SwagCreateSecretKey	true	"Fields need to create a Secret key"
// @Success 200 {object} models.SecretKey
// @Failure 400 {object} utils.ErrorResponse
// @Router /v1/secretkey [post]

func (h *SecretKeyHandler) CreateSecretKey(c *gin.Context) {
	csk := &models.SecretKey{}
	err := c.ShouldBind(csk)
	if err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Invalid req body",
			ErrorMsg:   err.Error(),
		})
		return
	}

	_, err = h.deps.SvcOpts.SecretKeySvc.CreateSecretKey(c.Request.Context(), csk)
	if err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Create secret key failed",
			ErrorMsg:   err.Error(),
		})
		return
	}
	t := h.deps.MqttClient.Publish(mqttSvc.TOPIC_SV_SYSTEM_U, 1, false,
		mqttSvc.ServerUpdateSecretKeyPayload("0", csk.Secret))
	if err := mqttSvc.HandleMqttErr(t); err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Update secret key mqtt failed",
			ErrorMsg:   err.Error(),
		})
		return
	}

	utils.ResponseJson(c, http.StatusOK, csk)
}

// Update secret key
// @Summary Update Mifare Card Secret Key
// @Schemes
// @Description Update secret key, must have "secret" field. Send updated info to MQTT broker
// @Accept  json
// @Produce json
// @Param	data	body	models.UpdateSecretKey	true	"Fields need to update a secret key"
// @Success 200 {boolean} true
// @Failure 400 {object} utils.ErrorResponse
// @Router /v1/secretkey [patch]

func (h *SecretKeyHandler) UpdateSecretKey(c *gin.Context) {
	s := &models.SecretKey{}
	err := c.ShouldBind(s)
	if err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Invalid req body",
			ErrorMsg:   err.Error(),
		})
		return
	}

	isSuccess, err := h.deps.SvcOpts.SecretKeySvc.UpdateSecretKey(c.Request.Context(), s)
	if err != nil || !isSuccess {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Update secret key failed",
			ErrorMsg:   err.Error(),
		})
		return
	}

	t := h.deps.MqttClient.Publish(mqttSvc.TOPIC_SV_SYSTEM_U, 1, false,
		mqttSvc.ServerUpdateSecretKeyPayload("0", s.Secret))
	if err := mqttSvc.HandleMqttErr(t); err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Update secret key mqtt failed",
			ErrorMsg:   err.Error(),
		})
		return
	}

	utils.ResponseJson(c, http.StatusOK, isSuccess)
}
