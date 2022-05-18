package models

import (
	"context"
	"fmt"

	"github.com/ecoprohcm/DMS_BackendServer/utils"
	"gorm.io/gorm"
)

type SecretKey struct {
	GormModel
	Secret string `gorm:"varchar(255); unique;not null" json:"secret"`
}

// Struct defines HTTP request payload for updating secret key
type UpdateSecretKey struct {
	Secret string `json:"secret" binding:"required"`
}
type SecretKeySvc struct {
	db *gorm.DB
}

func NewSecretKeySvc(db *gorm.DB) *SecretKeySvc {
	return &SecretKeySvc{
		db: db,
	}
}

func (sks *SecretKeySvc) FindSecretKey(ctx context.Context) (sk *SecretKey, err error) {
	result := sks.db.Find(&sk)
	if err := result.Error; err != nil {
		err = utils.HandleQueryError(err)
		return nil, err
	}
	return sk, nil
}

func (sks *SecretKeySvc) CreateSecretKey(ctx context.Context, sk *SecretKey) (*SecretKey, error) {
	var cnt int64
	sks.db.Find(&sk).Count(&cnt)
	if cnt > 0 {
		return nil, fmt.Errorf("secret key already exist. Use update instead")
	}
	if err := sks.db.Create(&sk).Error; err != nil {
		err = utils.HandleQueryError(err)
		return nil, err
	}
	return sk, nil
}

func (sks *SecretKeySvc) UpdateSecretKey(ctx context.Context, sk *SecretKey) (bool, error) {
	result := sks.db.Model(&sk).Where("id = ?", sk.ID).Updates(sk)
	return utils.ReturnBoolStateFromResult(result)
}
