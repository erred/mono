#
# Auth types
#
resource "vault_auth_backend" "token" {
  type = "token"
  description = "token based credentials"
}
resource "vault_auth_backend" "userpass" {
  type = "userpass"
  description = "user/password based credentials"
}
resource "vault_auth_backend" "nomad" {
  type = "nomad"
  description = "auth to nomad"
}
resource "vault_nomad_secret_backend" "nomad" {
    backend                   = vault_auth_backend.nomad.id
    description               = vault_auth_backend.nomad.description
    address                   = "http://192.168.100.1:4646"
    token = "cf905516-a1f4-505c-9f83-2b71083101cc"
}


#
# Groups & policies
#
resource "vault_identity_group" "admin" {
  name = "admin"
  type = "internal"
  policies = [
    vault_policy.admin.name,
  ]
  member_entity_ids = [
    vault_identity_entity.eevee.id,
  ]
}
resource "vault_policy" "admin" {
  name = "admin"
  policy = file("policy.admin.vault.hcl")
}

#
# Users
#
resource "vault_identity_entity" "eevee" {
  name      = "eevee"
}
resource "vault_identity_entity_alias" "eevee_userpass" {
  name            = "eevee"
  mount_accessor  = vault_auth_backend.userpass.accessor
  canonical_id    = vault_identity_entity.eevee.id
}
