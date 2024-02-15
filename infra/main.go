package main

import (
	"github.com/pulumi/pulumi-azure-native-sdk/resources/v2"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		// Create an Azure Resource Group
		_, err := resources.NewResourceGroup(ctx, "OneCloud-MarioKart-rg", nil)
		if err != nil {
			return err
		}

		// Create an Azure resource (Storage Account)
		// account, err := storage.NewStorageAccount(ctx, "ocbracketsa", &storage.StorageAccountArgs{
		// 	ResourceGroupName: resourceGroup.Name,
		// 	Sku: &storage.SkuArgs{
		// 		Name: pulumi.String("Standard_LRS"),
		// 	},
		// 	Kind: pulumi.String("StorageV2"),
		// })
		// if err != nil {
		// 	return err
		// }
		//
		// // Create a blob container to mount
		// _, err = storage.NewBlobContainer(ctx, "container", &storage.BlobContainerArgs{
		// 	AccountName: account.Name,
		// })
		// if err != nil {
		// 	return err
		// }
		// primaryStorageKey := pulumi.All(resourceGroup.Name, account.Name).ApplyT(
		// 	func(args []interface{}) (string, error) {
		// 		resourceGroupName := args[0].(string)
		// 		accountName := args[1].(string)
		// 		accountKeys, err := storage.ListStorageAccountKeys(ctx, &storage.ListStorageAccountKeysArgs{
		// 			ResourceGroupName: resourceGroupName,
		// 			AccountName:       accountName,
		// 		})
		// 		if err != nil {
		// 			return "", err
		// 		}
		//
		// 		return accountKeys.Keys[0].Value, nil
		// 	})

		// Export the primary key of the Storage Account
		ctx.Export("primaryStorageKey", pulumi.StringPtr("test"))
		// Create an App Service Plan
		// _, err = web.NewAppServicePlan(ctx, "asp", &web.AppServicePlanArgs{
		// 	ResourceGroupName: resourceGroup.Name,
		// 	Kind:              pulumi.String("Linux"),
		// 	Reserved:          pulumi.Bool(true),
		// 	Sku: &web.SkuDescriptionArgs{
		// 		Name: pulumi.String("F1"),
		// 		Tier: pulumi.String("Free"),
		// 	},
		// })
		// if err != nil {
		// 	return err
		// }
		// Creates the web app with the mount to Azure Storage Account
		// _, err = web.NewWebApp(ctx, "webapp", &web.WebAppArgs{
		// 	ResourceGroupName: resourceGroup.Name,
		// 	ServerFarmId:      plan.ID(),
		// 	SiteConfig: &web.SiteConfigArgs{
		// 		LinuxFxVersion: pulumi.String("DOCKER|appsvcsample/static-site:latest"),
		// 		AppSettings: web.NameValuePairArray{
		// 			&web.NameValuePairArgs{
		// 				Name:  pulumi.String("WEBSITES_ENABLE_APP_SERVICE_STORAGE"),
		// 				Value: pulumi.String("false"),
		// 			},
		// 		},
		// 		// Map Azure Storage Account to /data
		// 		AzureStorageAccounts: web.AzureStorageInfoValueMap{
		// 			"data": &web.AzureStorageInfoValueArgs{
		// 				AccessKey:   primaryStorageKey.(pulumi.StringOutput).ToStringPtrOutput(),
		// 				AccountName: account.Name,
		// 				ShareName:   container.Name,
		// 				MountPath:   pulumi.StringPtr("/data"),
		// 				Type:        web.AzureStorageTypeAzureBlob,
		// 			},
		// 		},
		// 	},
		// })
		// if err != nil {
		// 	return err
		// }

		return nil
	})
}
