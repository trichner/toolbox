package jira

import (
	"fmt"

	gojira "gopkg.in/andygrunwald/go-jira.v1"
)

type CreateUser struct {
	Name   string
	Email  string
	Groups []string
}

func (j *JiraService) CreateUser(user *CreateUser) (string, error) {
	u, _, err := j.client.User.Create(&gojira.User{
		EmailAddress: user.Email,
		DisplayName:  user.Name,
	})
	if err != nil {
		return "", err
	}

	for _, group := range user.Groups {
		err := j.AddUserToGroup(u.AccountID, group)
		if err != nil {
			_, err2 := j.client.User.Delete(u.AccountID)
			if err2 != nil {
				err = fmt.Errorf("rollback failed for user %q: %w", user, err)
			}
			return "", err
		}
	}

	return u.AccountID, nil
}

func (j *JiraService) AddUserToGroup(accountId, group string) error {
	apiEndpoint := fmt.Sprintf("/rest/api/2/group/user?groupname=%s", group)
	var user struct {
		AccountId string `json:"accountId"`
	}
	user.AccountId = accountId
	req, err := j.client.NewRequest("POST", apiEndpoint, &user)
	if err != nil {
		return err
	}

	responseGroup := new(gojira.Group)
	_, err = j.client.Do(req, responseGroup)
	if err != nil {
		return err
	}

	return nil
}
