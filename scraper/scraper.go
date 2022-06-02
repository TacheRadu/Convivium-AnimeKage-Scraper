package scraper

import (
	types "ak/types"
	"fmt"
	"strings"

	"github.com/gocolly/colly"
)

func GetRecentAnime(pageNumber int) []types.AnimeEpisode {
	var res []types.AnimeEpisode
	c := colly.NewCollector(
		colly.AllowedDomains("ak476.anime-kage.eu"),
	)
	c.OnHTML("body", func(e *colly.HTMLElement) {
		e.ForEach("div.news-item", func(i int, div *colly.HTMLElement) {
			if strings.HasPrefix(div.ChildAttr("a", "href"), "/ak/anime") {
				nameAndEpisode := div.ChildText(".news-title")
				id := strings.Index(nameAndEpisode, "Ep.")
				if id == -1 {
					res = append(res, types.AnimeEpisode{
						ImageLink: div.ChildAttr("img", "data-src"),
						AnimeName: nameAndEpisode,
						Link:      e.Request.AbsoluteURL(div.ChildAttr("a", "href")),
					})
				} else {
					res = append(res, types.AnimeEpisode{
						ImageLink:     div.ChildAttr("img", "data-src"),
						AnimeName:     nameAndEpisode[:id],
						EpisodeNumber: nameAndEpisode[id:],
						Link:          e.Request.AbsoluteURL(div.ChildAttr("a", "href")),
					})
				}
			}
		})
	})
	c.Visit(fmt.Sprintf("https://ak476.anime-kage.eu/?page=%d", pageNumber))
	return res
}

func GetAnime(url string, pageNumber int) types.Anime {
	var res types.Anime
	return res
}

func GetPlayerData(url string) types.PlayerData {
	var res types.PlayerData
	return res
}
