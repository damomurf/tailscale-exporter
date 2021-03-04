package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/pkg/errors"
)

type Tailnet struct {
	Devices []Device `json:"devices"`
}

type Device struct {
	Addresses                 []string `json:"addresses"`
	Authorized                bool     `json:"authorized"`
	BlocksIncomingConnections bool
	ClientVersion             string
	Expires                   time.Time
	Hostname                  string
	Name                      string
	ID                        string
	External                  bool `json:"isExternal"`
	KeyExpiryDisabled         bool
	LastSeen                  time.Time
	OS                        string
	UpdateAvailable           bool
	User                      string
	// This can be empty in responses and causes issues for JSON parsing:
	//Created                   *time.Time `json:"created,omitEmpty"`
}

const (
	DeviceURL = "https://api.tailscale.com/api/v2/tailnet/%s/devices"
)

var (
	addr     = flag.String("listen-address", ":8080", "The address to listen on for HTTP requests.")
	tailnet  = flag.String("tailnet", "", "The Tailscale network to query.")
	token    = flag.String("token", "", "The Tailscale API token.")
	interval = flag.Duration("interval", 5*time.Minute, "The interval to poll the Tailscale API.")
	debug    = flag.Bool("debug", false, "Debug output")
)

func b2f(b bool) float64 {
	set := 1.0
	if b {
		set = 0.0
	}
	return set

}

func debugf(fmt string, args ...interface{}) {
	if *debug {
		log.Printf(fmt, args)
	}
}

func fetchDevices(tailnet, key string) error {

	client := &http.Client{
		Timeout: time.Second * 10,
	}

	req, err := http.NewRequest("GET", fmt.Sprintf(DeviceURL, tailnet), nil)
	if err != nil {
		return errors.Wrap(err, "creating request")
	}
	req.SetBasicAuth(key, "")

	response, err := client.Do(req)
	if err != nil {
		return errors.Wrap(err, "executing request")
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusOK {

		buf, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return errors.Wrap(err, "reading response body")
		}

		tnet := &Tailnet{}

		if err := json.Unmarshal(buf, tnet); err != nil {
			return errors.Wrap(err, "parsing JSON")
		}

		for _, device := range tnet.Devices {
			generateMetrics(device)
			debugf("%+v\n", device)
		}
	} else {
		return errors.Errorf("Tailscale API request returned unexpected status code: %d - %s", response.StatusCode, response.Status)
	}

	return nil

}

func main() {

	flag.Parse()

	if *tailnet == "" || *token == "" {
		flag.Usage()
		os.Exit(1)
	}

	registerMetrics()

	if err := fetchDevices(*tailnet, *token); err != nil {
		log.Printf("Error fetching devices: %+v", err)
	}

	go func() {
		c := time.Tick(*interval)
		for range c {
			if err := fetchDevices(*tailnet, *token); err != nil {
				log.Printf("Error fetching devices: %+v", err)
			}
		}
	}()

	http.Handle("/metrics", metricsHandler())
	log.Fatal(http.ListenAndServe(*addr, nil))

}
