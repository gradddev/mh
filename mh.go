package mh

import (
	"net"
	"time"
)

type Config struct {
	IP      net.IP
	Timeout time.Duration
}

type Controller struct {
	config Config
}

func NewController(config Config) *Controller {
	controller := Controller{
		config: config,
	}
	return &controller
}

func (c *Controller) GetState() ([]byte, error) {
	request := []byte{0x81, 0x8a, 0x8b}
	response := make([]byte, 14)
	err := c.sendRequest(request, response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (c *Controller) GetPower() (bool, error) {
	state, err := c.GetState()
	if err != nil {
		return false, err
	}
	power := state[2] == 0x23
	return power, nil
}

func (c *Controller) SetPower(power bool) error {
	request := []byte{0x71, 0x23, 0x0f}
	if !power {
		request[1] = 0x24
	}
	err := c.sendRequest(request, nil)
	if err != nil {
		return err
	}
	return nil
}

type RGBW struct {
	Red, Green, Blue, White float64
}

func (c *Controller) GetRGBW() (*RGBW, error) {
	state, err := c.GetState()
	if err != nil {
		return nil, err
	}
	return &RGBW{
		Red:   float64(state[6]),
		Green: float64(state[7]),
		Blue:  float64(state[8]),
		White: float64(state[9]),
	}, nil
}

func (c *Controller) SetRGBW(rgbw *RGBW) error {
	err := c.SetPower(true)
	if err != nil {
		return err
	}
	request := []byte{
		0x31,
		uint8(rgbw.Red),
		uint8(rgbw.Green),
		uint8(rgbw.Blue),
		uint8(rgbw.White),
		0x00,
		0x0f,
	}
	err = c.sendRequest(request, nil)
	if err != nil {
		return err
	}
	return nil
}

func (c *Controller) sendRequest(request []byte, response []byte) error {
	connection, err := net.DialTimeout(
		"tcp",
		c.config.IP.String()+":5577",
		c.config.Timeout,
	)
	if err != nil {
		return err
	}

	checksum := calculateChecksum(request)
	request = append(request, checksum)
	_, err = connection.Write(request)
	if err != nil {
		return err
	}

	if response != nil {
		_, err = connection.Read(response)
		if err != nil {
			return err
		}
	}

	err = connection.Close()
	if err != nil {
		return err
	}

	return nil
}

func calculateChecksum(data []byte) byte {
	checksum := uint8(0)
	for _, b := range data {
		checksum += b
	}
	checksum &= 0xFF
	return checksum
}
