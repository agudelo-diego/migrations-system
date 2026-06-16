package migration

import (
	"fmt"
	"regexp"
	"strconv"
)

// Formato esperado:
// 20260616101500_table_create_bank_accounts.sql
// 20260616101501_sp_bank_accounts_create_v1.sql
// 20260616101502_seed_payment_methods.sql

var filenameRe = regexp.MustCompile(
	`^(\d{14})_(table|sp|seed)_(.+)\.sql$`,
)

func ParseFilename(filename string) (version int64, migType Type, name string, err error) {
	matches := filenameRe.FindStringSubmatch(filename)

	if len(matches) < 4 {
		return 0, "", "", fmt.Errorf(
			"nombre inválido: '%s'\n"+
				"  formato:       YYYYMMDDHHMMSS_tipo_nombre.sql\n"+
				"  tipos válidos: table | sp | seed",
			filename,
		)
	}

	version, err = strconv.ParseInt(matches[1], 10, 64)
	if err != nil {
		return 0, "", "", fmt.Errorf("timestamp inválido en '%s': %w", filename, err)
	}

	return version, Type(matches[2]), matches[3], nil
}
