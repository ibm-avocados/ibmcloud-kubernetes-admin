package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/moficodes/ibmcloud-kubernetes-admin/pkg/ibmcloud"
)

func (s *Server) DeleteTagHandler(w http.ResponseWriter, r *http.Request) {
	session, err := getCloudSessions(r)
	if err != nil {
		handleError(w, http.StatusUnauthorized, "could not get session", err.Error())
		return
	}

	var body ibmcloud.UpdateTag
	decoder := json.NewDecoder(r.Body)

	err = decoder.Decode(&body)
	if err != nil {
		handleError(w, http.StatusBadRequest, "could not decode", err.Error())
		return
	}

	res, err := session.DeleteTag(body)
	if err != nil {
		handleError(w, http.StatusInternalServerError, "could not delete tag", err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"status": "ok, deleted %d tags"}`, len(res.Results))
}

func (s *Server) GetTagHandler(w http.ResponseWriter, r *http.Request) {
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

	crn, ok := body["crn"]
	if !ok {
		handleError(w, http.StatusBadRequest, "no crn attached to body", err.Error())
		return
	}

	clusterCRN := fmt.Sprintf("%v", crn)

	tags, err := session.GetTags(clusterCRN)
	if err != nil {
		handleError(w, http.StatusUnauthorized, "could not get tags", err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	e := json.NewEncoder(w)
	e.Encode(tags)
}

func (s *Server) SetClusterTagHandler(w http.ResponseWriter, r *http.Request) {
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

	tag, ok := body["tag"]
	if !ok {
		handleError(w, http.StatusBadRequest, "no tag attached to body", err.Error())
		return
	}

	resourceGroup, ok := body["resourceGroup"]
	if !ok {
		handleError(w, http.StatusBadRequest, "no resourceID attached to body", err.Error())
		return
	}

	vars := mux.Vars(r)

	clusterID, ok := vars["clusterID"]

	if !ok {
		handleError(w, http.StatusBadRequest, "could not get clusterID")
		return
	}
	clusterTag := fmt.Sprintf("%v", tag)
	clusterResourceGroup := fmt.Sprintf("%v", resourceGroup)

	res, err := session.SetClusterTag(clusterTag, clusterID, clusterResourceGroup)
	if err != nil {
		log.Println("set cluster tag : ", err)
		handleError(w, http.StatusInternalServerError, "could not update tag", err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"status": "ok, wrote %d tags"}`, len(res.Results))

}

func (s *Server) SetTagHandler(w http.ResponseWriter, r *http.Request) {
	session, err := getCloudSessions(r)
	if err != nil {
		handleError(w, http.StatusUnauthorized, "could not get session", err.Error())
		return
	}

	var body ibmcloud.UpdateTag
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&body)
	if err != nil {
		handleError(w, http.StatusBadRequest, "could not decode", err.Error())
		return
	}

	res, err := session.SetTag(body)
	if err != nil {
		handleError(w, http.StatusInternalServerError, "could not update tag", err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"status": "ok, wrote %d tags"}`, len(res.Results))
}
