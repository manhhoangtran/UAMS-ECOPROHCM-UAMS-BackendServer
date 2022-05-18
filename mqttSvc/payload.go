package mqttSvc

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/ecoprohcm/DMS_BackendServer/models"
)

type UserIDPassword struct {
	UserId     string `json:"user_id"`
	RfidPass   string `json:"rfid_pw"`
	KeypadPass string `json:"keypad_pw"`
}
type DoorlockBootUp struct {
	DoorlockAddress string `json:"doorlock_address"`
	ActiveState     string `json:"doorlock_active_state"`
}

type SchedulerBootUp struct {
	SchedulerId     string `json:"register_id"`
	UserId          string `json:"user_id"`
	RfidPass        string `json:"rfid_pw"`
	KeypadPass      string `json:"keypad_pw"`
	DoorlockAddress string `json:"doorlock_address"`
	StartDate       string `json:"start_date"`
	EndDate         string `json:"end_date"`
	WeekDay         string `json:"week_day"`
	StartClass      string `json:"start_class"`
	EndClass        string `json:"end_class"`
}

func ServerCreateDoorlockPayload(doorlock *models.Doorlock) string {
	msg := fmt.Sprintf(`{"doorlock_address":"%s"}`, doorlock.DoorlockAddress)
	return PayloadWithGatewayId(doorlock.GatewayID, msg)
}

func ServerUpdateDoorlockPayload(doorlock *models.Doorlock) string {
	msg := fmt.Sprintf(`{"doorlock_address":"%s","doorlock_active_state":"%s"}`,
		doorlock.DoorlockAddress, doorlock.ActiveState)
	return PayloadWithGatewayId(doorlock.GatewayID, msg)
}

func ServerDeleteDoorlockPayload(doorlock *models.Doorlock) string {
	msg := fmt.Sprintf(`{"doorlock_address":"%s"}`, doorlock.DoorlockAddress)
	return PayloadWithGatewayId(doorlock.GatewayID, msg)
}

func ServerCmdDoorlockPayload(gwId string, doorlockAddress string, cmd *models.DoorlockCmd) string {
	var duration string = ""
	if cmd.Duration != "" {
		duration = fmt.Sprintf(`,"duration":"%s"`, cmd.Duration)
	}
	msg := fmt.Sprintf(`{"doorlock_address":"%s","action":"%s"%s}`, doorlockAddress, cmd.State, duration)
	return PayloadWithGatewayId(gwId, msg)
}

func ServerUpdateGatewayPayload(gw *models.Gateway) string {
	return fmt.Sprintf(`{"gateway_id":"%s","area_id":"%s","name":"%s"}`,
		gw.GatewayID, gw.AreaID, gw.Name)
}

func ServerDeleteGatewayPayload(gwID string) string {
	msg := `{}`
	return PayloadWithGatewayId(gwID, msg)
}

func ServerCreateRegisterPayload(
	gwId string,
	doorlockAddress string,
	sche *models.Scheduler,
	uP *UserIDPassword,
) string {

	loc, _ := time.LoadLocation("Asia/Ho_Chi_Minh")
	startDmySlice := getDayMonthYearSlice(sche.StartDate)
	start := time.Date(startDmySlice[2], time.Month(startDmySlice[1]), startDmySlice[0], 0, 0, 0, 0, loc).Unix()
	endDmySlice := getDayMonthYearSlice(sche.EndDate)
	end := time.Date(endDmySlice[2], time.Month(endDmySlice[1]), endDmySlice[0], 23, 59, 59, 0, loc).Unix()

	msg := fmt.Sprintf(`{"register_id":"%d",
	"user_id":"%s",
	"doorlock_address":"%s",
	"rfid_pw":"%s",
	"keypad_pw":"%s",
	"start_date":"%d",
	"end_date":"%d",
	"week_day":"%d",
	"start_class":"%d",
	"end_class":"%d"}`,
		sche.ID, uP.UserId, doorlockAddress, uP.RfidPass, uP.KeypadPass,
		start, end, sche.WeekDay, sche.StartClassTime, sche.EndClassTime)

	return PayloadWithGatewayId(gwId, msg)
}

func ServerUpdateRegisterPayload(gwId string, uSche *models.UpdateScheduler) string {
	sche := uSche.Scheduler
	loc, _ := time.LoadLocation("Asia/Ho_Chi_Minh")
	startDmySlice := getDayMonthYearSlice(sche.StartDate)
	start := time.Date(startDmySlice[2], time.Month(startDmySlice[1]), startDmySlice[0], 0, 0, 0, 0, loc).Unix()
	endDmySlice := getDayMonthYearSlice(sche.EndDate)

	end := time.Date(endDmySlice[2], time.Month(endDmySlice[1]), endDmySlice[0], 23, 59, 59, 0, loc).Unix()
	msg := fmt.Sprintf(`{"register_id":"%d",
	"user_id":"%s",
	"doorlock_address":"%s",
	"start_date":"%d",
	"end_date":"%d",
	"week_day":"%d",
	"start_class":"%d",
	"end_class":"%d"}`,
		sche.ID, uSche.UserID, uSche.DoorlockAddress,
		start, end, sche.WeekDay, sche.StartClassTime, sche.EndClassTime)
	return PayloadWithGatewayId(gwId, msg)
}

