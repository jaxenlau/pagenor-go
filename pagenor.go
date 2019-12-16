package pagenor

// FrontMatter Yaml front matter template
type FrontMatter struct {
	Layout            string   `json:"layout" yaml:"layout"`
	Title             string   `json:"title" yaml:"title"`
	Date              string   `json:"date" yaml:"date"`
	Category          string   `json:"category" yaml:"category"`
	Tags              []string `json:"tags" yaml:"tags"`
	TyporaRootURL     string   `json:"typora-root-url" yaml:"typora-root-url"`
	TyporaCopyImageTo string   `json:"typora-copy-image-to" yaml:"typora-copy-image-to"`
}

type Pagenor interface {
	Generate() error
}
