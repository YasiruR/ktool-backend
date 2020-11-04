package domain

var ClusterTopicMap map[int][]TopicMetrics

type KTopic struct {
	Name 		string
	Partitions 	[]int32
}

type TopicMetrics struct {
	Name 				string				`json:"name"`
	ClusterName			string				`json:"cluster_name"`
	WritablePartitions	[]int32				`json:"writable_partitions"`
	Partitions 			[]TopicPartition	`json:"partitions"`
	//Messages 			[]int				`json:"messages"`
}

type TopicPartition struct {
	ID 					int32				`json:"id"`
	Replicas 			[]int32				`json:"replicas"`
	InSyncReplicas 		[]int32				`json:"in_sync_replicas"`
	UnderReplicated 	bool				`json:"under_replicated"`
	OfflineReplicas		[]int32				`json:"offline_replicas"`
	FirstOffset 		int64				`json:"first_offset"`				//use get offset method
	NextOffset 			int64				`json:"last_offset"`
	HighWaterMark 		int64				`json:"last_repl_offset"`			//use high water mark function of
}