package java

import (
	"fmt"
)

const (
	AccessRuleContainerPath = "/paas/api/v1.1/instancemgmt/%s/services/jaas/instances/%s/accessrules/"
	AccessRuleResourcePath  = "/paas/api/v1.1/instancemgmt/%s/services/jaas/instances/%s/accessrules/%s"
)

// AccessRulesClient is a client for the AccessRules functions of the Compute API.
type AccessRulesClient struct {
	ResourceClient
}

// AccessRules obtains a AccessRulesClient which can be used to access to the
// AccessRules functions of the Compute API
func (c *JavaClient) AccessRules(serviceInstance string) *AccessRulesClient {
	return &AccessRulesClient{
		ResourceClient: ResourceClient{
			JavaClient:       c,
			ContainerPath:    AccessRuleContainerPath,
			ResourceRootPath: AccessRuleResourcePath,
			ServiceInstance:  serviceInstance,
		}}
}

type AccessRuleDestination string

const (
	AccessRuleDestinationWL  AccessRuleDestination = "WLS_ADMIN_SERVER"
	AccessRuleDestinationOTD AccessRuleDestination = "OTD"
)

type AccessRuleProtocol string

const (
	AccessRuleProtocolTCP AccessRuleProtocol = "tcp"
	AccessRuleProtocolUDP AccessRuleProtocol = "udp"
)

type AccessRuleType string

const (
	AccessRuleTypeSystem  AccessRuleType = "SYSTEM"
	AccessRuleTypeDefault AccessRuleType = "DEFAULT"
	AccessRuleTypeUser    AccessRuleType = "USER"
)

type AccessRuleSource string

const (
	AccessRuleSourceWLSAdmin   AccessRuleSource = "WLS_ADMIN_SERVER"
	AccessRuleSourceWLSManaged AccessRuleSource = "WLS_MANAGED_SERVER"
	AccessRuleOTD              AccessRuleSource = "OTD"
	AccessRuleDB               AccessRuleSource = "DB"
)

type AccessRuleStatus string

const (
	AccessRuleStatusEnabled  AccessRuleStatus = "enabled"
	AccessRuleStatusDisabled AccessRuleStatus = "disabled"
)

type AccessRuleOperation string

const (
	AccessRuleOperationUpdate AccessRuleOperation = "update"
	AccessRuleOperationDelete AccessRuleOperation = "delete"
)

type AccessRulesInfo struct {
	AccessRules []AccessRuleInfo `json:"accessRules"`
}

// AccessRuleInfo describes an existing AccessRule.
type AccessRuleInfo struct {
	// Description of the rule.
	Description string `json:"description"`
	// The service component to allow traffic to. For example, WLS_ADMIN_SERVER or OTD.
	Destination AccessRuleDestination `json:"destination"`
	// Ports for the rule. This can be a single port or a port range.
	Ports string `json:"ports"`
	// The name of the AccessRule
	Name string `json:"ruleName"`
	// Type of rule. For example, SYSTEM, DEFAULT, or USER.
	Type AccessRuleType `json:"ruleType"`
	// The hosts from which traffic is allowed. For example, PUBLIC-INTERNET for any host on the Internet, a
	// single IP address or a comma-separated list of subnets (in CIDR format) or IPv4 addresses, or a service
	// component name such as WLS_ADMIN_SERVER, WLS_MANAGED_SERVER, OTD, or DB.
	Source AccessRuleSource `json:"source"`
	// Status of the rule
	Status AccessRuleStatus `json:"status"`
}

