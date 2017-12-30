package main

import (
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/codingconcepts/env"
	"github.com/aotd1/argparse"
	"github.com/jacobsa/go-serial/serial"
)

type config struct {
	PortName 		string	`env:"UNIEL_PORT" default:"/dev/ttyUSB0" required:"true"`
	BaudRate 		uint	`env:"UNIEL_BAUD_RATE" default:"9600"`
	DataBits 		uint	`env:"UNIEL_DATA_BITS" default:"8"`
	StopBits 		uint	`env:"UNIEL_STOP_BITS" default:"1"`
	FlowControl 	bool	`env:"UNIEL_FLOW_CONTROL" default:"false"`
	DeviceAddress	int		`env:"UNIEL_DEVICE_ADDRESS" required:"true"`
}

func main() {
	godotenv.Load()

	config := config{}
	if err := env.Set(&config); err != nil {
		log.Fatal(err)
	}
	deviceAddress, err := getByteFromInt(config.DeviceAddress, 16)
	if err != nil {
		err := fmt.Errorf("bad address, please check environment: %v", err)
		log.Fatal(err)
	}

	p := argparse.NewParser("uniel", "Uniel UCH-* modules controlling program")
	offCmd := p.NewCommand("off", "Turn channel off")
	onCmd := p.NewCommand("on", "Turn channel on")
	channelFlag := p.String("c", "channel", &argparse.Options{Required: true, Help: "Channel: 1-8"})

	if err := p.Parse(os.Args); err != nil {
		fmt.Print(p.Usage(err))
		return
	}

	channel, err := getByteFromString(*channelFlag, 8)
	if err != nil {
		err := fmt.Errorf("bad channel, please check usage: %v", err)
		log.Fatal(p.Usage(err))
	}
	port := openPort(config)
	defer port.Close()

	if onCmd.Happened() {
		run(port, deviceAddress, 0x0A, [3]byte{0xFF, channel, 0x00})
	} else if offCmd.Happened() {
		run(port, deviceAddress, 0x0A, [3]byte{0x00, channel, 0x00})
	} else {
		err := fmt.Errorf("bad arguments, please check usage")
		log.Fatal(p.Usage(err))
	}

}

func getByteFromString(value string, max int) (byte, error) {
	intValue, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("argument not convertable to integer")
	} else {
		return getByteFromInt(intValue, max)
	}
}

func getByteFromInt(value int, max int) (byte, error) {
	if value <= 0 || value > max {
		return 0, fmt.Errorf("argument not in range 1 - %v", max)
	} else {
		return byte(value - 1), nil
	}
}

func openPort(cfg config) (io.ReadWriteCloser) {
	options := serial.OpenOptions{
		PortName: cfg.PortName,
		BaudRate: cfg.BaudRate,
		DataBits: cfg.DataBits,
		StopBits: cfg.StopBits,
		RTSCTSFlowControl: cfg.FlowControl,
		MinimumReadSize: 8,
		InterCharacterTimeout: 300,
	}
	port, err := serial.Open(options)
	if err != nil {
		log.Fatalf("serial.Open: %v", err)
	}
	return port
}

func run(port io.ReadWriteCloser, address byte, command byte, payload [3]byte) {
	//FF FF CC MM XX YY ZZ SS
	// CC код команды
	// MM Адрес модуля – адресата команды (00 – ответ модуля)
	// XX YY ZZ – байты данных, специфические для команд.
	// SS = CC + MM + XX + YY + ZZ
	checksum := uint8(command + address + payload[0] + payload[1] + payload[2])
	content := []byte{0xFF, 0xFF, command, address, payload[0], payload[1], payload[2], checksum}
	_, err := port.Write(content)
	if err != nil {
		log.Fatalf("port.Write: %v", err)
	}
	fmt.Printf("Wrote %s bytes", hex.Dump(content))
}