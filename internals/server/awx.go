package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/moficodes/ibmcloud-kubernetes-admin/pkg/awx"
)

func (s *Server) GetAWXWorkflowJobTemplates(w http.ResponseWriter, r *http.Request) {
	token := os.Getenv("AWX_ACCESS_TOKEN")
	templates, err := awx.GetWorkflowJobTemplates(token)
	if err != nil {
		handleError(w, http.StatusUnauthorized, "token invalid", err.Error())
		return
	}

	query := r.URL.Query().Get("labels")
	fmt.Println("QUERY: ", query)

	var result []awx.ResultsWorkflowTemplate
	if query == "" {
		result = templates.Results
	} else {
		for _, res := range templates.Results {
			for _, label := range res.SummaryFields.Labels.Results {
				if label.Name == query {
					result = append(result, res)
					break
				}
			}
		}
	}

	w.WriteHeader(http.StatusOK)
	e := json.NewEncoder(w)
	e.Encode(result)
}

func (s *Server) GetAWXJobTemplates(w http.ResponseWriter, r *http.Request) {
	token := os.Getenv("AWX_ACCESS_TOKEN")
	templates, err := awx.GetJobTemplates(token)
	if err != nil {
		handleError(w, http.StatusUnauthorized, "token invalid", err.Error())
		return
	}

	result := templates.Results

	w.WriteHeader(http.StatusOK)
	e := json.NewEncoder(w)
	e.Encode(result)
}
