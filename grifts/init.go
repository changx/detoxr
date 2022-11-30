package grifts

import (
	"git.kilokb.com/jon/detoxr/actions"

	"github.com/gobuffalo/buffalo"
)

func init() {
	buffalo.Grifts(actions.App())
}
