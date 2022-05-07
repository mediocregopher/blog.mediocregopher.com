package cfg

// this file contains functionality specific to the mediocre blog.

// NewBlogCfg returns a Cfg specifically configured for mediocre blog processes.
func NewBlogCfg(params Params) *Cfg {
	params.EnvPrefix = "MEDIOCRE_BLOG"
	return New(params)
}
