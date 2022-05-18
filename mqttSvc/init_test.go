//go:build unit
// +build unit

package mqttSvc

import (
	"fmt"
	"testing"

	"github.com/ecoprohcm/DMS_BackendServer/models"
)

func TestParseDoorlockPayload(t *testing.T) {
	payload := fmt.Sprintf(`{
		"gateway_id":"%s",
		"message": {
		"doorlock_address":"%s",
		"doorlock_active_state":"%s",
	}
		
	}`, "test", "test", "test")
	dl := parseDoorlockPayload(payload)
	expected := &models.Doorlock{
		GatewayID:       "test",
		DoorlockAddress: "test",
		ActiveState:     "test",
	}

	if expected.GatewayID != dl.GatewayID {
		t.Errorf("got %+v, wanted %+v", dl.GatewayID, expected.GatewayID)
	}

	if expected.DoorlockAddress != dl.DoorlockAddress {
		t.Errorf("got %+v, wanted %+v", dl.DoorlockAddress, expected.DoorlockAddress)
	}

	if expected.ActiveState != dl.ActiveState {
		t.Errorf("got %+v, wanted %+v", dl.ActiveState, expected.ActiveState)
	}
}
