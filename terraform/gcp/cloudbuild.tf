resource "google_cloudbuild_trigger" "vanity" {
  name        = "vanity"
  description = "go.seankhliao.com ci/cd"

  filename = "cmd/go.seankhliao.com/cloudbuild.yaml"
  included_files = [
    "go.*",
    "internal/envconf/**",
    "internal/runhttp/**",
    "internal/go.seankhliao.com/**",
    "internal/web/**",
    "vanity/**",
  ]

  github {
    owner = "seankhliao"
    name  = "mono"
    push {
      branch = "^main$"
    }
  }
}

resource "google_cloudbuild_trigger" "blog" {
  name        = "blog"
  description = "seankhliao.com ci/cd"

  filename = "cmd/seankhliao.com/cloudbuild.yaml"
  included_files = [
    "cmd/seankhliao.com/**",
    "go.*",
    "internal/envconf/**",
    "internal/runhttp/**",
    "internal/seankhliao.com/**",
    "internal/web/**",
  ]

  github {
    owner = "seankhliao"
    name  = "mono"
    push {
      branch = "^main$"
    }
  }
}
