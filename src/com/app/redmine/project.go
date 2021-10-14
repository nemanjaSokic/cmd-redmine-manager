package redmine

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"
)

type projectResult struct {
	Project Project `json:"project"`
}

type projectsResult struct {
	Projects []Project `json:"projects"`
}

type Project struct {
	Id           int            `json:"id"`
	Parent       IdName         `json:"parent"`
	Name         string         `json:"name"`
	Identifier   string         `json:"identifier"`
	Description  string         `json:"description"`
	CreatedOn    string         `json:"created_on"`
	UpdatedOn    string         `json:"updated_on"`
	CustomFields []*CustomField `json:"custom_fields,omitempty"`
}

func (c *Client) Project(id int) (*Project, error) {
	res, err := c.Get(c.endpoint + "/projects/" + strconv.Itoa(id) + ".json")
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)
	var r projectResult
	if res.StatusCode != 200 {
		var er errorsResult
		err = decoder.Decode(&er)
		if err == nil {
			err = errors.New(strings.Join(er.Errors, "\n"))
		}
	} else {
		err = decoder.Decode(&r)
	}
	if err != nil {
		return nil, err
	}
	return &r.Project, nil
}

func (c *Client) Projects() ([]Project, error) {
	res, err := c.Get(c.endpoint + "/projects.json" + c.getPaginationClause())
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)
	var r projectsResult
	if res.StatusCode != 200 {
		var er errorsResult
		err = decoder.Decode(&er)
		if err == nil {
			err = errors.New(strings.Join(er.Errors, "\n"))
		}
	} else {
		err = decoder.Decode(&r)
	}
	if err != nil {
		return nil, err
	}
	return r.Projects, nil
}
