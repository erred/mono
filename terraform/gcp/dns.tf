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
resource "google_dns_record_set" "MX__medea" {
  managed_zone = google_dns_managed_zone.medea.name
  name         = google_dns_managed_zone.medea.dns_name
  type         = "MX"
  rrdatas      = ["10 ${google_dns_record_set.AAAA__mx1_medea.name}"]
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
resource "google_dns_record_set" "TXT__medea" {
  managed_zone = google_dns_managed_zone.medea.name
  name         = google_dns_managed_zone.medea.dns_name
  type         = "TXT"
  rrdatas = [
    "\"v=spf1 mx -all\"",
  ]
  ttl = 21600
}

resource "google_dns_record_set" "TXT___dmarc_medea" {
  managed_zone = google_dns_managed_zone.medea.name
  name         = "_dmarc.${google_dns_managed_zone.medea.dns_name}"
  type         = "TXT"
  rrdatas = [
    "\"v=DMARC1;p=reject;adkim=r;aspf=r;ruf=mailto:dmarc-forensic@seankhliao.com;rua=mailto:dmarc-aggregate@seankhliao.com;fo=1;ri=86400;pct=100;rf=afrf;\"",
  ]
  ttl = 1
}

locals {
  medea_dkim = "v=DKIM1; k=rsa; p=MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA3R2wlknuJGSUG6St2MBKoD2/BlCt0yA1LpYBRy7rw+uaiePMJHEUM5LZPTsgM6uz0PaRN2u+wOg0ulPdpKhdn5LylX5mEtM+kGBIya2QTsBVDEzgoecOj+sdufVB43sPRXSEdzav+bMv4nvMYtMPbNX1hlk8GEvnMooHB85tDL7LipK26rdc/gIy39kiMqHJavPae3CsMIZiNG6D4oMtePFz9yPlQmm9LVVvCqPTqKvR6Rva3nFTLVBUrO7U4FlKWa+/4VdE89SNDzrZshkSTq6fJ75eA8TRzi0jwwT4silfNXpMnloy4hMC3NHibr9ftAncuRLm1zbWy4LLfFlbmQIDAQAB"
}

resource "google_dns_record_set" "TXT__default__domainkey_medea" {
  managed_zone = google_dns_managed_zone.medea.name
  name         = "default._domainkey.${google_dns_managed_zone.medea.dns_name}"
  type         = "TXT"
  rrdatas = [
    "\"${substr(local.medea_dkim, 0, 253)}\"",
    "\"${substr(local.medea_dkim, 253, -1)}\"",
  ]
  ttl = 1
}

resource "google_dns_record_set" "A__mx1_medea" {
  managed_zone = google_dns_managed_zone.medea.name
  name         = "mx1.${google_dns_managed_zone.medea.dns_name}"
  type         = "A"
  rrdatas      = ["65.21.73.144"]
  ttl          = 1
}
resource "google_dns_record_set" "AAAA__mx1_medea" {
  managed_zone = google_dns_managed_zone.medea.name
  name         = "mx1.${google_dns_managed_zone.medea.dns_name}"
  type         = "AAAA"
  rrdatas      = ["2a01:4f9:3b:4e2f::1"]
  ttl          = 1
}
