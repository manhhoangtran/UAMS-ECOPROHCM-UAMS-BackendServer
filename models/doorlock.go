package models

import (
	"context"
	"fmt"
	"time"

	"github.com/ecoprohcm/DMS_BackendServer/utils"
	"gorm.io/gorm"
)

type Doorlock struct {
	GormModel
	DoorSerialID    string      `gorm:"type:varchar(256);unique;not null" json:"doorSerialId"`
	Location        string      `json:"location"`
	Description     string      `json:"description"`
	GatewayID       string      `gorm:"type:varchar(256);" json:"gatewayId"`
	LastOpenTime    uint        `json:"lastOpenTime"`
	ConnectState    string      `json:"connectState"`
	BlockId         string      `json:"blockId"`
	FloorId         string      `json:"floorId"`
	RoomId          string      `json:"roomId"`
	DoorState       string      `json:"doorState"`
	LockState       string      `json:"lockState"`
	DoorlockAddress string      `json:"doorlockAddress"`
	ActiveState     string      `json:"activeState"`
	Schedulers      []Scheduler `gorm:"foreignKey:DoorID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"schedulers"`
}

// Struct defines HTTP request payload for openning doorlock
type DoorlockCmd struct {
	ID       string `json:"id"`
	State    string `json:"state"`
	Duration string `json:"duration"`
}

// Struct defines HTTP request payload for deleting doorlock
type DoorlockDelete struct {
	ID string `json:"id" binding:"required"`
}

// Struct defines HTTP request payload for getting doorlock status
type DoorlockStatus struct {
	ID              string `json:"id"`
	GatewayID       string `json:"gatewayId"`
	DoorlockAddress string `json:"doorlockAddress"`
	ConnectState    string `json:"connectState"`
	DoorState       string `json:"doorState"`
	LockState       string `json:"lockState"`
}

type DoorlockSvc struct {
	db *gorm.DB
}

func NewDoorlockSvc(db *gorm.DB) *DoorlockSvc {
	return &DoorlockSvc{
		db: db,
	}
}

func (dls *DoorlockSvc) FindAllDoorlock(ctx context.Context) (dlList []Doorlock, err error) {
	result := dls.db.Preload("Schedulers").Find(&dlList)
	if err := result.Error; err != nil {
		err = utils.HandleQueryError(err)
		return nil, err
	}
	return dlList, nil
}

func (dls *DoorlockSvc) FindDoorlockByID(ctx context.Context, id string) (dl *Doorlock, err error) {
	result := dls.db.Preload("Schedulers").First(&dl, id)
	if err := result.Error; err != nil {
		err = utils.HandleQueryError(err)
		return nil, err
	}
	return dl, nil
}

func (dls *DoorlockSvc) FindDoorlockByAddress(ctx context.Context, address string, gwID string) (dl *Doorlock, err error) {
	var cnt int64
	result := dls.db.Preload("Schedulers").Where("doorlock_address = ? AND gateway_id = ?", address, gwID).Find(&dl).Count(&cnt)
	if err := result.Error; err != nil {
		err = utils.HandleQueryError(err)
		return nil, err
	}

	if cnt <= 0 {
		return nil, fmt.Errorf("find no records")
	}

	return dl, nil
}

func (dls *DoorlockSvc) CreateDoorlock(ctx context.Context, dl *Doorlock) (*Doorlock, error) {
	if err := dls.db.Create(&dl).Error; err != nil {
		err = utils.HandleQueryError(err)
		return nil, err
	}
	return dl, nil
}

func (dls *DoorlockSvc) UpdateDoorlock(ctx context.Context, dl *Doorlock) (bool, error) {
	result := dls.db.Model(&dl).Where("doorlock_address = ? AND gateway_id = ?", dl.DoorlockAddress, dl.GatewayID).Updates(dl)
	return utils.ReturnBoolStateFromResult(result)
}

