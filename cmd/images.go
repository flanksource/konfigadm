package cmd

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/moshloop/konfigadm/pkg/build/ova"
	"github.com/moshloop/konfigadm/pkg/types"

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

var driverName, outputDir, outputFilename, resize, image, outputFormat, captureLogs string
var inline bool
var cfg *types.Config
var driver Driver

func copyImage(image string) string {
	cachedImage := image
	if outputFilename != "" {
		image = utils.GetBaseName(outputFilename) + path.Ext(image)
	} else {
		image = utils.GetBaseName(image) + "-" + utils.ShortTimestamp() + path.Ext(image)
	}

	if outputDir != "" {
		image = path.Join(outputDir, image)
	}

	log.Infof("Creating new base image: %s", image)
	if err := utils.FileCopy(cachedImage, image); err != nil {
		log.Fatalf("Failed to create new base image %s, %s", image, err)
	}
	log.Infof("Created new base image")
	if resize != "" {
		log.Infof("Resizing %s to %s\n", image, resize)
		if err := utils.Exec("qemu-img resize \"%s\" %s", image, resize); err != nil {
			log.Fatalf("Error resizing disk  %s", err)
		}
	}
	return image
}

func downloadImage(image string) string {
	if !strings.HasPrefix(image, "http") {
		return image
	}
	home, _ := os.UserHomeDir()
	imageCache := home + "/.konfigadm/images"
	basename := path.Base(image)
	cachedImage := imageCache + "/" + basename
	if utils.FileExists(cachedImage) {
		// TODO(moshloop) verify SHASUM
		log.Infof("Image found in cache: %s", basename)
	} else {
		log.Infof("Downloading image %s", image)
		if err := os.MkdirAll(imageCache, 0755); err != nil {
			log.Fatalf("Failed to create cache dir %s", imageCache)
		}
		if err := utils.Exec(fmt.Sprintf("wget --no-check-certificate -O %s %s", cachedImage, image)); err != nil {
			log.Fatalf("Failed to download image %s, %s", image, err)
		}
	}
	return cachedImage
}

func cloneImage(image string) string {
	if image == "" {
		log.Fatalf("Must specify --image")
	}

	if strings.HasPrefix(image, "http") {
		image = downloadImage(image)
	}
	if !inline {
		image = copyImage(image)
	}
	if !utils.FileExists(image) {
		log.Fatalf("%s does not exists", image)
	}
	return image
}

func buildImage(image string) string {
	image = cloneImage(image)

	if !strings.HasPrefix(image, "/") {
		cwd, _ := os.Getwd()
		image = path.Join(cwd, image)
	}
	driver.Build(image, cfg)
	return image
}

func testImage(image string) string {
	if image == "" {
		log.Fatalf("Must specify --image")
	}

	if !inline {
		image = cloneImage(image)
	}

	driver.Test(image, cfg)
	return image
}

var (
	//Images command
	Images = cobra.Command{
		Use: "images",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			driverName, _ = cmd.Flags().GetString("driver")
			outputDir, _ = cmd.Flags().GetString("output-dir")
			outputFilename, _ = cmd.Flags().GetString("output-filename")
			outputFormat, _ = cmd.Flags().GetString("output-format")
			resize, _ = cmd.Flags().GetString("resize")
			image, _ = cmd.Flags().GetString("image")
			inline, _ = cmd.Flags().GetBool("inline")
			captureLogs, _ = cmd.Flags().GetString("capture-logs")
			if url, ok := images[image]; ok {
				log.Infof("%s is an alias for %s", image, url)
				image = url
			}
			if driverName != "" {
				var ok bool
				if driver, ok = drivers[driverName]; !ok {
					log.Fatalf("Invalid driver name: %s ", driverName)
				}
			}
			cfg = GetConfig(cmd, args)

			cfg.Context.CaptureLogs = captureLogs

		},
	}
	upload = cobra.Command{
		Use:   "upload",
		Short: "Upload an image into a cloud provider",
	}

	list = cobra.Command{
		Use:   "list",
		Short: "List all available images",
		Run: func(cmd *cobra.Command, args []string) {
			for k, v := range images {
				fmt.Printf("%s: %s\n", k, v)
			}
		},
	}

	build = cobra.Command{
		Use:   "build",
		Short: "Build an image ",
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {

			image = buildImage(image)
			var err error
			if outputFormat != "" {
				if outputFormat == "ova" || outputFormat == "ovf" {
					image, err = ova.Create(image, image, make(map[string]string))
					if err != nil {
						log.Fatalf("Failed to create OVA %s", err)
					}
				} else {
					log.Fatalf("Unsupported format %s", outputFormat)
				}
			}
			fmt.Println(image)
		},
	}
)

func init() {
	cwd, _ := os.Getwd()
	upload.AddCommand(&ova.Ova)

	Images.AddCommand(&list, &build, &upload)

	Images.PersistentFlags().String("image", "", "A local or remote path to a disk image")
	Images.PersistentFlags().Bool("inline", false, "If true do not make a copy of the image and work on it directly")
	Images.PersistentFlags().String("capture-logs", "", "Attach a scratch drive to copy logs onto for debugging purposes ")
	Images.PersistentFlags().String("driver", "qemu", "The image build driver to use:  Supported options are: qemu,libguestfs")
	Images.PersistentFlags().String("resize", "", "Resize the image")
	Images.PersistentFlags().String("output-filename", "", "Output filename of image")
	Images.PersistentFlags().String("output-dir", cwd, "Output directory for new images")

}
