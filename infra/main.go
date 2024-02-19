package main

import (
	"github.com/pulumi/pulumi-azure-native-sdk/resources/v2"
	"github.com/pulumi/pulumi-azure-native-sdk/storage/v2"
	"github.com/pulumi/pulumi-azure-native-sdk/web"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		// Create an Azure Resource Group
		resourceGroup, err := resources.NewResourceGroup(ctx, "resourceGroup", &resources.ResourceGroupArgs{
			ResourceGroupName: pulumi.String("OneCloud-MarioKart-rg"),
		})
		if err != nil {
			return err
		}

		// Create an Azure resource (Storage Account)
		account, err := storage.NewStorageAccount(ctx, "sa", &storage.StorageAccountArgs{
			ResourceGroupName: resourceGroup.Name,
			AccessTier:        storage.AccessTierHot,
			Sku: &storage.SkuArgs{
				Name: storage.SkuName_Standard_LRS,
			},
			Kind: storage.KindStorageV2,
		})
		if err != nil {
			return err
		}
		_, err = storage.NewBlobContainer(ctx, "container", &storage.BlobContainerArgs{
			AccountName:       account.Name,
			ResourceGroupName: resourceGroup.Name,
		})
		if err != nil {
			return err
		}
		fileshare, err := storage.NewFileShare(ctx, "fileshare", &storage.FileShareArgs{
			AccountName:       account.Name,
			ResourceGroupName: resourceGroup.Name,
		})
		if err != nil {
			return err
		}

		primaryStorageKey := pulumi.All(resourceGroup.Name, account.Name).ApplyT(
			func(args []interface{}) (string, error) {
				accountName := args[1].(string)
				accountKeys, err := storage.ListStorageAccountKeys(ctx, &storage.ListStorageAccountKeysArgs{
					ResourceGroupName: "OneCloud-MarioKart-rg",
					AccountName:       accountName,
				})
				if err != nil {
					return "", err
				}

				return accountKeys.Keys[0].Value, nil
			})
		// Export the primary key of the Storage Account
		ctx.Export("primaryStorageKey", primaryStorageKey)
		plan, err := web.NewAppServicePlan(ctx, "asp", &web.AppServicePlanArgs{
			ResourceGroupName: resourceGroup.Name,
			Kind:              pulumi.String("Linux"),
			Reserved:          pulumi.Bool(true),
			Sku: &web.SkuDescriptionArgs{
				Name: pulumi.String("F1"),
				Tier: pulumi.String("Free"),
			},
		})
		if err != nil {
			return err
		}

		_, err = web.NewWebApp(ctx, "webapp", &web.WebAppArgs{
			ResourceGroupName: resourceGroup.Name,
			ServerFarmId:      plan.ID(),
			HttpsOnly:         pulumi.Bool(true),
			SiteConfig: &web.SiteConfigArgs{
				LinuxFxVersion: pulumi.String("DOCKER|bgaechter/bracket:latest"),
				AppSettings: web.NameValuePairArray{
					&web.NameValuePairArgs{
						Name:  pulumi.String("WEBSITES_ENABLE_APP_SERVICE_STORAGE"),
						Value: pulumi.String("false"),
					},
					&web.NameValuePairArgs{
						Name:  pulumi.String("WEBSITES_PORT"),
						Value: pulumi.String("8080"),
					},
				},

				// Map Azure Storage Account to /data
				AzureStorageAccounts: web.AzureStorageInfoValueMap{
					"data": &web.AzureStorageInfoValueArgs{
						AccessKey:   primaryStorageKey.(pulumi.StringOutput).ToStringPtrOutput(),
						AccountName: account.Name,
						ShareName:   fileshare.Name,
						MountPath:   pulumi.StringPtr("/data"),
						Type:        web.AzureStorageTypeAzureFiles,
					},
				},
			},
		})
		if err != nil {
			return err
		}
		return nil
	})
}