// CreateAccessRuleInput defines a AccessRule to be created.
type CreateAccessRuleInput struct {
	// Description of the rule.
	// Required.
	Description string `json:"description"`
	// Destination network. Specify the service component to allow traffic to. For example, WLS_ADMIN_SERVER for the
	// virtual machine where the WebLogic Administration Server is running, or OTD for the virtual machine that contains
	// the Oracle Traffic Director administration server.
	// Required
	Destination AccessRuleDestination `json:"destination"`
	// Network port. Specify a single port or a port range. For example, 8989 or 7000-8000.
	// Required.``
	Ports string `json:"ports"`
	// Communication protocol. Valid values are: tcp or udp.
	// Default is tcp
	// Optional
	Protocol AccessRuleProtocol `json:"string,omitempty"`
	// Name of the access rule
	// Required
	Name string `json:"ruleName"`
	// Network address of source. Specify the hosts from which traffic is allowed. Valid values include:
	// PUBLIC-INTERNET for any host on the Internet
	// A single IP address or a comma-separated list of subnets (in CIDR format) or IPv4 addresses
	// A service component name. Valid values include WLS_ADMIN_SERVER, WLS_MANAGED_SERVER, OTD, DB
	// Required
	Source AccessRuleSource `json:"source"`
	// Status of the rule. Specify whether the status should be enabled or disabled. Valid value: disabled or enabled.
	// Required
	Status AccessRuleStatus `json:"status"`
}

// CreateAccessRule creates a new AccessRule.
func (c *AccessRulesClient) CreateAccessRule(createInput *CreateAccessRuleInput) (*AccessRuleInfo, error) {
	var accessRuleInfo AccessRuleInfo
	if err := c.createResource(createInput, &accessRuleInfo); err != nil {
		return nil, err
	}

	getInput := &GetAccessRuleInput{
		Name: createInput.Name,
	}

	return c.GetAccessRule(getInput)
}

// GetAccessRuleInput describes the AccessRule to get
type GetAccessRuleInput struct {
	// The name of the AccessRule to query for
	// Required
	Name string `json:"name"`
}

// GetAccessRule retrieves the AccessRule with the given name.
func (c *AccessRulesClient) GetAccessRule(getInput *GetAccessRuleInput) (*AccessRuleInfo, error) {
	// The Oracle API does not have a way to return a specific access rule so we need to retrieve all of the access rules
	// and iterate through until we find the one we need
	var (
		accessRulesInfo AccessRulesInfo
		accessRuleInfo  AccessRuleInfo
	)
	if err := c.getResource("", &accessRulesInfo); err != nil {
		return nil, err
	}

	for _, accessRule := range accessRulesInfo.AccessRules {
		if accessRule.Name == getInput.Name {
			accessRuleInfo = accessRule
			break
		}
	}

	if accessRuleInfo.Name == "" {
		return nil, fmt.Errorf("Could not find access rule: %s", getInput.Name)
	}

	return &accessRuleInfo, nil
}

// UpdateAccessRuleInput describes a secruity rule to update
type UpdateAccessRuleInput struct {
	// Name of the access rule to update
	// Required
	Name string `json:"-"`
	// Type of operation to perform on the access rule. Valid values are: update (to disable or enable a rule)
	// and delete (to delete a rule).
	// Required
	Operation AccessRuleOperation `json:"operation"`
	// State of the access rule to update to. This attribute is required only when you disable or enable a rule.
	// Valid value is disable or enable.
	// Required
	Status AccessRuleStatus `json:"status"`
}

// UpdateAccessRule modifies the properties of the AccessRule with the given name.
func (c *AccessRulesClient) UpdateAccessRule(updateInput *UpdateAccessRuleInput) (*AccessRuleInfo, error) {
	var accessRuleInfo AccessRuleInfo
	if err := c.updateResource(updateInput.Name, updateInput, &accessRuleInfo); err != nil {
		return nil, err
	}

	return &accessRuleInfo, nil
}

type DeleteAccessRuleInput struct {
	// Name of the access rule to delete
	// Required
	Name string `json:"-"`
}

func (c *AccessRulesClient) DeleteAccessRule(deleteInput *DeleteAccessRuleInput) error {
	// The Oracle API does not have a DELETE call and so we must go through the UPDATE call
	updateAccessRule := &UpdateAccessRuleInput{
		Name:      deleteInput.Name,
		Operation: AccessRuleOperationDelete,
		Status:    AccessRuleStatusDisabled,
	}
	_, err := c.UpdateAccessRule(updateAccessRule)
	if err != nil {
		return err
	}
	return nil
}
