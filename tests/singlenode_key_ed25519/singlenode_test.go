//go:build integration
// +build integration

package test

import (
	"os"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/threefoldtech/terraform-provider-grid/tests"
)

func TestSingleNodeDeployment(t *testing.T) {
	/* Test case for deployeng a singlenode.

	   **Test Scenario**

	   - Deploy a singlenode.
	   - Check that the outputs not empty.
	   - Up wireguard.
	   - Check that containers is reachable
	   - Verify the VMs ips
	   - Check that env variables set successfully.
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

	node1Container1IP := terraform.Output(t, terraformOptions, "node1_container1_ip")
	assert.NotEmpty(t, node1Container1IP)

	node1Container2IP := terraform.Output(t, terraformOptions, "node1_container2_ip")
	assert.NotEmpty(t, node1Container2IP)

	// Up wireguard
	wgConfig := terraform.Output(t, terraformOptions, "wg_config")
	assert.NotEmpty(t, wgConfig)

	pIP := strings.Split(publicIP, "/")[0]

	status := false
	status = tests.Wait(pIP, "22")
	if status == false {
		t.Errorf("public ip not reachable")
	}

	verifyIPs := []string{publicIP, node1Container1IP, node1Container1IP}
	tests.VerifyIPs(wgConfig, verifyIPs)
	defer tests.DownWG()

	// Check that env variables set successfully
	res, _ := tests.RemoteRun("root", pIP, "cat /proc/1/environ")
	assert.Contains(t, string(res), "TEST_VAR=this value for test")
}
