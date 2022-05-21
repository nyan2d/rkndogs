package offline

import (
	"log"
	"net"
	"os"
	"regexp"

	"github.com/nyan2d/rkndogs/dns"
	"github.com/nyan2d/rkndogs/ssh"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Device ssh.DeviceConfig `yaml:"device"`
	Gate   string           `yaml:"gate"`
	Routes []RouteConfig    `yaml:"routes"`
}

type RouteConfig struct {
	Address string `yaml:"addr"`
	Mask    int    `yaml:"mask"`
}

func ReadConfig(path string) (Config, error) {
	var obj Config

	buf, err := os.ReadFile(path)
	if err != nil {
		return obj, err
	}

	err = yaml.Unmarshal(buf, &obj)
	if err != nil {
		return obj, err
	}

	ips := make([]RouteConfig, 0)
	for _, v := range obj.Routes {
		if v.isIP() {
			ips = append(ips, v)
		} else {
			addresses, err := dns.Resolve(v.Address)
			if err != nil {
				log.Printf("resolving %v: %v\n", v.Address, err)
				continue
			}
			for _, address := range addresses {
				ips = append(ips, RouteConfig{
					Address: address.String(),
					Mask:    v.Mask,
				})
			}
		}
	}
	obj.Routes = ips

	return obj, nil
}

func (rc RouteConfig) isIP() bool {
	rg := regexp.MustCompile(`^((25[0-5]|(2[0-4]|1\d|[1-9]|)\d)\.){3}(25[0-5]|(2[0-4]|1\d|[1-9]|)\d)$`)
	return rg.MatchString(rc.Address)
}

func (rc RouteConfig) IP() net.IP {
	return net.ParseIP(rc.Address)
}

func (rc RouteConfig) IPMask() net.IPMask {
	return net.CIDRMask(rc.Mask, 32)
}

func (rc RouteConfig) Contains(c RouteConfig) bool {
	cidr := net.IPNet{
		IP:   rc.IP(),
		Mask: rc.IPMask(),
	}
	return cidr.Contains(c.IP())
}
