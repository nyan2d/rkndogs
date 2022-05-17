package main

import (
	"fmt"
	"log"
	"net"
	"regexp"
)

type Address struct {
	IP   net.IP
	Mask net.IPMask
}

func main() {
	config, err := ReadConfig("test.yaml")
	if err != nil {
		log.Fatal(err)
	}

	entries := ParseEntries(config)

	for _, ip := range entries {
		fmt.Println(ip)
	}
}

func Simplify(adresses []Address) {
}

func IpToInt(ip net.IP) int {
	return (int(ip[0]) << 24) | (int(ip[1]) << 16) | (int(ip[2]) << 8) | int(ip[3])
}

func IntToIp(ip int) net.IP {
	return net.IP([]byte{
		byte((ip >> 24) & 0xFF),
		byte((ip >> 16) & 0xFF),
		byte((ip >> 8) & 0xFF),
		byte(ip & 0xFF),
	})
}

func ParseEntries(config []Config) []Address {
	result := []Address{}
	for _, conf := range config {
		if IsIP(conf.Address) {
			x := Address{
				IP:   net.ParseIP(conf.Address).To4(),
				Mask: net.CIDRMask(conf.Mask, 32),
			}
			result = append(result, x)
			continue
		}
		ips, err := Resolve(conf.Address)
		if err != nil {
			log.Printf("problem with resolving: %s\n", conf.Address)
			continue
		}
		for _, ip := range ips {
			x := Address{
				IP:   ip.To4(),
				Mask: net.CIDRMask(conf.Mask, 32),
			}
			result = append(result, x)
		}
	}
	return result
}

func IsIP(address string) bool {
	rg := regexp.MustCompile(`^((25[0-5]|(2[0-4]|1\d|[1-9]|)\d)\.){3}(25[0-5]|(2[0-4]|1\d|[1-9]|)\d)$`)
	return rg.MatchString(address)
}
