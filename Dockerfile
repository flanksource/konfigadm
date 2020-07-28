FROM ubuntu:bionic
ARG SYSTOOLS_VERSION=3.6

RUN apt-get update && \
  apt-get install -y  checkinstall qemu genisoimage  gnupg-agent curl apt-transport-https wget jq git sudo python-setuptools python-pip python-dev build-essential xz-utils ca-certificates unzip zip software-properties-common && \
  rm -Rf /var/lib/apt/lists/*  && \
  rm -Rf /usr/share/doc && rm -Rf /usr/share/man  && \
  apt-get clean

RUN wget --no-check-certificate https://github.com/moshloop/systools/releases/download/${SYSTOOLS_VERSION}/systools.deb && dpkg -i systools.deb
RUN curl -fksSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -
RUN sudo add-apt-repository \
  "deb [arch=amd64] https://download.docker.com/linux/ubuntu \
  $(lsb_release -cs) \
  stable"
RUN apt-get update && \
  apt-get install -y docker-ce docker-ce-cli containerd.io  && \
  rm -Rf /var/lib/apt/lists/*  && \
  rm -Rf /usr/share/doc && rm -Rf /usr/share/man  && \
  apt-get clean

RUN alias curl="curl -k" && curl https://sdk.cloud.google.com | bash -s --  --disable-prompts --install-dir=/opt
RUN wget --no-check-certificate https://dl.google.com/go/go1.13.1.linux-amd64.tar.gz && \
  tar -C /usr/local -xzf go1.13.1.linux-amd64.tar.gz && \
  rm go1.13.1.linux-amd64.tar.gz
RUN  pip install awscli azure-cli sshtunnel==0.1.3

ARG OVFTOOL_LOCATION
RUN wget $OVFTOOL_LOCATION && \
  OVFTOOL_BIN=$(basename $OVFTOOL_LOCATION) && \
  chmod +x $OVFTOOL_BIN && \
  ./$OVFTOOL_BIN --eulas-agreed --required && \
  rm $OVFTOOL_BIN

ARG PACKER_VERSION=1.2.4
RUN install_bin https://releases.hashicorp.com/packer/${PACKER_VERSION}/packer_${PACKER_VERSION}_linux_amd64.zip

ARG GOVC_VERSION=prerelease-v0.21.0-58-g8d28646
RUN install_bin https://github.com/vmware/govmomi/releases/download/${GOVC_VERSION}/govc_linux_amd64.gz

ARG SOPS_VERSION=3.4.0
RUN install_deb https://github.com/mozilla/sops/releases/download/${SOPS_VERSION}/sops_${SOPS_VERSION}_amd64.deb

RUN install_bin https://github.com/CrunchyData/postgres-operator/releases/download/v4.1.0/expenv
RUN install_bin https://github.com/hongkailiu/gojsontoyaml/releases/download/e8bd32d/gojsontoyaml

ARG KONFIGADM_VERSION=
RUN install_bin https://github.com/flanksource/konfigadm/releases/download/${KONFIGADM_VERSION}/konfigadm
