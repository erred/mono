resource "google_dns_managed_zone" "medea" {
  description = "medea subdomain for API access"
  dns_name    = "medea.seankhliao.com."

  dnssec_config {
    default_key_specs {
      algorithm  = "rsasha256"
      key_length = 2048
      key_type   = "keySigning"
      kind       = "dns#dnsKeySpec"
    }

    default_key_specs {
      algorithm  = "rsasha256"
      key_length = 1024
      key_type   = "zoneSigning"
      kind       = "dns#dnsKeySpec"
    }

    kind          = "dns#managedZoneDnsSecConfig"
    non_existence = "nsec3"
    state         = "on"
  }

  force_destroy = false

  name       = "medea"
  visibility = "public"
}

resource "google_dns_record_set" "A__medea" {
  managed_zone = google_dns_managed_zone.medea.name
  name         = google_dns_managed_zone.medea.dns_name
  type         = "A"
  rrdatas      = ["65.21.73.144"]
  ttl          = 1
}
resource "google_dns_record_set" "AAAA__medea" {
  managed_zone = google_dns_managed_zone.medea.name
  name         = google_dns_managed_zone.medea.dns_name
  type         = "AAAA"
  rrdatas      = ["2a01:4f9:3b:4e2f::1"]
  ttl          = 1
}
resource "google_dns_record_set" "NS__medea" {
  managed_zone = google_dns_managed_zone.medea.name
  name         = google_dns_managed_zone.medea.dns_name
  type         = "NS"
  rrdatas      = formatlist("ns-cloud-b%d.googledomains.com.", range(1, 5))
  ttl          = 21600
}
resource "google_dns_record_set" "SOA__medea" {
  managed_zone = google_dns_managed_zone.medea.name
  name         = google_dns_managed_zone.medea.dns_name
  type         = "SOA"
  rrdatas      = ["ns-cloud-b1.googledomains.com. cloud-dns-hostmaster.google.com. 1 21600 3600 259200 300"]
  ttl          = 21600
}
