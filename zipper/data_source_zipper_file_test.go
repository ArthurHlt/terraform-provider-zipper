package zipper

import (
	"fmt"
	"regexp"
	"testing"

	r "github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceZipperFile_Basic(t *testing.T) {
	var fileSize string
	r.Test(t, r.TestCase{
		Providers: testProviders,
		Steps: []r.TestStep{
			r.TestStep{
				Config: testAccDsZipperFileContent,
				Check: r.ComposeTestCheckFunc(
					testAccZipFileExists(zipFilePath, &fileSize),
					r.TestCheckResourceAttrPtr("data.zipper_file.foo", "output_size", &fileSize),

					r.TestMatchResourceAttr(
						"data.zipper_file.foo", "output_sha", regexp.MustCompile(`^[0-9a-f]{40}$`),
					),
				),
			},
			r.TestStep{
				Config: testAccDsZipperFileAutoDetect,
				Check: r.ComposeTestCheckFunc(
					testAccZipFileExists(zipFilePath, &fileSize),
					r.TestCheckResourceAttrPtr("data.zipper_file.foo", "output_size", &fileSize),
				),
			},
		},
	})
}

var testAccDsZipperFileContent = fmt.Sprintf(`
data "zipper_file" "foo" {
  type        = "local"
  source      = "%s"
  output_path = "%s"
}
`, folderToZip, zipFilePath)

var testAccDsZipperFileAutoDetect = fmt.Sprintf(`
data "zipper_file" "foo" {
  source      = "%s"
  output_path = "%s"
}
`, folderToZip, zipFilePath)
