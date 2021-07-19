package proxy_server

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"net/url"
	"os"
	"strings"
)

type HTTP struct {
	logger  *log.Logger
	address string
}

func NewHTTP() *HTTP {
	return &HTTP{
		address: ":8080",
		logger:  log.New(os.Stdout, "", 0),
	}
}
func (h *HTTP) Run() error{
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	if l, err := net.Listen("tcp", h.address); err == nil {
		for {
			client, err := l.Accept()
			if err != nil {
				log.Panic(err)
			}

			go h.handleClientRequest(client)
		}
	} else {
		log.Panic(err)
		return err
	}
	return nil

}
func (h *HTTP) handleClientRequest(client net.Conn) {
	if client == nil {
		return
	}
	defer client.Close()

	var b [1024]byte
	n, err := client.Read(b[:])
	if err != nil {
		log.Println(err)
		return
	}
	var method, host, address string
	fmt.Sscanf(string(b[:bytes.IndexByte(b[:], '\n')]), "%s%s", &method, &host)
	hostPortURL, err := url.Parse(host)
	if err != nil {
		log.Println(err)
		return
	}

	if hostPortURL.Opaque == "443" {
		address = hostPortURL.Scheme + ":443"
	} else {
		if strings.Index(hostPortURL.Host, ":") == -1 {
			address = hostPortURL.Host + ":80"
		} else {
			address = hostPortURL.Host
		}
	}

	if server, err := net.Dial("tcp", address); err == nil {
		if method == "CONNECT" {
			client.Write([]byte("HTTP/1.1 200 Connection established\r\n\r\n"))
		} else {
			server.Write(b[:n])
		}
		go io.Copy(server, client)
		io.Copy(client, server)
	} else {
		log.Println(err)
		return
	}

}
