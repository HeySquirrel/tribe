package rally

import (
	"encoding/json"
	"fmt"
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

func New(apikey string) *Rally {
	rally := new(Rally)
	rally.apikey = apikey

	return rally
}

func (r *Rally) GetByFormattedId(formattedIds ...string) ([]Artifact, error) {
	artifacts := make([]Artifact, 0)

	client := &http.Client{}

	for _, formattedId := range formattedIds {
		req, _ := http.NewRequest("GET", "https://rally1.rallydev.com/slm/webservice/v2.0/Artifact", nil)
		req.Header.Set("zsessionid", r.apikey)

		q := req.URL.Query()
		q.Add("query", fmt.Sprintf("(FormattedID = %s)", formattedId))
		q.Add("fetch", "FormattedID,Name,Description")
		req.URL.RawQuery = q.Encode()

		res, err := client.Do(req)
		if err != nil {
			return artifacts, err
		}
		defer res.Body.Close()

		var queryResult Result
		json.NewDecoder(res.Body).Decode(&queryResult)

		artifacts = append(artifacts, queryResult.QueryResult.Results...)
	}

	return artifacts, nil
}
