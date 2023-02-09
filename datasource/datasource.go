// SPDX-FileCopyrightText: (c) Mauve Mailorder Software GmbH & Co. KG, 2020. Licensed under [MIT](LICENSE) license
//
// SPDX-License-Identifier: MIT

package datasource

// DataSource defines the interface to get report files
type DataSource interface {
	// NewestReport returns the newest XML report
	NewestReport() (*Report, error)
}
