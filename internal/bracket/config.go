package bracket

import (
	"os"
)

type BracketConfig struct {
	DataDir     string
	TemplateDir string
}

func newConfig() BracketConfig {
	c := BracketConfig{
		DataDir:     "/data",
		TemplateDir: "/templates/*",
	}

	if os.Getenv("BRACKET_DEV") == "true" {
		c.DataDir = "./data"
		c.TemplateDir = "./templates"
	}

	return c
}
