package internal

type (
	SiteConfig struct {
		Url string
	}

	SystemConfig struct {
		Site    SiteConfig
		Request RequestConfig
	}

	RequestConfig struct {
		TimeoutPerRequest int
	}
)
