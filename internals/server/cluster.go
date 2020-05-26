package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/moficodes/ibmcloud-kubernetes-admin/pkg/ibmcloud"
)

func (s *Server) ClusterDeleteHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	session, err := getCloudSessions(r)
	if err != nil {
		handleError(w, http.StatusUnauthorized, "could not get session", err.Error())
		return
	}
	var body map[string]interface{}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&body)
	if err != nil {
		handleError(w, http.StatusBadRequest, "could not decode", err.Error())
		return
	}

	id := fmt.Sprintf("%v", body["id"])
	resoueceGroup := fmt.Sprintf("%v", body["resourceGroup"])
	deleteResources := fmt.Sprintf("%v", body["deleteResources"])
	err = session.DeleteCluster(id, resoueceGroup, deleteResources)
	if err != nil {
		handleError(w, http.StatusInternalServerError, "could not delete", err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, statusOkMessage)
}

func (s *Server) ClusterCreateHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	session, err := getCloudSessions(r)
	if err != nil {
		handleError(w, http.StatusUnauthorized, "could not get session", err.Error())
		return
	}

	var body ibmcloud.CreateClusterRequest

	decoder := json.NewDecoder(r.Body)

	err = decoder.Decode(&body)
	if err != nil {
		handleError(w, http.StatusBadRequest, "could not decode ", err.Error())
		return
	}

	createResponse, err := session.CreateCluster(body)

	if err != nil {
		handleError(w, http.StatusUnauthorized, "could not create cluster", err.Error())
		return
	}

	log.Println("cluster created :", createResponse.ID)

	w.WriteHeader(http.StatusOK)
	e := json.NewEncoder(w)
	e.Encode(createResponse)
}

func (s *Server) ClusterHandler(w http.ResponseWriter, r *http.Request) {
	session, err := getCloudSessions(r)
	if err != nil {
		handleError(w, http.StatusUnauthorized, "could not get session", err.Error())
		return
	}

	vars := mux.Vars(r)

	clusterID, ok := vars["clusterID"]

	if !ok {
		handleError(w, http.StatusBadRequest, "could not get clusterID")
		return
	}

	var body map[string]interface{}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&body)
	if err != nil {
		handleError(w, http.StatusBadRequest, "could not decode", err.Error())
		return
	}

	resourceGroup, ok := body["resourceGroup"]
	if !ok {
		handleError(w, http.StatusBadRequest, "no tag attached to body", err.Error())
		return
	}
	clusterResourceGroup := fmt.Sprintf("%v", resourceGroup)
	cluster, err := session.GetCluster(clusterID, clusterResourceGroup)
	if err != nil {
		handleError(w, http.StatusUnauthorized, "could not get cluster", err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	e := json.NewEncoder(w)
	e.Encode(cluster)
}

func (s *Server) ClusterListHandler(w http.ResponseWriter, r *http.Request) {
	session, err := getCloudSessions(r)
	if err != nil {
		handleError(w, http.StatusUnauthorized, "could not get session", err.Error())
		return
	}

	clusters, err := session.GetClusters("")
	if err != nil {
		handleError(w, http.StatusUnauthorized, "could not get clusters", err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	e := json.NewEncoder(w)
	e.Encode(clusters)
}

func (s *Server) ClusterWorkerListHandler(w http.ResponseWriter, r *http.Request) {
	session, err := getCloudSessions(r)
	if err != nil {
		handleError(w, http.StatusUnauthorized, "could not get session", err.Error())
		return
	}

	vars := mux.Vars(r)

	clusterID, ok := vars["clusterID"]

	if !ok {
		handleError(w, http.StatusBadRequest, "could not get clusterID")
		return
	}

	workers, err := session.GetWorkers(clusterID)

	if err != nil {
		handleError(w, http.StatusUnauthorized, "could not get workers for cluster : ", clusterID, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	e := json.NewEncoder(w)
	e.Encode(workers)
}

func (s *Server) VlanEndpointHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	session, err := getCloudSessions(r)
	if err != nil {
		handleError(w, http.StatusUnauthorized, "could not get session", err.Error())
		return
	}
	vars := mux.Vars(r)

	datacenters, ok := vars["datacenter"]

	if !ok {
		handleError(w, http.StatusBadRequest, "could not get datacenter")
		return
	}

	vlans, err := session.GetDatacenterVlan(datacenters)
	if err != nil {
		handleError(w, http.StatusNotFound, "could not get vlans")
	}
	w.WriteHeader(http.StatusOK)
	e := json.NewEncoder(w)
	e.Encode(vlans)
}
