package customerimporter

import (
	"database/sql"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCountDomains(t *testing.T) {
	db, err := sql.Open("mysql", "root:rootpassword@tcp(localhost:3306)/mydatabase")
	require.NoError(t, err)
	defer db.Close()

	err = CreateAndPopulateTable(db, "test.csv")
	require.NoError(t, err)

	rows, err := CountDomains(db)
	require.NoError(t, err)
	defer rows.Close()

	//require proper columns
	columns, err := rows.Columns()
	require.NoError(t, err)
	require.Equal(t, "domain", columns[0])
	require.Equal(t, "count", columns[1])

	//require first row is domain and count
	require.True(t, rows.Next())
	var domain string
	var count int
	err = rows.Scan(&domain, &count)
	require.NoError(t, err)
	require.Equal(t, "360.cn", domain)
	require.Equal(t, 1, count)

	//require no more rows
	require.False(t, rows.Next())
}
