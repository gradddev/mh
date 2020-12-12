package mh

import (
	"net"
)

type Request []byte
type Response []byte
type Controller struct {
	Address    string
	Connection *net.TCPConn
}

func NewController(address string) *Controller {
	controller := Controller{
		Address: address,
	}
	return &controller
}

func (c *Controller) OpenConnection() error {
	tcpAddr, err := net.ResolveTCPAddr("tcp", c.Address)
	if err != nil {
		return err
	}
	connection, err := net.DialTCP("tcp", nil, tcpAddr)
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
	Red   uint8
	Green uint8
	Blue  uint8
	White uint8
}

func (c *Controller) GetRGBW() (*RGBW, error) {
	state, err := c.GetState()
	if err != nil {
		return nil, err
	}
	return &RGBW{
		Red:   state[6],
		Green: state[7],
		Blue:  state[8],
		White: state[9],
	}, nil
}

func (c *Controller) SetRGBW(rgbw *RGBW) error {
	request := []byte{
		0x31,
		rgbw.Red,
		rgbw.Green,
		rgbw.Blue,
		rgbw.White,
		0x00,
		0x0f,
	}
	err := c.SendRequest(request, nil)
	if err != nil {
		return err
	}
	return nil
}
