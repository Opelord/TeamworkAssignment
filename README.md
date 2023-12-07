# Teamwork.com assignment

### Task description
Package customerimporter reads from a CSV file and returns a sorted (data
structure of your choice) of email domains along with the number of customers
with e-mail addresses for each domain.
This should be able to be run from the CLI and output the sorted domains to the terminal or to a file.
Any errors should be logged (or handled).
Performance matters (this is only ~3k lines, but could be 1m lines or run on a small machine).

### Prerequisites
- makefile
- go 1.20
- docker

### Before
```
cd TeamworkAssignment
make start-db
```

### How to run
```
go build -o ./bin/assignment ./cmd/main.go
./bin/assignment -inputfile=customers.csv -outputfile=customers.txt
```

### Test
```
go test ./...
```

### Cleanup
```
make stop-db
rm -r ./bin
```