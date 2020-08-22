package phases

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"

	. "github.com/flanksource/konfigadm/pkg/types" // nolint: golint
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

const (
	certificateHeader = "-----BEGIN CERTIFICATE-----"
)

var (
	caCertificateFiles = []string{
		"/etc/ssl/certs/ca-certificates.crt",
		"/etc/ssl/certs/%s.pem",
		"/etc/pki/ca-trust/extracted/pem/tls-ca-bundle.pem",
		"/etc/pki/tls/certs/%s.pem",
		"/etc/pki/tls/certs/ca-bundle.crt",
		"/usr/lib/ssl/certs/%s.pem",
		"/usr/lib/python3.6/site-packages/pip/_vendor/requests/cacert.pem",
		"/usr/lib/python3.8/site-packages/pip/_vendor/certifi/cacert.pem",
	}
)

var TrustedCA Phase = trustedCA{}

type trustedCA struct{}

func (p trustedCA) ApplyPhase(sys *Config, ctx *SystemContext) ([]Command, Filesystem, error) {
	files := Filesystem{}
	commands := make([]Command, len(sys.TrustedCA))

	if len(sys.TrustedCA) == 0 {
		return commands, files, nil
	}

	scriptFilename := "/tmp/install_certs"
	scriptFile := installCertificatesScript()

	files[scriptFilename] = File{
		Content:     scriptFile,
		Permissions: "0700",
		Owner:       "root",
	}

	for i, caFile := range sys.TrustedCA {
		tmpFile := fmt.Sprintf("/tmp/konfigadm-trusted-%d.pem", i)

		file, err := certificateToPem(string(caFile))
		if err != nil {
			return nil, files, errors.Wrapf(err, "failed to parse certificate %s", caFile)
		}

		files[tmpFile] = *file
		cmd := fmt.Sprintf("%s %s", scriptFilename, tmpFile)
		commands[i] = Command{Cmd: cmd}
	}

	rmCertsCommand := Command{Cmd: "rm -r /tmp/konfigadm-trusted-*.pem"}
	rmScriptCommand := Command{Cmd: "rm -r /tmp/install_certs"}
	commands = append(commands, rmCertsCommand, rmScriptCommand)

	return commands, files, nil
}

func (p trustedCA) Verify(cfg *Config, results *VerifyResults, flags ...Flag) bool {
	if len(cfg.TrustedCA) == 0 {
		return true
	}

	verify := false

	for _, caFile := range cfg.TrustedCA {
		file, err := certificateToPem(string(caFile))
		if err != nil {
			return false
		}

		var content string

		if file.Content != "" {
			content = file.Content
		} else {
			resp, err := http.Get(file.ContentFromURL)
			if err != nil {
				results.Fail("certificate %s download error: %v", file.ContentFromURL, err)
				continue
			}
			defer resp.Body.Close()
			b, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				results.Fail("certificate %s read body error: %v", file.ContentFromURL, err)
				continue
			}
			content = string(b)
		}

		block, _ := pem.Decode([]byte(content))
		if block == nil {
			results.Fail("could not read certificate")
			continue
		}

		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			results.Fail("could not parse certificate: %v", err)
			continue
		}

		if found := findCertificateInFiles(results, cert.Subject.CommonName, content); found {
			verify = true
		} else {
			results.Fail("certificate %s does not exist in cert file paths", cert.Subject.CommonName)
		}
	}

	return verify
}

func findCertificateInFiles(results *VerifyResults, certName string, certBytes string) bool {
	found := false
	for _, fp := range caCertificateFiles {
		filePath := fmt.Sprintf(fp, certName)
		bytes, err := ioutil.ReadFile(filePath)
		if err != nil {
			log.Printf("Error reading: %s : %s", filePath, err)
			continue
		}
		if strings.Contains(string(bytes), certBytes) {
			found = true
			results.Pass("certificate %s found in path %s", certName, filePath)
		}
	}
	return found
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

	fullPath, err := filepath.Abs(certificate)
	if err != nil {
		return nil, errors.Wrap(err, "failed to expand path")
	}

	body, err := ioutil.ReadFile(fullPath)
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
    cp $cert "/usr/share/pki/ca-trust-source/$desc.crt"
    update-ca-trust extract
  fi

  if which update-ca-certificates 2>&1 > /dev/null; then
    echo "Updating ca certs via update-ca-certificates"
    cp $cert "/usr/local/share/ca-certificates/$desc.crt"
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
