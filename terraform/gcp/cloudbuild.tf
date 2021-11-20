resource "google_cloudbuild_trigger" "vanity" {
  name        = "vanity"
  description = "go.seankhliao.com ci/cd"

  filename = "ci/cloudbuild/vanity.yaml"
  included_files = [
    "ci/cloudbuild/vanity.yaml",
    "go.*",
    "svc/cmd/vanity/**",
    "svc/webserver/**",
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
    "ci/cloudbuild/w16.yaml",
    "internal/web/**",
    "go.*",
    "static/**",
    "svc/cmd/w16/**",
    "svc/webserver/*",
  ]

  github {
    owner = "seankhliao"
    name  = "mono"
    push {
      branch = "^main$"
    }
  }
}
