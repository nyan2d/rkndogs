package main

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"regexp"
	"strings"

	"golang.org/x/crypto/ssh"
)

const (
	routePattern = `(?m)^(\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3})\s+(\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3})\s+(\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3})`
)

type TomatoRoute struct {
	Host string
	Mask string
	Gate string
}

func LoadRoutes(d DeviceConfig) []TomatoRoute {
	result := make([]TomatoRoute, 0)
	resp, err := ExecuteSshCommand(d, "route -n")
	if err != nil {
		return result
	}

	rg := regexp.MustCompile(routePattern)
	lines := strings.Split(resp, "\n")
	for _, line := range lines {
		if !rg.MatchString(line) {
			continue
		}
		submatches := rg.FindStringSubmatch(line)
		result = append(result, TomatoRoute{
			Host: submatches[1],
			Mask: submatches[3],
			Gate: submatches[2],
		})
	}
	return result
}

func ClearRoutes(d DeviceConfig, routes []TomatoRoute) {
	commands := []string{}
	for _, route := range routes {
		command := fmt.Sprintf("route del -net %v netmask %v gw %v", route.Host, route.Mask, route.Gate)
		commands = append(commands, command)
	}

	if _, err := ExecuteSshCommands(d, commands); err != nil {
		log.Println(err)
	} else {
		log.Println("routing table cleared")
	}
}

func PushRoutes(d DeviceConfig, routes []TomatoRoute) {
	commands := []string{}
	for _, v := range routes {
		command := fmt.Sprintf("route add -net %v netmask %v gw %v", v.Host, v.Mask, v.Gate)
		commands = append(commands, command)
	}

	if _, err := ExecuteSshCommands(d, commands); err != nil {
		log.Println(err)
	} else {
		log.Println("routing table successfully updated")
	}
}

func ExecuteSshCommand(c DeviceConfig, command string) (string, error) {
	result, err := ExecuteSshCommands(c, []string{command})
	if err != nil {
		return "", err
	}
	return result[0], nil
}

func ExecuteSshCommands(c DeviceConfig, commands []string) ([]string, error) {
	config := &ssh.ClientConfig{
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		User:            c.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(c.Password),
		},
	}

	result := []string{}

	client, err := ssh.Dial("tcp", net.JoinHostPort(c.Host, "22"), config)
	if err != nil {
		return result, err
	}
	defer client.Close()

	for _, v := range commands {
		session, err := client.NewSession()
		if err != nil {
			return result, err
		}
		defer session.Close()

		var outBuffer bytes.Buffer
		session.Stdout = &outBuffer

		err = session.Run(v)
		if err != nil {
			return result, err
		}

		result = append(result, outBuffer.String())
	}

	return result, nil
}
