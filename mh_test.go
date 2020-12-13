package mh_test

import (
	"os"
	"testing"
	"time"

	"github.com/AlexeySemigradsky/mh"
)

var address = os.Getenv("DEVICE_ADDRESS")
var timeout = 3 * time.Second
var controller = mh.NewController(address, timeout)

func TestController_OpenConnection(t *testing.T) {
	err := controller.OpenConnection()
	if err != nil {
		t.Error(err)
	}
}

func TestController_CloseConnection(t *testing.T) {
	err := controller.CloseConnection()
	if err != nil {
		t.Error(err)
	}
}

func TestController_SendRequest(t *testing.T) {
	request := []byte{0x81, 0x8a, 0x8b}
	response := make([]byte, 14)
	err := controller.SendRequest(request, response)
	if err != nil {
		t.Error(err)
	}
}

func TestController_GetState(t *testing.T) {
	_, err := controller.GetState()
	if err != nil {
		t.Error(err)
	}
}

func TestController_GetPower(t *testing.T) {
	_, err := controller.GetPower()
	if err != nil {
		t.Error(err)
	}
}

func TestController_GetRGBW(t *testing.T) {
	_, err := controller.GetRGBW()
	if err != nil {
		t.Error(err)
	}
}

func TestController_SetPower(t *testing.T) {
	power, err := controller.GetPower()
	if err != nil {
		t.Error(err)
	}

	err = controller.SetPower(!power)
	if err != nil {
		t.Error(err)
	}

	time.Sleep(time.Second)

	err = controller.SetPower(power)
	if err != nil {
		t.Error(err)
	}
}

func TestController_SetRGBW(t *testing.T) {
	rgbw, err := controller.GetRGBW()
	if err != nil {
		t.Error(err)
	}

	err = controller.SetRGBW(&mh.RGBW{Red: 255, Green: 0, Blue: 0, White: 0})
	if err != nil {
		t.Error(err)
	}

	time.Sleep(time.Second)

	err = controller.SetRGBW(&mh.RGBW{Red: 0, Green: 255, Blue: 0, White: 0})
	if err != nil {
		t.Error(err)
	}

	time.Sleep(time.Second)

	err = controller.SetRGBW(&mh.RGBW{Red: 0, Green: 0, Blue: 255, White: 0})
	if err != nil {
		t.Error(err)
	}

	time.Sleep(time.Second)

	err = controller.SetRGBW(&mh.RGBW{Red: 0, Green: 0, Blue: 0, White: 255})
	if err != nil {
		t.Error(err)
	}

	time.Sleep(time.Second)

	err = controller.SetRGBW(rgbw)
	if err != nil {
		t.Error(err)
	}
}
