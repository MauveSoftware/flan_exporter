package datasource

import "time"

type Report struct {
	Date  time.Time
	Files []*ReportFile
}
