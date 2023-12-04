package config

import (
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Aws struct {
		IdentityStoreArn               string `json:"identityStoreArn"`
		SsoInstanceArn                 string `json:"ssoInstanceArn"`
		CapabilityPermissionSetArn     string `json:"capabilityPermissionSetArn"`
		CapabilityLogsPermissionSetArn string `json:"capabilityLogsPermissionSetArn"`
		CapabilityLogsAwsAccountAlias  string `json:"capabilityLogsAwsAccountAlias"`
		SharedEcrPullPermissionSetArn  string `json:"sharedEcrPullPermissionSetArn"`
		SharedEcrPullAwsAccountAlias   string `json:"sharedEcrPullAwsAccountAlias"`
		AccountNamePrefix              string `json:"accountNamePrefix"`
		SsoRegion                      string `json:"ssoRegion"`
		AssumableRoles                 struct {
			SsoManagementArn          string `json:"ssoManagementArn"`
			CapabilityAccountRoleName string `json:"capabilityAccountRoleName"`
		} `json:"assumableRoles"`
		Scim struct {
			Endpoint string `json:"endpoint"`
			Token    string `json:"token"`
		}
		OrganizationsParentId     string `json:"organizationsParentId"`
		RootOrganizationsParentId string `json:"rootOrganizationsParentId"`
	} `json:"aws"`
	Azure struct {
		TenantId            string `json:"tenantId"`
		ClientId            string `json:"clientId"`
		ClientSecret        string `json:"clientSecret"`
		ApplicationId       string `json:"applicationId"`
		ApplicationObjectId string `json:"applicationObjectId"`
	} `json:"azure"`
	CapSvc struct { // Capability-Service
		Host         string `json:"host"`
		TokenScope   string `json:"tokenScope"`
		ClientId     string `json:"clientId"`
		ClientSecret string `json:"clientSecret"`
	} `json:"capSvc"`
	Finout struct {
		Username     string `json:"username"`
		Password     string `json:"password"`
		MfaUrl       string `json:"mfaUrl"`
		ClientId     string `json:"clientId"`
		ClientSecret string `json:"clientSecret"`
	}
	Log struct {
		Level string `json:"level"`
		Debug bool   `json:"debug"`
	}
	EventHandling struct {
		Enabled bool `json:"enable"`
	}
	Scheduler struct {
		Frequency               string `json:"scheduleFrequency" default:"30m"`
		EnableAzure2Finout      bool   `json:"enableAzure2Finout"`
		EnableCostCentre2Finout bool   `json:"enableCostCentre2Finout"`
	}
}

const APP_CONF_PREFIX = "AFS"

func LoadConfig() (Config, error) {
	var conf Config
	err := envconfig.Process(APP_CONF_PREFIX, &conf)

	return conf, err
}
