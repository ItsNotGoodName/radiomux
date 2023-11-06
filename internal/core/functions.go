package core

import (
	"crypto/rand"
	"encoding/base64"
	"net"
	"strings"
)

func GenerateToken() (string, error) {
	secret := make([]byte, 32)
	_, err := rand.Read(secret)
	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(secret), nil
}

func PossiblePublicIPs() ([]net.IP, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return nil, err
	}

	var ips []net.IP
	for _, addr := range addrs {
		ip := net.ParseIP(strings.Split(addr.String(), "/")[0])
		if ip.To4() != nil && ip.String() != "127.0.0.1" {
			ips = append(ips, ip)
		}
	}

	return ips, nil
}
