resource "google_cloudbuild_trigger" "vanity" {
  name        = "vanity"
  description = "go.seankhliao.com ci/cd"

  filename = "ci/cloudbuild/vanity.yaml"
  included_files = [
    "ci/Dockerfile",
    "go.*",
    "go/cmd/vanity/**",
    "go/webserver/**",
  ]

  github {
    owner = "seankhliao"
    name  = "mono"
    push {
      branch = "^main$"
    }
  }
}

resource "google_cloudbuild_trigger" "w16" {
  name        = "w16"
  description = "seankhliao.com ci/cd"

  filename = "ci/cloudbuild/w16.yaml"
  included_files = [
    "blog/**",
    "ci/Dockerfile",
    "go.*",
    "go/cmd/w16/**",
    "go/cmd/webrender/**",
    "go/internal/w16/**",
    "go/render/**",
    "go/webserver/**",
  ]

  github {
    owner = "seankhliao"
    name  = "mono"
    push {
      branch = "^main$"
    }
  }
}
