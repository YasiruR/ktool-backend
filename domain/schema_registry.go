package domain

type SchemaRegistry struct {
	Host 			string
	Port 			int
	NumOfSchemas	int
	Schemas 		[]Schema
}

type Schema struct {
	ID 			float64
	Name 		string
	NameSpace 	string
	MaxVersion 	int
}