package phases

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/flanksource/konfigadm/pkg/types"
	get "github.com/hashicorp/go-getter"
	log "github.com/sirupsen/logrus"
)

//getDir tmp dir to save downloaded dir
const getDir = "/tmp/konfigadm/"

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
			if strings.HasSuffix(k, "/") && strings.HasPrefix(k, "/") {
				err = GetDir(v)
				if err != nil {
					log.Printf("Error downloading dir: %s", err)
					continue
				}
				LookUpDIR(k, &files)
				CleanUpDIR()
				continue
			}
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

// GetDir : Function to download dir to tmp folder
func GetDir(path string) error {
	currDir, err := os.Getwd()
	if err != nil {
		log.Printf("Error getting current directory: %s", err)
	}
	client := get.Client{
		Src:  path,
		Dst:  getDir,
		Pwd:  currDir,
		Mode: 3, // This ensures src is an archive or directory
	}
	err = client.Get()
	return err
}

// LookUpDIR : Fuction to read contents of files in dir
func LookUpDIR(repPath string, files *map[string]types.File) {
	err := filepath.Walk(getDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		// Replace getDir path with dst path
		newpath := strings.Replace(path, getDir, repPath, 1)
		(*files)[newpath] = types.File{Content: Lookup(path)}
		return nil
	})
	if err != nil {
		log.Printf("Error reading dir %s", err)
	}
}

// CleanUpDIR : Function to delete tmp dir
func CleanUpDIR() {
	err := os.RemoveAll(getDir)
	if err != nil {
		log.Printf("Unable to cleanup dir %s: %s", getDir, err)
	}
}
