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
