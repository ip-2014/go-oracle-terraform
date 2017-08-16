package database

const (
	DBAccessContainerPath = "/paas/api/v1.1/instancemgmt/%s/services/dbaas/instances/%s/accessrules"
	DBAccessRootPath      = "/paas/api/v1.1/instancemgmt/%s/services/dbaas/instances/%s/accessrules/%s"
)

type UtilityClient struct {
	UtilityResourceClient
}

func (c *DatabaseClient) AccessRules() *UtilityClient {
	return &UtilityClient{
		UtilityResourceClient: UtilityResourceClient{
			DatabaseClient:   c,
			ContainerPath:    DBAccessContainerPath,
			ResourceRootPath: DBAccessRootPath,
		},
	}
}

type AccessRuleInfo struct {
	ServiceInstanceID string `json:"-"`
	Description       string `json:"description"`
}

type CreateAccessRuleInput struct {
	ServiceInstanceID string `json:"-"`
	Description       string `json:"description"`
}

func (c *UtilityClient) CreateAccessRule(input *CreateAccessRuleInput) (*AccessRuleInfo, error) {
	c.ServiceInstance = input.ServiceInstanceID
	var accessRule AccessRuleInfo
	if err := c.createResource(input, &accessRule); err != nil {
		return nil, err
	}

	return &accessRule, nil
}
