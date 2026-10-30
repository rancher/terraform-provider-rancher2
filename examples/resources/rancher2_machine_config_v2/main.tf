terraform {
  required_version = "~> v1.14.8"

  required_providers {
    rancher2 = {
      source  = "terraform.local/local/rancher2"
      version = "0.0.0-dev"
    }
  }
}

provider "rancher2" {
  api_url    = "https://rancher.example.com"
  insecure   = true
  access_key = "token-*****"
  secret_key = "*************************************************"
}

resource "rancher2_machine_config_v2" "foo-harvester-v2" {
  generate_name = "foo-harvester-v2"

  harvester_config {
    vm_namespace = "default"

    cpu_count               = "2"
    cpu_pinning             = "true"
    isolate_emulator_thread = "true"

    memory_size = "4"

    disk_info = <<EOF
    {
        "disks": [{
            "imageName": "harvester-public/image-57hzg",
            "size": 40,
            "bootOrder": 1
        }]
    }
    EOF

    network_info = <<EOF
    {
        "interfaces": [{
            "networkName": "harvester-public/vlan1"
        }]
    }
    EOF

    ssh_user = "ubuntu"

    user_data = <<EOF
    package_update: true
    packages:
      - qemu-guest-agent
      - iptables
    runcmd:
      - - systemctl
        - enable
        - '--now'
        - qemu-guest-agent.service
    EOF
  }
}
