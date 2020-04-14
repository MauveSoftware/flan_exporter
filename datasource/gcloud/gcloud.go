package gcloud

import (
	"github.com/MauveSoftware/flan_exporter/datasource"
)

type gcloudDataSource struct {
}

// New returns a datasource implementation using Google Cloud Storage
func New(bucket, credentialsFile string) (datasource.DataSource, error) {
	return &gcloudDataSource{}
}

func (d *gcloudDataSource) NewestReport() ([]*datasource.ReportFile, error) {
	return nil, nil
}
