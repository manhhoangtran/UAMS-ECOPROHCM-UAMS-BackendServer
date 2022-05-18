// Package utils provides utility functions using in entire app
package utils

import (
	"errors"
	"fmt"

	logger "github.com/ecoprohcm/DMS_BackendServer/logs"
	"gorm.io/gorm"
)

// Error handler for ORM query
func HandleQueryError(err error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		logger.LogWithoutFields(logger.SQLSERVER, logger.ErrorLevel, "Can't find any record", err.Error())
		return fmt.Errorf("can't find any record")
	}
	logger.LogWithoutFields(logger.SQLSERVER, logger.ErrorLevel, err.Error())
	return err
}

// Define return values for update, delete query
func ReturnBoolStateFromResult(result *gorm.DB) (bool, error) {
	err := result.Error
	ra := result.RowsAffected
	if err != nil {
		err = HandleQueryError(err)
		return false, err
	}
	if ra > 0 {
		return true, nil
	} else {
		logger.LogWithoutFields(logger.SQLSERVER, logger.ErrorLevel, "No record affected", err.Error())
		return false, fmt.Errorf("no record affected")
	}
}
