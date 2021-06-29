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

type Database struct {
	data []SatelliteConnection
}
type SatelliteConnection struct {
	name     string
	address  string
	port     int
	data     data
	commChan chan data
}
type data struct {
	packets []Packet
}

type Packet struct {
	UnixTimestamp int64
	telemetryID   uint16
	value         float32
}

func main() {
	log.Print("Hello log")
	log.SetLevel(log.DebugLevel)

	sat1 := SatelliteConnection{"StringSatellite", "localhost", 8000, data{}, make(chan data)}
	go connectToSatellite(sat1)

	sat2 := SatelliteConnection{"BinarySatellite", "localhost", 8001, data{}, make(chan data)}
	go connectToSatellite(sat2)
	db := Database{}
	db.data = append(db.data, sat1)
	db.data = append(db.data, sat2)
	Chart(db)
	// for {
	// 	select {
	// 	case boop := <-sat1.commChan:
	// 		log.Debug("Booped", boop)

	// 	case boop := <-sat2.commChan:
	// 		log.Debug("Booped", boop)
	// 	}
	// }
}

func handleErr(e error) {
	log.Error(e)
}

func (p *Packet) proccessString(s string) {
	if utf8.ValidString(s) {
		s = strings.Trim(s, "[]")
		values := strings.Split(s, ":")

		UnixTimestamp, err := strconv.ParseInt(values[0], 10, 64)
		if err != nil {
			handleErr(err)
			return
		}
		p.UnixTimestamp = UnixTimestamp

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
		// println(time.Unix(p.UnixTimestamp, 0).String(), p.telemetryID, p.value)
	} else {
		log.Debug("Not valid utf-8")
	}
}

func connectToSatellite(s SatelliteConnection) {
	socketAddress := s.address + ":" + strconv.Itoa(s.port)
	log.Debug("Attempting to contact ", socketAddress)

	conn, err := net.Dial("tcp", socketAddress)
	if err != nil {
		log.Error(err)
		return
	}

	defer conn.Close()
	if s.name == "StringSatellite" {
		for {
			message, err := bufio.NewReader(conn).ReadString(']')
			if err != nil {
				handleErr(err)
				return
			}
			var p Packet
			p.proccessString(message)
			s.data.packets = append(s.data.packets, p)
			// log.Debug(s.data.packets)
			s.commChan <- s.data
			// Channel <- (dataStore.packets)
		}
	} else if s.name == "BinarySatellite" {
		for {
			buff := bufio.NewReader(conn)
			var p Packet
			p.decodeBinary(buff)
			s.data.packets = append(s.data.packets, p)
			// log.Debug(s.data.packets)
			s.commChan <- s.data
		}
	}
}

func (p *Packet) decodeBinary(buff *bufio.Reader) {
	header := make([]byte, 4)
	fields := []interface{}{header, &p.UnixTimestamp, &p.telemetryID, &p.value}
	for _, field := range fields {
		err := binary.Read(buff, binary.LittleEndian, field)
		if err != nil {
			handleErr(err)
			return
		}
	}
}
