package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/ecoprohcm/DMS_BackendServer/models"
	"github.com/ecoprohcm/DMS_BackendServer/utils"
	"github.com/gin-gonic/gin"
)

type GatewayLogHandler struct {
	deps *HandlerDependencies
}

func NewGatewayLogHandler(deps *HandlerDependencies) *GatewayLogHandler {
	return &GatewayLogHandler{
		deps,
	}
}

// Find all gateway logs info
// @Summary Find All GatewayLog
// @Schemes
// @Description find all gateway logs info
// @Produce json
// @Success 200 {array} []models.GatewayLog
// @Failure 400 {object} utils.ErrorResponse
// @Router /v1/gatewayLogs [get]
func (h *GatewayLogHandler) FindAllGatewayLog(c *gin.Context) {
	glList, err := h.deps.SvcOpts.LogSvc.FindAllGatewayLog(c)
	if err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Get all gateway logs failed",
			ErrorMsg:   err.Error(),
		})
		return
	}
	utils.ResponseJson(c, http.StatusOK, glList)
}

// Find gateway log info by id
// @Summary Find GatewayLog By ID
// @Schemes
// @Description find gateway log info by id
// @Produce json
// @Param        id	path	string	true	"GatewayLog ID"
// @Success 200 {object} models.GatewayLog
// @Failure 400 {object} utils.ErrorResponse
// @Router /v1/gatewayLog/{id} [get]
func (h *GatewayLogHandler) FindGatewayLogByID(c *gin.Context) {
	id := c.Param("id")
	gl, err := h.deps.SvcOpts.LogSvc.FindGatewayLogByID(c, id)
	if err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Get gateway log failed",
			ErrorMsg:   err.Error(),
		})
		return
	}
	utils.ResponseJson(c, http.StatusOK, gl)
}

// Update GatewayLogs Cleaner time period (Default: 1 week)
// @Summary Update GatewayLogs Cleaner time period (Default: 1 week)
// @Schemes
// @Description Change time period for GatewayLogs Cleaner
// @Produce json
// @Param        period	path	string	true	"Period (HOURS)"
// @Success 200 {object} boolean
// @Failure 400 {object} utils.ErrorResponse
// @Router /v1/gatewayLogs/period [post]
func (h *GatewayLogHandler) UpdateGatewayLogCleanPeriod(c *gin.Context) {
	period := models.GatewayLogTime{}
	err := c.ShouldBind(&period)
	if err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Invalid req body",
			ErrorMsg:   err.Error(),
		})
		return
	}
	err = h.deps.SvcOpts.LogSvc.UpdateGatewayLogCleanPeriod(&period)
	if err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Incorrect period",
			ErrorMsg:   err.Error(),
		})
	}
	utils.ResponseJson(c, http.StatusOK, true)
}

// Find Gateway logs by period of time
// @Summary Find Gateway logs by period of time
// @Schemes
// @Description find Gateway logs by period of time
// @Produce json
// @Param        id	path	string	true	"GatewayLog ID"
// @Param 		 from path  string  true    "From Unix time"
// @Param 		 to path    string  true    "To Unix time"
// @Success 200 {object} []models.GatewayLog
// @Failure 400 {object} utils.ErrorResponse
// @Router /v1/gatewayLogs/period/{id}/date/:from/:to [get]
func (h *GatewayLogHandler) FindGatewayLogsByTime(c *gin.Context) {
	gatewayId := c.Param("id")
	from := c.Param("from")
	to := c.Param("to")
	fromInt, _ := strconv.ParseInt(from, 10, 64)
	toInt, _ := strconv.ParseInt(to, 10, 64)
	fromFormatted := time.Unix(fromInt, 0).Format(models.DEFAULT_TIME_FORMAT)
	toFormatted := time.Unix(toInt, 0).Format(models.DEFAULT_TIME_FORMAT)
	glList, err := h.deps.SvcOpts.LogSvc.FindGatewayLogsByTime(gatewayId, fromFormatted, toFormatted)

	if err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Failed to get gateway logs",
			ErrorMsg:   err.Error(),
		})
	}
	utils.ResponseJson(c, http.StatusOK, glList)
}

// Find Gateway logs by log type in period of time
// @Summary Find Gateway logs by log type in period of time
// @Schemes
// @Description Get gateway's logs by log type (DEBUG/INFO) in a period of time (UNIX)
// @Produce json
// @Param        id		path	string	true	"GatewayLog ID"
// @Param 		 type  	query   string  no    	"Log's type (DEBUG/INFO), default is DEBUG"
// @Param 		 from  	query	string  true    "From Unix time"
// @Param 		 to    	query	string  no    	"To Unix time, default use current time"
// @Success 200 {object} []models.GatewayLog
// @Failure 400 {object} utils.ErrorResponse
// @Router /v1/gatewayLogs/:id/period [get]
func (h *GatewayLogHandler) FindGatewayLogsTypeByTime(c *gin.Context) {
	var fromFormatted, toFormatted string
	gatewayId := c.Param("id")
	logType := c.DefaultQuery("type", "DEBUG")
	from, isExist := c.GetQuery("from")
	if !isExist {
		/* If from parameter not exist, return StatusBadRequest */
		utils.ResponseJson(c, http.StatusBadRequest, utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Invalid request body",
			ErrorMsg:   fmt.Errorf("Missing <from> parameters").Error(),
		})
	} else {
		fromInt, _ := strconv.ParseInt(from, 10, 64)
		fromFormatted = time.Unix(fromInt, 0).Format(models.DEFAULT_TIME_FORMAT)
	}

	to, isExist := c.GetQuery("to")
	if !isExist || to == "" {
		/* If to parameter not exist or to empty, use current time instead */
		toFormatted = time.Now().Format(models.DEFAULT_TIME_FORMAT)
	} else {
		toInt, _ := strconv.ParseInt(to, 10, 64)
		toFormatted = time.Unix(toInt, 0).Format(models.DEFAULT_TIME_FORMAT)
	}
	glList, err := h.deps.SvcOpts.LogSvc.FindGatewayLogsTypeByTime(gatewayId, logType, fromFormatted, toFormatted)
	if err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        fmt.Sprintf("Failed to get gateway logs"),
			ErrorMsg:   err.Error(),
		})
	}
	utils.ResponseJson(c, http.StatusOK, glList)
}
