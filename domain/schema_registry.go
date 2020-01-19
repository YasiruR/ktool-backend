package domain

type SchemaRegistry struct {
	Host 			string
	Port 			int64
	NumOfSchemas	int64
	Schemas 		[]Schema
}

type Schema struct {
	ID 			float64
	Name 		string
	NameSpace 	string
	MaxVersion 	int64
}