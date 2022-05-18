package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	fads "github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	armresources "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	_ "github.com/joho/godotenv/autoload"
)

var (
	azureClientId     string
	ctx               = context.TODO()
	subscriptionId    = os.Getenv("AZURE_SUBSCRIPTION_ID")
	location          = "westus2"
	resourceGroupName = "resourceGroupName"
	interval          = 5 * time.Second
)

func ParseEnvironment() error {
	azureClientId = os.Getenv("AZURE_CLIENT_ID")
	fmt.Println(azureClientId)
	return nil
}

func handleErr(err error) {
	if err != nil {
		log.Panicf(fmt.Sprintf("Error during processing: %s", err.Error()))
	}
}

func createResourceGroup(ctx context.Context, credential azcore.TokenCredential) (*armresources.ResourceGroupsClientCreateOrUpdateResponse, error) {
	rgClient, err := armresources.NewResourceGroupsClient(subscriptionId, credential, nil)
	if err != nil {
		return nil, err
	}

	param := armresources.ResourceGroup{
		Location: to.Ptr(location),
	}

	resp, err := rgClient.CreateOrUpdate(ctx, resourceGroupName, param, nil)

	return &resp, err
}

func main() {
	err := ParseEnvironment()
	handleErr(err)
	cred, err := fads.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("authentication failure: %+v", err)
	}

	resourceGroup, err := createResourceGroup(ctx, cred)
	if err != nil {
		log.Fatalf("cannot create resource group: %+v", err)
	}
	log.Printf("Resource Group %s created", *resourceGroup.ResourceGroup.ID)
}
