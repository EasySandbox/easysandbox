package main

import (
	"io"
	"net"
	"net/http"
	"os"
	"strings"
	"fmt"
	"time"
)

func main() {

	if len(os.Args) != 2 {
		fmt.Println("Usage:", os.Args[0], "<command>")
		os.Exit(1)
	}

	server := "hostsystem"

	hostname, _ := os.Hostname()
	ip, _ := getOurIP()

	url := fmt.Sprintf("http://%s:8080/add?domain=%s&ip=%s", server, hostname, ip)

	for {

		resp, err := http.Post(url, "text/plain", nil)

		if err != nil {
			time.Sleep(1 * time.Second)
			continue
		}

		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			body, err := io.ReadAll(resp.Body)
			if err == nil {
				if strings.HasPrefix(string(body), "Added IP") {
					time.Sleep(5 * time.Minute)
					continue
				}
			}
		}
		time.Sleep(1 * time.Second)

	}

}

func getOurIP() (string, error) {
	conn, err := net.Dial("udp", "1.1.1.1:80") // placeholder IP address, no network activity is actually done
	if err != nil {
		return "", err
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String(), nil
}
