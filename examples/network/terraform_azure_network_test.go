package test

import (
    "strings"
    "testing"

    "github.com/allanore/aztest/modules/azure"
    "github.com/gruntwork-io/terratest/modules/random"
    "github.com/gruntwork-io/terratest/modules/terraform"
    "github.com/stretchr/testify/assert"
)

// An example of how to test the Terraform module in examples/terraform-azure-example using Terratest.
func TestTerraformAzureNetworkingExample(t *testing.T) {
    t.Parallel()

    var regions = []string{
        "centralus",
        "eastus",
        "eastus2",
        "northcentralus",
        "southcentralus",
        "westcentralus",
        "westus",
        "westus2",
    }

    // Pick a random Azure region to test in.
    azureRegion := random.RandomString(regions)

    // Network Settings for Vnet and Subnet
    systemName := strings.ToLower(random.UniqueId())
    vnetAddress := "10.0.0.0/16"
    subnetPrefix := "10.0.0.0/24"

    terraformOptions := &terraform.Options{

        // The path to where our Terraform code is located
        TerraformDir: "../examples/network",

        // Variables to pass to our Terraform code using -var options
        Vars: map[string]interface{}{
            "system":             systemName,
            "location":           azureRegion,
            "vnet_address_space": vnetAddress,
            "subnet_prefix":      subnetPrefix,
        },
    }

    // At the end of the test, run `terraform destroy` to clean up any resources that were created
    defer terraform.Destroy(t, terraformOptions)

    // This will run `terraform init` and `terraform apply` and fail the test if there are any errors
    terraform.InitAndApply(t, terraformOptions)

    // Run `terraform output` to get the value of an output variable
    vnetRG := terraform.Output(t, terraformOptions, "vnet_rg")
    subnetID := terraform.Output(t, terraformOptions, "subnet_id")
    nsgName := terraform.Output(t, terraformOptions, "nsg_name")

    // Look up Subnet and NIC ID associations of NSG
    nsgAssociations := azure.GetAssociationsforNSG(t, vnetRG, nsgName, "")

    //Check if subnet is associated with NSG
    assert.Contains(t, nsgAssociations, subnetID)

}