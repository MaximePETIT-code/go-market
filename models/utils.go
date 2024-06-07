package models

type CSVExportable interface {
    CSVHeaders() []string
    CSVRecord() []string
    TableName() string
}