package zipper

import (
	"testing"

	"fmt"
	r "github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"io/ioutil"
	"os"
	"path/filepath"
)

var wd, _ = os.Getwd()
var folderToZip = filepath.Join(wd, "..", "fixtures", "zip")
var zipFile, _ = ioutil.TempFile("", "provider-zipper")
var zipFilePath = zipFile.Name()

func TestMain(m *testing.M) {
	retCode := m.Run()
	zipFile.Close()
	os.Remove(zipFilePath)
	os.Exit(retCode)
}

var testProviders = map[string]*schema.Provider{
	"zipper": Provider(),
}

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func testAccZipFileExists(filename string, fileSize *string) r.TestCheckFunc {
	return func(s *terraform.State) error {
		*fileSize = ""
		fi, err := os.Stat(filename)
		if err != nil {
			return err
		}
		*fileSize = fmt.Sprintf("%d", fi.Size())
		return nil
	}
}
