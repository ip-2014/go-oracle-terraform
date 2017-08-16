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
	ContentType      string
}

func (c *ResourceClient) createResource(serviceInstance string, requestBody interface{}, responseBody interface{}) error {
	_, err := c.executeRequest("POST", c.getContainerPath(c.ContainerPath, serviceInstance), c.ContentType, requestBody)
	if err != nil {
		return err
	}

	return nil
}

func (c *ResourceClient) updateResource(name string, serviceInstance string, requestBody interface{}, responseBody interface{}) error {
	_, err := c.executeRequest("PUT", c.getObjectPath(c.ResourceRootPath, serviceInstance, name), c.ContentType, requestBody)
	if err != nil {
		return err
	}

	return nil
}

func (c *ResourceClient) getResource(name string, serviceInstance string, responseBody interface{}) error {
	var objectPath string
	if name != "" {
		objectPath = c.getObjectPath(c.ResourceRootPath, serviceInstance, name)
	} else {
		objectPath = c.getContainerPath(c.ContainerPath, serviceInstance)
	}
	resp, err := c.executeRequest("GET", objectPath, c.ContentType, nil)
	if err != nil {
		return err
	}

	return c.unmarshalResponseBody(resp, responseBody)
}

func (c *ResourceClient) deleteResource(name, serviceInstance string) error {
	var objectPath string
	if name != "" {
		objectPath = c.getObjectPath(c.ResourceRootPath, serviceInstance, name)
	} else {
		objectPath = c.ResourceRootPath
	}
	_, err := c.executeRequest("DELETE", objectPath, c.ContentType, nil)
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
	_, err := c.executeRequest("PUT", objectPath, c.ContentType, requestBody)
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
