package routes

import (
	"fmt"

	"github.com/jpmchia/ip2location-pfsense/backend/service/controller"

	"github.com/labstack/echo/v4"
)

type (
	Home struct {
		controller.Controller
	}

	stat struct {
		Title       string
		MetricsJson string
	}
)

func (c *Home) Get(ctx echo.Context) error {
	page := controller.NewPage(ctx)
	page.Layout = "main"
	page.Name = "home"
	page.Metatags.Description = "Backend IP2Location service for pfSense"
	page.Metatags.Keywords = []string{"pfSense", "IP2Location", "geolocation", "firewall", "security"}
	page.Pager = controller.NewPager(ctx, 4)
	page.Data = c.fetchStats(&page.Pager)

	return c.RenderPage(ctx, page)
}

// fetchPosts is an mock example of fetching posts to illustrate how paging works
func (c *Home) fetchStats(pager *controller.Pager) []stat {
	pager.SetItems(20)
	stats := make([]stat, 20)

	for k := range stats {
		stats[k] = stat{
			Title:       fmt.Sprintf("Post example #%d", k+1),
			MetricsJson: fmt.Sprintf("Lorem ipsum example #%d ddolor sit amet, consectetur adipiscing elit. Nam elementum vulputate tristique.", k+1),
		}
	}
	return stats[pager.GetOffset() : pager.GetOffset()+pager.ItemsPerPage]
}
