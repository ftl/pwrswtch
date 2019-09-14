package trx

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/ftl/rigproxy/pkg/protocol"
	"github.com/pkg/errors"
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

func (t *TRX) SetPowerLevel(value float32) error {
	valueStr := fmt.Sprintf("%.6f", value)
	log.Printf("Setting power to %s", valueStr)
	return t.sendCommand(
		protocol.LongCommand("set_level"),
		"RFPOWER",
		valueStr,
	)
}

func (t *TRX) SetTx(enabled bool) error {
	valueStr := boolTo01(enabled)
	log.Printf("Setting the Tx to %s", valueStr)
	return t.sendCommand(
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

func (t *TRX) sendCommand(command protocol.Command, args ...string) error {
	out, err := net.Dial("tcp", t.address)
	if err != nil {
		return errors.Wrap(err, "cannot connect to TRX")
	}
	trx := protocol.NewTransceiver(out)
	defer trx.Close()
	trx.WhenDone(func() {
		out.Close()
	})

	request := protocol.Request{
		Command: command,
		Args: args,
	}
	_, err = trx.Send(context.Background(), request)
	if err != nil {
		return errors.Wrapf(err, "cannot send %v command with args %v", command.Long, args)
	}
	return nil
}