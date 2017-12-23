package main

import (
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/aotd1/argparse"
	"github.com/joho/godotenv"
	"github.com/codingconcepts/env"
	"github.com/jacobsa/go-serial/serial"
)

type config struct {
	PortName 		string	`env:"UNIEL_PORT" required:"true"`
	BaudRate 		uint	`env:"UNIEL_BAUD_RATE" default:"9600"`
	DataBits 		uint	`env:"UNIEL_DATA_BITS" default:"8"`
	StopBits 		uint	`env:"UNIEL_STOP_BITS" default:"1"`
	FlowControl 	bool	`env:"UNIEL_FLOW_CONTROL" default:"false"`
}

func main() {
	godotenv.Load()

	config := config{}
	if err := env.Set(&config); err != nil {
		log.Fatal(err)
	}

	p := argparse.NewParser("uniel", "Uniel UCH-* modules controlling program")
	offCmd := p.NewCommand("off", "Turn channel off")
	onCmd := p.NewCommand("on", "Turn channel on")
	address := byte(0x01)
	channel := byte(0x01)
	//address, err1 := getByte(p.String("a", "address", &argparse.Options{Required: true, Help: "Device address: 1-16"}))
	//channel, err2 := getByte(p.String("c", "channel", &argparse.Options{Required: true, Help: "Channel: 1-8"}))
	//if err1 != nil {
	//	fmt.Print(p.Usage(err1))
	//	return
	//}
	//
	//if err2 != nil {
	//	fmt.Print(p.Usage(err2))
	//	return
	//}

	if err := p.Parse(os.Args); err != nil {
		fmt.Print(p.Usage(err))
		return
	}

	port := openPort(config)
	defer port.Close()

	if onCmd.Happened() {
		run(port, address, 0x0A, [3]byte{0xFF, channel, 0x00})
	} else if offCmd.Happened() {
		run(port, address, 0x0A, [3]byte{0x00, channel, 0x00})
	} else {
		err := fmt.Errorf("bad arguments, please check usage")
		fmt.Print(p.Usage(err))
	}

}

func getByte(address *string) (byte, error) {
	switch *address {
		case "1": return 0x00, nil
		case "2": return 0x01, nil
		case "3": return 0x02, nil
		case "4": return 0x03, nil
		case "5": return 0x04, nil
		case "6": return 0x05, nil
		case "7": return 0x06, nil
		case "8": return 0x07, nil
		default: return 0x00, fmt.Errorf("bad byte argument, please check usage")
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