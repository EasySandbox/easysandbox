package getourip

// this is in its own package for C requirements reasons.
// when I tried to put it in the domains package, client was missing a libvirt lib
import (
	"net"
)

func GetOurIP() (string, error) {
	conn, err := net.Dial("udp", "1.1.1.1:80") // placeholder IP address, no network activity is actually done
	if err != nil {
		return "", err
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String(), nil
}
