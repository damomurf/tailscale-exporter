package main

import (
	"net/http"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	tsDeviceInfo = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "tailscale_device_info",
			Help: "Information about the tailscale device.",
		},
		[]string{"name", "id", "external", "hostname"},
	)

	tsUpgradeAvailable = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "tailscale_upgrade_available",
			Help: "Whether this device has an update available.",
		},
		[]string{"name", "id"},
	)

	tsLastSeen = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "tailscale_last_seen",
			Help: "The last time this device was seen.",
		},
		[]string{"name", "id"},
	)

	tsExpires = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "tailscale_expires",
			Help: "When this device's key will expire.",
		},
		[]string{"name", "id"},
	)

	tsBlocksIncoming = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "tailscale_blocks_incoming",
			Help: "Whether this device blocks incoming connections.",
		},
		[]string{"name", "id"},
	)

	tsExternal = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "tailscale_external",
			Help: "Whether this device is external to this Tailnet.",
		},
		[]string{"name", "id"},
	)
)

func registerMetrics() {
	prometheus.MustRegister(tsDeviceInfo)
	prometheus.MustRegister(tsUpgradeAvailable)
	prometheus.MustRegister(tsLastSeen)
	prometheus.MustRegister(tsExpires)
	prometheus.MustRegister(tsBlocksIncoming)
	prometheus.MustRegister(tsExternal)
}

func generateMetrics(device Device) {

	tsDeviceInfo.WithLabelValues(device.Name, device.ID, strconv.FormatBool(device.External), device.Hostname).Set(1)
	tsUpgradeAvailable.WithLabelValues(device.Name, device.ID).Set(b2f(device.UpdateAvailable))
	tsLastSeen.WithLabelValues(device.Name, device.ID).Set(float64(device.LastSeen.Unix()))
	tsExpires.WithLabelValues(device.Name, device.ID).Set(float64(device.Expires.Unix()))
	tsBlocksIncoming.WithLabelValues(device.Name, device.ID).Set(b2f(device.BlocksIncomingConnections))
	tsExternal.WithLabelValues(device.Name, device.ID).Set(b2f(device.External))
}

func metricsHandler() http.Handler {
	return promhttp.Handler()
}
