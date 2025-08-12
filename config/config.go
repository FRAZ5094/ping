package config

type Host struct {
	Addr string
	Name string
}

func Parse(path string) ([]Host, error) {

	hosts := []Host{
		{
			Name: "Google",
			Addr: "www.google.com",
		},
		{
			Name: "Cloudflare",
			Addr: "1.1.1.1",
		},
	}

	return hosts, nil
}
