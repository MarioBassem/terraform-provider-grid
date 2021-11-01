package provider

import (
	"context"
	"crypto/ed25519"
	"log"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/pkg/errors"
	substrate "github.com/threefoldtech/substrate-client"
	client "github.com/threefoldtech/terraform-provider-grid/internal/node"
	"github.com/threefoldtech/zos/pkg/rmb"
)

// providerFactories are used to instantiate a provider during acceptance testing.
// The factory function will be invoked for every Terraform CLI command executed
// to create a provider server to which the CLI can reattach.
var providerFactories = map[string]func() (*schema.Provider, error){
	"scaffolding": func() (*schema.Provider, error) {
		return New("dev")(), nil
	},
}
var testAccProviders map[string]*schema.Provider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = New("test")()
	testAccProvider.ConfigureContextFunc = mockproviderConfigure
	testAccProviders = map[string]*schema.Provider{
		"example": testAccProvider,
	}
}
func mockproviderConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	var err error

	apiClient := apiClient{}
	apiClient.mnemonics = d.Get("mnemonics").(string)
	apiClient.use_rmb_proxy = d.Get("use_rmb_proxy").(bool)

	apiClient.rmb_redis_url = d.Get("rmb_redis_url").(string)

	identity, err := substrate.IdentityFromPhrase(string(apiClient.mnemonics))
	if err != nil {
		return nil, diag.FromErr(errors.Wrap(err, "error getting identity"))
	}
	sk, err := identity.SecureKey()
	apiClient.userSK = sk
	if err != nil {
		return nil, diag.FromErr(errors.Wrap(err, "error getting user secret"))
	}
	apiClient.identity = &identity
	network := d.Get("network").(string)
	if network != "dev" && network != "test" {
		return nil, diag.Errorf("network must be one of dev and test")
	}
	apiClient.substrate_url = SUBSTRATE_URL[network]
	apiClient.graphql_url = GRAPHQL_URL[network]
	apiClient.rmb_proxy_url = RMB_PROXY_URL[network]
	substrate_url := d.Get("substrate_url").(string)
	graphql_url := d.Get("graphql_url").(string)
	rmb_proxy_url := d.Get("rmb_proxy_url").(string)
	if substrate_url != "" {
		log.Printf("substrate url is not null %s", substrate_url)
		apiClient.substrate_url = substrate_url
	}
	if graphql_url != "" {
		apiClient.graphql_url = graphql_url
	}
	if rmb_proxy_url != "" {
		apiClient.rmb_proxy_url = rmb_proxy_url
	}
	log.Printf("substrate url: %s %s\n", apiClient.substrate_url, substrate_url)
	apiClient.sub, err = NewSubstrateMock(*apiClient.identity)
	if err != nil {
		return nil, diag.FromErr(errors.Wrap(err, "couldn't create substrate client"))
	}

	pub := sk.Public().(ed25519.PublicKey)
	twin, err := apiClient.sub.GetTwinByPubKey(pub)
	if err != nil && errors.Is(err, substrate.ErrNotFound) {
		return nil, diag.Errorf("no twin associated with the accound with the given mnemonics")
	}
	if err != nil {
		return nil, diag.FromErr(errors.Wrap(err, "failed to get twin for the given mnemonics"))
	}
	apiClient.twin_id = twin
	var cl rmb.Client
	if apiClient.use_rmb_proxy {
		cl = client.NewProxyBus(apiClient.rmb_proxy_url, apiClient.twin_id)
	} else {
		cl, err = rmb.NewClient(apiClient.rmb_redis_url)
	}
	if err != nil {
		return nil, diag.FromErr(errors.Wrap(err, "couldn't create rmb client"))
	}
	apiClient.rmb = cl
	return &apiClient, nil

}
func TestProvider(t *testing.T) {
	if err := New("dev")().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func testAccPreCheck(t *testing.T) {
	// You can add code here to run prior to any test case execution, for example assertions
	// about the appropriate environment variables being set are common to see in a pre-check
	// function.
}
