package redmine

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"net/url"
	"fmt"
)

type userResult struct {
	User User `json:"user"`
}

type usersResult struct {
	Users []User `json:"users"`
}

type User struct {
	Id          int          `json:"id"`
	Login       string       `json:"login"`
	Firstname   string       `json:"firstname"`
	Lastname    string       `json:"lastname"`
	Mail        string       `json:"mail"`
	CreatedOn   string       `json:"created_on"`
	LatLoginOn  string       `json:"last_login_on"`
	Memberships []Membership `json:"memberships"`
	ApiKey      string       `json:"api_key"`
}

func (c *Client) UsersByFilter(filter url.Values) ([]User, error) {
	var filterStr string;
	if (len(filter) > 0) {
		filterStr = "&"+filter.Encode()
	}

	res, err := c.Get(c.endpoint + "/users.json?key=" + c.apikey+filterStr)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)
	var r usersResult
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
	return r.Users, nil
}

func (c *Client) Users() ([]User, error) {	// For backward compatibility
	return c.UsersByFilter(url.Values{})
}

func (c *Client) UserByLogin(login string) (*User, error) {
	users,err := c.UsersByFilter( url.Values { "name": []string{login} } )
	if (err != nil) {
		return nil, err
	}

	if (len(users) == 0) {
		return nil, fmt.Errorf("Cannot find user with login \"%v\"", login)
	}

	return &users[0], nil
}

func (c *Client) User(id int) (*User, error) {
	res, err := c.Get(c.endpoint + "/users/" + strconv.Itoa(id) + ".json?key=" + c.apikey)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)
	var r userResult
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
	return &r.User, nil
}
