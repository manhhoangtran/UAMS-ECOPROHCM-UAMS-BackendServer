package models

import (
	"context"

	"github.com/ecoprohcm/DMS_BackendServer/utils"
	"gorm.io/gorm"
)

type GwNetwork struct {
	GormModel
	GatewayID          string `gorm:"type:varchar(256);"`
	InterfaceName      string `gorm:"type:varchar(50);not null"`
	PrimaryIpAddress   string `gorm:"type:varchar(20);"`
	SecondaryIpAddress string `gorm:"type:varchar(20);"`
	MacAddress         string `gorm:"type:varchar(20);not null;"`
}

type GwNetworkSvc struct {
	db *gorm.DB
}

func NewGwNetworkSvc(db *gorm.DB) *GwNetworkSvc {
	return &GwNetworkSvc{
		db: db,
	}
}

func (gwns *GwNetworkSvc) CreateGwNetwork(ctx context.Context, gwNet *GwNetwork) (*GwNetwork, error) {
	if err := gwns.db.Create(&gwNet).Error; err != nil {
		err = utils.HandleQueryError(err)
		return nil, err
	}
	return gwNet, nil
}

func (gwns *GwNetworkSvc) FindGwNetworkByName(ctx context.Context, gwId string, ifName string) (gwNet *GwNetwork, err error) {
	result := gwns.db.Where("gateway_id = ? AND interface_name = ?", gwId, ifName).Find(&gwNet)
	if err := result.Error; err != nil {
		err = utils.HandleQueryError(err)
		return nil, err
	}
	return gwNet, err
}

func (gwns *GwNetworkSvc) UpdateGwNetwork(ctx context.Context, gwNet *GwNetwork) (bool, error) {
	gw, err := gwns.FindGwNetworkByName(ctx, gwNet.GatewayID, gwNet.InterfaceName)
	if err != nil || gw == nil {
		err = utils.HandleQueryError(err)
		return false, err
	}
	if err = gwns.db.Model(&gw).Select("primary_ip_address", "secondary_ip_address", "mac_address").Where("gateway_id = ? AND interface_name = ?", gwNet.GatewayID, gwNet.InterfaceName).Updates(GwNetwork{
		PrimaryIpAddress:   gwNet.PrimaryIpAddress,
		SecondaryIpAddress: gwNet.SecondaryIpAddress,
		MacAddress:         gwNet.MacAddress,
	}).Error; err != nil {
		err = utils.HandleQueryError(err)
		return false, err
	}
	return true, nil
}

func (gwns *GwNetworkSvc) DeleteGwNetwork(ctx context.Context, gwNet *GwNetwork) (bool, error) {
	// Delete all rows match specific gateway_id
	if err := gwns.db.Where("gateway_id = ?", gwNet.GatewayID).Delete(&GwNetwork{}).Error; err != nil {
		err = utils.HandleQueryError(err)
		return false, err
	}
	return true, nil
}
