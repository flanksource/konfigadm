# Contributing

## Signing the CLA

Please sign the CLA over [here](https://cla-assistant.io/flanksource/konfigadm)
before contributing. Link -

https://cla-assistant.io/flanksource/konfigadm

## Developing

Konfigadm is designed to make changes to the machine it runs on by default. This
means that testing changes locally requires running within some sort of
container or VM. This guide will use a systemd-nspawn container for rapid
feedback. For running tests in the same way as the CI pipeline does, see 
[Testing](#testing).

```
sudo debootstrap buster /var/lib/machines/konfigadm-debian-buster/
sudo systemd-nspawn -D /var/lib/machines/konfigadm-debian-buster/ --machine konfigadm
```

This should present a console prompt on the container. Here we will configure
the root password so that we can boot the container later. You can close the
container by pressing `Ctrl + ]]]`.

```
passwd
```

Back on the host, run the following from within the Kubeadm directory to boot
the container so that Konfigadm can run:

```
make linux
sudo systemd-nspawn -D /var/lib/machines/konfigadm-debian-stable --machine konfigadm --bind-ro $(pwd):/opt/konfigadm
cd /opt/konfigadm
.bin/konfigadm --help
```

Now you can make changes and build locally and then run the resulting binary
within the container to verify progress.

## Testing

Make sure both unit and integration tests pass:

```bash
make
```

You can run unit tests only via:

```bash
make test
```

And only integration tests via:

```bash
make e2e-all
```
