package main

import (
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	tsUpdated = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "tailscale_last_updated",
			Help: "Timestamp of when Tailscale data was last updated from the API.",
		},
		[]string{},
	)

	tsDeviceInfo = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "tailscale_device_info",
			Help: "Information about the tailscale device.",
		},
		[]string{"name", "id", "external", "hostname"},
	)

	tsUpgradeAvailable = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "tailscale_device_upgrade_available",
			Help: "Whether this device has an update available.",
		},
		[]string{"name", "id"},
	)

	tsLastSeen = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "tailscale_device_last_seen",
			Help: "The last time this device was seen.",
		},
		[]string{"name", "id"},
	)

	tsExpires = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "tailscale_device_expires",
			Help: "When this device's key will expire.",
		},
		[]string{"name", "id"},
	)

	tsBlocksIncoming = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "tailscale_device_blocks_incoming",
			Help: "Whether this device blocks incoming connections.",
		},
		[]string{"name", "id"},
	)

	tsExternal = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "tailscale_device_external",
			Help: "Whether this device is external to this Tailnet.",
		},
		[]string{"name", "id"},
	)
)

type metrics struct {
	registry *prometheus.Registry
}

func (metrics *metrics) registerMetrics() {
	metrics.registry.MustRegister(tsUpdated)
	metrics.registry.MustRegister(tsDeviceInfo)
	metrics.registry.MustRegister(tsUpgradeAvailable)
	metrics.registry.MustRegister(tsLastSeen)
	metrics.registry.MustRegister(tsExpires)
	metrics.registry.MustRegister(tsBlocksIncoming)
	metrics.registry.MustRegister(tsExternal)
}

func (metrics *metrics) generateMetrics(device Device) {

	// create a new registry instance to ensure old data is cleaned up
	metrics.registry = prometheus.NewRegistry()
	metrics.registerMetrics()

	tsUpdated.WithLabelValues([]string{}...).Set(float64(time.Now().Unix()))

	tsDeviceInfo.WithLabelValues(device.Name, device.ID, strconv.FormatBool(device.External), device.Hostname).Set(1)
	tsUpgradeAvailable.WithLabelValues(device.Name, device.ID).Set(b2f(device.UpdateAvailable))
	tsLastSeen.WithLabelValues(device.Name, device.ID).Set(float64(device.LastSeen.Unix()))
	tsExpires.WithLabelValues(device.Name, device.ID).Set(float64(device.Expires.Unix()))
	tsBlocksIncoming.WithLabelValues(device.Name, device.ID).Set(b2f(device.BlocksIncomingConnections))
	tsExternal.WithLabelValues(device.Name, device.ID).Set(b2f(device.External))

}

func (metrics *metrics) metricsHandler() http.Handler {
	return promhttp.HandlerFor(metrics.registry, promhttp.HandlerOpts{})
}

func b2f(b bool) float64 {
	set := 1.0
	if !b {
		set = 0.0
	}
	return set

}
