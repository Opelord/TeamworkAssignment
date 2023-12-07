package main

import (
	"TeamworkAssignment/pkg/customerimporter"
	"database/sql"
	"flag"
	"log"
)

func main() {
	input := flag.String(
		"inputfile",
		"customers.csv",
		"CSV file with input data.",
	)
	output := flag.String(
		"outputfile",
		"",
		"optional - name of csv output file (without '.csv'). If left empty output will be print to terminal.")
	flag.Parse()

	db, err := sql.Open("mysql", "root:rootpassword@tcp(localhost:3306)/mydatabase")
	if err != nil {
		log.Fatalf("couldn't setup db connection: %v", err)
	}
	defer db.Close()

	if err = customerimporter.CreateAndPopulateTable(db, *input); err != nil {
		log.Fatalf("couldn't setup db: %v", err)
	}

	rows, err := customerimporter.CountDomains(db)
	if err != nil {
		log.Fatalf("couldn't fetch data from database: %v", err)
	}
	defer rows.Close()

	//store the result in Go object
	var result customerimporter.DomainCountList
	for rows.Next() {
		var domainCount customerimporter.DomainCount
		if err := rows.Scan(&domainCount.Domain, &domainCount.Count); err != nil {
			log.Fatalf("could scan row %v", err)
		}
		result = append(result, domainCount)
	}

	//write to output
	if *output != "" {
		err = result.WriteToCSV(*output)
		log.Fatalf("couldn't write to csv: %v", err)
	} else {
		result.Print()
	}
}
