package models

import (
	"context"

	"github.com/ecoprohcm/DMS_BackendServer/utils"
	"gorm.io/gorm"
)

type DoorlockStatusLogTime struct {
	Day    int `json:"day"`
	Hour   int `json:"hour"`
	Minute int `json:"minute"`
	Second int `json:"second"`
}
type DoorlockStatusLog struct {
	GormModel
	DoorID     string `json:"doorId"`
	StateType  string `json:"statusType"` // ConnectState, DoorState, LockState
	StateValue string `json:"stateValue"` // corresponding statetype
}
type DoorlockStatusLogSvc struct {
	db *gorm.DB
}

func NewDoorlockStatusLogSvc(db *gorm.DB) *DoorlockStatusLogSvc {
	doorlockStatusLogSvc := &DoorlockStatusLogSvc{
		db: db,
	}
	return doorlockStatusLogSvc
}

func (dlsls *DoorlockStatusLogSvc) GetAllDoorlockStatusLogs(ctx context.Context) (dlslList []DoorlockStatusLog, err error) {
	result := dlsls.db.Find(&dlslList)
	if err := result.Error; err != nil {
		err = utils.HandleQueryError(err)
		return nil, err
	}
	return dlslList, nil
}

func (dlsls *DoorlockStatusLogSvc) GetDoorlockStatusLogByDoorID(ctx context.Context, doorId string) (dlslList []DoorlockStatusLog, err error) {
	result := dlsls.db.Where("door_id = ?", doorId).Find(&dlslList)
	if err := result.Error; err != nil {
		err = utils.HandleQueryError(err)
		return nil, err
	}
	return dlslList, nil
}

func (dlsls *DoorlockStatusLogSvc) CreateDoorlockStatusLog(ctx context.Context, dlsl *DoorlockStatusLog) (*DoorlockStatusLog, error) {
	if err := dlsls.db.Create(&dlsl).Error; err != nil {
		err = utils.HandleQueryError(err)
		return nil, err
	}
	return dlsl, nil
}

func (dlsls *DoorlockStatusLogSvc) GetDoorlockStatusLogInTimeRange(from string, to string) (dlslList *[]DoorlockStatusLog, err error) {
	result := dlsls.db.Where("created_at >= ? AND created_at <= ?", from, to).Find(&dlslList)
	if err := result.Error; err != nil {
		err = utils.HandleQueryError(err)
		return nil, err
	}
	return dlslList, nil
}

func (dlsls *DoorlockStatusLogSvc) DeleteDoorlockStatusLogInTimeRange(from string, to string) (bool, error) {
	result := dlsls.db.Unscoped().Where("created_at >= ? AND created_at <= ?", from, to).Delete(&DoorlockStatusLog{})
	return utils.ReturnBoolStateFromResult(result)
}

func (dlsls *DoorlockStatusLogSvc) DeleteDoorlockStatusLogByDoorID(doorId string) (bool, error) {
	result := dlsls.db.Unscoped().Where("door_id = ?", doorId).Delete(&DoorlockStatusLog{})
	return utils.ReturnBoolStateFromResult(result)
}
