resource "nomad_job" "podinfo" {
  jobspec = jsonencode({
    Job = {
      ID = "podinfo"
      Datacenters = [
        "finland-hetzner",
      ]
      Name = "podinfo"
      TaskGroups = [{
        Name = "podinfo"
        Networks = [{
          DynamicPorts = [{
            Label = "http"
            To = 9898
          }]
          Mode = "bridge"
        }]
      }]
      Tasks = [{
        Config = {
          image = "index.docker.io/stefanprodan/podinfo"
        }
        Driver = "containerd-driver"
      }]
      Type = "service"
    }
  })
  json = true
}

data "nomad_job_parser" "my_job" {
  hcl = <<EOF
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
  EOF
  canonicalize = true
}

output "json" {
  value = data.nomad_job_parser.my_job.json
}
