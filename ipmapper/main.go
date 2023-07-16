package main

import (
	//"crypto/hmac"
	//"crypto/sha256"
	//"io/ioutil"
	"fmt"
	"net"
	"net/http"
	"strings"
	"sync"

	"github.com/adrg/xdg"

	"git.voidnet.tech/kev/easysandbox/filesystem"
)

var hmacKeyPath = xdg.DataHome + "/easysandbox/hmac.key"

var ipMap struct {
	sync.RWMutex
	m map[string]net.IP
}

var hmacKey []byte

func init() {

	if exists, _ := filesystem.PathExists(hmacKeyPath); exists {
		// generate hmac key

		hmacKey = make([]byte, 32)
		// generate dummy key of length 32
		key := strings.Repeat("a", 32)
		hmacKey = []byte(key)

	}
	//hmacKey, err := ioutil.ReadFile(hmacKeyPath)

	// if err != nil {

	// }

	ipMap.m = make(map[string]net.IP)
}

func addIP(w http.ResponseWriter, r *http.Request) {
	ipStr := r.URL.Query().Get("ip")
	domainName := r.URL.Query().Get("domain")

	// // Calculate HMAC
	// h := hmac.New(sha256.New, hmacKey)
	// h.Write([]byte(ipStr))
	// mac := h.Sum(nil)

	// // Compare HMAC with request HMAC
	// providedHMAC := r.Header.Get("X-HMAC")
	// if !hmac.Equal(mac, []byte(providedHMAC)) {
	// 	w.WriteHeader(http.StatusUnauthorized)
	// 	w.Write([]byte("Invalid HMAC"))
	// 	return
	// }

	ip := net.ParseIP(ipStr)
	if ip == nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid IP"))
		return
	}

	ipMap.Lock()
	ipMap.m[domainName] = ip
	ipMap.Unlock()

	fmt.Println("Added IP:", ipStr)
	w.Write([]byte("Added IP: " + ipStr))
}

func getIP(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")

	ipMap.RLock()
	ip, ok := ipMap.m[key]
	ipMap.RUnlock()

	if !ok {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("No IP found for key: " + key))
		return
	}

	w.Write([]byte(ip.String()))
}

func main() {
	http.HandleFunc("/add", addIP)
	http.HandleFunc("/get", getIP) // add this line
	http.ListenAndServe(":8080", nil)
}
