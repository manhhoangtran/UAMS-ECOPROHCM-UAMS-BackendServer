package models

import (
	"context"

	"github.com/ecoprohcm/DMS_BackendServer/utils"
	"gorm.io/gorm"
)

type Area struct {
	GormModel
	Name    string `gorm:"unique;not null" json:"name"`
	Manager string `gorm:"not null" json:"manager"`
}
type AreaSvc struct {
	db *gorm.DB
}

func NewAreaSvc(db *gorm.DB) *AreaSvc {
	return &AreaSvc{
		db: db,
	}
}

func (as *AreaSvc) FindAllArea(ctx context.Context) (aList []Area, err error) {
	result := as.db.Find(&aList)
	if err := result.Error; err != nil {
		err = utils.HandleQueryError(err)
		return nil, err
	}
	return aList, nil
}

func (as *AreaSvc) FindAreaByID(ctx context.Context, id string) (a *Area, err error) {
	result := as.db.First(&a, id)
	if err := result.Error; err != nil {
		err = utils.HandleQueryError(err)
		return nil, err
	}
	return a, nil
}

func (as *AreaSvc) CreateArea(a *Area, ctx context.Context) (*Area, error) {
	if err := as.db.Create(&a).Error; err != nil {
		err = utils.HandleQueryError(err)
		return nil, err
	}
	return a, nil
}

func (as *AreaSvc) UpdateArea(ctx context.Context, a *Area) (bool, error) {
	result := as.db.Model(&a).Where("id = ?", a.ID).Updates(a)
	return utils.ReturnBoolStateFromResult(result)
}

func (as *AreaSvc) DeleteArea(ctx context.Context, areaId uint) (bool, error) {
	result := as.db.Unscoped().Where("id = ?", areaId).Delete(&Area{})
	return utils.ReturnBoolStateFromResult(result)
}
