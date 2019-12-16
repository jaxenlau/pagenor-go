package services

import (
	"bytes"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/yaml.v2"

	"github.com/jaxenlau/pagenor-go"
	"github.com/jaxenlau/pagenor-go/log"
	"gopkg.in/AlecAivazis/survey.v1"
)

var _ pagenor.Pagenor = (*Pagenor)(nil)

const (
	dateTimeFormat = "2006-01-02 15:04:05"
	dateFormat     = "2006-01-02"
)

type PagenorOptions struct {
	Path              string   `json:"path" yaml:"path" mapstructure:"path"`
	Layout            string   `json:"layout" yaml:"layout" mapstructure:"layout"`
	Categories        []string `json:"categories" yaml:"categories" mapstructure:"categories"`
	TyporaRootURL     string   `json:"typora-root-url" yaml:"typora-root-url" mapstructure:"typora-root-url"`
	TyporaCopyImageTo string   `json:"typora-copy-image-to" yaml:"typora-copy-image-to" mapstructure:"typora-copy-image-to"`
}

const (
	defaultLayout            = "post"
	defaultTyporaRootURL     = "../"
	defaultTyporaCopyImageTo = "../images"
)

func (opts *PagenorOptions) loadDefault() {
	if opts.Layout == "" {
		opts.Layout = defaultLayout
	}

	if opts.TyporaRootURL == "" {
		opts.TyporaRootURL = defaultTyporaRootURL
	}
	if opts.TyporaCopyImageTo == "" {
		opts.TyporaCopyImageTo = defaultTyporaCopyImageTo
	}
}

func NewPagenor(opts *PagenorOptions) *Pagenor {
	opts.loadDefault()
	return &Pagenor{
		cfg: opts,
		questions: []*survey.Question{
			{
				Name:      "title",
				Prompt:    &survey.Input{Message: "请输入标题:"},
				Validate:  survey.Required,
				Transform: survey.Title,
			},
			{
				Name: "category",
				Prompt: &survey.Select{
					Message: "选择分类:",
					Options: opts.Categories,
					Default: "Other",
				},
			},
			{
				Name:   "tags",
				Prompt: &survey.Input{Message: `输入标签(多个标签以 ", " 号间隔):`},
			},
		},
	}
}

type Pagenor struct {
	cfg       *PagenorOptions
	questions []*survey.Question
}

func (p *Pagenor) Generate() error {
	logger := log.DefaultLogger.WithField("scope", "Pagenor.Generate")

	answers := struct {
		Title    string `survey:"title"`
		Category string `survey:"category"`
		Tags     string `survey:"tags"`
	}{}

	err := survey.Ask(p.questions, &answers)
	if err != nil {
		return err
	}

	nowTime := time.Now()

	frontMatter := pagenor.FrontMatter{
		Layout:            p.cfg.Layout,
		Title:             answers.Title,
		Date:              nowTime.Format(dateTimeFormat),
		Category:          answers.Category,
		Tags:              strings.Split(answers.Tags, ","),
		TyporaRootURL:     p.cfg.TyporaRootURL,
		TyporaCopyImageTo: p.cfg.TyporaCopyImageTo,
	}

	result, err := yaml.Marshal(frontMatter)
	if err != nil {
		return err
	}

	output := bytes.Buffer{}
	output.WriteString("---\n")
	output.Write(result)
	output.WriteString("---\n")

	fileName := strings.Join([]string{nowTime.Format(dateFormat), frontMatter.Title}, "-")

	if err := ioutil.WriteFile(filepath.Join(p.cfg.Path, fileName), output.Bytes(), 0644); err != nil {
		return err
	}

	var errBuffer bytes.Buffer

	cmd := exec.Command("open", "-a", "typora", fileName)
	cmd.Stdout = &errBuffer
	cmd.Stderr = &errBuffer

	if err := cmd.Run(); err != nil {
		logger.WithError(err).WithField("filename", fileName).Error("open file failed")

	}

	return nil
}
