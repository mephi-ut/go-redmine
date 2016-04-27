package redmine

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"fmt"
)

type membershipsResult struct {
	Memberships []Membership `json:"memberships"`
}

type membershipResult struct {
	Membership Membership `json:"membership"`
}

type membershipRequest struct {
	Membership Membership `json:"membership"`
}

type Membership struct {
	Id      int      `json:"id,omitempty"`
	Project IdName   `json:"project,omitempty"`
	User    IdName   `json:"user,omitempty"`
	UserId  int      `json:"user_id,omitempty"`
	Roles   []IdName `json:"roles,omitempty"`
	RoleIds []int    `json:"role_ids,omitempty"`
	Groups  []IdName `json:"groups,omitempty"`
}

func (c *Client) Memberships(projectId int) ([]Membership, error) {
	res, err := c.Get(c.endpoint + "/projects/" + strconv.Itoa(projectId) + "/memberships.json?key=" + c.apikey)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)
	var r membershipsResult
	if res.StatusCode == 404 {
		return nil, errors.New("Not Found")
	}
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
	return r.Memberships, nil
}

func (c *Client) Membership(id int) (*Membership, error) {
	res, err := c.Get(c.endpoint + "/memberships/" + strconv.Itoa(id) + ".json?key=" + c.apikey)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)
	var r membershipResult
	if res.StatusCode == 404 {
		return nil, errors.New("Not Found")
	}
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
	return &r.Membership, nil
}

func (c *Client) CreateMembership(membership Membership) (*Membership, error, int) {
	var ir membershipRequest
	ir.Membership = membership
	s, err := json.Marshal(ir)
	if err != nil {
		return nil, fmt.Errorf("Got error from json.Marshal(ir): %v", err.Error()), 0
	}

	project,err := c.Project(membership.Project.Id);
	if err != nil {
		return nil, fmt.Errorf("Got error from c.Project(%v): %v", membership.Project.Id, err.Error()), 0
	}

	url      := c.endpoint+"/projects/"+project.Identifier+"/memberships.json?"
	req, err := http.NewRequest("POST", url+"key="+c.apikey, strings.NewReader(string(s)))
	if err != nil {
		return nil, fmt.Errorf("Got error from http.NewRequest(\"POST\", \"%vkey=********\", \"%v\"): %v", url, s, err.Error()), 0
	}
	req.Header.Set("Content-Type", "application/json")
	res, err := c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Got error from c.Do(req): %v", err.Error()), 0
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)
	var r membershipRequest
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
		return nil, fmt.Errorf("Got error from decoder.Decode() [res.StatusCode == %v]: %v", res.StatusCode, err.Error()), res.StatusCode
	}
	return &r.Membership, nil, res.StatusCode
}

func (c *Client) UpdateMembership(membership Membership) error {
	var ir membershipRequest
	ir.Membership = membership
	s, err := json.Marshal(ir)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("PUT", c.endpoint+"/memberships/"+strconv.Itoa(membership.Id)+".json?key="+c.apikey, strings.NewReader(string(s)))
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

func (c *Client) DeleteMembership(id int) error {
	req, err := http.NewRequest("DELETE", c.endpoint+"/memberships/"+strconv.Itoa(id)+".json?key="+c.apikey, strings.NewReader(""))
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
