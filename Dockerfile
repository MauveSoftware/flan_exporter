FROM golang:1.27.1-alpine3.24@sha256:cf6fca6641884b8433441b2b0652976f975e1d0fdd26d177eaaf8596087f3125 AS builder
ADD . /go/flan_exporter/
WORKDIR /go/flan_exporter
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /go/bin/flan_exporter

FROM alpine:3.24.1@sha256:28bd5fe8b56d1bd048e5babf5b10710ebe0bae67db86916198a6eec434943f8b
RUN apk --no-cache add ca-certificates bash
COPY --from=builder /go/bin/flan_exporter /app/flan_exporter
ENV FLAN_DATASOURCE="fs"
ENV FLAN_FS_PATH="/reports"
ENV FLAN_GCLOUD_CREDENTIAL_FILE="/app/gcloud_credentials.json"
ENV FLAN_GCLOUD_BUCKET_NAME="mauve-flan"
EXPOSE 9711
ENTRYPOINT /app/flan_exporter -datasource.provider=$FLAN_DATASOURCE -datasource.fs.report-path=$FLAN_FS_PATH -datasource.gcloud.credentials-path=$FLAN_GCLOUD_CREDENTIAL_FILE -datasource.gcloud.bucket-name=$FLAN_GCLOUD_BUCKET_NAME
