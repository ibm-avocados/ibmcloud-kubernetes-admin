package server

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/moficodes/ibmcloud-kubernetes-admin/pkg/ibmcloud"
)

func (s *Server) VersionEndpointHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	versions, err := ibmcloud.GetVersions()
	if err != nil {
		handleError(w, http.StatusNotFound, "could not get locations")
	}
	w.WriteHeader(http.StatusOK)
	e := json.NewEncoder(w)
	e.Encode(versions)
}

func (s *Server) LocationEndpointHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	locations, err := ibmcloud.GetLocations()
	if err != nil {
		handleError(w, http.StatusNotFound, "could not get locations")
	}
	w.WriteHeader(http.StatusOK)
	e := json.NewEncoder(w)
	e.Encode(locations)
}

func (s *Server) LocationGeoEndpointHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	vars := mux.Vars(r)

	geo, ok := vars["geo"]

	if !ok {
		handleError(w, http.StatusBadRequest, "could not get clusterID")
		return
	}

	locations, err := ibmcloud.GetGeoLocations(geo)
	if err != nil {
		handleError(w, http.StatusNotFound, "could not get locations")
	}
	w.WriteHeader(http.StatusOK)
	e := json.NewEncoder(w)
	e.Encode(locations)
}

func (s *Server) ZonesEndpointHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	showFlavors := r.FormValue("showFlavors")
	location := r.FormValue("location")
	zones, err := ibmcloud.GetZones(showFlavors, location)
	if err != nil {
		handleError(w, http.StatusNotFound, "could not load zones")
	}
	w.WriteHeader(http.StatusOK)
	e := json.NewEncoder(w)
	e.Encode(zones)
}

func (s *Server) MachineTypeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	serverType := r.FormValue("type")
	os := r.FormValue("os")
	cpuLimitStr := r.FormValue("cpuLimit")
	cpuLimit, err := strconv.Atoi(cpuLimitStr)
	if err != nil {
		handleError(w, http.StatusBadRequest, "cpuLimit should be a number")
		return
	}
	memoryLimitStr := r.FormValue("memoryLimit")
	memoryLimit, err := strconv.Atoi(memoryLimitStr)
	if err != nil {
		handleError(w, http.StatusBadRequest, "memoryLimit should be a number")
		return
	}

	vars := mux.Vars(r)

	datacenter, ok := vars["datacenter"]

	if !ok {
		handleError(w, http.StatusBadRequest, "could not get clusterID")
		return
	}
	flavors, err := ibmcloud.GetMachineType(datacenter, serverType, os, cpuLimit, memoryLimit)
	if err != nil {
		handleError(w, http.StatusNotFound, "could not load flavor")
	}
	w.WriteHeader(http.StatusOK)
	e := json.NewEncoder(w)
	e.Encode(flavors)
}
