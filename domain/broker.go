package domain

type Broker struct {
	ID 					int
	Host 				string
	Port 				int
	Version 			float64
	Processors 			int
	NumOfTopics			int
	TotalPartitions 	int
	OfflinePartitions 	int
	ActiveControllers 	int
	Status 				string
	CreatedAt 			string
	ClusterID 			int
}

type Server struct {
	Host 			string 		`json:"host"`
	Port 			int			`json:"port"`
	MetricsPort 	int			`json:"metrics_port"`
}
