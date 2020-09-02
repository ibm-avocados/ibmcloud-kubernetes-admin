package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/moficodes/ibmcloud-kubernetes-admin/pkg/ibmcloud"
)

func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	apiKey := os.Getenv("IBMCLOUD_API_KEY")
	if apiKey == "" {
		panic("api key is needed")
	}
	session, err := ibmcloud.IAMAuthenticate(apiKey)
	if err != nil {
		panic(err)
	}

	resourceGroup := os.Getenv("RESOURCE_GROUP")
	if resourceGroup == "" {
		panic("resource group is needed")
	}
	dataCenter := os.Getenv("DATACENTER")
	if dataCenter == "" {
		panic("datacenter is needed")
	}
	masterVersion := os.Getenv("KUBE_VERSION")
	if masterVersion == "" {
		panic("master version is needed")
	}

	entitlement := os.Getenv("ENTITLEMENT")
	machineType := os.Getenv("MACHINE_TYPE")
	if machineType == "" {
		panic("machine type is needed")
	}
	name := os.Getenv("CLUSTER_NAME")
	if name == "" {
		panic("name is needed")
	}

	privateVlan := os.Getenv("PRIVATE_VLAN")
	publicVlan := os.Getenv("PUBLIC_VLAN")

	_workerNum := os.Getenv("WORKER_NUMBER")
	if _workerNum == "" {
		panic("worker num is needed")
	}

	workerNum, err := strconv.Atoi(_workerNum)
	if err != nil {
		panic("worker num needs to be a valid integer")
	}

	request := ibmcloud.CreateClusterRequest{
		ResourceGroup: resourceGroup,
		ClusterRequest: ibmcloud.ClusterRequest{
			PublicServiceEndpoint: true,
			DataCenter:            dataCenter,
			MachineType:           machineType,
			MasterVersion:         masterVersion,
			Name:                  name,
			WorkerNum:             workerNum,
		},
	}

	if strings.Contains(masterVersion, "_openshift") {
		if entitlement == "" {
			panic("entitlement needed for openshift cluster")
		} else {
			request.ClusterRequest.DefaultWorkerPoolEntitlement = entitlement
		}
	}

	if privateVlan != "" && publicVlan != "" {
		request.ClusterRequest.PrivateVlan = privateVlan
		request.ClusterRequest.PublicVlan = publicVlan
	}

	var response *ibmcloud.CreateClusterResponse
	alreadyExists := false
	// janky retry logic
	for i := 0; i < 3; i++ {
		response, err = session.CreateCluster(request)
		if err != nil {
			if strings.Contains(err.Error(), "A cluster with the same name already exists. Choose another name.") {
				alreadyExists = true
				break
			}
			// sleep 5 second before retrying
			time.Sleep(5 * time.Second)
		} else {
			// error was nil
			// we succeded
			break
		}
	}

	// we are out of the loop
	// if err is nil, we are good.
	// else there was an error

	if err != nil && !alreadyExists {
		panic("could not create cluster")
	}

	var id string

	if alreadyExists {
		cluster, err := session.GetCluster(name, resourceGroup)
		if err != nil {
			panic("cannot get cluster")
		}
		id = cluster.ID
	} else {
		id = response.ID
	}

	_tags := os.Getenv("TAGS")

	tags := strings.Split(_tags, ",")
	for _, tag := range tags {
		_, err = session.SetClusterTag(tag, id, resourceGroup)
		if err != nil {
			fmt.Println(err)
		}
	}

	for {
		session, err = session.RenewSession()
		if err != nil {
			fmt.Println(err)
		}
		cluster, err := session.GetCluster(id, resourceGroup)
		if err != nil {
			fmt.Println(err)
		}
		if clusterProvisionComplete(cluster) {
			break
		}
		fmt.Println("sleeping for 5 minute before checking again")
		time.Sleep(5 * time.Minute)
	}

	fmt.Println("done!")
}

func clusterProvisionComplete(cluster *ibmcloud.Cluster) bool {
	return cluster.State == "normal" &&
		cluster.IngressHostname != "" &&
		cluster.IngressSecretName != "" &&
		cluster.MasterHealth == "normal" &&
		cluster.MasterState == "deployed" &&
		cluster.MasterStatus == "Ready"
}
