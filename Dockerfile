FROM --platform=$BUILDPLATFORM golang:1.27.4-alpine3.24@sha256:70bd7ddea8b6cbe5d6943cde279d6b67feb10cc924fb9100d1d8757c1575b280 AS builder
ADD . /go/flan_exporter/
WORKDIR /go/flan_exporter
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /go/bin/flan_exporter

FROM alpine:3.24.4@sha256:f59abe656b692d901e184be3445e084786602108599d2a51d229c0f5a48c71bf
RUN apk --no-cache add ca-certificates bash
COPY --from=builder /go/bin/flan_exporter /app/flan_exporter
ENV FLAN_DATASOURCE="fs"
ENV FLAN_FS_PATH="/reports"
ENV FLAN_GCLOUD_CREDENTIAL_FILE="/app/gcloud_credentials.json"
ENV FLAN_GCLOUD_BUCKET_NAME="mauve-flan"
EXPOSE 9711
ENTRYPOINT /app/flan_exporter -datasource.provider=$FLAN_DATASOURCE -datasource.fs.report-path=$FLAN_FS_PATH -datasource.gcloud.credentials-path=$FLAN_GCLOUD_CREDENTIAL_FILE -datasource.gcloud.bucket-name=$FLAN_GCLOUD_BUCKET_NAME
