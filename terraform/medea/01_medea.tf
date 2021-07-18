resource "null_resource" "medea" {
  connection {
    host        = "medea.seankhliao.com"
    private_key = file(pathexpand("~/.ssh/id_ed25519"))
    agent       = false
  }

  provisioner "remote-exec" {
    inline = [
      "rm /etc/sysctl.d/* || true",
      "rm /etc/ssh/ssh_host_{dsa,rsa,ecdsa}_key* || true",
      "pacman -Rns --noconfirm btrfs-progs gptfdisk haveged xfsprogs wget vim net-tools cronie",
      "pacman -Syu --noconfirm neovim",
      "systemctl enable --now systemd-timesyncd",
      "mkdir -p /etc/rancher/k3s",
    ]
  }
  provisioner "file" {
    destination = "/root/.ssh/authorized_keys"
    content     = <<-EOT
      ${file(pathexpand("~/.ssh/id_ed25519.pub"))}
      ${file(pathexpand("~/.ssh/id_ed25519_sk.pub"))}
      ${file(pathexpand("~/.ssh/id_ecdsa_sk.pub"))}
    EOT
  }
  provisioner "file" {
    destination = "/etc/sysctl.d/30-ipforward.conf"
    content     = <<-EOT
      net.ipv4.ip_forward=1
      net.ipv4.conf.lxc*.rp_filter=0
      net.ipv6.conf.default.forwarding=1
      net.ipv6.conf.all.forwarding=1
    EOT
  }
  provisioner "file" {
    destination = "/etc/modules-load.d/br_netfilter.conf"
    content     = "br_netfilter"
  }
  provisioner "file" {
    destination = "/etc/systemd/network/40-wg0.netdev"
    content     = file("secrets/40-wg0.netdev")
  }
  provisioner "file" {
    destination = "/etc/systemd/network/41-wg0.network"
    content     = file("41-wg0.network")
  }
  provisioner "file" {
    destination = "/etc/rancher/k3s/config.yaml"
    content     = file("k3s/config.yaml")
  }
  provisioner "file" {
    destination = "/etc/rancher/k3s/registries.yaml"
    content     = file("secrets/registries.yaml")
  }
  # provisioner "remote-exec" {
  #   inline = [
  #     "curl -sfL https://get.k3s.io | sh -"
  #   ]
  # }
}
