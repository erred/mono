resource "google_cloudbuild_trigger" "vanity" {
  name        = "vanity"
  description = "go.seankhliao.com ci/cd"

  filename = "ci/cloudbuild/vanity.yaml"
  included_files = [
    "ci/cloudbuild/vanity.yaml",
    "go.*",
    "content/**",
    "internal/stdlog/**",
    "internal/web/**",
    "svc/cmd/vanity/**",
    "svc/runsvr/**",
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
    "go.*",
    "content/**",
    "internal/o11y/**",
    "internal/web/picture/**",
    "internal/web/render/**",
    "static/**",
    "internal/stdlog/**",
    "svc/runsvr/**",
    "svc/cmd/w16/**",
  ]

  github {
    owner = "seankhliao"
    name  = "mono"
    push {
      branch = "^main$"
    }
  }
}
