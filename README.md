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
        Type of the database if it cannot be guessed right from the connection string (default "postgres"), options: postgres, sqlserver
  -schema string
        Schema name for the database, defaults to dbo for MSSQL, public for postgres
```

Example:

```bash
mermaider -connection-string "host=localhost port=5432 user=postgres password=postgres dbname=postgres sslmode=disable" -schema public
```

### Supported databases

| Database        | Supported |
| --------------- | --------- |
| Postgres        | Yes       |
| MSSQL           | Yes       |
| MySQL / MariaDB | Planned   |
| SQLite          | Planned   |

### Limitations

Currently app designed only for usage with a single schema. It is relatively easy to change that and will be done at the later
point.

## Contribution

Feel free to add output formats or additional databases support. I'm open to pull requests. Make sure to have tests for your code _including_ integration tests using test containers.

## License

BSD 2-Clause License

## Plans

- [ ] Add support for Graphviz output
