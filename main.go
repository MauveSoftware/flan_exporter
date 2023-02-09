// SPDX-FileCopyrightText: (c) Mauve Mailorder Software GmbH & Co. KG, 2020. Licensed under [MIT](LICENSE) license
//
// SPDX-License-Identifier: MIT

package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"

	"github.com/MauveSoftware/flan_exporter/datasource"
	"github.com/MauveSoftware/flan_exporter/datasource/filesystem"
	"github.com/MauveSoftware/flan_exporter/datasource/gcloud"
)

const version string = "0.2.3"

var (
	showVersion           = flag.Bool("version", false, "Print version information.")
	listenAddress         = flag.String("web.listen-address", ":9711", "Address on which to expose metrics and web interface.")
	metricsPath           = flag.String("web.telemetry-path", "/metrics", "Path under which to expose metrics.")
	dataSourceProvider    = flag.String("datasource.provider", "fs", "Data source provider (gcloud for Google Cloud Storage or fs for local filesystem)")
	fsPath                = flag.String("datasource.fs.report-path", "", "Path to report files")
	gcloudCredentialsFile = flag.String("datasource.gcloud.credentials-path", "", "Path to Google Cloud Credentials JSON file")
	gcloudBucketName      = flag.String("datasource.gcloud.bucket-name", "flan-reports", "Name ")
	tlsEnabled            = flag.Bool("tls.enabled", false, "Enables TLS")
	tlsCertChainPath      = flag.String("tls.cert-file", "", "Path to TLS cert file")
	tlsKeyPath            = flag.String("tls.key-file", "", "Path to TLS key file")
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
	ds, err := datasourceProvider()
	if err != nil {
		panic(err)
	}

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

	logrus.Infof("Listening for %s on %s (TLS: %v)", *metricsPath, *listenAddress, *tlsEnabled)
	if *tlsEnabled {
		logrus.Fatal(http.ListenAndServeTLS(*listenAddress, *tlsCertChainPath, *tlsKeyPath, nil))
		return
	}

	logrus.Fatal(http.ListenAndServe(*listenAddress, nil))
}

func datasourceProvider() (datasource.DataSource, error) {
	switch *dataSourceProvider {
	case "fs":
		return filesystem.New(*fsPath), nil
	case "gcloud":
		return gcloud.New(*gcloudBucketName, *gcloudCredentialsFile)
	default:
		return nil, errors.Errorf("data source is unknown")
	}
}
