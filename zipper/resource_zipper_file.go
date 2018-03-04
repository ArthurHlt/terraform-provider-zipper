package zipper

import (
	"github.com/ArthurHlt/zipper"
	"github.com/hashicorp/terraform/helper/schema"
	"io"
	"os"
)

func resourceFile() *schema.Resource {
	return &schema.Resource{
		Create: resourceFileCreate,
		Read:   resourceFileRead,
		Update: resourceFileUpdate,
		Delete: resourceFileDelete,
		CustomizeDiff: func(diff *schema.ResourceDiff, v interface{}) error {
			if diff.Id() != diff.Get("output_sha") {
				return diff.SetNewComputed("output_sha")
			}
			return nil
		},
		Schema: map[string]*schema.Schema{
			"type": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Source type to use to create zip, e.g.: http, local or git. (if omitted type will be auto-detected)",
			},
			"source": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Target source for zipper",
			},
			"output_path": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The output of the archive file.",
			},
			"not_when_nonexists": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Set to true to not create zip when not exists at output_path if sources files didn't change. (to earn time if not necessary)",
			},
			"output_sha": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				ForceNew:    true,
				Description: "SHA1 checksum made by zipper.",
			},
			"output_size": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Size of the zip file.",
			},
		},
	}
}
func resourceFileCreate(d *schema.ResourceData, meta interface{}) error {
	z := meta.(*zipper.Manager)
	s, err := z.CreateSession(d.Get("source").(string), d.Get("type").(string))
	if err != nil {
		return err
	}
	sha1, err := s.Sha1()
	if err != nil {
		return err
	}
	d.SetId(sha1)
	d.Set("output_sha", sha1)

	size, err := createZip(s, d.Get("output_path").(string))
	if err != nil {
		return err
	}
	d.Set("output_size", int(size))
	return nil
}

func resourceFileRead(d *schema.ResourceData, meta interface{}) error {
	z := meta.(*zipper.Manager)
	s, err := z.CreateSession(d.Get("source").(string), d.Get("type").(string))
	if err != nil {
		return err
	}
	newSha1, err := s.Sha1()
	if err != nil {
		return err
	}
	d.Set("output_sha", newSha1)

	onlyWhen := d.Get("not_when_nonexists").(bool)
	if onlyWhen {
		return nil
	}
	currentSize := d.Get("output_size").(int)
	if currentSize == 0 {
		currentSize = -1
	}
	outputPath := d.Get("output_path").(string)
	fstat, err := os.Stat(outputPath)
	if err == nil && int(fstat.Size()) == currentSize {
		return nil
	}
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	size, err := createZip(s, d.Get("output_path").(string))
	if err != nil {
		return err
	}
	d.Set("output_size", int(size))
	return nil
}

func resourceFileDelete(d *schema.ResourceData, meta interface{}) error {
	outputPath := d.Get("output_path").(string)
	_, err := os.Stat(outputPath)
	if err != nil && os.IsNotExist(err) {
		return nil
	}
	return os.Remove(outputPath)
}

func resourceFileUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func createZip(s *zipper.Session, outputPath string) (int64, error) {
	zipFile, err := s.Zip()
	if err != nil {
		return 0, err
	}
	defer zipFile.Close()

	f, err := os.Create(outputPath)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	_, err = io.Copy(f, zipFile)
	if err != nil {
		return 0, err
	}
	return zipFile.Size(), nil

}
