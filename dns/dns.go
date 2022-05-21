package dns

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
)

const (
	googleDnsUri = "https://dns.google/resolve?name=%v&type=A"
)

type dnsResponse struct {
	Status int64       `json:"Status"`
	Answer []dnsAnswer `json:"Answer"`
}

type dnsAnswer struct {
	Type int64  `json:"type"`
	Data string `json:"data"`
}

func Resolve(domain string) ([]net.IP, error) {
	ips := []net.IP{}

	resp, err := http.Get(fmt.Sprintf(googleDnsUri, domain))
	if err != nil {
		return ips, err
	}
	defer resp.Body.Close()

	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		return ips, err
	}

	var response dnsResponse
	err = json.Unmarshal(buf, &response)
	if err != nil {
		return ips, err
	}

	if response.Status != 0 {
		return ips, fmt.Errorf("dns: wrong status")
	}

	for _, answer := range response.Answer {
		if answer.Type == 1 {
			ips = append(ips, net.ParseIP(answer.Data).To4())
		}
	}

	return ips, err
}
