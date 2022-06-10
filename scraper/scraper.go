package scraper

import (
	types "ak/types"
	"fmt"
	"strings"

	"github.com/gocolly/colly"
)

func GetRecentAnime(pageNumber int) []types.AnimeEpisode {
	var res []types.AnimeEpisode
	c := colly.NewCollector()
	c.DetectCharset = true
	c.OnHTML("body", func(e *colly.HTMLElement) {
		e.ForEach("div.news-item", func(i int, div *colly.HTMLElement) {
			if strings.HasPrefix(div.ChildAttr("a", "href"), "/ak/anime") {
				nameAndEpisode := div.ChildText(".news-title")
				id := strings.Index(nameAndEpisode, "Ep.")
				if id == -1 {
					if div.ChildAttr("img", "src") == "" {
						res = append(res, types.AnimeEpisode{
							ImageLink: div.ChildAttr("img", "data-src"),
							AnimeName: nameAndEpisode,
							Link:      e.Request.AbsoluteURL(div.ChildAttr("a", "href")),
						})
					} else {
						res = append(res, types.AnimeEpisode{
							ImageLink: div.ChildAttr("img", "src"),
							AnimeName: nameAndEpisode,
							Link:      e.Request.AbsoluteURL(div.ChildAttr("a", "href")),
						})
					}
				} else {
					if div.ChildAttr("img", "src") == "" {
						res = append(res, types.AnimeEpisode{
							ImageLink:     div.ChildAttr("img", "data-src"),
							AnimeName:     nameAndEpisode,
							EpisodeNumber: nameAndEpisode[id:],
							Link:          e.Request.AbsoluteURL(div.ChildAttr("a", "href")),
						})
					} else {
						res = append(res, types.AnimeEpisode{
							ImageLink:     div.ChildAttr("img", "src"),
							AnimeName:     nameAndEpisode,
							EpisodeNumber: nameAndEpisode[id:],
							Link:          e.Request.AbsoluteURL(div.ChildAttr("a", "href")),
						})
					}
				}
			}
		})
	})
	c.Visit(fmt.Sprintf("https://ak476.anime-kage.eu/?page=%d", pageNumber))
	return res
}

func GetAnime(url string, pageNumber int) types.Anime {
	var res types.Anime
	c := colly.NewCollector()
	c.DetectCharset = true
	c.OnHTML(".episode-list-picture", func(e *colly.HTMLElement) {
		res.ImageLink = e.ChildAttr("img", "src")
	})
	c.OnHTML(".page-title", func(e *colly.HTMLElement) {
		res.Title = strings.Trim(e.Text, " \n")
	})
	c.OnHTML("body", func(e *colly.HTMLElement) {
		e.ForEach(".col-12.col-lg-6", func(i int, h *colly.HTMLElement) {
			// first is episode list, this is with details:
			if i == 1 {
				h.ForEach(".row", func(index int, row *colly.HTMLElement) {
					if strings.Trim(row.ChildText(".right-left-desktop"), " \n") == "Genuri:" {
						res.Genres = strings.Split(row.ChildText(".left"), ", ")
					} else if strings.Trim(row.ChildText(".right-left-desktop"), " \n") == "Descriere:" {
						res.Summary = row.ChildText(".left")
					} else if strings.Trim(row.ChildText(".right-left-desktop"), " \n") == "Data lansării:" {
						res.Year = row.ChildText(".left")
					}
				})
			}
		})
	})

	c.OnHTML(".episode-list", func(e *colly.HTMLElement) {
		e.ForEach("a", func(_ int, a *colly.HTMLElement) {
			if strings.HasPrefix(a.Attr("href"), "/ak/anime") {
				res.Episodes = append(res.Episodes, types.AnimeEpisode{
					ImageLink:     res.ImageLink,
					AnimeName:     res.Title,
					EpisodeNumber: strings.Trim(a.Text, " \n"),
					Link:          e.Request.AbsoluteURL(a.Attr("href")),
				})
			}
		})
		if (len(res.Episodes) - (pageNumber+1)*50) > 0 {
			res.HasNextPage = true
		}
		if (len(res.Episodes) - (pageNumber+1)*50) >= 0 {
			res.Episodes = res.Episodes[len(res.Episodes)-(pageNumber+1)*50 : len(res.Episodes)-pageNumber*50]
		} else if len(res.Episodes)-pageNumber*50 > 0 {
			res.Episodes = res.Episodes[0 : len(res.Episodes)-pageNumber*50]
		}
	})
	c.Visit(url)
	return res
}

func GetPlayerData(url string) types.PlayerData {
	var res types.PlayerData
	c := colly.NewCollector()
	c.DetectCharset = true
	c.OnHTML(".col-12.col-md-4.left-center-desktop", func(div *colly.HTMLElement) {
		res.PrevEpisode = div.Request.AbsoluteURL(div.ChildAttr("a", "href"))
	})
	c.OnHTML(".col-12.col-md-4.center", func(div *colly.HTMLElement) {
		if div.ChildAttr("a", "href") != "" {
			res.AnimeLink = div.Request.AbsoluteURL(div.ChildAttr("a", "href"))
		}
	})
	c.OnHTML(".col-12.col-md-4.right-center-desktop", func(div *colly.HTMLElement) {
		if div.ChildAttr("a", "href") != "" {
			res.NextEpisode = div.Request.AbsoluteURL(div.ChildAttr("a", "href"))
		}
	})
	c.OnHTML(".news-title", func(div *colly.HTMLElement) {
		i := strings.Index(div.Text, "Ep.")
		if i != -1 {
			res.EpisodeNumber = strings.Trim(div.Text[i+3:], " \n")
		}
	})
	c.OnHTML("#source15", func(div *colly.HTMLElement) {
		d := colly.NewCollector()
		d.DetectCharset = true
		d.OnHTML("source", func(h *colly.HTMLElement) {
			res.Servers = append(res.Servers, types.Server{
				Title: "Source 15",
				Link:  h.Attr("src"),
			})
		})
		d.Visit(div.ChildAttr("iframe", "data-src"))
	})
	c.Visit(url)
	return res
}
