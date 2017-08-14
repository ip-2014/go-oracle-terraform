package java

import (
	"fmt"
	"net/http"

	"github.com/hashicorp/go-oracle-terraform/client"
	"github.com/hashicorp/go-oracle-terraform/opc"
)

const AUTH_HEADER = "Authorization"
const TENANT_HEADER = "X-ID-TENANT-NAME"
const JAVA_QUALIFIED_NAME = "%s%s/%s"

// Client represents an authenticated java client, with compute credentials and an api client.
type JavaClient struct {
	client     *client.Client
	authHeader *string
}

func NewJavaClient(c *opc.Config) (*JavaClient, error) {
	javaClient := &JavaClient{}
	client, err := client.NewClient(c)
	if err != nil {
		return nil, err
	}
	javaClient.client = client

	javaClient.authHeader = javaClient.getAuthenticationHeader()

	return javaClient, nil
}

func (c *JavaClient) executeRequest(method, path string, contentType string, body interface{}) (*http.Response, error) {
	req, err := c.client.BuildRequest(method, path, body)
	if err != nil {
		return nil, err
	}

	debugReqString := fmt.Sprintf("HTTP %s Path (%s)", method, path)
	if body != nil {
		req.Header.Set("Content-Type", contentType)
	}
	// Log the request before the authentication header, so as not to leak credentials
	c.client.DebugLogString(debugReqString)
	c.client.DebugLogString(fmt.Sprintf("Req (%+v)", req))

	// Set the authentiation headers
	req.Header.Add(AUTH_HEADER, *c.authHeader)
	req.Header.Add(TENANT_HEADER, *c.client.IdentityDomain)

	resp, err := c.client.ExecuteRequest(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// Each resource except service instance needs a service instance to attach to.
// This is checking whether or not it's a service instance or not and creating the approriate path
func (c *JavaClient) getContainerPath(root, serviceInstance string) string {
	if serviceInstance == "" {
		return fmt.Sprintf(root, *c.client.IdentityDomain)
	}
	return fmt.Sprintf(root, *c.client.IdentityDomain, serviceInstance)
}

// Each resource except service instance needs a service instance to attach to.
// This is checking whether or not it's a service instance or not and creating the approriate path
func (c *JavaClient) getObjectPath(root, serviceInstance, name string) string {
	if serviceInstance == "" {
		return fmt.Sprintf(root, *c.client.IdentityDomain, name)
	}
	return fmt.Sprintf(root, *c.client.IdentityDomain, serviceInstance, name)
}