func ServerDeleteRegisterPayload(gwId string, registerId uint) string {
	msg := fmt.Sprintf(`{"register_id":"%d"}`, registerId)
	return PayloadWithGatewayId(gwId, msg)
}

func ServerBootuptHPEmployeePayload(gwId string, emps []models.Employee) string {
	bootupEmps := []UserIDPassword{}
	for _, emp := range emps {
		buEmp := UserIDPassword{
			UserId:     emp.MSNV,
			RfidPass:   emp.RfidPass,
			KeypadPass: emp.KeypadPass,
		}
		bootupEmps = append(bootupEmps, buEmp)
	}
	bootupEmpsJson, _ := json.Marshal(bootupEmps)
	return PayloadWithGatewayId(gwId, string(bootupEmpsJson))
}

func ServerUpdateUserPayload(gwId string, userId string, rfidPw string, keypadPw string) string {
	msg := fmt.Sprintf(`{"user_id":"%s","rfid_pw":"%s", "keypad_pw":"%s"}`,
		userId, rfidPw, keypadPw)
	return PayloadWithGatewayId(gwId, msg)
}

func ServerDeleteUserPayload(gwId string, msnv string) string {
	msg := fmt.Sprintf(`{"user_id":"%s"}`, msnv)
	return PayloadWithGatewayId(gwId, msg)
}

func PayloadWithGatewayId(gwId string, msg string) string {
	return fmt.Sprintf(`{"gateway_id":"%s","message":%s}`, gwId, msg)
}

func getDayMonthYearSlice(str string) []int {
	strs := strings.Split(str, "/")
	var dmySlice = []int{}
	for _, s := range strs {
		number, err := strconv.ParseUint(s, 10, 32)
		if err != nil {
			panic(err)
		}
		dmySlice = append(dmySlice, int(number))
	}
	return dmySlice
}

func ServerUpdateSecretKeyPayload(gwId string, secretKey string) string {
	msg := fmt.Sprintf(`{"secret_key":"%s"}`, secretKey)
	return PayloadWithGatewayId(gwId, msg)
}

func ServerUpdateGatewayCmd(gwId string, action string) string {
	msg := fmt.Sprintf(`{"action":"%s"}`, action)
	return PayloadWithGatewayId(gwId, msg)
}

func ServerBootupDoorlocksPayload(gwId string, dls []models.Doorlock) string {
	bootupDls := []DoorlockBootUp{}
	for _, dl := range dls {
		buDl := DoorlockBootUp{
			DoorlockAddress: dl.DoorlockAddress,
			ActiveState:     dl.ActiveState,
		}
		bootupDls = append(bootupDls, buDl)
	}
	bootupDlsJson, _ := json.Marshal(bootupDls)
	return PayloadWithGatewayId(gwId, string(bootupDlsJson))
}

func ServerBootupRegisterPayload(
	gwId string,
	scheBoUpListPointer []*SchedulerBootUp,
) string {
	scheBoUpList := []SchedulerBootUp{}
	for _, sche := range scheBoUpListPointer {
		loc, _ := time.LoadLocation("Asia/Ho_Chi_Minh")
		startDmySlice := getDayMonthYearSlice(sche.StartDate)
		start := time.Date(startDmySlice[2], time.Month(startDmySlice[1]), startDmySlice[0], 0, 0, 0, 0, loc).Unix()
		sche.StartDate = strconv.FormatInt(start, 10)
		endDmySlice := getDayMonthYearSlice(sche.EndDate)
		end := time.Date(endDmySlice[2], time.Month(endDmySlice[1]), endDmySlice[0], 23, 59, 59, 0, loc).Unix()
		sche.EndDate = strconv.FormatInt(end, 10)

		if !isPastTime(end) {
			scheBoUpList = append(scheBoUpList, *sche)
		}
	}
	bootupScheJson, _ := json.Marshal(scheBoUpList)
	return PayloadWithGatewayId(gwId, string(bootupScheJson))
}

func isPastTime(t_compared int64) bool {

	t_now := time.Now().Unix()
	if t_compared >= t_now {
		return false
	}

	return true
}

func ServerBootupSystemPayload(gwId string, srKey string) string {
	msg := fmt.Sprintf(`{"secret_key":"%s"}`, srKey)
	return PayloadWithGatewayId(gwId, msg)
}
