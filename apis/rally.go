package apis

import (
	"encoding/json"
	"fmt"
	"github.com/heysquirrel/tribe/config"
	"net/http"
	"strings"
)

type Rally struct {
	host   string
	apikey string
}

type RallyResult struct {
	QueryResult QueryResult `json:"QueryResult"`
}

type QueryResult struct {
	Artifacts []Artifact `json:"Results"`
}

type Artifact struct {
	ObjectType  string `json:"_type"`
	Name        string `json:"Name"`
	Description string `json:"Description"`
	Id          string `json:"_refObjectUUID"`
	FormattedID string `json:"FormattedID"`
}

func (a *Artifact) GetType() string        { return a.ObjectType }
func (a *Artifact) GetName() string        { return a.Name }
func (a *Artifact) GetDescription() string { return a.Description }
func (a *Artifact) GetId() string          { return a.FormattedID }

func NewRallyFromConfig(servername string) (*Rally, error) {
	serverconfig := config.WorkItemServer(config.ServerName(servername))

	return NewRally(serverconfig["host"], serverconfig["apikey"])
}

func NewRally(host, apikey string) (*Rally, error) {
	if strings.TrimSpace(host) == "" {
		return nil, fmt.Errorf("Invalid hostname: '%s'", host)
	}

	if strings.TrimSpace(apikey) == "" {
		return nil, fmt.Errorf("Invalid apikey: '%s'", apikey)
	}

	return &Rally{host, apikey}, nil
}

func (r *Rally) GetWorkItem(id string) (WorkItem, error) {
	client := &http.Client{}
	url := fmt.Sprintf("%s/slm/webservice/v2.0/Artifact", r.host)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("zsessionid", r.apikey)

	q := req.URL.Query()
	q.Add("query", fmt.Sprintf("(FormattedID = %s)", id))
	q.Add("fetch", "FormattedID,Name,Description")
	req.URL.RawQuery = q.Encode()

	res, err := client.Do(req)
	if err != nil {
		return NullWorkItem(id), err
	}
	defer res.Body.Close()

	var queryResult RallyResult
	json.NewDecoder(res.Body).Decode(&queryResult)

	for _, result := range queryResult.QueryResult.Artifacts {
		if result.FormattedID == id {
			return &result, nil
		}
	}

	return NullWorkItem(id), ItemNotFoundError(id)
}
