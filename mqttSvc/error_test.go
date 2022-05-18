//go:build unit
// +build unit

package mqttSvc

import (
	"errors"
	"testing"
	"time"
)

type MockToken struct {
}

func (mt *MockToken) Wait() bool {
	return false
}

func (mt *MockToken) WaitTimeout(td time.Duration) bool {
	return false
}

func (mt *MockToken) Done() <-chan struct{} {
	ch := make(chan struct{})
	return ch
}

func (mt *MockToken) Error() error {
	return errors.New("MQTT Error")
}
func TestHandleMqttErr(t *testing.T) {

	mockToken := &MockToken{}
	err := HandleMqttErr(mockToken)

	if err.Error() != "MQTT Error" {
		t.Errorf("got %v, wanted %v", err.Error(), "MQTT Error")
	}
}

func TestHandleMqttNilErr(t *testing.T) {

	err := HandleMqttErr(nil)

	if err != nil {
		t.Errorf("got %v, wanted %v", err, nil)
	}
}
