package gcloud

import (
	"context"
	"io/ioutil"
	"strings"
	"sync"
	"time"

	"cloud.google.com/go/storage"
	"github.com/MauveSoftware/flan_exporter/datasource"
	"github.com/pkg/errors"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

type gcloudDataSource struct {
	client      *storage.Client
	bucketName  string
	ctx         context.Context
	cached      *datasource.Report
	cachedMutex sync.Mutex
}

// New returns a datasource implementation using Google Cloud Storage
func New(bucketName, credentialsFile string) (datasource.DataSource, error) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx, option.WithCredentialsFile(credentialsFile))
	if err != nil {
		return nil, errors.Wrap(err, "could not create storage client for Google Cloud Storage")
	}

	return &gcloudDataSource{
		ctx:        ctx,
		client:     client,
		bucketName: bucketName,
	}, nil
}

func (d *gcloudDataSource) NewestReport() (*datasource.Report, error) {
	bucket := d.client.Bucket(d.bucketName)
	objs, err := d.listReportObjects(bucket)
	if err != nil {
		return nil, errors.Wrap(err, "could not get list of report objects")
	}

	if len(objs) == 0 {
		return nil, errors.Errorf("could not find report files in bucket")
	}

	return d.newestReportFromObjects(bucket, objs)
}

func (d *gcloudDataSource) listReportObjects(bucket *storage.BucketHandle) ([]*storage.ObjectAttrs, error) {
	objs := bucket.Objects(d.ctx, &storage.Query{Prefix: "xml_files/"})
	list := make([]*storage.ObjectAttrs, 0)

	for {
		attrs, err := objs.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}

		list = append(list, attrs)
	}

	return list, nil
}

func (d *gcloudDataSource) newestReportFromObjects(bucket *storage.BucketHandle, objs []*storage.ObjectAttrs) (*datasource.Report, error) {
	reports := d.reportMap(objs)

	var newestKey string
	var newestDate time.Time
	for key := range reports {
		isReportDir, date := datasource.DateOfReportDir(key)
		if !isReportDir {
			continue
		}

		if newestKey == "" || date.After(newestDate) {
			newestKey = key
			newestDate = date
		}
	}

	d.cachedMutex.Lock()
	defer d.cachedMutex.Unlock()

	if d.cached != nil && d.cached.Date == newestDate {
		return d.cached, nil
	}

	r, err := d.reportFromObjects(newestDate, bucket, reports[newestKey])
	if err != nil {
		return nil, err
	}
	d.cached = r

	return r, nil
}

func (d *gcloudDataSource) reportMap(objs []*storage.ObjectAttrs) map[string][]*storage.ObjectAttrs {
	reports := make(map[string][]*storage.ObjectAttrs)

	for _, obj := range objs {
		t := strings.Split(obj.Name, "/")
		reportName := t[1]

		reports[reportName] = append(reports[reportName], obj)
	}

	return reports
}

func (d *gcloudDataSource) reportFromObjects(reportDate time.Time, bucket *storage.BucketHandle, objs []*storage.ObjectAttrs) (*datasource.Report, error) {
	files, err := d.reportFiles(bucket, objs)
	if err != nil {
		return nil, err
	}

	return &datasource.Report{
		Date:  reportDate,
		Files: files,
	}, nil
}

func (d *gcloudDataSource) reportFiles(bucket *storage.BucketHandle, objs []*storage.ObjectAttrs) ([]*datasource.ReportFile, error) {
	files := make([]*datasource.ReportFile, len(objs))

	for i, obj := range objs {
		f, err := d.reportFileFromObject(bucket, obj)
		if err != nil {
			return nil, errors.Wrapf(err, "could not read object %s", obj.Name)
		}

		files[i] = f
	}

	return files, nil
}

func (d *gcloudDataSource) reportFileFromObject(bucket *storage.BucketHandle, obj *storage.ObjectAttrs) (*datasource.ReportFile, error) {
	r, err := bucket.Object(obj.Name).NewReader(d.ctx)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	return &datasource.ReportFile{
		Name:    obj.Name,
		Content: b,
	}, nil
}
