package zipper

import (
	"crypto/tls"
	"github.com/ArthurHlt/zipper"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"net/http"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"skip_ssl_validation": &schema.Schema{
				Type:        schema.TypeBool,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("ZIPPER_SKIP_SSL_VALIDATION", "true"),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"zipper_file": resourceFile(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"zipper_file": dataSourceFile(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {

	m, err := zipper.NewManager(zipper.NewGitHandler(), &zipper.HttpHandler{}, &zipper.LocalHandler{})
	if err != nil {
		return nil, err
	}
	m.SetHttpClient(&http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: d.Get("skip_ssl_validation").(bool),
			},
		},
	})
	return m, nil
}
