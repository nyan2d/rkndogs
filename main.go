package main

import (
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
	config, err := ReadConfig("config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	entries := ParseEntries(config.Routes)
	entries = Reduce(entries)

	routes := LoadRoutes(config.Device)
	filteredRoutes := make([]TomatoRoute, 0)
	for _, v := range routes {
		if v.Gate == config.Gate {
			filteredRoutes = append(filteredRoutes, v)
		}
	}
	ClearRoutes(config.Device, filteredRoutes)

	newRoutes := []TomatoRoute{}
	for _, v := range entries {
		mask := net.CIDRMask(v.Mask, 32)
		newRoutes = append(newRoutes, TomatoRoute{
			Host: v.IP.Mask(mask).String(),
			Mask: net.IP(mask).String(),
			Gate: config.Gate,
		})
	}
	PushRoutes(config.Device, newRoutes)
}

// TODO: there is quadratic complexity rn, and it would be great to reduce it
func Reduce(adr []Address) []Address {
	if len(adr) < 2 {
		return adr
	}
	result := make([]Address, len(adr))
	copy(result, adr)

	done := false
	for !done {
		done = true
		for i := 0; i < len(result); i++ {
			isbreak := false
			for j := len(result) - 1; j >= 0; j-- {
				if i == j {
					continue
				}
				if result[i].Contains(result[j]) {
					done = false
					result = delete(result, j)
					isbreak = true
					break
				}
			}
			if isbreak {
				break
			}
		}
	}

	return result
}

func ParseEntries(config []RouteConfig) []Address {
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

func delete[S ~[]T, T any](s S, i int) S {
	return append(s[:i], s[i+1:]...)
}
