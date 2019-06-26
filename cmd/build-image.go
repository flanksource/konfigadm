package cmd

import (
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	. "github.com/moshloop/konfigadm/pkg/build"

	"github.com/moshloop/konfigadm/pkg/utils"
	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

var IMAGE_CACHE string

var drivers = map[string]Driver{
	"libguestfs": Libguestfs{},
	"qemu":       Qemu{},
}

func buildImage(cmd *cobra.Command, args []string, image, outputDir string) {
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
		if resize, _ := cmd.Flags().GetString("resize"); resize != "" {
			log.Infof("Reizing %s to %s\n", image, resize)
			if err := utils.Exec("qemu-img resize \"%s\" %s", image, resize); err != nil {
				log.Fatalf("Error resizing disk  %s", err)
			}
		}

	}
	if !utils.FileExists(image) {
		log.Fatalf("%s does not exists", image)
	}

	cfg := GetConfig(cmd, args)

	var driver Driver
	var ok bool
	if driver, ok = drivers[driverName]; !ok {
		log.Fatalf("Invalid driver name: %s ", driverName)
	}

	if !strings.HasPrefix(image, "/") {
		cwd, _ := os.Getwd()
		image = path.Join(cwd, image)
	}
	driver.Build(image, cfg)
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
			buildImage(cmd, args, image, outputDir)
		},
	}
)

func init() {
	cwd, _ := os.Getwd()
	BuildImg.Flags().String("image", "", "A local or remote path to a disk image")
	BuildImg.Flags().Bool("list-images", false, "Print a list of available public images to use")
	BuildImg.Flags().String("driver", "qemu", "The image build driver to use:  Supported options are: qemu,libguestfs")
	BuildImg.Flags().String("resize", "", "Resize the image")
	BuildImg.Flags().String("output-dir", cwd, "Output directory for new images")
}
