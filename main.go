package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/MauveSoftware/flan_exporter/datasource/filesystem"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

const version string = "0.1.0"

var (
	showVersion   = flag.Bool("version", false, "Print version information.")
	listenAddress = flag.String("web.listen-address", ":9999", "Address on which to expose metrics and web interface.")
	metricsPath   = flag.String("web.telemetry-path", "/metrics", "Path under which to expose metrics.")
)

func main() {
	flag.Parse()

	if *showVersion {
		printVersion()
		os.Exit(0)
	}

	startServer()
}

func printVersion() {
	fmt.Println("flan_exporter")
	fmt.Printf("Version: %s\n", version)
	fmt.Println("Author(s): Daniel Czerwonk")
	fmt.Println("Copyright: 2020, Mauve Mailorder Software GmbH & Co. KG, Licensed under MIT license")
	fmt.Println("Metric exporter for Cloudflare Flan Scan report results (https://github.com/cloudflare/flan)")
}

func startServer() {
	ds := filesystem.New("/Users/daniel/xml_files")

	logrus.Infof("Starting Cloudflare Flan Scan exporter (Version: %s)", version)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
			<head><title>Flan Scan Result Exporter (Version ` + version + `)</title></head>
			<body>
			<h1>Senderscore Exporter by Mauve Mailorder Software</h1>
			<h2>Metrics</h2>
			<p><a href="/metrics">here</a></p>
			<h2>More information</h2>
			<p><a href="https://github.com/MauveSoftware/flan_exporter">github.com/MauveSoftware/flan_exporter</a></p>
			</body>
			</html>`))
	})

	c := &collector{dataSource: ds}
	prometheus.MustRegister(c)
	http.Handle("/metrics", promhttp.Handler())

	logrus.Infof("Listening for %s on %s", *metricsPath, *listenAddress)
	logrus.Fatal(http.ListenAndServe(*listenAddress, nil))
}
