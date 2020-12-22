package domain

import (
	"github.com/Azure/azure-sdk-for-go/profiles/latest/containerservice/mgmt/containerservice"
	"github.com/aws/aws-sdk-go/service/eks"
)

const (
	SUBMITTED = "REQUEST SUBMITTED"
	COMPLETED = "RUNNING"
	STOPPED   = "STOPPED"
	FAILED    = "FAILED"
	DELETED   = "DELETED"
	// gke states
	GKE_CREATING = "CREATING CLUSTER"
	GKE_DELETING = "DELETING"
	// eks states
	EKS_MASTER_CREATING     = "CREATING CONTROL PLANE"
	EKS_MASTER_CREATED      = "CONTROL PLANE CREATED"
	EKS_MASTER_FAILED       = "CONTROL PLANE CREATION FAILED"
	EKS_NODE_GROUP_CREATING = "CREATING NODE GROUP"
	EKS_NODE_GROUP_CREATED  = "NODE GROUP CREATED"
	EKS_NODE_GROUP_FAILED   = "NODE GROUP CREATION FAILED"

	EKS_NODE_GROUP_DELETING      = "DELETING NODE GROUP"
	EKS_NODE_GROUP_DELETED       = "NODE GROUP DELETED"
	EKS_NODE_GROUP_DELETE_FAILED = "DELETING NODE GROUP FAILED"
	EKS_MASTER_DELETING          = "DELETING CONTROL PLANE"
	EKS_MASTER_DELETED           = "CONTROL PLANE DELETED"
	EKS_MASTER_DELETE_FAILED     = "DELETING CONTROL PLANE FAILED"
	EKS_SUBMITTED_FOR_DELETION   = "TO BE DELETED"
	// aks state
	AKS_CREATING               = "CREATING CLUSTER"
	AKS_SUBMITTED              = "REQUEST SUBMITTED"
	AKS_SUBMITTED_FOR_DELETION = "TO BE DELETED"
)

type ClusterResponse struct {
	Clusters []KubCluster `json:"clusters"`
	Error    error        `json:"error"`
	Status   int          `json:"status"`
	Message  string       `json:"message"`
}

type KubCluster struct {
	Id              int    `json:"id"`
	Reference       string `json:"reference"`
	ClusterId       string `json:"cluster_id"`
	ClusterName     string `json:"cluster_name"`
	ServiceProvider string `json:"service_providers"`
	Status          string `json:"status"`
	CreatedOn       string `json:"created_on"`
	ProjectName     string `json:"project_name"`
	Location        string `json:"location"`
	SecretId        string `json:"secret_id"`
	ResourceGroup   string `json:"resource_group"`
}

type ClusterOptions struct {
	Provider          string    `json:"provider"`
	UserId            int       `json:"user_id"` //todo: remove this, doesnt make sense
	SecretId          int       `json:"secret_id"`
	ResourceGroupName string    `json:"resource_group"`
	Name              string    `json:"name"`
	ClusterId         string    `json:"cluster_id"`
	Description       string    `json:"description"`
	Location          string    `json:"location"`
	Zone              string    `json:"zone"`
	InstanceCount     int32     `json:"instances"`
	ImageType         string    `json:"image_type"`
	MachineType       string    `json:"machine_type"`
	MachineFamily     []*string `json:"machine_family"`
	DiskSize          int       `json:"disk_size"`
	DiskType          string    `json:"disk_type"`
	KubVersion        string    `json:"kub_version"`
}

type ValidationResponse struct {
	Error   error  `json:"error"`
	Status  int    `json:"status"`
	Message string `json:"message"`
}

//GKE specific structs

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

//AKS specific structures
type AksClusterContext struct {
	ClusterResponse AksClusterStatus `json:"cluster_status"`
	ClusterRequest  ClusterOptions   `json:"cluster_request"`
	SecretID        int              `json:"secret_id"`
}

type AksClusterStatus struct {
	Name          string `json:"name"`
	ResourceGroup string `json:"resource_group"`
	KubVersion    string `json:"kub_version"`
	UserName      string `json:"user_name"`
	SSHPvtKey     string `json:"ssh_key_pvt"`
	SSHPubKey     string `json:"ssh_key_pub"`
	Status        string `json:"status"`
	Error         string `json:"error"`
}

type AksAsyncJobParams struct {
	ClusterOptions ClusterOptions
	CreateRequest  containerservice.ManagedCluster
	Client         containerservice.ManagedClustersClient
}

type AksResourceGroup struct {
	Groups []string `json:"groups"`
	Status string   `json:"status"`
	Error  string   `json:"error"`
}

//EKS specific structs

//type EksClusterStatus struct {
//	CreateClusterOutput  eks.CreateClusterOutput   `json:"eks_create_cluster_output"`
//	CreateNodGroupOutput eks.CreateNodegroupOutput `json:"eks_create_node_group_output"`
//}

type EksClusterContext struct {
	ClusterStatus  EksClusterStatus `json:"cluster_status"`
	ClusterRequest ClusterOptions   `json:"cluster_request"`
	SecretID       int              `json:"secret_id"`
}

type EksNodeGroupContext struct {
	SecretId int           `json:"secret_id"`
	Response eks.Nodegroup `json:"response"`
	Region   string        `json:"region"`
}

type EksClusterStatus struct {
	Name         string     `json:"name"`
	ClusterArn   string     `json:"cluster_arn"`
	RequestToken string     `json:"request_token"`
	RoleArn      string     `json:"role_arn"`
	SubnetIds    *[]*string `json:"subnet_ids"`
	KubVersion   string     `json:"kub_version"`
	Status       string     `json:"status"`
	Error        string     `json:"error"`
}

type EksAsyncJobParams struct {
	NodeGroupName string
	Client        eks.EKS
}

// async job processing
type AsyncCloudJob struct {
	Provider    string
	Status      string
	Reference   string
	Information interface{} //this could be any struct that wraps provider specific information
}

const (
	IsRunning = 1
	IsStopped = 0
)
