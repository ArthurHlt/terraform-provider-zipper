package zipper

import (
	"github.com/ArthurHlt/zipper"
	"github.com/hashicorp/terraform/helper/schema"
	"os"
)

func dataSourceFile() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceZipRead,
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
			"not_when_exists": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Set to true to not create zip when already exists at output path. (to earn time if not necessary)",
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

func dataSourceZipRead(d *schema.ResourceData, meta interface{}) error {
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

	currentSize := d.Get("output_size").(int)
	if currentSize == 0 {
		currentSize = -1
	}
	outputPath := d.Get("output_path").(string)
	fstat, err := os.Stat(outputPath)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	if err == nil {
		d.Set("output_size", int(fstat.Size()))
	}

	onlyWhen := d.Get("not_when_exists").(bool)
	if onlyWhen && err == nil {
		return nil
	}

	size, err := createZip(s, d.Get("output_path").(string))
	if err != nil {
		return err
	}
	d.Set("output_size", int(size))
	return nil
}
