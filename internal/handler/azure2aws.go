package handler

import (
	"context"
	"fmt"

	"go.dfds.cloud/aad-finout-sync/internal/azure"
	"go.dfds.cloud/aad-finout-sync/internal/config"
	"go.dfds.cloud/aad-finout-sync/internal/util"
	"go.uber.org/zap"
)

const AzureAdToAwsName = "aadToAws"

func Azure2AwsHandler(ctx context.Context) error {
	conf, err := config.LoadConfig()
	if err != nil {
		return err
	}

	azClient := azure.NewAzureClient(azure.Config{
		TenantId:     conf.Azure.TenantId,
		ClientId:     conf.Azure.ClientId,
		ClientSecret: conf.Azure.ClientSecret,
	})

	appRoles, err := azClient.GetApplicationRoles(conf.Azure.ApplicationId)
	if err != nil {
		return err
	}

	appRoleId, err := appRoles.GetRoleId("User")
	if err != nil {
		return err
	}

	appAssignments, err := azClient.GetAssignmentsForApplication(conf.Azure.ApplicationObjectId)
	if err != nil {
		return err
	}

	groups, err := azClient.GetGroups(azure.AZURE_CAPABILITY_GROUP_PREFIX)
	if err != nil {
		return err
	}

	for _, group := range groups.Value {
		select {
		case <-ctx.Done():
			util.Logger.Info("Job cancelled", zap.String("jobName", AzureAdToAwsName))
			return nil
		default:
		}

		util.Logger.Debug(group.DisplayName, zap.String("jobName", AzureAdToAwsName))

		// If group is not already assigned to enterprise application, assign them.
		if !appAssignments.ContainsGroup(group.DisplayName) {
			util.Logger.Info(fmt.Sprintf("Group %s has not been assigned to application yet, assigning", group.DisplayName), zap.String("jobName", AzureAdToAwsName))
			_, err := azClient.AssignGroupToApplication(conf.Azure.ApplicationObjectId, group.ID, appRoleId)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
