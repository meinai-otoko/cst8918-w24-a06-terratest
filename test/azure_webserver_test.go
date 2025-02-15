package test

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/azure"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

// Your Azure subscription ID
var subscriptionID string = "2893229a-aac8-47a3-a415-71b6625a9ad3"

func TestAzureLinuxVMCreation(t *testing.T) {
	terraformOptions := &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: "../",
		// Override the default terraform variables
		Vars: map[string]interface{}{
			"labelPrefix": "kc000004",
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	// Run Terraform apply
	terraform.InitAndApply(t, terraformOptions)

	// Retrieve Terraform outputs
	vmName := terraform.Output(t, terraformOptions, "vm_name")
	resourceGroupName := terraform.Output(t, terraformOptions, "resource_group_name")
	nicName := terraform.Output(t, terraformOptions, "nic_name") // Ensure you have this output in Terraform

	// ✅ Test 1: Confirm VM Exists
	assert.True(t, azure.VirtualMachineExists(t, vmName, resourceGroupName, subscriptionID), "VM should exist in Azure")

	// ✅ Test 2: Confirm NIC Exists
	assert.True(t, azure.NetworkInterfaceExists(t, nicName, resourceGroupName, subscriptionID), "NIC should exist and be attached to the VM")

	// ✅ Test 3: Confirm the VM is Running Ubuntu
	vm := azure.GetVirtualMachine(t, resourceGroupName, vmName, subscriptionID)
	actualUbuntuVersion := vm.StorageProfile.ImageReference.Offer
	assert.Contains(t, actualUbuntuVersion, "UbuntuServer", "VM should be running Ubuntu")

	// ✅ Test 4: Confirm NIC is Attached to the VM
	nics := vm.NetworkProfile.NetworkInterfaces
	assert.NotEmpty(t, nics, "NIC should be attached to the VM")
}
