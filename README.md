[![CircleCI](https://circleci.com/gh/MauveSoftware/flan_exporter.svg?style=shield)](https://circleci.com/gh/MauveSoftware/flan_exporter)
[![Go Report Card](https://goreportcard.com/badge/github.com/mauvesoftware/flan_exporter)](https://goreportcard.com/report/github.com/mauvesoftware/flan_exporter)
[![Docker Build Statu](https://img.shields.io/docker/build/MauveSoftware/flan_exporter.svg)](https://hub.docker.com/r/MauveSoftware/flan_exporter/builds)

# flan_exporter
Metrics exporter for Cloudflares Flan Scan vulnerability / network scanner

This is not an offical Cloudflare project.

## Install
```
go get -u github.com/MauveSoftware/flan_exporter
```

## Usage

### Binary
```bash
./flan_exporter -datasource.fs.report-path=/opt/flan/shared/reports
```

### Docker
```bash
docker run -d --restart always --name flan_exporter -v /reports:/opt/flan/shared/reports -p 9711:9711 mauvesoftware/flan_exporter
```

or with an Google Cloud Storage backend:

```bash
docker run -d --restart always --name flan_exporter -e FLAN_DATASOURCE=gcloud -e FLAN_GCLOUD_BUCKET_NAME=my-bucket -v /app/gcloud_credentials.json:/path/to/credentials.json -p 9711:9711 mauvesoftware/flan_exporter
```

## License
(c) Mauve Mailorder Software GmbH & Co. KG, 2020. Licensed under [MIT](LICENSE) license.

## Prometheus
see https://prometheus.io/

## Cloudflare Flan Scan
see https://github.com/cloudflare/flan
