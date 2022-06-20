terraform {
  required_providers {
    grid = {
      source = "threefoldtechdev.com/providers/grid"
    }
  }
}

provider "grid" {
}

resource "grid_network" "net1" {
    nodes = [34]
    ip_range = "10.1.0.0/16"
    name = "network"
    description = "newer network"
    add_wg_access = true
}
resource "grid_deployment" "d1" {
  node = 34
  network_name = grid_network.net1.name
  ip_range = lookup(grid_network.net1.nodes_ip_range, 34, "")
  disks{
    name = "rootfs"
    size = 9
  }
  vms {
    name = "vm1"
    flist = "https://hub.grid.tf/tf-official-apps/base:latest.flist"
    cpu = 2 
    memory = 1024
    entrypoint = "/sbin/zinit init"
    env_vars = {
      SSH_KEY = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQCtXFGwsXo3PR2iR1ex8Ko5d/+qj2oeMp6CBcfNwR8k03kPfEeMQVi3BU5dvx/O+fOSrN1vluEh0T+x9NmN7HM3APb5sWevCKDSAMTvr7oy2g6OyRLUbXdGpYEpy8Ne4F0X/7CrSBuX3oJjCARnCN8YL33g+YVMIhOuQK0ZU2dFufRQdog/0LmEuLE1Zvb52KyIRe/FB7J0u98bmxR/DBCPmJ3+QuAtNyxV972m1M0tF9fnI3u4bvQmSViSuiL6KqqrRu7VUS4v5MYGF9hBmIMzQfDXClTwlk7GamiZwAQi1aTKzyUnkEFMX51tL6mQ2rTu1hdEtuU/DhF9cjFYCoDLBkV8I8rKHR/X2kDmQqS2g6B5eSth/Fn3NQyiZ87h31SDSh4Rr/HmBSjiRXLP3nN/I2OajjclWB+3ECdIFs2J6gnQtLM5mWOLojGbhQPjhOwaUCKQadXBJDNnsENbboPtSGg20Mil4GulXEcD+qSeF6+ZVbZXWIHSVmGms0gTVus= mariobassem@Mario"
    }
    planetary = true
    mounts{
      disk_name = "rootfs"
      mount_point = "/mnt/nonUsedMountToCopyRootFileSystemTo"
    }
  }
  vms {
    name = "vm2"
    flist = "https://hub.grid.tf/tf-official-apps/base:latest.flist"
    cpu = 2 
    memory = 1024
    entrypoint = "/sbin/zinit init"
    env_vars = {
      SSH_KEY = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQCtXFGwsXo3PR2iR1ex8Ko5d/+qj2oeMp6CBcfNwR8k03kPfEeMQVi3BU5dvx/O+fOSrN1vluEh0T+x9NmN7HM3APb5sWevCKDSAMTvr7oy2g6OyRLUbXdGpYEpy8Ne4F0X/7CrSBuX3oJjCARnCN8YL33g+YVMIhOuQK0ZU2dFufRQdog/0LmEuLE1Zvb52KyIRe/FB7J0u98bmxR/DBCPmJ3+QuAtNyxV972m1M0tF9fnI3u4bvQmSViSuiL6KqqrRu7VUS4v5MYGF9hBmIMzQfDXClTwlk7GamiZwAQi1aTKzyUnkEFMX51tL6mQ2rTu1hdEtuU/DhF9cjFYCoDLBkV8I8rKHR/X2kDmQqS2g6B5eSth/Fn3NQyiZ87h31SDSh4Rr/HmBSjiRXLP3nN/I2OajjclWB+3ECdIFs2J6gnQtLM5mWOLojGbhQPjhOwaUCKQadXBJDNnsENbboPtSGg20Mil4GulXEcD+qSeF6+ZVbZXWIHSVmGms0gTVus= mariobassem@Mario"
    }
    planetary = true

  }
  vms {
    name = "vm4"
    flist = "https://hub.grid.tf/tf-official-apps/base:latest.flist"
    
    cpu = 2 
    memory = 1024
    entrypoint = "/sbin/zinit init"
    env_vars = {
      SSH_KEY = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQCtXFGwsXo3PR2iR1ex8Ko5d/+qj2oeMp6CBcfNwR8k03kPfEeMQVi3BU5dvx/O+fOSrN1vluEh0T+x9NmN7HM3APb5sWevCKDSAMTvr7oy2g6OyRLUbXdGpYEpy8Ne4F0X/7CrSBuX3oJjCARnCN8YL33g+YVMIhOuQK0ZU2dFufRQdog/0LmEuLE1Zvb52KyIRe/FB7J0u98bmxR/DBCPmJ3+QuAtNyxV972m1M0tF9fnI3u4bvQmSViSuiL6KqqrRu7VUS4v5MYGF9hBmIMzQfDXClTwlk7GamiZwAQi1aTKzyUnkEFMX51tL6mQ2rTu1hdEtuU/DhF9cjFYCoDLBkV8I8rKHR/X2kDmQqS2g6B5eSth/Fn3NQyiZ87h31SDSh4Rr/HmBSjiRXLP3nN/I2OajjclWB+3ECdIFs2J6gnQtLM5mWOLojGbhQPjhOwaUCKQadXBJDNnsENbboPtSGg20Mil4GulXEcD+qSeF6+ZVbZXWIHSVmGms0gTVus= mariobassem@Mario"
    }
    planetary = true

  }
}