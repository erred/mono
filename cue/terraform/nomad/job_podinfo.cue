package terraform

import (
    "encoding/json"
)

resource: nomad_job: podinfo: {
    "json": true
    jobspec: json.Marshal({Job: {
        ID: "podinfo"
        Datacenters: ["finland-hetzner"]
        Name: "podinfo"
        TaskGroups: [{
            Name: "podinfo"
            Networks: [{
                DynamicPorts: [{
                    Label: "http"
                    To: 9898
                }]
                Mode: "bridge"
            }]
            Tasks: [{
                Config: image: "index.docker.io/stefanprodan/podinfo"
                Driver: "containerd-driver"
                Name: "podinfo"
            }]
        }]
        Type: "service"
    }})
}
