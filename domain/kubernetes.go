package domain

import "github.com/aws/aws-sdk-go/service/eks"

type ClusterResponse struct {
	Clusters []KubCluster `json:"clusters"`
	Error    error        `json:"error"`
	Status   int          `json:"status"`
	Message  string       `json:"message"`
}

type KubCluster struct {
	Id              int    `json:"id"`
	ClusterId       string `json:"cluster_id"`
	ClusterName     string `json:"cluster_name"`
	ServiceProvider string `json:"service_providers"`
	Status          string `json:"status"`
	CreatedOn       string `json:"created_on"`
	ProjectName     string `json:"project_name"`
	Location        string `json:"location"`
}

//GKE specific structs
type GkeClusterOptions struct {
	UserId        int    `json:"user_id"` //todo: remove this, doesnt make sense
	SecretId      int    `json:"secret_id"`
	Name          string `json:"name"`
	ClusterId     string `json:"cluster_id"`
	Description   string `json:"description"`
	Location      string `json:"location"`
	Zone          string `json:"zone"`
	InstanceCount int32  `json:"instances"`
	ImageType     string `json:"image_type"`
	MachineType   string `json:"machine_type"`
	MachineFamily string `json:"machine_family"`
	DiskSize      int    `json:"disk_size"`
	DiskType      string `json:"disk_type"`
	KubVersion    string `json:"kub_version"`
}

type GkeClusterStatus struct {
	Name      string `json:"name"`
	ClusterId string `json:"cluster_id"`
	OpId      string `json:"operation_id"`
	Status    string `json:"status"`
	Error     string `json:"error"`
}

type GkeLROperation struct {
	Name      string
	Id        int
	Zone      string
	ProjectId string
	Status    string
	Error     error
}

type GkeOperationStatusCheck struct {
	OperationName string `json:"operation_name"`
	Status        string `json:"status"`
	Detail        string `json:"detail"`
	Error         error  `json:"error"`
}

type GkeRecommendations struct {
	Nodes  []Node `json:"nodes"`
	Status int    `json:"status"`
	Detail string `json:"detail"`
}

type Node struct {
	Type        string `json:"type""`
	Region      string `json:"region"`
	Processor   string `json:"processor"`
	Memory      string `json:"memory"`
	Network     string `json:"network"`
	StartupTime string `json:"startup_time"`
	NodeCount   string `json:"node_count"`
	Cost        string `json:"cost"`
}

type GkeResources struct {
	Continents  []string `json:"continents"`
	RegionNames []string `json:"region_names"`
	RegionIds   []string `json:"region_ids"`
	Status      int      `json:"status"`
	Detail      string   `json:"detail"`
}

type ResourceLocation struct {
	Continent  string `json:"continent"`
	RegionName string `json:"region_name"`
	RegionId   string `json:"region_id"`
}

//EKS specific structs
type EksClusterStatus struct {
	CreateClusterOutput  eks.CreateClusterOutput   `json:"eks_create_cluster_output"`
	CreateNodGroupOutput eks.CreateNodegroupOutput `json:"eks_create_node_group_output"`
}
