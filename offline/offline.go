package offline

import (
	"log"
	"net"

	"github.com/nyan2d/rkndogs/ssh"
	"github.com/nyan2d/rkndogs/util"
)

func Do(configPath string) {
	config, err := ReadConfig(configPath)
	if err != nil {
		log.Fatal("config reading:", config)
	}

	routes := ssh.LoadRoutes(config.Device)
	log.Println("routes loaded")
	oldRoutes := util.Filter(routes, func(item ssh.SshRoute) bool {
		return item.Gate == config.Gate
	})
	err = ssh.RemoveRoutes(config.Device, oldRoutes)
	if err != nil {
		log.Fatal("removing routes:", err)
	}
	log.Println("routes removed")

	config.Routes = reduceRoutes(config.Routes)
	sshRoutes := routesToSshRoutes(config.Routes, config.Gate)
	err = ssh.PushRoutes(config.Device, sshRoutes)
	if err != nil {
		log.Fatal("updating routes:", err)
	}
	log.Println("routes updated")
}

// TODO: there is quadratic complexity rn, and it would be great to reduce it
func reduceRoutes(routes []RouteConfig) []RouteConfig {
	if len(routes) < 2 {
		return routes
	}
	r := make([]RouteConfig, len(routes))
	copy(r, routes)

	done := false
	for !done {
		done = true
		for i := 0; i < len(r); i++ {
			isbreak := false
			for j := len(r) - 1; j >= 0; j-- {
				if i == j {
					continue
				}
				if r[i].Contains(r[j]) {
					done = false
					r = util.Delete(r, j)
					isbreak = true
					break
				}
			}
			if isbreak {
				break
			}
		}
	}

	return r
}

func routesToSshRoutes(routes []RouteConfig, gate string) []ssh.SshRoute {
	r := make([]ssh.SshRoute, 0)
	for _, v := range routes {
		route := ssh.SshRoute{
			Host: v.IP().String(),
			Mask: net.IP(v.IPMask()).String(),
			Gate: gate,
		}
		r = append(r, route)
	}
	return r
}
