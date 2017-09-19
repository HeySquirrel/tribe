package rally

import (
	"encoding/json"
	"fmt"
	"github.com/heysquirrel/tribe/apis"
	"net/http"
)

type Rally struct {
	apikey string
}

type Result struct {
	QueryResult QueryResult `json:"QueryResult"`
}

type QueryResult struct {
	Results []Artifact `json:"Results"`
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

func New(apikey string) *Rally {
	rally := new(Rally)
	rally.apikey = apikey

	return rally
}

func (r *Rally) GetWorkItem(id string) (apis.WorkItem, error) {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", "https://rally1.rallydev.com/slm/webservice/v2.0/Artifact", nil)
	req.Header.Set("zsessionid", r.apikey)

	q := req.URL.Query()
	q.Add("query", fmt.Sprintf("(FormattedID = %s)", id))
	q.Add("fetch", "FormattedID,Name,Description")
	req.URL.RawQuery = q.Encode()

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var queryResult Result
	json.NewDecoder(res.Body).Decode(&queryResult)

	for _, result := range queryResult.QueryResult.Results {
		if result.FormattedID == id {
			return &result, nil
		}
	}

	return &Artifact{FormattedID: id}, nil
}

func (r *Rally) queryForArtifact(formattedid string) (*Artifact, error) {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", "https://rally1.rallydev.com/slm/webservice/v2.0/Artifact", nil)
	req.Header.Set("zsessionid", r.apikey)

	q := req.URL.Query()
	q.Add("query", fmt.Sprintf("(FormattedID = %s)", formattedid))
	q.Add("fetch", "FormattedID,Name,Description")
	req.URL.RawQuery = q.Encode()

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var queryResult Result
	json.NewDecoder(res.Body).Decode(&queryResult)

	for _, result := range queryResult.QueryResult.Results {
		if result.FormattedID == formattedid {
			return &result, nil
		}
	}

	return nil, nil
}

func (r *Rally) GetByFormattedId(formattedId string) (*Artifact, error) {
	client := &http.Client{}

	req, _ := http.NewRequest("GET", "https://rally1.rallydev.com/slm/webservice/v2.0/Artifact", nil)
	req.Header.Set("zsessionid", r.apikey)

	q := req.URL.Query()
	q.Add("query", fmt.Sprintf("(FormattedID = %s)", formattedId))
	q.Add("fetch", "FormattedID,Name,Description")
	req.URL.RawQuery = q.Encode()

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var queryResult Result
	json.NewDecoder(res.Body).Decode(&queryResult)

	for _, result := range queryResult.QueryResult.Results {
		if result.FormattedID == formattedId {
			return &result, nil
		}
	}

	return nil, nil
}

func (r *Rally) GetByFormattedIds(formattedIds ...string) ([]*Artifact, error) {
	artifacts := make([]*Artifact, 0)

	for _, formattedId := range formattedIds {
		artifact, err := r.GetByFormattedId(formattedId)
		if err != nil {
			return artifacts, err
		}
		artifacts = append(artifacts, artifact)
	}

	return artifacts, nil
}
