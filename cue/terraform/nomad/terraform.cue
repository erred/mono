package terraform

terraform: required_providers: nomad: {
    source: "hashicorp/nomad"
    version: "1.4.16"
}

provider: nomad: {
    address: "http://192.168.100.1:4646"
}
