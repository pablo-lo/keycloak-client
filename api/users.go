package api

import (
	"errors"

	"github.com/cloudtrust/keycloak-client"
	"gopkg.in/h2non/gentleman.v2/plugins/body"
	"gopkg.in/h2non/gentleman.v2/plugins/url"
)

const (
	userPath                   = "/auth/admin/realms/:realm/users"
	usersAdminExtensionAPIPath = "/auth/realms/:realmReq/api/admin/realms/:realm/users"
	userCountPath              = userPath + "/count"
	userIDPath                 = userPath + "/:id"
	userGroupsPath             = userIDPath + "/groups"
	userGroupIDPath            = userGroupsPath + "/:groupId"
	executeActionsEmailPath    = userIDPath + "/execute-actions-email"
	sendReminderEmailPath      = "/auth/realms/:realm/onboarding/sendReminderEmail"
	smsAPI                     = "/auth/realms/:realm/smsApi"
	sendNewEnrolmentCode       = smsAPI + "/sendNewCode"
	shadowUser                 = userIDPath + "/federated-identity/:provider"
)

// GetUsers returns a list of users, filtered according to the query parameters.
// Parameters: email, first (paging offset, int), firstName, lastName, username,
// max (maximum result size, default = 100),
// search (string contained in username, firstname, lastname or email)
func (c *Client) GetUsers(accessToken string, reqRealmName, targetRealmName string, paramKV ...string) (keycloak.UsersPageRepresentation, error) {
	var resp keycloak.UsersPageRepresentation
	if len(paramKV)%2 != 0 {
		return resp, errors.New(keycloak.MsgErrInvalidParam + "." + keycloak.EvenParams)
	}

	var plugins = append(createQueryPlugins(paramKV...), url.Path(usersAdminExtensionAPIPath), url.Param("realmReq", reqRealmName), url.Param("realm", targetRealmName))
	var err = c.get(accessToken, &resp, plugins...)
	return resp, err
}

// CreateUser creates the user from its UserRepresentation. The username must be unique.
func (c *Client) CreateUser(accessToken string, reqRealmName, targetRealmName string, user keycloak.UserRepresentation) (string, error) {
	return c.post(accessToken, nil, url.Path(usersAdminExtensionAPIPath), url.Param("realmReq", reqRealmName), url.Param("realm", targetRealmName), body.JSON(user))
}

// CountUsers returns the number of users in the realm.
func (c *Client) CountUsers(accessToken string, realmName string) (int, error) {
	var resp = 0
	var err = c.get(accessToken, &resp, url.Path(userCountPath), url.Param("realm", realmName))
	return resp, err
}

// GetUser gets the represention of the user.
func (c *Client) GetUser(accessToken string, realmName, userID string) (keycloak.UserRepresentation, error) {
	var resp = keycloak.UserRepresentation{}
	var err = c.get(accessToken, &resp, url.Path(userIDPath), url.Param("realm", realmName), url.Param("id", userID))
	return resp, err
}

// GetGroupsOfUser gets the groups of the user.
func (c *Client) GetGroupsOfUser(accessToken string, realmName, userID string) ([]keycloak.GroupRepresentation, error) {
	var resp = []keycloak.GroupRepresentation{}
	var err = c.get(accessToken, &resp, url.Path(userGroupsPath), url.Param("realm", realmName), url.Param("id", userID))
	return resp, err
}

// AddGroupToUser adds a group to the groups of the user.
func (c *Client) AddGroupToUser(accessToken string, realmName, userID, groupID string) error {
	return c.put(accessToken, url.Path(userGroupIDPath), url.Param("realm", realmName), url.Param("id", userID), url.Param("groupId", groupID))
}

// DeleteGroupFromUser adds a group to the groups of the user.
func (c *Client) DeleteGroupFromUser(accessToken string, realmName, userID, groupID string) error {
	return c.delete(accessToken, url.Path(userGroupIDPath), url.Param("realm", realmName), url.Param("id", userID), url.Param("groupId", groupID))
}

// UpdateUser updates the user.
func (c *Client) UpdateUser(accessToken string, realmName, userID string, user keycloak.UserRepresentation) error {
	return c.put(accessToken, url.Path(userIDPath), url.Param("realm", realmName), url.Param("id", userID), body.JSON(user))
}

// DeleteUser deletes the user.
func (c *Client) DeleteUser(accessToken string, realmName, userID string) error {
	return c.delete(accessToken, url.Path(userIDPath), url.Param("realm", realmName), url.Param("id", userID))
}

// ExecuteActionsEmail sends an update account email to the user. An email contains a link the user can click to perform a set of required actions.
func (c *Client) ExecuteActionsEmail(accessToken string, realmName string, userID string, actions []string, paramKV ...string) error {
	if len(paramKV)%2 != 0 {
		return errors.New(keycloak.MsgErrInvalidParam + "." + keycloak.EvenParams)
	}

	var plugins = append(createQueryPlugins(paramKV...), url.Path(executeActionsEmailPath), url.Param("realm", realmName), url.Param("id", userID), body.JSON(actions))

	return c.put(accessToken, plugins...)
}

// SendNewEnrolmentCode sends a new enrolment code and return it
func (c *Client) SendNewEnrolmentCode(accessToken string, realmName string, userID string) (keycloak.SmsCodeRepresentation, error) {
	var paramKV []string
	paramKV = append(paramKV, "userid", userID)
	var plugins = append(createQueryPlugins(paramKV...), url.Path(sendNewEnrolmentCode), url.Param("realm", realmName))
	var resp = keycloak.SmsCodeRepresentation{}

	_, err := c.post(accessToken, &resp, plugins...)

	return resp, err
}

// SendReminderEmail sends a reminder email to a user
func (c *Client) SendReminderEmail(accessToken string, realmName string, userID string, paramKV ...string) error {
	if len(paramKV)%2 != 0 {
		return errors.New(keycloak.MsgErrInvalidParam + "." + keycloak.EvenParams)
	}
	var newParamKV = append(paramKV, "userid", userID)

	var plugins = append(createQueryPlugins(newParamKV...), url.Path(sendReminderEmailPath), url.Param("realm", realmName))

	_, err := c.post(accessToken, nil, plugins...)
	return err
}

// LinkShadowUser links shadow user to a realm in the context of brokering
func (c *Client) LinkShadowUser(accessToken string, reqRealmName string, userID string, provider string, fedIDKC keycloak.FederatedIdentityRepresentation) error {
	_, err := c.post(accessToken, nil, url.Path(shadowUser), url.Param("realm", reqRealmName), url.Param("id", userID), url.Param("provider", provider), body.JSON(fedIDKC))
	return err
}
