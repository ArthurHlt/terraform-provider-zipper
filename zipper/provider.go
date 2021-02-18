package zipper

import (
	"context"
	"crypto/tls"
	"github.com/ArthurHlt/zipper"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"net/http"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"skip_ssl_validation": &schema.Schema{
				Type:        schema.TypeBool,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("ZIPPER_SKIP_SSL_VALIDATION", "false"),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"zipper_file": resourceFile(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"zipper_file": dataSourceFile(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {

	m, err := zipper.NewManager(zipper.NewGitHandler(), &zipper.HttpHandler{}, &zipper.LocalHandler{})
	if err != nil {
		return nil, diag.FromErr(err)
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
