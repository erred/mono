FROM alpine

RUN apk add git openssh \
    && adduser -h /srv -s /usr/bin/git-shell -D git \
    && passwd -u git \
    && echo "Port 22223" >> /etc/ssh/sshd_config \
    && echo "PasswordAuthentication no" >> /etc/ssh/sshd_config \
    && echo "HostKey /etc/ssh/data/id_ed25519" >> /etc/ssh/sshd_config \
    && echo "HostCertificate /etc/ssh/data/id_ed25519-cert.pub" >> /etc/ssh/sshd_config \
    && echo "HostKeyAlgorithms ssh-ed25519-cert-v01@openssh.com,ssh-ed25519" >> /etc/ssh/sshd_config \
    && echo "TrustedUserCAKeys /etc/ssh/data/user_ca_keys" >> /etc/ssh/sshd_config

ENTRYPOINT ["/usr/sbin/sshd", "-D", "-e"]
