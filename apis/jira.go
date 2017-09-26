package apis

import (
	"encoding/json"
	"fmt"
	"github.com/heysquirrel/tribe/config"
	"net/http"
	"strings"
)

type jira struct {
	host     string
	username string
	password string
}

type JiraResult struct {
	Issue Issue `json:"fields"`
}

type IssueType struct {
	Name string `json:"name"`
}

type Issue struct {
	Key         string    `json:"key"`
	Summary     string    `json:"summary"`
	Description string    `json:"description"`
	IssueType   IssueType `json:"issuetype"`
}

func (i *Issue) GetType() string        { return i.IssueType.Name }
func (i *Issue) GetName() string        { return i.Summary }
func (i *Issue) GetDescription() string { return i.Description }
func (i *Issue) GetId() string          { return i.Key }

func NewJiraFromConfig(servername string) (*jira, error) {
	serverconfig := config.WorkItemServer(config.ServerName(servername))

	return NewJira(serverconfig["host"], serverconfig["username"], serverconfig["password"])
}

func NewJira(host, username, password string) (*jira, error) {
	if strings.TrimSpace(host) == "" {
		return nil, fmt.Errorf("Invalid hostname: '%s'", host)
	}

	if strings.TrimSpace(username) == "" {
		return nil, fmt.Errorf("Invalid username: '%s'", username)
	}

	if strings.TrimSpace(password) == "" {
		return nil, fmt.Errorf("Invalid password: '%s'", password)
	}

	return &jira{host, username, password}, nil
}

func (j *jira) GetWorkItem(id string) (WorkItem, error) {
	client := &http.Client{}
	url := fmt.Sprintf("%s/rest/api/2/issue/%s", j.host, id)

	req, _ := http.NewRequest("GET", url, nil)
	req.SetBasicAuth(j.username, j.password)

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == 404 {
		return NullWorkItem(id), ItemNotFoundError(id)
	}

	var result JiraResult
	json.NewDecoder(res.Body).Decode(&result)

	return &result.Issue, nil
}
