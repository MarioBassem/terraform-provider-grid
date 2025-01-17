//go:build integration
// +build integration

package test

import (
	"os"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/threefoldtech/terraform-provider-grid/tests"
)

func TestSingleNodeDeployment(t *testing.T) {
	/* Test case for deployeng a VM with planetary network.

	   **Test Scenario**

	   - Deploy a VM in single node.
	   - Check that the outputs not empty.
	   - Verify the VMs ips
	   - Check that ygg ip is reachable.
	   - ssh to VM and check if yggdrasil is active
	   - Destroy the deployment
	*/

	// retryable errors in terraform testing.
	// generate ssh keys for test
	tests.SSHKeys()
	publicKey := os.Getenv("PUBLICKEY")
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "./",
		Vars: map[string]interface{}{
			"public_key": publicKey,
		},
		Parallelism: 1,
	})
	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	// Check that the outputs not empty
	publicIP := terraform.Output(t, terraformOptions, "computed_public_ip")
	assert.NotEmpty(t, publicIP)

	yggIP := terraform.Output(t, terraformOptions, "ygg_ip")
	assert.NotEmpty(t, yggIP)

	status := false
	status = tests.Wait(yggIP, "22")
	if status == false {
		t.Errorf("Yggdrasil IP not reachable")
	}

	verifyIPs := []string{publicIP, yggIP}
	tests.VerifyIPs("", verifyIPs)
	defer tests.DownWG()

	// ssh to VM by ygg_ip
	res, _ := tests.RemoteRun("root", yggIP, "cat /proc/1/environ")
	assert.Contains(t, string(res), "TEST_VAR=this value for test")
}
