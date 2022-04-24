job "podinfo" {
  datacenters = ["finland-hetzner"]
  type = "service"
  group "podinfo" {
    network {
      mode = "bridge"
      port "http" {
        # static = 9898
        to = 9898
      }
    }
    task "podinfo" {
      driver = "containerd-driver"
      config {
        image = "index.docker.io/stefanprodan/podinfo"
        # host_network = true
      }
    }
  }
}
