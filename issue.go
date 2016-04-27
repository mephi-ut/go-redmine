package redmine

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type issueRequest struct {
	Issue Issue `json:"issue"`
}

type issueResult struct {
	Issue Issue `json:"issue"`
}

type issuesResult struct {
	Issues []Issue `json:"issues"`
}

type Issue struct {
	Id          int     `json:"id,omitempty"`
	Subject     string  `json:"subject,omitempty"`
	Description string  `json:"description,omitempty"`
	ProjectId   int     `json:"project_id,omitempty"`
	Project     *IdName `json:"project,omitempty"`
	Tracker     *IdName `json:"tracker,omitempty"`
	StatusId    int     `json:"status_id,omitempty"`
	Status      *IdName `json:"status,omitempty"`
	Priority    *IdName `json:"priority,omitempty"`
	Author      *IdName `json:"author,omitempty"`
	AssignedTo  *IdName `json:"assigned_to,omitempty"`
	Notes       string  `json:"notes,omitempty"`
	StatusDate  string  `json:"status_date,omitempty"`
	CreatedOn   string  `json:"created_on,omitempty"`
	UpdatedOn   string  `json:"updated_on,omitempty"`
	CustomFields []*CustomField `json:"custom_fields,omitempty"`
}

type CustomField struct {
	Id      int     `json:"id,omitempty"`
	Name	string  `json:"name,omitempty"`
	Value	string  `json:"value"`
}

func (c *Client) IssuesOf(projectId int) ([]Issue, error) {
	res, err := c.Get(c.endpoint + "/issues.json?project_id=" + strconv.Itoa(projectId) + "&key=" + c.apikey)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)
	var r issuesResult
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
	return r.Issues, nil
}

func (c *Client) Issue(id int) (*Issue, error) {
	res, err := c.Get(c.endpoint + "/issues/" + strconv.Itoa(id) + ".json?key=" + c.apikey)
	if res.StatusCode == 404 {
		return nil, errors.New("Not Found")
	}
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)
	var r issueRequest
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
	return &r.Issue, nil
}

func (c *Client) IssuesByQuery(query_id int) ([]Issue, error) {
	res, err := http.Get(c.endpoint + "/issues.json?query_id=" + strconv.Itoa(query_id) + "&key=" + c.apikey)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)
	var r issuesResult
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
	return r.Issues, nil
}

func (c *Client) Issues() ([]Issue, error) {
	res, err := c.Get(c.endpoint + "/issues.json?key=" + c.apikey)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)
	var r issuesResult
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
	return r.Issues, nil
}

func (c *Client) CreateIssue(issue Issue) (*Issue, error) {
	var ir issueRequest
	ir.Issue = issue
	s, err := json.Marshal(ir)
	if err != nil {
		return nil, fmt.Errorf("Got error from json.Marshal(ir): %v", err.Error())
	}
	uri      := c.endpoint+"/issues.json?";
	req, err := http.NewRequest("POST", uri+"key="+c.apikey, strings.NewReader(string(s)))
	if err != nil {
		return nil, fmt.Errorf("Got error from http.NewRequest(\"POST\", \"%vkey=********\", \"%v\")", uri, s, err.Error())
	}
	req.Header.Set("Content-Type", "application/json")
	res, err := c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Got error from c.Do(req): %v", err.Error())
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)
	var r issueRequest
	if res.StatusCode != 201 {
		var er errorsResult
		err = decoder.Decode(&er)
		if err == nil {
			err = errors.New(strings.Join(er.Errors, "\n"))
		}
	} else {
		err = decoder.Decode(&r)
	}
	if err != nil {
		return nil, fmt.Errorf("Got error from decoder.Decode() [res.StatusCode == %v]: %v", res.StatusCode, err.Error())
	}
	return &r.Issue, nil
}

func (c *Client) UpdateIssue(issue Issue) error {
	var ir issueRequest
	ir.Issue = issue
	s, err := json.Marshal(ir)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("PUT", c.endpoint+"/issues/"+strconv.Itoa(issue.Id)+".json?key="+c.apikey, strings.NewReader(string(s)))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	res, err := c.Do(req)
	if res.StatusCode == 404 {
		return errors.New("Not Found")
	}
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		decoder := json.NewDecoder(res.Body)
		var er errorsResult
		err = decoder.Decode(&er)
		if err == nil {
			err = errors.New(strings.Join(er.Errors, "\n"))
		}
	}
	if err != nil {
		return err
	}
	return err
}

func (c *Client) DeleteIssue(id int) error {
	req, err := http.NewRequest("DELETE", c.endpoint+"/issues/"+strconv.Itoa(id)+".json?key="+c.apikey, strings.NewReader(""))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	res, err := c.Do(req)
	if res.StatusCode == 404 {
		return errors.New("Not Found")
	}
	if err != nil {
		return err
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)
	if res.StatusCode != 200 {
		var er errorsResult
		err = decoder.Decode(&er)
		if err == nil {
			err = errors.New(strings.Join(er.Errors, "\n"))
		}
	}
	return err
}

func (issue *Issue) GetTitle() string {
	return issue.Tracker.Name + " #" + strconv.Itoa(issue.Id) + ": " + issue.Subject
}
