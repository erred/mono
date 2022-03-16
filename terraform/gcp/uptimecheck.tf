variable "https_checks" {
  description = "services with standard https uptime checks"
  type = map(object({
    host     = string
    path     = string
    contains = string
  }))
  default = {
    "seankhliao.com" = {
      host     = "seankhliao.com"
      path     = "/"
      contains = "seankhliao"
    }
    "go.seankhliao.com" = {
      host     = "go.seankhliao.com"
      path     = "/uptime-check"
      contains = "go-import"
    }
  }
}


resource "google_monitoring_uptime_check_config" "https_checks" {
  for_each = var.https_checks

  display_name = each.key
  period       = "600s"
  timeout      = "10s"

  content_matchers {
    content = each.value.contains
    matcher = "CONTAINS_STRING"
  }

  http_check {
    path           = each.value.path
    port           = 443
    request_method = "GET"
    use_ssl        = true
    validate_ssl   = true
  }

  monitored_resource {
    labels = {
      "host"       = each.value.host
      "project_id" = data.google_project.default.project_id
    }
    type = "uptime_url"
  }
}

resource "google_monitoring_alert_policy" "https_check_missing" {
  for_each = var.https_checks

  combiner     = "OR"
  display_name = "missing: ${each.key}"
  notification_channels = [
    local.notification_default_app,
  ]
  conditions {
    display_name = "Missing check"
    condition_monitoring_query_language {
      duration = "600s"
      query    = <<-EOT
        fetch uptime_url
        | metric 'monitoring.googleapis.com/uptime_check/check_passed'
        | filter (metric.check_id == '${google_monitoring_uptime_check_config.https_checks[each.key].uptime_check_id}')
        | align next_older(1m)
        | every 1m
        | group_by [resource.host],
            [check_present: count(value.check_passed)]
        | absent_for 600s
      EOT
      trigger {
        count = 1
      }
    }
  }
}

resource "google_monitoring_alert_policy" "https_check_failed" {
  for_each = var.https_checks

  combiner     = "OR"
  display_name = "failed: ${each.key}"
  notification_channels = [
    local.notification_default_app,
  ]
  conditions {
    display_name = "Check failed"
    condition_monitoring_query_language {
      duration = "600s"
      query    = <<-EOT
        fetch uptime_url
        | metric 'monitoring.googleapis.com/uptime_check/check_passed'
        | filter (metric.check_id == '${google_monitoring_uptime_check_config.https_checks[each.key].uptime_check_id}')
        | align next_older(5m)
        | every 5m
        | group_by [resource.host],
            [check_failed: count_true(not(value.check_passed))]
        | condition check_failed > 1 '1'
      EOT
      trigger {
        count = 2
      }
    }
  }
}

resource "google_monitoring_alert_policy" "https_check_latency" {
  for_each = var.https_checks

  combiner     = "OR"
  display_name = "latency: ${each.key}"
  notification_channels = [
    local.notification_default_app,
  ]
  conditions {
    display_name = "Latency"
    condition_monitoring_query_language {
      duration = "600s"
      query    = <<-EOT
        fetch uptime_url
        | metric 'monitoring.googleapis.com/uptime_check/request_latency'
        | filter (metric.check_id == '${google_monitoring_uptime_check_config.https_checks[each.key].uptime_check_id}')
        | group_by 5m, [max_latency: max(value.request_latency)]
        | every 5m
        | condition max_latency > 1500 'ms'
      EOT
      trigger {
        count = 2
      }
    }
  }
}

resource "google_monitoring_alert_policy" "https_check_cert" {
  for_each = var.https_checks

  combiner     = "OR"
  display_name = "cert: ${each.key}"
  notification_channels = [
    local.notification_default_app,
  ]
  conditions {
    display_name = "Cert expiry"
    condition_monitoring_query_language {
      duration = "600s"
      query    = <<-EOT
        fetch uptime_url
        | metric 'monitoring.googleapis.com/uptime_check/time_until_ssl_cert_expires'
        | filter (metric.check_id == '${google_monitoring_uptime_check_config.https_checks[each.key].uptime_check_id}')
        | align next_older(5m)
        | every 5m
        | group_by [resource.host],
            [min_cert_expiry_time:
              min(value.time_until_ssl_cert_expires)]
        | condition min_cert_expiry_time < 15 'd'
      EOT
      trigger {
        count = 2
      }
    }
  }
}
