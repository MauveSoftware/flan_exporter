package datasource

// DataSource defines the interface to get report files
type DataSource interface {
	// NewestReport returns the newest XML report
	NewestReport() (*Report, error)
}
