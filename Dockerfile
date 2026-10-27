FROM --platform=$BUILDPLATFORM golang:1.27.3-alpine3.24@sha256:d530cf60723abf1ffae0cd16bb546b379ec0550d2f3815e5b16ed69dbc00fba4 AS builder
ADD . /go/flan_exporter/
WORKDIR /go/flan_exporter
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /go/bin/flan_exporter

FROM alpine:3.24.3@sha256:be2656061b536c1883c4034790e21d0a09c5aaaf51968b89368a3c0371ea2c9a
RUN apk --no-cache add ca-certificates bash
COPY --from=builder /go/bin/flan_exporter /app/flan_exporter
ENV FLAN_DATASOURCE="fs"
ENV FLAN_FS_PATH="/reports"
ENV FLAN_GCLOUD_CREDENTIAL_FILE="/app/gcloud_credentials.json"
ENV FLAN_GCLOUD_BUCKET_NAME="mauve-flan"
EXPOSE 9711
ENTRYPOINT /app/flan_exporter -datasource.provider=$FLAN_DATASOURCE -datasource.fs.report-path=$FLAN_FS_PATH -datasource.gcloud.credentials-path=$FLAN_GCLOUD_CREDENTIAL_FILE -datasource.gcloud.bucket-name=$FLAN_GCLOUD_BUCKET_NAME
