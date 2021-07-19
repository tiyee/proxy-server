package proxy_server

import (
	"io"
	"log"
	"net"
	"os"
	"strconv"
)

// Socket5 doc see: https://www.ietf.org/rfc/rfc1928.txt
type Socket5 struct {
	logger  *log.Logger
	address string
}

func NewSocket5() *Socket5 {
	return &Socket5{
		address: ":8081",
		logger:  log.New(os.Stdout, "", 0),
	}
}
func (s *Socket5) Run() error {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	if l, err := net.Listen("tcp", s.address); err == nil {
		for {
			if client, err := l.Accept(); err == nil {
				go s.handleClientRequest(client)
			}
		}
	} else {
		log.Panic(err)
		return err
	}
	return nil

}
func (s *Socket5) handleClientRequest(client net.Conn) {
	if client == nil {
		return
	}
	defer client.Close()

	var b [1024]byte
	if n, err := client.Read(b[:]); err == nil {
		if b[0] == 0x05 {
			client.Write([]byte{0x05, 0x00})
			n, err = client.Read(b[:])
			var host, port string
			switch b[3] {
			case 0x01: //IP V4 address: X'01'
				host = net.IPv4(b[4], b[5], b[6], b[7]).String()
			case 0x03: // DOMAINNAME: X'03'
				host = string(b[5 : n-2])
			case 0x04: // IP V6 address: X'04'
				host = net.IP{b[4], b[5], b[6], b[7], b[8], b[9], b[10], b[11], b[12], b[13], b[14], b[15], b[16], b[17], b[18], b[19]}.String()
			}
			port = strconv.Itoa(int(b[n-2])<<8 | int(b[n-1]))

			server, err := net.Dial("tcp", net.JoinHostPort(host, port))
			if err != nil {
				log.Println(err)
				return
			}
			defer server.Close()
			client.Write([]byte{0x05, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})
			go io.Copy(server, client)
			io.Copy(client, server)
		}
	}

}
