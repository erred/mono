resource "google_cloudbuild_trigger" "vanity" {
  name        = "vanity"
  description = "go.seankhliao.com ci/cd"

  filename = "vanity/cloudbuild.yaml"
  included_files = [
    "go.*",
    "internal/envconf/**",
    "internal/runhttp/**",
    "internal/vanity/**",
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

  filename = "blog/cloudbuild.yaml"
  included_files = [
    "blog/**",
    "go.*",
    "internal/blog/**",
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
