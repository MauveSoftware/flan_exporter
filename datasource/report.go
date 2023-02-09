// SPDX-FileCopyrightText: (c) Mauve Mailorder Software GmbH & Co. KG, 2020. Licensed under [MIT](LICENSE) license
//
// SPDX-License-Identifier: MIT

package datasource

import "time"

type Report struct {
	Date  time.Time
	Files []*ReportFile
}
