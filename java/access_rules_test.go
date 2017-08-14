package java

import (
	"fmt"
	"testing"

	"github.com/hashicorp/go-oracle-terraform/database"
	"github.com/hashicorp/go-oracle-terraform/helper"
	"github.com/hashicorp/go-oracle-terraform/opc"
	"github.com/kylelemons/godebug/pretty"
)

const (
	_AccessRuleDescription  = "Acceptance Test Description"
	_AccessRuleDestination  = AccessRuleDestinationWL
	_AccessRulePorts        = "8989"
	_AccessRuleName         = "test-acc-access-rule-1"
	_AccessRuleSource       = AccessRuleSourceWLSManaged
	_AccessRuleStatus       = AccessRuleStatusEnabled
	_AccessRuleUpdateStatus = AccessRuleStatusDisabled
	_AccessRuleOperation    = AccessRuleOperationUpdate
)

func TestAccAccessRuleLifeCycle(t *testing.T) {
	helper.Test(t, helper.TestCase{})

	arClient, siClient, dClient, err := getAccessRuleTestClients()
	if err != nil {
		t.Fatal(err)
	}

	databaseParameter := database.ParameterInput{
		AdminPassword:                   _ServiceInstanceDBAPassword,
		BackupDestination:               _ServiceInstanceBackupDestinationBoth,
		SID:                             _ServiceInstanceDBSID,
		Type:                            _ServiceInstanceDBType,
		UsableStorage:                   _ServiceInstanceUsableStorage,
		CloudStorageContainer:           _ServiceInstanceDBCloudStorageContainer,
		CreateStorageContainerIfMissing: _ServiceInstanceCloudStorageCreateIfMissing,
	}

	createDatabaseServiceInstance := &database.CreateServiceInstanceInput{
		Name:             _ServiceInstanceDatabaseName,
		Edition:          _ServiceInstanceEdition,
		Level:            _ServiceInstanceLevel,
		Shape:            _ServiceInstanceShape,
		SubscriptionType: _ServiceInstanceSubscriptionType,
		Version:          _ServiceInstanceDBVersion,
		VMPublicKey:      _ServiceInstancePubKey,
		Parameter:        databaseParameter,
	}

	_, err = dClient.CreateServiceInstance(createDatabaseServiceInstance)
	if err != nil {
		t.Fatal(err)
	}
	defer destroyDatabaseServiceInstance(t, dClient, _ServiceInstanceDatabaseName)

	parameter := Parameter{
		Type:          _ServiceInstanceType,
		DBAName:       _ServiceInstanceDBAUser,
		DBAPassword:   _ServiceInstanceDBAPassword,
		DBServiceName: _ServiceInstanceDatabaseName,
		Shape:         _ServiceInstanceShape,
		Version:       _ServiceInstanceVersion,
		AdminUsername: _ServiceInstanceAdminUsername,
		AdminPassword: _ServiceInstanceAdminPassword,
		VMsPublicKey:  _ServiceInstancePubKey,
	}

	createServiceInstance := &CreateServiceInstanceInput{
		CloudStorageContainer:           _ServiceInstanceCloudStorageContainer,
		CreateStorageContainerIfMissing: _ServiceInstanceCloudStorageCreateIfMissing,
		ServiceName:                     _ServiceInstanceName,
		Level:                           _ServiceInstanceLevel,
		SubscriptionType:                _ServiceInstanceSubscriptionType,
		Parameters:                      []Parameter{parameter},
	}

	_, err = siClient.CreateServiceInstance(createServiceInstance)
	if err != nil {
		t.Fatal(err)
	}
	defer destroyServiceInstance(t, siClient, _ServiceInstanceName)

	createAccessRule := &CreateAccessRuleInput{
		Description: _AccessRuleDescription,
		Destination: _AccessRuleDestination,
		Ports:       _AccessRulePorts,
		Name:        _AccessRuleName,
		Source:      _AccessRuleSource,
		Status:      _AccessRuleStatus,
	}

	createdAccessRule, err := arClient.CreateAccessRule(createAccessRule)
	if err != nil {
		t.Fatal(err)
	}
	defer destroyAccessRule(t, arClient, _AccessRuleName)

	getAccessRuleInput := &GetAccessRuleInput{
		Name: _AccessRuleName,
	}
	receivedAccessRule, err := arClient.GetAccessRule(getAccessRuleInput)
	if err != nil {
		t.Fatal(err)
	}

	if diff := pretty.Compare(receivedAccessRule, createdAccessRule); diff != "" {
		t.Fatalf(fmt.Sprintf("Result Diff (-got +want)\n%s", diff))
	}

	updateAccessRuleInput := &UpdateAccessRuleInput{
		Name:      _AccessRuleName,
		Operation: _AccessRuleOperation,
		Status:    _AccessRuleUpdateStatus,
	}
	updatedAccessRule, err := arClient.UpdateAccessRule(updateAccessRuleInput)
	if err != nil {
		t.Fatal(err)
	}
	receivedAccessRule, err = arClient.GetAccessRule(getAccessRuleInput)
	if err != nil {
		t.Fatal(err)
	}

	if diff := pretty.Compare(receivedAccessRule, updatedAccessRule); diff != "" {
		t.Fatalf(fmt.Sprintf("Result Diff (-got +want)\n%s", diff))
	}

}

func getAccessRuleTestClients() (*AccessRulesClient, *ServiceInstanceClient, *database.ServiceInstanceClient, error) {
	client, err := getJavaTestClient(&opc.Config{})
	if err != nil {
		return &AccessRulesClient{}, &ServiceInstanceClient{}, &database.ServiceInstanceClient{}, err
	}

	dClient, err := database.GetDatabaseTestClient(&opc.Config{})
	if err != nil {
		return &AccessRulesClient{}, &ServiceInstanceClient{}, &database.ServiceInstanceClient{}, err
	}

	return client.AccessRules(_ServiceInstanceName), client.ServiceInstanceClient(), dClient.ServiceInstanceClient(), nil
}

func destroyAccessRule(t *testing.T, client *AccessRulesClient, name string) {
	input := &DeleteAccessRuleInput{
		Name: name,
	}

	if err := client.DeleteAccessRule(input); err != nil {
		t.Fatal(err)
	}
}
