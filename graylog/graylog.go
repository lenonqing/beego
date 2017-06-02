package graylog

import (
	"bytes"
	"compress/zlib"
	"crypto/rand"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net"
	"os"
	"strconv"

	"github.com/astaxie/beego"
)

const (
	defaultGraylogPort     = 12201
	defaultGraylogHost     = "127.0.0.1"
	defaultConnectionType  = "wan"
	defaultMaxWanChunkSize = 1420
	defaultMaxLanChunkSize = 8154

	infoLevel  = 6
	errorLevel = 3
)

// Config Graylog配置信息
type Config struct {
	GraylogPort     int
	GraylogHost     string
	AppName         string
	ConnectionType  string
	MaxLanChunkSize int
	MaxWanChunkSize int
}

type gelf struct {
	config Config
}

var g *gelf

// Init 初始化.
func Init() {
	port, _ := beego.AppConfig.Int(`graylogPort`)
	config := &Config{
		GraylogPort: port,
		GraylogHost: beego.AppConfig.String(`graylogHost`),
		AppName:     beego.AppConfig.String(`appname`),
	}
	if config.GraylogPort == 0 {
		config.GraylogPort = defaultGraylogPort
	}
	if config.GraylogHost == "" {
		config.GraylogHost = defaultGraylogHost
	}
	if config.ConnectionType == "" {
		config.ConnectionType = defaultConnectionType
	}
	if config.MaxLanChunkSize == 0 {
		config.MaxLanChunkSize = defaultMaxLanChunkSize
	}
	if config.MaxWanChunkSize == 0 {
		config.MaxWanChunkSize = defaultMaxWanChunkSize
	}
	g = &gelf{
		config: *config,
	}
}

func (g *gelf) compress(b []byte) bytes.Buffer {
	var buf bytes.Buffer
	comp := zlib.NewWriter(&buf)

	comp.Write(b)
	comp.Close()

	return buf
}

func (g *gelf) createChunkedMessage(index int, chunkCountInt int, id []byte, compressed *bytes.Buffer) bytes.Buffer {
	var packet bytes.Buffer

	chunksize := g.getChunksize()

	packet.Write(g.intToBytes(30))
	packet.Write(g.intToBytes(15))
	packet.Write(id)

	packet.Write(g.intToBytes(index))
	packet.Write(g.intToBytes(chunkCountInt))

	packet.Write(compressed.Next(chunksize))

	return packet
}

func (g *gelf) getChunksize() int {

	if g.config.ConnectionType == "wan" {
		return g.config.MaxWanChunkSize
	}

	if g.config.ConnectionType == "lan" {
		return g.config.MaxLanChunkSize
	}

	return g.config.MaxWanChunkSize
}

func (g *gelf) intToBytes(i int) []byte {
	buf := new(bytes.Buffer)

	err := binary.Write(buf, binary.LittleEndian, int8(i))
	if err != nil {
		log.Printf("Uh oh! %s", err)
	}
	return buf.Bytes()
}

func (g *gelf) send(b []byte) {
	var addr = g.config.GraylogHost + ":" + strconv.Itoa(g.config.GraylogPort)
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		log.Printf("Uh oh! %s", err)
		return
	}
	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		log.Printf("Uh oh! %s", err)
		return
	}
	conn.Write(b)
}

func (g *gelf) log(level uint, message string, data map[string]interface{}) {
	hostname, _ := os.Hostname()
	obj := map[string]interface{}{
		"version":       "1.0",
		"host":          g.config.AppName,
		"_server":       hostname,
		"short_message": message,
		"level":         level,
	}
	for k, v := range data {
		obj[fmt.Sprint("_", k)] = v
	}
	msg, err := json.Marshal(obj)
	if err != nil {
		log.Println(err)
	}
	compressed := g.compress(msg)

	chunksize := g.config.MaxWanChunkSize
	length := compressed.Len()

	if length > chunksize {
		chunkCountInt := length / chunksize
		chunkCount := math.Ceil(float64(chunkCountInt))
		chunkCountInt = int(chunkCount) + 1

		id := make([]byte, 8)
		rand.Read(id)

		for i, index := 0, 0; i < length; i, index = i+chunksize, index+1 {
			packet := g.createChunkedMessage(index, chunkCountInt, id, &compressed)
			g.send(packet.Bytes())
		}

	} else {
		g.send(compressed.Bytes())
	}
	log.Println(string(msg))
}

// Info 打印Info信息
func Info(message string, data map[string]interface{}) {
	g.log(infoLevel, message, data)
}

// Error 打印错误信息
func Error(message string, data map[string]interface{}) {
	g.log(errorLevel, message, data)
}
