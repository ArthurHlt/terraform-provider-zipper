package zipper

import "github.com/hashicorp/terraform/helper/schema"

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
	return resourceFileCreate(d, meta)
}
