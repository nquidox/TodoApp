package types

type DatabaseWorker interface {
	InitTable(model any) error

	CreateRecord(model any) error

	ReadOneRecord(model any, params map[string]any) error
	ReadRecordSubmodel(model any, submodel any, params map[string]any) error
	ReadManyRecords(model any, submodel any, params map[string]any) error
	ReadWithPagination(model any, params map[string]any) error

	UpdateRecord(model any, params map[string]any) error
	UpdateRecordSubmodel(model any, submodel any, params map[string]any) error

	DeleteRecord(model any, params map[string]any) error
	DeleteManyExceptOne(model any, params map[string]any) error
}
