package main

import (
	"bufio"
	"encoding/binary"
	"net"
	"strconv"
	"strings"
	"unicode/utf8"

	log "github.com/sirupsen/logrus"
)

type Data struct {
	ID   int
	data []packet
}
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

	sat1 := SatelliteConnection{"StringSatellite", "localhost", 8000}
	go connectToSatellite(sat1)

	sat2 := SatelliteConnection{"BinarySatellite", "localhost", 8001}
	connectToSatellite(sat2)
}

func handleErr(e error) {
	log.Error(e)
}

func (p *packet) proccessString(s string) {
	if utf8.ValidString(s) {
		s = strings.Trim(s, "[]")
		values := strings.Split(s, ":")

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
		// println(time.Unix(p.unixTimestamp, 0).String(), p.telemetryID, p.value)
	} else {
		log.Info("Not valid utf-8")
	}
}

func connectToSatellite(s SatelliteConnection) {
	socketAddress := s.address + ":" + strconv.Itoa(s.port)
	log.Info("Attempting to contact ", socketAddress)

	conn, err := net.Dial("tcp", socketAddress)
	if err != nil {
		log.Error(err)
		return
	}
	defer conn.Close()
	dataStore := Data{ID: 2}
	if s.name == "StringSatellite" {
		for {
			message, err := bufio.NewReader(conn).ReadString(']')
			if err != nil {
				handleErr(err)
				return
			}
			var p packet
			p.proccessString(message)
			dataStore.data = append(dataStore.data, p)
			log.Info(dataStore.data)
		}
	} else if s.name == "BinarySatellite" {
		for {
			buff := bufio.NewReader(conn)
			var p packet
			p.decodeBinary(buff)
			dataStore.data = append(dataStore.data, p)
			log.Info(dataStore.data)
		}
	}
}

func (p *packet) decodeBinary(buff *bufio.Reader) {
	header := make([]byte, 4)
	err := binary.Read(buff, binary.LittleEndian, header)
	err = binary.Read(buff, binary.LittleEndian, &p.unixTimestamp)
	err = binary.Read(buff, binary.LittleEndian, &p.telemetryID)
	err = binary.Read(buff, binary.LittleEndian, &p.value)
	if err != nil {
		handleErr(err)
		return
	}
	// println(time.Unix(p.unixTimestamp, 0).String(), p.telemetryID, p.value)
}

// // func proccessBinary(bytes []byte) {
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