func (dls *DoorlockSvc) UpdateDoorlockByAddress(ctx context.Context, dl *Doorlock) (bool, error) {
	result := dls.db.Model(&dl).Where("gateway_id = ? AND doorlock_address = ?", dl.GatewayID, dl.DoorlockAddress).Updates(dl)
	return utils.ReturnBoolStateFromResult(result)
}

func (dls *DoorlockSvc) UpdateDoorlockState(ctx context.Context, dl *DoorlockCmd) (bool, error) {
	result := dls.db.Model(&Doorlock{}).Where("id = ?", dl.ID).Update("last_open_time", time.Now().UnixMilli())
	return utils.ReturnBoolStateFromResult(result)
}

func (dls *DoorlockSvc) DeleteDoorlock(ctx context.Context, id string) (bool, error) {
	result := dls.db.Unscoped().Where("id = ?", id).Delete(&Doorlock{})
	return utils.ReturnBoolStateFromResult(result)
}

func (dls *DoorlockSvc) DeleteDoorlockByAddress(ctx context.Context, dl *Doorlock) (bool, error) {
	result := dls.db.Unscoped().Where("gateway_id = ? AND doorlock_address = ?", dl.GatewayID, dl.DoorlockAddress).Delete(&Doorlock{})
	return utils.ReturnBoolStateFromResult(result)
}

func (dls *DoorlockSvc) UpdateDoorlockStatus(ctx context.Context, dl *DoorlockStatus) (bool, error) {
	result := dls.db.Model(&Doorlock{}).Where("id = ?", dl.ID).Updates(Doorlock{DoorState: dl.DoorState, LockState: dl.LockState})
	return utils.ReturnBoolStateFromResult(result)
}

func (dls *DoorlockSvc) UpdateDoorState(ctx context.Context, dl *DoorlockStatus) (bool, error) {
	result := dls.db.Model(&Doorlock{}).Where("gateway_id = ? AND doorlock_address = ?", dl.GatewayID, dl.DoorlockAddress).Update("door_state", dl.DoorState)
	return utils.ReturnBoolStateFromResult(result)
}

func (dls *DoorlockSvc) UpdateLockState(ctx context.Context, dl *DoorlockStatus) (bool, error) {
	result := dls.db.Model(&Doorlock{}).Where("gateway_id = ? AND doorlock_address = ?", dl.GatewayID, dl.DoorlockAddress).Update("lock_state", dl.LockState)
	return utils.ReturnBoolStateFromResult(result)
}

func (dls *DoorlockSvc) GetDoorlockStatusByID(ctx context.Context, id string) (dl *DoorlockStatus, err error) {
	var cnt int64
	result := dls.db.Model(&Doorlock{}).Where("id = ?", id).Find(&dl).Count(&cnt)
	if err := result.Error; err != nil {
		err = utils.HandleQueryError(err)
		return nil, err
	}

	if cnt <= 0 {
		return nil, fmt.Errorf("find no records")
	}

	return dl, nil
}

func (dls *DoorlockSvc) UpdateDoorlockStateCmd(ctx context.Context, dl *DoorlockCmd) (bool, error) {
	result := dls.db.Model(&Doorlock{}).Where("id = ?", dl.ID).Update("lock_state", dl.State)
	return utils.ReturnBoolStateFromResult(result)
}

func (dls *DoorlockSvc) FindAllDoorlocksByRoomID(ctx context.Context, roomId string) (dl []*Doorlock, err error) {
	var cnt int64
	result := dls.db.Model(&Doorlock{}).Where("room_id = ?", roomId).Find(&dl).Count(&cnt)
	if err := result.Error; err != nil {
		err = utils.HandleQueryError(err)
		return nil, err
	}

	if cnt <= 0 {
		return nil, fmt.Errorf("find no doorlock records by roomId")
	}

	return dl, nil
}

func (dls *DoorlockSvc) FindAllDoorlockByGatewayID(ctx context.Context, gwId string) (dlList []Doorlock, err error) {
	result := dls.db.Preload("Schedulers").Where("gateway_id = ?", gwId).Find(&dlList)
	if err := result.Error; err != nil {
		err = utils.HandleQueryError(err)
		return nil, err
	}
	return dlList, nil
}
