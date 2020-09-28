package internal

import (
	"sync"

	"github.com/pinpt/agent/v4/sdk"
)

// JiraIntegration is an integration for Jira
// easyjson:skip
type JiraIntegration struct {
	logger      sdk.Logger
	config      sdk.Config
	manager     sdk.Manager
	httpmanager sdk.HTTPClientManager
	client      sdk.GraphQLClient
	lock        sync.Mutex
}

var _ sdk.Integration = (*JiraIntegration)(nil)

// Start is called when the integration is starting up
func (i *JiraIntegration) Start(logger sdk.Logger, config sdk.Config, manager sdk.Manager) error {
	i.logger = sdk.LogWith(logger, "pkg", "jira")
	i.config = config
	i.manager = manager
	i.httpmanager = manager.HTTPManager()
	sdk.LogInfo(i.logger, "starting")
	return nil
}

// Enroll is called when a new integration instance is added
func (i *JiraIntegration) Enroll(instance sdk.Instance) error {
	authConfig, err := i.createAuthConfig(&instance)
	if err != nil {
		return err
	}
	if err := i.installWebHookIfNecessary(i.logger, instance.Config(), instance.State(), authConfig, instance.CustomerID(), instance.IntegrationInstanceID()); err != nil {
		return err
	}
	return nil
}

// Dismiss is called when an existing integration instance is removed
func (i *JiraIntegration) Dismiss(instance sdk.Instance) error {
	authConfig, err := i.createAuthConfig(&instance)
	if err != nil {
		return err
	}
	if err := i.uninstallWebHookIfNecessary(i.logger, instance.Config(), instance.State(), authConfig, instance.CustomerID(), instance.IntegrationInstanceID()); err != nil {
		return err
	}
	return nil
}

// Stop is called when the integration is shutting down for cleanup
func (i *JiraIntegration) Stop() error {
	sdk.LogInfo(i.logger, "stopping")
	return nil
}
