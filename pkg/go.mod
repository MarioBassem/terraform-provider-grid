module github.com/threefoldtech/terraform-provider-grid/pkg

go 1.19

require (
	github.com/threefoldtech/substrate-client-dev v0.0.1
	github.com/threefoldtech/substrate-client-main v0.0.1
	github.com/threefoldtech/substrate-client-qa v0.0.1
	github.com/threefoldtech/substrate-client-test v0.0.1
)

require github.com/threefoldtech/zos v0.5.7 // indirect

replace github.com/centrifuge/go-substrate-rpc-client/v4 v4.0.5 => github.com/threefoldtech/go-substrate-rpc-client/v4 v4.0.6-0.20220927094755-0f0d22c73cc7

replace github.com/threefoldtech/substrate-client-dev v0.0.1 => ./substrates/substrate-dev

replace github.com/threefoldtech/substrate-client-test v0.0.1 => ./substrates/substrate-test

replace github.com/threefoldtech/substrate-client-qa v0.0.1 => ./substrates/substrate-qa

replace github.com/threefoldtech/substrate-client-main v0.0.1 => ./substrates/substrate-main
