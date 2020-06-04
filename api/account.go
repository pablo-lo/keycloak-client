package api

import (
	"github.com/cloudtrust/keycloak-client"
	"gopkg.in/h2non/gentleman.v2/plugins/body"
	"gopkg.in/h2non/gentleman.v2/plugins/headers"
	"gopkg.in/h2non/gentleman.v2/plugins/url"
)

const (
	accountExtensionAPIPath            = "/auth/realms/master/api/account/realms/:realm"
	accountExecuteActionsEmail         = accountExtensionAPIPath + "/execute-actions-email"
	accountCredentialsPath             = accountExtensionAPIPath + "/credentials"
	accountPasswordPath                = accountCredentialsPath + "/password"
	accountCredentialsRegistratorsPath = accountCredentialsPath + "/registrators"
	accountCredentialIDPath            = accountCredentialsPath + "/:credentialID"
	accountCredentialLabelPath         = accountCredentialIDPath + "/label"
	accountMoveFirstPath               = accountCredentialIDPath + "/moveToFirst"
	accountMoveAfterPath               = accountCredentialIDPath + "/moveAfter/:previousCredentialID"
)

var (
	hdrAcceptJSON           = headers.Set("Accept", "application/json")
	hdrContentTypeTextPlain = headers.Set("Content-Type", "text/plain")
)

// GetCredentials returns the list of credentials of the user
func (c *AccountClient) GetCredentials(accessToken string, realmName string) ([]keycloak.CredentialRepresentation, error) {
	var resp = []keycloak.CredentialRepresentation{}
	var err = c.client.get(accessToken, &resp, url.Path(accountCredentialsPath), url.Param("realm", realmName), hdrAcceptJSON)
	return resp, err
}

// GetCredentialRegistrators returns list of credentials types available for the user
func (c *AccountClient) GetCredentialRegistrators(accessToken string, realmName string) ([]string, error) {
	var resp = []string{}
	var err = c.client.get(accessToken, &resp, url.Path(accountCredentialsRegistratorsPath), url.Param("realm", realmName), hdrAcceptJSON)
	return resp, err
}

// UpdateLabelCredential updates the label of credential
func (c *AccountClient) UpdateLabelCredential(accessToken string, realmName string, credentialID string, label string) error {
	return c.client.put(accessToken, url.Path(accountCredentialLabelPath), url.Param("realm", realmName), url.Param("credentialID", credentialID), body.String(label), hdrAcceptJSON, hdrContentTypeTextPlain)
}

// DeleteCredential deletes the credential
func (c *AccountClient) DeleteCredential(accessToken string, realmName string, credentialID string) error {
	return c.client.delete(accessToken, url.Path(accountCredentialIDPath), url.Param("realm", realmName), url.Param("credentialID", credentialID), hdrAcceptJSON)
}

// MoveToFirst moves the credential at the top of the list
func (c *AccountClient) MoveToFirst(accessToken string, realmName string, credentialID string) error {
	_, err := c.client.post(accessToken, nil, url.Path(accountMoveFirstPath), url.Param("realm", realmName), url.Param("credentialID", credentialID), hdrAcceptJSON)
	return err
}

// MoveAfter moves the credential after the specified one into the list
func (c *AccountClient) MoveAfter(accessToken string, realmName string, credentialID string, previousCredentialID string) error {
	_, err := c.client.post(accessToken, nil, url.Path(accountMoveAfterPath), url.Param("realm", realmName), url.Param("credentialID", credentialID), url.Param("previousCredentialID", previousCredentialID), hdrAcceptJSON)
	return err
}

// UpdatePassword updates the user's password
// Parameters: realm, currentPassword, newPassword, confirmPassword
func (c *AccountClient) UpdatePassword(accessToken, realm, currentPassword, newPassword, confirmPassword string) (string, error) {
	var m = map[string]string{"currentPassword": currentPassword, "newPassword": newPassword, "confirmation": confirmPassword}
	return c.client.post(accessToken, nil, url.Path(accountPasswordPath), url.Param("realm", realm), body.JSON(m))
}

// GetAccount provides the user's information
func (c *AccountClient) GetAccount(accessToken string, realm string) (keycloak.UserRepresentation, error) {
	var resp = keycloak.UserRepresentation{}
	var err = c.client.get(accessToken, &resp, url.Path(accountExtensionAPIPath), url.Param("realm", realm), hdrAcceptJSON)
	return resp, err
}

// UpdateAccount updates the user's information
func (c *AccountClient) UpdateAccount(accessToken string, realm string, user keycloak.UserRepresentation) error {
	_, err := c.client.post(accessToken, nil, url.Path(accountExtensionAPIPath), url.Param("realm", realm), body.JSON(user))
	return err
}

// DeleteAccount delete current user
func (c *AccountClient) DeleteAccount(accessToken string, realmName string) error {
	return c.client.delete(accessToken, url.Path(accountExtensionAPIPath), url.Param("realm", realmName), hdrAcceptJSON)
}

// ExecuteActionsEmail send an email with required actions to the user
func (c *AccountClient) ExecuteActionsEmail(accessToken string, realmName string, actions []string) error {
	return c.client.put(accessToken, url.Path(accountExecuteActionsEmail), url.Param("realm", realmName), body.JSON(actions))
}
