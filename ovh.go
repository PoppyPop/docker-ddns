package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/ovh/go-ovh/ovh"
)

// DNSRecord holds the first name of the currently logged-in user.
// Visit https://api.ovh.com/console/#/me#GET for the full definition
type DNSRecord struct {
	Target    string `json:"target"`
	TTL       uint64 `json:"ttl"`
	Zone      string `json:"zone"`
	FieldType string `json:"fieldType"`
	ID        uint64 `json:"id"`
	SubDomain string `json:"subDomain"`
}

// ParamsTarget ...
type ParamsTarget struct {
	Target string `json:"target"`
}

const (
	// Time allowed to read the next pong message from the peer.
	ipFile = "/conf/ovh.ip"
)

var domain = flag.String("domain", "test.com", "Domain")
var subdomain = flag.String("subdomain", "ddns", "Subdomain")

var ak = flag.String("ak", "invalid", "Application Key")
var as = flag.String("as", "invalid", "Application Secret")
var ck = flag.String("ck", "invalid", "Consumer Key")

// GetCurrentIP ...
func GetCurrentIP() (ip string, err error) {
	resp, err := http.Get("http://ipv4.icanhazip.com")
	if err != nil {
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	return strings.Trim(string(body), "\n"), nil
}

// ControlIP ...
func ControlIP(ip string) (res bool, err error) {

	if _, errStat := os.Stat(ipFile); os.IsNotExist(errStat) {
		_, errCr := os.Create(ipFile)
		if errCr != nil {
			return false, errCr
		}
	}

	byteFile, err := ioutil.ReadFile(ipFile)
	if err != nil {
		return
	}

	return (string(byteFile) == ip), nil
}

// SaveIP ...
func SaveIP(ip string) (err error) {
	return ioutil.WriteFile(ipFile, []byte(ip), 0644)
}

func main() {
	flag.Parse()

	ip, errIP := GetCurrentIP()
	if errIP != nil {
		fmt.Println(errIP)
		return
	}
	fmt.Printf("Detected IP : %s\n", ip)

	ok, errCtrl := ControlIP(ip)
	if errCtrl != nil {
		fmt.Println(errCtrl)
		return
	}

	if !ok {
		var list []uint64

		// API key with rules :
		// GET /domain/*
		// PUT /domain/*
		client, _ := ovh.NewClient(
			"ovh-eu",
			*ak,
			*as,
			*ck,
		)

		req := "/domain/zone/" + *domain + "/record?subDomain=" + *subdomain + ""
		err := client.Get(req, &list)
		if err != nil {
			fmt.Println(err)
			return
		}

		for _, element := range list {
			var record DNSRecord
			reqDetails := "/domain/zone/" + *domain + "/record/" + strconv.FormatUint(element, 10) + ""
			err := client.Get(reqDetails, &record)

			if err != nil {
				fmt.Println(err)
				return
			}

			if record.FieldType == "A" {
				fmt.Printf("%s.%s : %s\n", record.SubDomain, record.Zone, record.Target)
				if record.Target != ip {
					params := &ParamsTarget{Target: ip}
					if err := client.Put(reqDetails, params, nil); err != nil {
						fmt.Printf("Error: %q\n", err)
						return
					}
				} else {
					fmt.Printf("IP (Remote): %s nothing to do\n", ip)
				}
			} else {
				fmt.Println("Not an A record!")
			}
		}

		// MAJ IP
		if err := SaveIP(ip); err != nil {
			fmt.Printf("Error: %q\n", err)
			return
		}
	} else {
		fmt.Printf("IP (local): %s nothing to do\n", ip)
	}
}
