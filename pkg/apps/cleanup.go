package apps

import (
	. "github.com/moshloop/konfigadm/pkg/types"
)

var Cleanup Phase = cleanup{}

type cleanup struct{}

func (c cleanup) ApplyPhase(sys *Config, ctx *SystemContext) ([]Command, Filesystem, error) {

	// tdnf clean all
	// /sbin/ldconfig
	// /usr/sbin/pwconv
	// /usr/sbin/grpconv
	// rm /etc/resolv.conf
	// ln -sf ../run/systemd/resolve/resolv.conf /etc/resolv.conf
	// rm -rf /tmp/*
	// rm -rf /usr/share/man/*
	// rm -rf /usr/share/doc/*
	// find /var/cache -type f -exec rm -rf {} \;
	// 	unset HISTFILE
	// echo -n > /root/.bash_history

	// find /var/log -type f | while read -r f; do echo -ne '' > "$f"; done;
	// echo -ne '' >/var/log/lastlog
	// echo -ne '' >/var/log/wtmp
	// echo -ne '' >/var/log/btmp

	// echo -ne '' > /root/.bashrc
	// echo -ne '' > /root/.bash_profile
	// echo 'shopt -s histappend' >> /root/.bash_profile
	// echo 'export PROMPT_COMMAND="history -a; history -c; history -r; $PROMPT_COMMAND"' >> /root/.bash_profile

	// rm -f /etc/ssh/{ssh_host_dsa_key,ssh_host_dsa_key.pub,ssh_host_ecdsa_key,ssh_host_ecdsa_key.pub,ssh_host_ed25519_key,ssh_host_ed25519_key.pub,ssh_host_rsa_key,ssh_host_rsa_key.pub}

	// dd if=/dev/zero of=/EMPTY bs=1M  2>/dev/null || echo "dd exit code $? is suppressed"
	// rm -f /EMPTY

	// > /etc/machine-id

	return []Command{}, Filesystem{}, nil
}
