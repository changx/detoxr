package grifts

import (
	"github.com/changx/detoxr/actions"

	"github.com/gobuffalo/buffalo"
)

func init() {
	buffalo.Grifts(actions.App())
}
