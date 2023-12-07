// Package customerimporter reads from a CSV file and returns a sorted (data
// structure of your choice) of email domains along with the number of customers
// with e-mail addresses for each domain. This should be able to be run from the
// CLI and output the sorted domains to the terminal or to a file. Any errors
// should be logged (or handled). Performance matters (this is only ~3k lines,
// but could be 1m lines or run on a small machine).
package customerimporter

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"os"
	"strconv"
)

type DomainCountList []DomainCount

type DomainCount struct {
	Domain string
	Count  int
}

func (dcl DomainCountList) Print() {
	for _, dc := range dcl {
		fmt.Printf("%-20v %d\n", dc.Domain, dc.Count)
	}
}

func (dcl DomainCountList) WriteToCSV(fileName string) error {
	f, err := os.Create(fileName + ".csv")
	if err != nil {
		return fmt.Errorf("couldn't create file: %v", err)
	}
	csvWriter := csv.NewWriter(f)
	for _, dc := range dcl {
		err = csvWriter.Write([]string{dc.Domain, strconv.Itoa(dc.Count)})
		if err != nil {
			return fmt.Errorf("couldn't write to file: %v", err)
		}
	}
	csvWriter.Flush()
	return nil
}

// CreateAndPopulateTable uses db connection to load data from csv to database using mysql queries.
func CreateAndPopulateTable(db *sql.DB, inputFile string) error {
	//drop table if another already exists e.g. from previous run
	_, err := db.Exec("DROP TABLE IF EXISTS customers")
	if err != nil {
		return err
	}

	//prepare create table query
	f, err := os.Open(inputFile)
	if err != nil {
		return fmt.Errorf("couldn't open file: %v", err)
	}
	csvReader := csv.NewReader(f)
	headers, err := csvReader.Read()
	if err != nil {
		return fmt.Errorf("couldn't read headers for table")
	}
	createTableQuery := "CREATE TABLE customers ("
	for _, header := range headers {
		createTableQuery += fmt.Sprintf(" %v TEXT ,", header)
	}
	createTableQuery = createTableQuery[:len(createTableQuery)-1] + ")"

	//create table 'customers'
	_, err = db.Exec(createTableQuery)
	if err != nil {
		return fmt.Errorf("couldn't create table customers: %v", err)
	}

	//populate data
	mysql.RegisterLocalFile(inputFile)
	_, err = db.Exec("LOAD DATA LOCAL INFILE '" + inputFile + "' INTO TABLE customers FIELDS TERMINATED BY ',' LINES TERMINATED BY '\\n' IGNORE 1 LINES;")
	if err != nil {
		return fmt.Errorf("couldn't load csv file into db: %v", err)
	}

	return nil
}

// CountDomains uses sql query to count domains in emails.
// SUBSTRING_INDEX uses column 'email' and splits the string whenever '@' appears.
// -1 specifies to use rightmost result of the split, which is the domain.
func CountDomains(db *sql.DB) (*sql.Rows, error) {
	sqlQuery := `SELECT SUBSTRING_INDEX(email, '@', -1) AS domain, COUNT(*) AS count
	FROM customers
	GROUP BY domain
	ORDER BY count DESC;`
	return db.Query(sqlQuery)
}
