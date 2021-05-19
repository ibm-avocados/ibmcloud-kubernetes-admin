package kubeadmin

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/labstack/echo/v4"

	"github.com/moficodes/ibmcloud-kubernetes-admin/pkg/ibmcloud"
)

func TestGithubCommentHandler(t *testing.T) {
	issue := ibmcloud.GithubIssueComment{
		IssueNumber: os.Getenv("TEST_GITHUB_REPO"),
		EventName:   "handlertest",
		Password:    "password",
		AccountID:   "1234",
		GithubUser:  "Mofizur-Rahman",
		GithubToken: os.Getenv("TEST_GITHUB_TOKEN"),
		ClusterRequest: ibmcloud.GithubIssueClusterRequest{
			Count:      10,
			Type:       "kubernetes",
			ErrorCount: 0,
			Regions:    "dal-10,dal-12",
		},
	}
	b, err := json.Marshal(issue)
	if err != nil {
		t.Fatal(err)
	}
	buf := bytes.NewBuffer(b)
	r := httptest.NewRequest(http.MethodPost, "/github/comment", buf)
	w := httptest.NewRecorder()
	e := echo.New()
	c := e.NewContext(r, w)

	err = GithubCommentHandler(c)
	if err != nil {
		t.Fatal(err)
	}
}
