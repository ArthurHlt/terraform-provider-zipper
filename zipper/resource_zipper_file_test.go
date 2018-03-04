package zipper

import (
	"fmt"
	"regexp"
	"testing"

	r "github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"os"
	"path/filepath"
)

func TestAccResourceZipperFile_Basic(t *testing.T) {
	var fileSize string
	r.Test(t, r.TestCase{
		Providers:    testProviders,
		CheckDestroy: testAccCheckZipFileDestroyed(zipFilePath),
		Steps: []r.TestStep{
			r.TestStep{
				Config: testAccZipperFileContent,
				Check: r.ComposeTestCheckFunc(
					testAccZipFileExists(zipFilePath, &fileSize),
					r.TestCheckResourceAttrPtr("zipper_file.foo", "output_size", &fileSize),

					r.TestMatchResourceAttr(
						"zipper_file.foo", "output_sha", regexp.MustCompile(`^[0-9a-f]{40}$`),
					),
				),
			},
		},
	})
}
func TestAccResourceZipperFile_SourceChange(t *testing.T) {
	var fileSize string
	newFile := filepath.Join(folderToZip, "new.txt")
	defer func() {
		os.Remove(newFile)
	}()
	currentSha1 := ""
	r.Test(t, r.TestCase{
		Providers:    testProviders,
		CheckDestroy: testAccCheckZipFileDestroyed(zipFilePath),
		Steps: []r.TestStep{
			r.TestStep{
				Config: testAccZipperFileContent,
				Check: r.ComposeTestCheckFunc(
					testAccZipFileExists(zipFilePath, &fileSize),
					r.TestCheckResourceAttrPtr("zipper_file.foo", "output_size", &fileSize),

					r.TestMatchResourceAttr(
						"zipper_file.foo", "output_sha", regexp.MustCompile(`^[0-9a-f]{40}$`),
					),
					func(s *terraform.State) error {
						currentSha1 = s.RootModule().Resources["zipper_file.foo"].Primary.Attributes["output_sha"]
						return nil
					},
				),
			},
			r.TestStep{
				PreConfig: func() {
					f, err := os.Create(newFile)
					if err != nil {
						panic(err)
					}
					f.Close()
				},
				Config: testAccZipperFileContent,
				Check: r.ComposeTestCheckFunc(
					testAccZipFileExists(zipFilePath, &fileSize),
					r.TestCheckResourceAttrPtr("zipper_file.foo", "output_size", &fileSize),

					r.TestMatchResourceAttr(
						"zipper_file.foo", "output_sha", regexp.MustCompile(`^[0-9a-f]{40}$`),
					),
					func(s *terraform.State) error {
						p := s.RootModule().Resources["zipper_file.foo"].Primary
						newSha1 := p.Attributes["output_sha"]
						if newSha1 == currentSha1 || p.ID == currentSha1 {
							return fmt.Errorf("Sha1 didn't change after updating file")
						}
						return nil
					},
				),
			},
		},
	})
}
func TestAccResourceZipperFile_NotOnNonExists(t *testing.T) {
	var fileSize string

	r.Test(t, r.TestCase{
		Providers:    testProviders,
		CheckDestroy: testAccCheckZipFileDestroyed(zipFilePath),
		Steps: []r.TestStep{
			r.TestStep{
				Config: testAccZipperFileNotOnNonExists,
				Check: r.ComposeTestCheckFunc(
					testAccZipFileExists(zipFilePath, &fileSize),
					r.TestCheckResourceAttrPtr("zipper_file.foo", "output_size", &fileSize),

					r.TestMatchResourceAttr(
						"zipper_file.foo", "output_sha", regexp.MustCompile(`^[0-9a-f]{40}$`),
					),
				),
			},
			r.TestStep{
				PreConfig: func() {
					os.Remove(zipFilePath)
				},
				Config: testAccZipperFileNotOnNonExists,
				Check: r.ComposeTestCheckFunc(
					func(s *terraform.State) error {
						_, err := os.Stat(zipFilePath)
						if err == nil {
							return fmt.Errorf("File %s has been recreated but shouldn't be.", zipFilePath)
						}
						if err != nil && os.IsNotExist(err) {
							return nil
						}
						return err
					},
					r.TestCheckResourceAttrPtr("zipper_file.foo", "output_size", &fileSize),

					r.TestMatchResourceAttr(
						"zipper_file.foo", "output_sha", regexp.MustCompile(`^[0-9a-f]{40}$`),
					),
				),
			},
		},
	})
}
func TestAccResourceZipperFile_AutoDetect(t *testing.T) {
	var fileSize string
	r.Test(t, r.TestCase{
		Providers:    testProviders,
		CheckDestroy: testAccCheckZipFileDestroyed(zipFilePath),
		Steps: []r.TestStep{
			r.TestStep{
				Config: testAccZipperFileAutoDetect,
				Check: r.ComposeTestCheckFunc(
					testAccZipFileExists(zipFilePath, &fileSize),
					r.TestCheckResourceAttrPtr("zipper_file.foo", "output_size", &fileSize),

					r.TestMatchResourceAttr(
						"zipper_file.foo", "output_sha", regexp.MustCompile(`^[0-9a-f]{40}$`),
					),
				),
			},
		},
	})
}

var testAccZipperFileContent = fmt.Sprintf(`
resource "zipper_file" "foo" {
  type        = "local"
  source      = "%s"
  output_path = "%s"
}
`, folderToZip, zipFilePath)

var testAccZipperFileAutoDetect = fmt.Sprintf(`
resource "zipper_file" "foo" {
  source      = "%s"
  output_path = "%s"
}
`, folderToZip, zipFilePath)

var testAccZipperFileNotOnNonExists = fmt.Sprintf(`
resource "zipper_file" "foo" {
  type               = "local"
  not_when_nonexists = true
  source             = "%s"
  output_path        = "%s"
}
`, folderToZip, zipFilePath)

func testAccCheckZipFileDestroyed(filePath string) r.TestCheckFunc {

	return func(s *terraform.State) error {
		_, err := os.Stat(filePath)
		if err == nil {
			return fmt.Errorf("File %s has not been deleted.", filePath)
		}
		if err != nil && os.IsNotExist(err) {
			return nil
		}
		return err
	}
}
