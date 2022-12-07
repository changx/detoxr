package actions

import (
	"github.com/changx/detoxr/ds"
	"github.com/gobuffalo/buffalo"
)

func GetDoHServiceURL(c buffalo.Context) error {
	q := SettingsQuery{
		Server: ds.GetDohService(),
	}
	return successResponse(c, q)
}

func TryDoHServiceWithURL(c buffalo.Context) error {
	return nil
}

func SaveDoHServiceURL(c buffalo.Context) error {
	r := SettingsQuery{}

	if err := c.Bind(&r); err != nil {
		return serverErrorWithShortMessage(c, "illegal request", err.Error())
	}

	if r.Server != "" {
		ds.SetDohService(r.Server)
	}

	return successResponse(c, r)
}
