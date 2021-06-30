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
	satellite []SatelliteConnection
}
type SatelliteConnection struct {
	name     string
	address  string
	port     int
	data     data
	commChan chan Packet
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
	log.SetLevel(log.InfoLevel)

	sat1 := SatelliteConnection{"StringSatellite", "localhost", 8000, data{}, make(chan Packet)}
	go connectToSatellite(sat1)

	sat2 := SatelliteConnection{"BinarySatellite", "localhost", 8001, data{}, make(chan Packet)}
	go connectToSatellite(sat2)

	db := Database{}
	db.satellite = append(db.satellite, sat1)
	db.satellite = append(db.satellite, sat2)
	go db.receiver()
	Chart(db)
}

func (d Database) receiver() {
	for {
		select {
		case newPacket := <-d.satellite[0].commChan:
			log.Debug("ONE", newPacket)
			d.satellite[0].data.packets = append(d.satellite[0].data.packets, newPacket)
		case newPacket := <-d.satellite[1].commChan:
			log.Debug("TWO", newPacket)
			d.satellite[1].data.packets = append(d.satellite[1].data.packets, newPacket)
		}
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
			log.Debug(s.data.packets)
			s.commChan <- p
		}
	} else if s.name == "BinarySatellite" {
		for {
			buff := bufio.NewReader(conn)
			var p Packet
			p.decodeBinary(buff)
			log.Debug(s.data.packets)
			s.commChan <- p
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

func handleErr(e error) {
	log.Error(e)
}
