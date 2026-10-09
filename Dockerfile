FROM golang:1.27.2-alpine3.24@sha256:85dc1069ac644ea3c527b177303a406eb3358192816cd7f9e5848eb658851673 AS builder
ADD . /go/flan_exporter/
WORKDIR /go/flan_exporter
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /go/bin/flan_exporter

FROM alpine:3.24.2@sha256:294b683cb724975bec92580e1e685676bd4b50bda910ddb8c51d4cabeaec77e6
RUN apk --no-cache add ca-certificates bash
COPY --from=builder /go/bin/flan_exporter /app/flan_exporter
ENV FLAN_DATASOURCE="fs"
ENV FLAN_FS_PATH="/reports"
ENV FLAN_GCLOUD_CREDENTIAL_FILE="/app/gcloud_credentials.json"
ENV FLAN_GCLOUD_BUCKET_NAME="mauve-flan"
EXPOSE 9711
ENTRYPOINT /app/flan_exporter -datasource.provider=$FLAN_DATASOURCE -datasource.fs.report-path=$FLAN_FS_PATH -datasource.gcloud.credentials-path=$FLAN_GCLOUD_CREDENTIAL_FILE -datasource.gcloud.bucket-name=$FLAN_GCLOUD_BUCKET_NAME
