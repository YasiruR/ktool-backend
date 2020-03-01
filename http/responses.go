package http

//-------------------Cluster--------------------------//

type clusterRes struct {
	Clusters 	[]clusterInfo	`json:"clusters"`
}

type clusterInfo struct {
	Id 					int 		`json:"id"`
	ClusterName			string 		`json:"cluster_name"`
	Brokers 			[]string	`json:"brokers"`
	Topics 				[]topic 	`json:"topics"`
	Available 			bool		`json:"available"`
}

type topic struct {
	Name 		string		`json:"name"`
	Partitions 	[]int32 	`json:"partitions"`
}

//type broker struct {
//	Host 		string			`json:"host"`
//	Port 		int				`json:"port"`
//}

type errorMessage struct {
	Mesg 	string		`json:"mesg"`
}