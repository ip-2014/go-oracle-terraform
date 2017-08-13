package java

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/mitchellh/mapstructure"
)

// ResourceClient is an AuthenticatedClient with some additional information about the resources to be addressed.
type ResourceClient struct {
	*JavaClient
	ContainerPath    string
	ResourceRootPath string
	ServiceInstance  string
}

func (c *ResourceClient) createResource(requestBody interface{}, responseBody interface{}) error {
	_, err := c.executeRequest("POST", c.getContainerPath(c.ContainerPath, c.ServiceInstance), requestBody)
	if err != nil {
		return err
	}

	return nil
}

func (c *ResourceClient) updateResource(name string, requestBody interface{}, responseBody interface{}) error {
	_, err := c.executeRequest("PUT", c.getObjectPath(c.ResourceRootPath, c.ServiceInstance, name), requestBody)
	if err != nil {
		return err
	}

	return nil
}

func (c *ResourceClient) getResource(name string, responseBody interface{}) error {
	var objectPath string
	if name != "" {
		objectPath = c.getObjectPath(c.ResourceRootPath, c.ServiceInstance, name)
	} else {
		objectPath = c.ResourceRootPath
	}
	resp, err := c.executeRequest("GET", objectPath, nil)
	if err != nil {
		return err
	}

	return c.unmarshalResponseBody(resp, responseBody)
}

func (c *ResourceClient) deleteResource(name string) error {
	var objectPath string
	if name != "" {
		objectPath = c.getObjectPath(c.ResourceRootPath, c.ServiceInstance, name)
	} else {
		objectPath = c.ResourceRootPath
	}
	_, err := c.executeRequest("DELETE", objectPath, nil)
	if err != nil {
		return err
	}

	// No errors and no response body to write
	return nil
}

// ServiceInstance needs a PUT and a body to be destroyed
func (c *ResourceClient) deleteInstanceResource(name string, requestBody interface{}) error {
	var objectPath string
	if name != "" {
		objectPath = c.getObjectPath(c.ResourceRootPath, "", name)
	} else {
		objectPath = c.ResourceRootPath
	}
	_, err := c.executeRequest("PUT", objectPath, requestBody)
	if err != nil {
		return err
	}

	// No errors and no response body to write
	return nil
}

func (c *ResourceClient) unmarshalResponseBody(resp *http.Response, iface interface{}) error {
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	c.client.DebugLogString(fmt.Sprintf("HTTP Resp (%d): %s", resp.StatusCode, buf.String()))
	// JSON decode response into interface
	var tmp interface{}
	dcd := json.NewDecoder(buf)
	if err := dcd.Decode(&tmp); err != nil {
		return fmt.Errorf("%+v", resp)
		return err
	}

	// Use mapstructure to weakly decode into the resulting interface
	msdcd, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		WeaklyTypedInput: true,
		Result:           iface,
		TagName:          "json",
	})
	if err != nil {
		return err
	}

	if err := msdcd.Decode(tmp); err != nil {
		return err
	}
	return nil
}
