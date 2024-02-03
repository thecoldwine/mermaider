# Mermaider

This is a CLI utility to make tedious task of documenting existing databases a bit easier.
Mermaider will generate ER diagrams using Mermaid notation for your table and you can just
paste it into your markdown file in MKDocs or wherever.

## Installation

```bash
go install github.com/thecoldwine/mermaider
```

## Usage

```bash
Usage of mermaider:
  -connection-string string
        Connection string to a target database, driver will be inferred automatically
  -db-type string
        Type of the database if it cannot be guessed right from the connection string (default "postgres")
  -schema string
        Schema name for the database, defaults to dbo for MSSQL, public for postgres
```

Example:

```bash
mermaider -connection-string "host=localhost port=5432 user=postgres password=postgres dbname=postgres sslmode=disable" -schema public
```

## Contribution

Feel free to add output formats or additional databases support. I'm open to pull requests. Make sure to have tests for your code _including_ integration tests using test containers.

## License

BSD 2-Clause License

## Plans

_before 1.0_
- [x] Add support for Mermaid output
- [x] Add support for Postgres
- [ ] Add support for MS SQL Server
  
_maybe_
- [ ] Add support for Graphviz output
- [ ] Add support for MySQL
