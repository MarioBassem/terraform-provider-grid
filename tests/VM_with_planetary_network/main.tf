variable "public_key" {
  type = string
}

terraform {
  required_providers {
    grid = {
      source = "threefoldtech/grid"
    }
  }
}

provider "grid" {
}

resource "grid_network" "net1" {
  nodes       = [2]
  ip_range    = "10.1.0.0/16"
  name        = "network"
  description = "newer network"
}
resource "grid_deployment" "d1" {
  node         = 2
  network_name = grid_network.net1.name
  vms {
    name       = "vm1"
    flist      = "https://hub.grid.tf/tf-official-apps/threefoldtech-ubuntu-20.04.flist"
    cpu        = 2
    publicip   = true
    memory     = 1024
    entrypoint = "/init.sh"
    env_vars = {
      SSH_KEY  = "${var.public_key}"
      TEST_VAR = "this value for test"
    }
    planetary = true
  }
}

output "computed_public_ip" {
  value = grid_deployment.d1.vms[0].computedip
}

output "ygg_ip" {
  value = grid_deployment.d1.vms[0].ygg_ip
}
