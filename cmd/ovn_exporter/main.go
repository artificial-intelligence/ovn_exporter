package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/go-kit/log/level"
	ovn "github.com/greenpau/ovn_exporter/pkg/ovn_exporter"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	var listenAddress string
	var metricsPath string
	var pollTimeout int
	var pollInterval int
	var isShowVersion bool
	var logLevel string
	var systemRunDir string
	var (
		disableVswitchMetrics   bool
		disableNorthboundMetrics bool
		disableSouthboundMetrics bool
	)

	flag.BoolVar(&disableVswitchMetrics, "disable-vswitch-metrics", false, "Disable collection of vswitch metrics.")
	flag.BoolVar(&disableNorthboundMetrics, "disable-northbound-metrics", false, "Disable collection of northbound metrics.")
	flag.BoolVar(&disableSouthboundMetrics, "disable-southbound-metrics", false, "Disable collection of southbound metrics.")

	opts := ovn.Options{
		Timeout: pollTimeout,
		Logger:  logger,
		Features: map[string]bool{
			"vswitch_metrics":   !disableVswitchMetrics,
			"northbound_metrics": !disableNorthboundMetrics,
			"southbound_metrics": !disableSouthboundMetrics,
		},
	}

	exporter, err := ovn.NewExporter(opts)
	if err != nil {
		level.Error(logger).Log(
			"msg", "failed to init properly",
			"error", err.Error(),
		)
		os.Exit(1)
	}
	prometheus.MustRegister(exporter)

	http.Handle(metricsPath, promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
             <head><title>OVN Exporter</title></head>
             <body>
             <h1>OVN Exporter</h1>
             <p><a href='` + metricsPath + `'>Metrics</a></p>
             </body>
             </html>`))
	})

	level.Info(logger).Log("listen_on ", listenAddress)
	if err := http.ListenAndServe(listenAddress, nil); err != nil {
		level.Error(logger).Log(
			"msg", "listener failed",
			"error", err.Error(),
		)
		os.Exit(1)
	}
}
