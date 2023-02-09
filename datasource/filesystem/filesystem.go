// SPDX-FileCopyrightText: (c) Mauve Mailorder Software GmbH & Co. KG, 2020. Licensed under [MIT](LICENSE) license
//
// SPDX-License-Identifier: MIT

package filesystem

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/MauveSoftware/flan_exporter/datasource"
	"github.com/pkg/errors"
)

type filesystemDataSource struct {
	directory string
}

// New returns a datasource implementation using local filesystem
func New(directory string) datasource.DataSource {
	return &filesystemDataSource{directory: directory}
}

func (d *filesystemDataSource) NewestReport() (*datasource.Report, error) {
	fs, err := ioutil.ReadDir(d.directory)
	if err != nil {
		return nil, errors.Wrap(err, "could not list reports directory")
	}

	if len(fs) == 0 {
		return nil, errors.Errorf("could not find any report directory in %s", d.directory)
	}

	date, dir := d.newestReport(fs)
	files, err := d.filesFromDir(dir)
	if err != nil {
		return nil, errors.Wrapf(err, "could not get files for report %s", dir.Name())
	}

	return &datasource.Report{
		Date:  date,
		Files: files,
	}, nil
}

func (d *filesystemDataSource) newestReport(fs []os.FileInfo) (date time.Time, dir os.FileInfo) {
	var newest os.FileInfo
	var newestTime time.Time

	for _, fi := range fs {
		if !fi.IsDir() {
			continue
		}

		isReport, t := datasource.DateOfReportDir(fi.Name())
		if !isReport {
			continue
		}

		if newest == nil || t.After(newestTime) {
			newestTime = t
			newest = fi
		}
	}

	return newestTime, newest
}

func (d *filesystemDataSource) filesFromDir(dir os.FileInfo) ([]*datasource.ReportFile, error) {
	p := filepath.Join(d.directory, dir.Name())
	fs, err := ioutil.ReadDir(p)
	if err != nil {
		return nil, errors.Wrapf(err, "could not list directory %s", dir.Name())
	}

	files := make([]*datasource.ReportFile, 0)
	for _, fi := range fs {
		if !strings.HasSuffix(fi.Name(), ".xml") {
			continue
		}

		f, err := d.fileFromInfo(dir, fi)
		if err != nil {
			return nil, err
		}

		files = append(files, f)
	}

	return files, nil
}

func (d *filesystemDataSource) fileFromInfo(dirInfo, fileInfo os.FileInfo) (*datasource.ReportFile, error) {
	p := filepath.Join(d.directory, dirInfo.Name(), fileInfo.Name())

	b, err := ioutil.ReadFile(p)
	if err != nil {
		return nil, errors.Errorf("could not read xml report file %s", fileInfo.Name())
	}

	return &datasource.ReportFile{
		Name:    fileInfo.Name(),
		Content: b,
	}, nil
}
