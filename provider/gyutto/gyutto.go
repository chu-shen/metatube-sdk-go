package gyutto

import (
	"fmt"
	"net/url"
	"path"
	"regexp"
	"strings"

	"github.com/gocolly/colly/v2"
	"golang.org/x/text/language"

	"github.com/metatube-community/metatube-sdk-go/common/parser"
	"github.com/metatube-community/metatube-sdk-go/model"
	"github.com/metatube-community/metatube-sdk-go/provider"
	"github.com/metatube-community/metatube-sdk-go/provider/internal/scraper"
)

var _ provider.MovieProvider = (*Gyutto)(nil)

const (
	Name     = "Gyutto"
	Priority = 1000
)

const (
	baseURL  = "https://gyutto.com/"
	movieURL = "https://gyutto.com/i/item%s"
)

type Gyutto struct {
	*scraper.Scraper
}

func New() *Gyutto {
	return &Gyutto{scraper.NewDefaultScraper(Name, baseURL, Priority, language.Japanese)}
}

func (gcu *Gyutto) NormalizeMovieID(id string) string {
	if ss := regexp.MustCompile(`^(?i)(?:GYUTTO[-_])?(\d+)$`).FindStringSubmatch(id); len(ss) == 2 {
		return ss[1]
	}
	return ""
}

func (gcu *Gyutto) GetMovieInfoByID(id string) (info *model.MovieInfo, err error) {
	return gcu.GetMovieInfoByURL(fmt.Sprintf(movieURL, id))
}

func (gcu *Gyutto) ParseMovieIDFromURL(rawURL string) (string, error) {
	homepage, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}
	return strings.TrimLeft(path.Base(homepage.Path), "item"), nil
}

func (gcu *Gyutto) GetMovieInfoByURL(rawURL string) (info *model.MovieInfo, err error) {
	id, err := gcu.ParseMovieIDFromURL(rawURL)
	if err != nil {
		return
	}

	info = &model.MovieInfo{
		ID:            id,
		Number:        fmt.Sprintf("GYUTTO-%s", id),
		Provider:      gcu.Name(),
		Homepage:      rawURL,
		Actors:        []string{},
		PreviewImages: []string{},
		Genres:        []string{},
	}

	c := gcu.ClonedCollector()

	// Title
	c.OnXML(`//div[contains(@class,"parts_Mds01")]//h1`, func(e *colly.XMLElement) {
		info.Title = strings.TrimSpace(e.Text)
	})

	// Cover
	c.OnXML(`//div[contains(@class,"unit_DojinMainPh")]//div[@class="ItemPh"]/a`, func(e *colly.XMLElement) {
		href := strings.TrimSpace(e.Attr("href"))
		if href != "" {
			info.CoverURL = e.Request.AbsoluteURL(href)
		}
	})

	// Preview Images
	c.OnXML(`//div[contains(@class,"unit_SamplePhSmall")]//div[@class="ItemPh"]/a`, func(e *colly.XMLElement) {
		href := strings.TrimSpace(e.Attr("href"))
		if href != "" {
			info.PreviewImages = append(info.PreviewImages,
				e.Request.AbsoluteURL(href))
		}
	})

	// Fields
	c.OnXML(`//div[contains(@class,"unit_DetailBasicInfo")][1]//dl[contains(@class,"BasicInfo")]`, func(e *colly.XMLElement) {
		key := strings.TrimSpace(e.ChildText(`.//dt`))
		switch key {
		case "サークル":
			info.Label = strings.TrimSpace(e.ChildText(`.//dd/a`))
			info.Maker = strings.TrimSpace(e.ChildText(`.//dd/a`))
		case "配信開始日":
			info.ReleaseDate = parser.ParseDate(e.ChildText(`.//dd`))
		case "ジャンル":
			texts := e.ChildTexts(`.//dd/a`)
			for _, t := range texts {
				t = strings.TrimSpace(t)
				if t != "" {
					info.Genres = append(info.Genres, t)
				}
			}
		}
	})

	// 作品概要
	c.OnXML(`//div[contains(@class,"unit_DetailSummary")]//p`, func(e *colly.XMLElement) {
		info.Summary = strings.TrimSpace(e.Text)
	})

	// Title (fallback)
	c.OnXML(`//meta[@property="og:title"]`, func(e *colly.XMLElement) {
		if info.Title != "" {
			return
		}
		info.Title = e.Attr("content")
	})

	// Cover (fallback)
	c.OnXML(`//meta[@property="og:image"]`, func(e *colly.XMLElement) {
		if info.CoverURL != "" {
			return
		}
		info.CoverURL = e.Request.AbsoluteURL(e.Attr("content"))
	})

	// Fallbacks
	c.OnScraped(func(_ *colly.Response) {
		if info.ThumbURL == "" {
			info.ThumbURL = info.CoverURL // same as cover.
		}
	})

	err = c.Visit(info.Homepage)
	return
}

func init() {
	provider.Register(Name, New)
}
