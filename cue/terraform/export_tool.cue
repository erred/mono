package terraform

import (
	"encoding/json"
	"tool/file"
)

// defaults
locals: {}
module: {}
output: {}
provider: {}
resource: {}
terraform: {}
variable: {}

config: {
    if len(locals) > 0 {
        "locals": locals
    }
    if len(module) > 0 {
        "module": module
    }
    if len(output) > 0 {
        "output": output
    }
    if len(provider) > 0 {
        "provider": provider
    }
    if len(resource) > 0 {
        "resource": resource
    }
    if len(terraform) > 0 {
        "terraform": terraform
    }
    if len(variable) > 0 {
        "variable": variable
    }
}

command: export: task: write: file.Create & {
    filename: "config.tf.json"
    contents: json.Indent(json.Marshal(config), "", "  ")
}
