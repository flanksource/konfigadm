package phases

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	"github.com/flanksource/konfigadm/pkg/types"
	log "github.com/sirupsen/logrus"
)

// Files variable defining the FilePhase
var Files types.Phase = filesPhase{}

type filesPhase struct{}

func (p filesPhase) ApplyPhase(sys *types.Config, ctx *types.SystemContext) ([]types.Command, types.Filesystem, error) {
	var commands []types.Command
	files := make(map[string]types.File)
	for k, v := range sys.Filesystem {
		files[k] = v
	}

	for k, v := range sys.Files {
		_, err := url.ParseRequestURI(v)
		if err == nil {
			files[k] = types.File{Content: LookupURL(v)}
			continue

		}
		files[k] = types.File{Content: Lookup(v)}
	}

	return commands, files, nil
}

func (p filesPhase) Verify(cfg *types.Config, results *types.VerifyResults, flags ...types.Flag) bool {
	verify := true
	for f := range cfg.Files {

		if _, err := os.Stat(f); err != nil {
			verify = false
			results.Fail("%s does not exist", f)
		} else {
			results.Pass("%s exists", f)
		}
	}

	for f := range cfg.Templates {
		if _, err := os.Stat(f); err != nil {
			verify = false
			results.Fail("%s does not exist", f)
		} else {
			results.Pass("%s exists", f)
		}
	}

	return verify
}

//Lookup Function to Fetch content of file from path
func Lookup(path string) string {
	_, err := os.Stat(path)
	if err != nil {
		return path
	}
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Printf("Error reading: %s : %s", path, err)
		return path
	}
	return string(data)
}

//LookupURL Function to Fetch content of file from url
func LookupURL(path string) string {
	_, err := url.ParseRequestURI(path)
	if err != nil {
		return path
	}
	resp, err := http.Get(path)
	if err != nil {
		log.Printf("GET Error %s : %s", path, err)
		return path
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Printf("Status error: %v: %s", resp.StatusCode, path)
		return path
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Read body: %s: %s", err, path)
		return path
	}
	return string(data)
}
