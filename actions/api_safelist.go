package actions

import (
	"time"

	"github.com/changx/detoxr/ds"
	"github.com/gobuffalo/buffalo"
	"github.com/miekg/dns"
)

func GetSafelist(c buffalo.Context) error {
	list := make([]*ListItem, 0)
	whitelist := ds.GetWhitelist()
	now := time.Now()
	whitelist.Each(func(item *ds.CacheItem) bool {
		c.Logger().Debugf("cache item: %+v", item)
		listItem := ListItem{
			Name:  item.Answer().Question[0].Name,
			QType: dns.TypeToString[item.Answer().Question[0].Qtype],
			TTL:   item.Expire() - now.Unix(),
		}
		list = append(list, &listItem)
		return true
	})
	return successResponse(c, List{List: list})
}
