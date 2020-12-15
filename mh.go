package mh

import (
	"net"
	"time"
)

type Request []byte
type Response []byte
type Controller struct {
	Address    string
	Timeout    time.Duration
	Connection net.Conn
}

func NewController(address string, timeout time.Duration) *Controller {
	controller := Controller{
		Address: address,
		Timeout: timeout,
	}
	return &controller
}

func (c *Controller) OpenConnection() error {
	connection, err := net.DialTimeout(
		"tcp",
		c.Address,
		c.Timeout,
	)
	if err != nil {
		return err
	}
	c.Connection = connection
	return nil
}

func (c *Controller) CloseConnection() error {
	err := c.Connection.Close()
	if err != nil {
		return err
	}
	c.Connection = nil
	return nil
}

func (c *Controller) SendRequest(request Request, response Response) error {
	if c.Connection == nil {
		err := c.OpenConnection()
		if err != nil {
			return err
		}
	}
	checksum := calculateChecksum(request)
	request = append(request, checksum)
	_, err := c.Connection.Write(request)
	if err != nil {
		return err
	}
	if response != nil {
		_, err = c.Connection.Read(response)
		if err != nil {
			return err
		}
	}
	err = c.CloseConnection()
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

func (c *Controller) GetState() (Response, error) {
	request := []byte{0x81, 0x8a, 0x8b}
	response := make([]byte, 14)
	err := c.SendRequest(request, response)
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
	err := c.SendRequest(request, nil)
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
	request := []byte{
		0x31,
		uint8(rgbw.Red),
		uint8(rgbw.Green),
		uint8(rgbw.Blue),
		uint8(rgbw.White),
		0x00,
		0x0f,
	}
	err := c.SendRequest(request, nil)
	if err != nil {
		return err
	}
	return nil
}
