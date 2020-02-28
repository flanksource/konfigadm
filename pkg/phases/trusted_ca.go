package phases

import (
	"fmt"
	"io/ioutil"
	"strings"

	. "github.com/flanksource/konfigadm/pkg/types"
	"github.com/pkg/errors"
)

const (
	certificateHeader = "-----BEGIN CERTIFICATE-----"
)

var TrustedCA Phase = trustedCA{}

type trustedCA struct{}

func (p trustedCA) ApplyPhase(sys *Config, ctx *SystemContext) ([]Command, Filesystem, error) {
	files := Filesystem{}
	commands := make([]Command, len(sys.TrustedCA))

	if len(sys.TrustedCA) == 0 {
		return commands, files, nil
	}

	scriptFilename := "/tmp/trustedCA/install_certs"
	scriptFile := installCertificatesScript()

	files[scriptFilename] = File{
		Content:     scriptFile,
		Permissions: "700",
		Owner:       "root",
	}

	for i, caFile := range sys.TrustedCA {
		tmpFile := fmt.Sprintf("/tmp/trustedCA/konfigadm-trusted-%d.pem", i)

		file, err := certificateToPem(string(caFile))
		if err != nil {
			return nil, files, errors.Wrapf(err, "failed to parse certificate %s", caFile)
		}

		files[tmpFile] = *file
		cmd := fmt.Sprintf("%s %s", scriptFilename, tmpFile)
		commands[i] = Command{Cmd: cmd}
	}

	rmCommand := Command{Cmd: "rm -r /tmp/trustedCA/"}
	commands = append(commands, rmCommand)

	return commands, files, nil
}

func certificateToPem(certificate string) (*File, error) {
	if strings.HasPrefix(certificate, certificateHeader) {
		file := &File{Content: certificate}
		return file, nil
	}

	if strings.HasPrefix(certificate, "http") || strings.HasPrefix(certificate, "https") {
		file := &File{ContentFromURL: certificate}
		return file, nil
	}

	body, err := ioutil.ReadFile(certificate)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read certificate %s from disk", certificate)
	}

	file := &File{Content: string(body)}
	return file, nil
}

func installCertificatesScript() string {
	script := `
#!/bin/bash

install_cert() {
  cert=$1
  desc=$(openssl x509 -in $1 -text -noout | grep Subject | grep CN | sed 's/.*=//' | sed 's/^ //')
  echo "Importing $desc"
  if which update-ca-trust  2>&1 > /dev/null; then
    echo "Updating ca certs via update-ca-trust"
    cp $cert /usr/share/pki/ca-trust-source/$desc.crt
    update-ca-trust extract
  fi

  if which update-ca-certificates 2>&1 > /dev/null; then
    echo "Updating ca certs via update-ca-certificates"
    cp $cert /usr/local/share/ca-certificates/$desc.crt
    update-ca-certificates
  fi

  if [[ -e $JAVA_HOME/jre/lib/security/cacerts ]]; then
       echo "Installing into Java cacerts"
       $JAVA_HOME/bin/keytool -import -noprompt -trustcacerts \
                -keystore  $JAVA_HOME/jre/lib/security/cacerts \
                -storepass changeit -keypass changeit \
                -alias "$desc" \
                -file $1
  fi

  for python in python python2 python3; do
    if which $python 2>&1 > /dev/null ; then
      for site in $($python -c "import site; print('\n'.join(site.getsitepackages()))"); do
          if [[ -e $site ]]; then
             site=$(realpath $site)
             for certs in $(find $site -name "certs.py"); do
                  for pem in $($python $certs); do
                    roots="$roots $pem"
                  done
             done
             roots="$roots $(find $site -name cacerts.txt)"
             roots="$roots $(find $site -name cacert.pem)"
          fi
      done
    fi
  done

  if [[ -e /usr/local/Cellar/ ]]; then
    for site in $(find /usr/local/Cellar/ -type d -name "site-packages"); do
        roots="$roots $(find $site -name cacerts.txt)"
        roots="$roots $(find $site -name cacert.pem)"
    done
  fi

  roots="$roots $(openssl version -d | cut -d":" -f2 | sed 's|"||g' | sed 's| ||')"
  roots=$(echo $roots | tr " " "\n" | sort | uniq)
  for root in $roots; do
    if [[ -d "$root/certs" ]]; then
      echo "Copying to $root/certs/$desc.pem"
      cp $cert "$root/certs/$desc.pem"
    elif [[ -e "$root" ]]; then
      echo "Appending to $root"
      cat $cert >> "$root"
    fi
  done

}


name=$(basename $1)
cert=$1
certname=cert
echo "Installing cert: $1"
roots="/etc/ssl/certs/ca-certificates.crt /etc/pki/tls/certs/ca-bundle.crt /etc/ssl/ca-bundle.pem /etc/pki/tls/cacert.pem /etc/pki/ca-trust/extracted/pem/tls-ca-bundle.pem /usr/local/etc/openssl/cert.pem"

tmp=$(mktemp -d)

if [[ "$1" == "/"* ]]; then
  install_cert $1
elif [[ "$1" == *".pem" ]]; then
  echo "Downloading certificate from $1"
  curl $1 > /tmp/certs
  install_cert $cert
elif [[ "$1" == *":"* ]]; then
  echo "Extracting certs from $1"
  openssl s_client   -showcerts -connect $1 </dev/null 2>&1 | sed -ne '/-BEGIN CERTIFICATE-/,/-END CERTIFICATE-/p' > $tmp/cert.pem
  cd $tmp
  if [[ "$(uname)" == "Darwin" ]]; then
      csplit -k  -f "cert.pem." $tmp/cert.pem  "/END CERTIFICATE/+1" {10}
  else
      csplit -z -k   -f "" -b  $tmp/%02d.pem $tmp/cert.pem "/END CERTIFICATE/+1" {10}
  fi
  rm $tmp/cert.pem
  ls $tmp
  for pem in $(ls $tmp/*.pem*); do
    install_cert $pem
  done
fi
`
	return script
}
