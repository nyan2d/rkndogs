package main

import (
	"fmt"
	"log"
	"net"
	"regexp"
	"strconv"
)

type Address struct {
	IP   net.IP
	Mask int
}

func (a *Address) Contains(addr Address) bool {
	_, netw, err := net.ParseCIDR(a.IP.String() + "/" + strconv.Itoa(a.Mask))
	if err != nil {
		return false
	}
	return netw.Contains(addr.IP)
}

func main() {
	config, err := ReadConfig("test.yaml")
	if err != nil {
		log.Fatal(err)
	}

	entries := ParseEntries(config)

	Simplify(entries)

	for _, ip := range entries {
		fmt.Println(ip)
	}
}

func Simplify(addresses []Address) {
	if len(addresses) < 2 {
		return
	}

	// sort.Slice(addresses, func(i, j int) bool {
	// 	return bytes.Compare(addresses[i].IP, addresses[j].IP) > 0
	// })
	done := false
	for !done {
		done = true
		for i := 0; i < len(addresses); i++ {
			isbreak := false
			for j := 0; j < len(addresses); i++ {
				if i == j {
					continue
				}
				if addresses[i].Contains(addresses[j]) {
					done = false
					addresses = remove(addresses, j)
					isbreak = true
					break
				}
			}
			if isbreak {
				break
			}
		}
	}
}

func ParseEntries(config []Config) []Address {
	result := []Address{}
	for _, conf := range config {
		if IsIP(conf.Address) {
			x := Address{
				IP:   net.ParseIP(conf.Address).To4(),
				Mask: conf.Mask,
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
				Mask: conf.Mask,
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

func remove[T any](index []T, s int) []T {
	return append(index[:s], index[s+1:]...)
}
