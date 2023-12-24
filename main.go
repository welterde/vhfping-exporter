package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	flags "github.com/jessevdk/go-flags"
	"github.com/VictoriaMetrics/metrics"
)

var opts struct {
	Listen  string `short:"l" long:"listen" description:"Listen address" value-name:"[HOST]:PORT" default:":9605"`
	Period  uint   `short:"p" long:"period" description:"Period in seconds, should match Prometheus scrape interval" value-name:"SECS" default:"60"`
	Fping   string `short:"f" long:"fping"  description:"Fping binary path" value-name:"PATH" default:"/usr/bin/fping"`
	Count   uint   `short:"c" long:"count"  description:"Number of pings to send at each period" value-name:"N" default:"20"`
	Version bool   `long:"version" description:"Show version"`
	DumpRaw bool   `long:"dump-raw" description:"Supply raw latency metrics in addition to histogram"`
}

var (
	buildVersion = "dev"
	buildCommit  = "none"
	buildDate    = "unknown"
)

func probeHandler(w http.ResponseWriter, r *http.Request) {
	targetParam := r.URL.Query().Get("target")
	if targetParam == "" {
		w.Header().Set("Content-Type", "text/html")
		_, _ = w.Write([]byte(`<html>
		    <head><title>VM-H FPing Exporter</title></head>
			<body>
			<b>ERROR: missing target parameter</b>
			</body>`))
		return
	}

	target := GetTarget(
		WorkerSpec{
			period: time.Second * time.Duration(opts.Period),
		},
		TargetSpec{
			host: targetParam,
		},
	)

	target.vset.WritePrometheus(w)
}

func main() {
	if _, err := flags.Parse(&opts); err != nil {
		os.Exit(0)
	}
	if opts.Version {
		fmt.Printf("vhfping-exporter %v (commit %v, built %v)\n", buildVersion, buildCommit, buildDate)
		os.Exit(0)
	}
	if _, err := os.Stat(opts.Fping); os.IsNotExist(err) {
		fmt.Printf("could not find fping at %q\n", opts.Fping)
		os.Exit(1)
	}
	http.HandleFunc("/probe", probeHandler)
	http.HandleFunc("/metrics", func(w http.ResponseWriter, req *http.Request) {
		metrics.WritePrometheus(w, true)
	})
	log.Fatal(http.ListenAndServe(opts.Listen, nil))
}
