// SPDX-FileCopyrightText: (c) Mauve Mailorder Software GmbH & Co. KG, 2020. Licensed under [MIT](LICENSE) license
//
// SPDX-License-Identifier: MIT

package datasource

import "time"

// DateOfReportDir parses the name of the report dir and extracts the date of the report from it
func DateOfReportDir(name string) (isReportDir bool, reportDate time.Time) {
	t, err := time.Parse("2006.01.02-15.04", name)
	if err != nil {
		return false, time.Time{}
	}

	return true, t
}
