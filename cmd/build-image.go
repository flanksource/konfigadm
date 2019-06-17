package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"time"

	. "github.com/moshloop/konfigadm/pkg/build"

	"github.com/moshloop/konfigadm/pkg/utils"
	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v3"
)

var IMAGE_CACHE string

var drivers = map[string]Driver{
	"libguestfs": Libguestfs{},
}

func buildImage(cmd *cobra.Command, image, outputDir string) {
	HOME, _ := os.UserHomeDir()
	IMAGE_CACHE := HOME + "/.konfigadm/images"
	driverName, _ := cmd.Flags().GetString("driver")

	if image == "" {
		log.Fatalf("Must specify --image or --list-images")
	}

	if url, ok := images[image]; ok {
		log.Infof("%s is an alias for %s", image, url)
		image = url
	}

	if strings.HasPrefix(image, "http") {
		basename := path.Base(image)
		cachedImage := IMAGE_CACHE + "/" + basename
		if utils.FileExists(cachedImage) {
			// TODO(moshloop) verify SHASUM
			log.Infof("Image found in cache: %s", basename)
		} else {
			log.Infof("Downloading image %s", image)
			if err := os.MkdirAll(IMAGE_CACHE, 0755); err != nil {
				log.Fatalf("Failed to create cache dir %s", IMAGE_CACHE)
			}
			if err := utils.Exec(fmt.Sprintf("wget -O %s %s", cachedImage, image)); err != nil {
				log.Fatalf("Failed to download image %s, %s", image, err)
			}
		}
		timestamp := time.Now().Format("-20060102150405")
		image = outputDir + "/" + strings.Split(basename, ".")[0] + timestamp + "." + strings.Split(basename, ".")[1]
		log.Infof("Creating new base image: %s", image)
		if err := utils.FileCopy(cachedImage, image); err != nil {
			log.Fatalf("Failed to create new base image %s, %s", image, err)
		}
		log.Infof("Created new base image")

	}
	if !utils.FileExists(image) {
		log.Fatalf("%s does not exists", image)
	}

	tmpfile, err := ioutil.TempFile("", "konfigadm.*.yml")
	if err != nil {
		log.Fatalf("Cannot create tempfile %s", err)
	}
	cfg := GetConfig(cmd)
	_, _, err = cfg.ApplyPhases()
	if err != nil {
		log.Fatalf("Error applying phases %s\n", err)
	}
	data, _ := yaml.Marshal(cfg)

	if _, err := tmpfile.Write(data); err != nil {
		log.Fatalf("Error writing tmp file %s", err)
	}

	var driver Driver
	var ok bool
	if driver, ok = drivers[driverName]; !ok {
		log.Fatalf("Invalid driver name: %s ", driverName)
	}
	driver.Build(image, tmpfile)

	fmt.Println(image)
}

var (
	BuildImg = cobra.Command{
		Use:   "build-image",
		Short: "Build an image ",
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {

			if list, _ := cmd.Flags().GetBool("list-images"); list {
				for k, v := range images {
					fmt.Printf("%s: %s\n", k, v)
				}
				return
			}
			outputDir, _ := cmd.Flags().GetString("output-dir")
			image, _ := cmd.Flags().GetString("image")
			buildImage(cmd, image, outputDir)
		},
	}
)

func init() {
	cwd, _ := os.Getwd()
	BuildImg.Flags().String("image", "", "A local or remote path to a disk image")
	BuildImg.Flags().Bool("list-images", false, "Print a list of available public images to use")
	BuildImg.Flags().String("driver", "libguestfs", "The image build driver to use, currently only libguestfs supported")
	BuildImg.Flags().String("output-dir", cwd, "Output directory for new images")
}
