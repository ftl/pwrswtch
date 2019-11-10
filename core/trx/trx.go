package trx

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/pkg/errors"

	"github.com/ftl/rigproxy/pkg/protocol"
)

func Open(address string) *TRX {
	if address == "" {
		address = "localhost:4532"
	}

	return &TRX{
		address: address,
	}
}

type TRX struct {
	address string
}

func (t *TRX) SetPowerLevel(value string) error {
	log.Printf("Setting power to %s", value)
	return t.setValue(
		protocol.LongCommand("set_level"),
		"RFPOWER",
		value,
	)
}

func (t *TRX) GetPowerLevel() (string, error) {
	data, err := t.sendSingleCommand(
		protocol.LongCommand("get_level"),
		"RFPOWER",
	)
	if err != nil {
		return "", err
	}
	if len(data) != 1 {
		return "", fmt.Errorf("unexpected RFPOWER level %v", data)
	}
	log.Printf("Got RFPOWER level %s", data[0])
	return data[0], nil
}

func (t *TRX) SetMode(mode, passband string) error {
	log.Printf("Setting mode to %s %s", mode, passband)
	return t.setValue(
		protocol.LongCommand("set_mode"),
		mode,
		passband,
	)
}

func (t *TRX) GetMode() (string, string, error) {
	data, err := t.sendSingleCommand(
		protocol.LongCommand("get_mode"),
	)
	if err != nil {
		return "", "", err
	}
	if len(data) != 2 {
		return "", "", fmt.Errorf("unexpeced mode response %v", data)
	}
	log.Printf("Got mode %s %s", data[0], data[1])
	return data[0], data[1], nil
}

func (t *TRX) SetTx(enabled bool) error {
	valueStr := boolTo01(enabled)
	log.Printf("Setting the Tx to %s", valueStr)
	return t.setValue(
		protocol.LongCommand("set_ptt"),
		valueStr,
	)
}

func boolTo01(b bool) string {
	if b {
		return "1"
	}
	return "0"
}

func (t *TRX) setValue(command protocol.Command, args ...string) error {
	_, err := t.sendSingleCommand(command, args...)
	return err
}

func (t *TRX) sendSingleCommand(command protocol.Command, args ...string) ([]string, error) {
	out, err := net.Dial("tcp", t.address)
	if err != nil {
		return nil, errors.Wrap(err, "cannot connect to TRX")
	}
	trx := protocol.NewTransceiver(out)
	defer trx.Close()
	trx.WhenDone(func() {
		out.Close()
	})

	request := protocol.Request{
		Command: command,
		Args:    args,
	}
	response, err := trx.Send(context.Background(), request)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot send %v command with args %v", command.Long, args)
	}
	return response.Data, nil
}
