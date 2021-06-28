package main

import (
	"bufio"
	"encoding/binary"
	"net"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	log "github.com/sirupsen/logrus"
)

type SatelliteConnection struct {
	name    string
	address string
	port    int
}

type packet struct {
	unixTimestamp int64
	telemetryID   uint16
	value         float32
}

func main() {
	log.Print("Hello log")

	sat1 := SatelliteConnection{"BinarySatellite", "localhost", 8001}
	go tcpExampleBinary(sat1)

	sat2 := SatelliteConnection{"StringSatellite", "localhost", 8000}
	tcpExampleString(sat2)

}

func handleErr(e error) {
	log.Error(e)
}

func tcpExampleString(s SatelliteConnection) {
	socketAddress := s.address + ":" + strconv.Itoa(s.port)
	log.Info("Attempting to contact ", socketAddress)
	c, err := net.Dial("tcp", socketAddress)
	if err != nil {
		handleErr(err)
		return
	}
	defer c.Close()

	for {
		message, err := bufio.NewReader(c).ReadString(']')
		if err != nil {
			handleErr(err)
			return
		}
		proccessString(message)
	}
}

func tcpExampleBinary(s SatelliteConnection) {
	socketAddress := s.address + ":" + strconv.Itoa(s.port)
	log.Info("Attempting to contact ", socketAddress)

	conn, err := net.Dial("tcp", socketAddress)
	if err != nil {
		log.Error(err)
		return
	}
	defer conn.Close()

	for {
		buff := bufio.NewReader(conn)

		var p packet
		header := make([]byte, 4)
		err := binary.Read(buff, binary.LittleEndian, header)
		err = binary.Read(buff, binary.LittleEndian, &p.unixTimestamp)
		err = binary.Read(buff, binary.LittleEndian, &p.telemetryID)
		err = binary.Read(buff, binary.LittleEndian, &p.value)
		println(time.Unix(p.unixTimestamp, 0).String(), p.telemetryID, p.value)

		if err != nil {
			log.Error(err)
			return
		}

	}
}

func proccessString(s string) {
	if utf8.ValidString(s) {
		s = strings.Trim(s, "[]")
		values := strings.Split(s, ":")
		var p packet

		unixTimestamp, err := strconv.ParseInt(values[0], 10, 64)
		if err != nil {
			handleErr(err)
			return
		}
		p.unixTimestamp = unixTimestamp

		telIDSigned, err := strconv.ParseInt(values[1], 10, 16)
		if err != nil {
			handleErr(err)
			return
		}
		p.telemetryID = uint16(telIDSigned)

		valueFloat64, err := strconv.ParseFloat(values[2], 32)
		if err != nil {
			handleErr(err)
			return
		}
		p.value = float32(valueFloat64)
		println(time.Unix(p.unixTimestamp, 0).String(), p.telemetryID, p.value)
	} else {
		log.Info("Not valid utf-8")
	}
}

// func proccessBinary(bytes []byte) {
// 	// print(len(bytes))
// 	// for _, byte := range binaryBytes {
// 	// 	print(byte)
// 	// }
// 	var p packet
// 	header := bytes[0:4]
// 	unixTimestamp := bytes[4:12]
// 	// log.Info(unixTimestamp)
// 	telemetryID := bytes[12:14]
// 	value := bytes[14:18]
// 	println(header, unixTimestamp, telemetryID, value)

// 	unixTimestampUint := binary.LittleEndian.Uint64(unixTimestamp)
// 	p.unixTimestamp = int64(unixTimestampUint)

// 	p.telemetryID = binary.LittleEndian.Uint16(telemetryID)

// 	// valueUint32 := binary.LittleEndian.Uint32(value)

// 	// p.value = (valueFloat32)
// 	println(binary.LittleEndian.Uint16(header), p.unixTimestamp, time.Unix(p.unixTimestamp, 0).String(), p.telemetryID, p.value)

// 	println()
// }
