package handler

import (
	"context"
	"go.dfds.cloud/aad-finout-sync/internal/config"
	"go.dfds.cloud/aad-finout-sync/internal/finout"
)

const AzureAdToFinoutName = "aadToFinout"

func Azure2FinoutHandler(ctx context.Context) error {
	conf, err := config.LoadConfig()
	if err != nil {
		return err
	}

	//finoutClientAuth := finout.NewFinoutClient()
	finoutClientApp := finout.NewFinoutClient()
	//finoutClientAuth.SetAuthMethod(finout.AuthUserMethod(conf.Finout.Username, conf.Finout.Password, &conf.Finout.MfaUrl))
	finoutClientApp.SetAuthMethod(finout.AuthClientSecretMethod(finout.Config{ClientId: conf.Finout.ClientId, ClientSecret: conf.Finout.ClientSecret}))

	tags, err := finoutClientApp.ApiApp().ListViews(ctx)
	if err != nil {
		return err
	}

	tags.GetByName("dfds.cost.centre")

	return nil

	//azClient := azure.NewAzureClient(azure.Config{
	//	TenantId:     conf.Azure.TenantId,
	//	ClientId:     conf.Azure.ClientId,
	//	ClientSecret: conf.Azure.ClientSecret,
	//})
	//
	//groups, err := azClient.GetGroups(azure.AZURE_CAPABILITY_GROUP_PREFIX)
	//if err != nil {
	//	return err
	//}
	//
	//departments := make(map[string]int)
	//
	//for _, group := range groups.Value {
	//	select {
	//	case <-ctx.Done():
	//		util.Logger.Info("Job cancelled", zap.String("jobName", AzureAdToFinoutName))
	//		return nil
	//	default:
	//	}
	//
	//	util.Logger.Debug(group.DisplayName, zap.String("jobName", AzureAdToFinoutName))
	//	members, err := azClient.GetGroupMembers(group.ID)
	//	if err != nil {
	//		return err
	//	}
	//
	//	for _, member := range members.Value {
	//		util.Logger.Debug(member.UserPrincipalName, zap.String("jobName", AzureAdToFinoutName))
	//		if _, ok := departments[member.Department]; !ok {
	//			departments[member.Department] = 0
	//		}
	//	}
	//}

	return nil
}
