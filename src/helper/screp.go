package helper

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gocolly/colly/v2"
)

type Dye struct {
	Name  string `json:"name"`
	Color string `json:"color"`
}

func DyeView() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		data, err := scrapeDye()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"data": data,
		})
	}
}

func scrapeDye() ([]Dye, error) {
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36"),
	)

	c.SetRequestTimeout(90 * time.Second)

	data := []Dye{}

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8")
		r.Headers.Set("Accept-Language", "en-US,en;q=0.9")
		r.Headers.Set("Referer", "https://tanaka0.work/")
	})

	c.OnHTML("table", func(h *colly.HTMLElement) {
		h.ForEach("tr", func(_ int, el *colly.HTMLElement) {
			td := el.ChildTexts("td")

			if len(td) < 2 {
				return
			}

			name := cleanText(td[0])
			color := cleanText(td[1])

			if name == "" || color == "" {
				return
			}

			data = append(data, Dye{
				Name:  name,
				Color: color,
			})
		})
	})

	err := c.Visit("https://tanaka0.work/AIO/en/DyePredictor/ColorWeapon")
	if err != nil {
		return nil, err
	}

	return data, nil
}

func cleanText(text string) string {
	text = strings.TrimSpace(text)
	text = strings.ReplaceAll(text, "■", "")
	text = strings.Join(strings.Fields(text), " ")

	return text
}
