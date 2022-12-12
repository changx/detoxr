package actions

import (
	"time"

	"github.com/changx/detoxr/ds"
	"github.com/gobuffalo/buffalo"
	"github.com/miekg/dns"
)

type ListItem struct {
	Name  string `json:"name"`
	QType string `json:"qtype"`
	TTL   int64  `json:"ttl"`
}

type List struct {
	List []*ListItem `json:"list"`
}

func GetVictims(c buffalo.Context) error {
	list := make([]*ListItem, 0)
	blacklist := ds.GetBlacklist()
	now := time.Now()
	blacklist.Each(func(item *ds.CacheItem) bool {
		c.Logger().Debugf("cache item: %+v", item)
		listItem := ListItem{
			Name:  item.Answer().Question[0].Name,
			QType: dns.TypeToString[item.Answer().Question[0].Qtype],
			TTL:   item.Expire() - now.Unix(),
		}
		if item.Expire() < 0 {
			listItem.TTL = -999999
		}
		list = append(list, &listItem)
		return true
	})
	return successResponse(c, List{List: list})
}
