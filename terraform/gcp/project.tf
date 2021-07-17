locals {
  # https://cloud.google.com/iam/docs/service-accounts
  cloudbuild_service_account = "${data.google_project.default.number}@cloudbuild.gserviceaccount.com"
  # https://cloud.google.com/iam/docs/service-agents
  googleapis_service_agent = "${data.google_project.default.number}@cloudservices.gserviceaccount.com"

}

resource "google_project_iam_policy" "default" {
  project     = data.google_project.default.project_id
  policy_data = data.google_iam_policy.project.policy_data
}
data "google_iam_policy" "project" {
  audit_config {
    audit_log_configs {
      log_type = "ADMIN_READ"
    }
    audit_log_configs {
      log_type = "DATA_READ"
    }
    audit_log_configs {
      log_type = "DATA_WRITE"
    }
    service = "allServices"
  }

  binding {
    members = [
      "serviceAccount:${google_service_account.hetzner_medea.email}"
    ]
    role = "roles/artifactregistry.reader"
  }
  binding {
    members = [
      "serviceAccount:${google_service_account.cluster30_kaniko.email}",
    ]
    role = "roles/artifactregistry.writer"
  }
  binding {
    members = [
      "serviceAccount:${local.cloudbuild_service_account}",
    ]
    role = "roles/cloudbuild.builds.builder"
  }
  binding {
    members = [
      "serviceAccount:${google_service_account.cluster30_traefik.email}",
    ]
    role = "roles/dns.admin"
  }
  binding {
    members = [
      "serviceAccount:${local.googleapis_service_agent}",
    ]
    role = "roles/editor"
  }
  binding {
    members = [
      "user:seankhliao@gmail.com",
    ]
    role = "roles/owner"
  }
  binding {
    members = [
      "serviceAccount:${google_project_service_identity.svc["cloudbuild"].email}"
    ]
    role = "roles/run.admin"
  }

  dynamic "binding" {
    for_each = toset([
      "artifactregistry",
      "cloudasset",
      "cloudbuild",
      "run",
    ])
    iterator = svc

    content {
      members = [
        "serviceAccount:${google_project_service_identity.svc["${svc.value}"].email}",
      ]
      role = "roles/${svc.value}.serviceAgent"
    }
  }
}
