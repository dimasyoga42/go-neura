package helper

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
)

// Endpoint & action hash otakudesu untuk resolve mirror streaming lewat AJAX.
// CATATAN: action hash ini (aa1208d27f29ca340c92c66d1926f13f dan
// 2a3505c93b0035d3f455df82bf976b84) diambil dari script inline halaman
// episode saat ini. Otakudesu kadang merotasi hash ini saat update situs,
// jadi kalau suatu saat resolveStreamURL berhenti mengembalikan hasil,
// action hash ini perlu diambil ulang dari <script> di halaman episode
// (cari pemanggilan admin-ajax.php di source halaman).
const (
	otakudesuAjaxURL   = "https://otakudesu.blog/wp-admin/admin-ajax.php"
	otakudesuNonceAct  = "aa1208d27f29ca340c92c66d1926f13f"
	otakudesuMirrorAct = "2a3505c93b0035d3f455df82bf976b84"
	otakudesuUA        = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36"
)

type ajaxResponse struct {
	Data string `json:"data"`
}

type mirrorPayload struct {
	ID int    `json:"id"`
	I  int    `json:"i"`
	Q  string `json:"q"`
}

// fetchOtakudesuNonce mengambil nonce yang dibutuhkan sebelum bisa resolve
// mirror streaming manapun.
func fetchOtakudesuNonce(client *http.Client, referer string) (string, error) {
	form := url.Values{"action": {otakudesuNonceAct}}

	req, err := http.NewRequest(http.MethodPost, otakudesuAjaxURL, strings.NewReader(form.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Referer", referer)
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("User-Agent", otakudesuUA)

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var ajaxResp ajaxResponse
	if err := json.Unmarshal(body, &ajaxResp); err != nil {
		return "", err
	}

	return ajaxResp.Data, nil
}

// resolveOtakudesuStreamURL mendecode data-content (base64 -> {id,i,q}),
// lalu memanggil admin-ajax.php untuk mendapatkan HTML iframe embed asli,
// dan mengembalikan src iframe tersebut.
func resolveOtakudesuStreamURL(client *http.Client, referer, dataContent, nonce string) (string, error) {
	decoded, err := base64.StdEncoding.DecodeString(dataContent)
	if err != nil {
		return "", err
	}

	var payload mirrorPayload
	if err := json.Unmarshal(decoded, &payload); err != nil {
		return "", err
	}

	form := url.Values{
		"id":     {strconv.Itoa(payload.ID)},
		"i":      {strconv.Itoa(payload.I)},
		"q":      {payload.Q},
		"nonce":  {nonce},
		"action": {otakudesuMirrorAct},
	}

	req, err := http.NewRequest(http.MethodPost, otakudesuAjaxURL, strings.NewReader(form.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Referer", referer)
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("User-Agent", otakudesuUA)

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var ajaxResp ajaxResponse
	if err := json.Unmarshal(body, &ajaxResp); err != nil {
		return "", err
	}

	htmlBytes, err := base64.StdEncoding.DecodeString(ajaxResp.Data)
	if err != nil {
		return "", err
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(htmlBytes))
	if err != nil {
		return "", err
	}

	src, _ := doc.Find("iframe").Attr("src")
	return strings.TrimSpace(src), nil
}

type Anime struct {
	Name      string `json:"Name"`
	Image_Url string `json:"Image_Url"`
	Link      string `json:"Link"`
}

type Episode struct {
	Name string `json:"Name"`
	Link string `json:"Link"`
	Date string `json:"Date"`
}

type VideoMirror struct {
	Provider string `json:"Provider"`
	Url      string `json:"Url"`
}

type DownloadQuality struct {
	Resolution string        `json:"Resolution"`
	Size       string        `json:"Size"`
	Mirrors    []VideoMirror `json:"Mirrors"`
}

type StreamMirror struct {
	Quality   string `json:"Quality"`
	Provider  string `json:"Provider"`
	Embed_Url string `json:"Embed_Url"`
}

type EpisodeDetail struct {
	Title        string            `json:"Title"`
	Stream_Url   string            `json:"Stream_Url"`
	Streams      []StreamMirror    `json:"Streams"`
	Prev_Episode string            `json:"Prev_Episode"`
	Next_Episode string            `json:"Next_Episode"`
	Anime_Link   string            `json:"Anime_Link"`
	Downloads    []DownloadQuality `json:"Downloads"`
}

type AnimeDetail struct {
	Name         string    `json:"Name"`
	Japanese     string    `json:"Japanese"`
	Score        string    `json:"Score"`
	Producer     string    `json:"Producer"`
	Type         string    `json:"Type"`
	Status       string    `json:"Status"`
	TotalEpisode string    `json:"Total_Episode"`
	Aired        string    `json:"Aired"`
	Duration     string    `json:"Duration"`
	Studio       string    `json:"Studio"`
	Genre        []string  `json:"Genre"`
	Image_Url    string    `json:"Image_Url"`
	Description  string    `json:"Description"`
	Episodes     []Episode `json:"Episodes"`
}

func GetanimeHelper() ([]Anime, error) {
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36"),
	)

	c.SetRequestTimeout(90 * time.Second)

	data := make([]Anime, 0)

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8")
		r.Headers.Set("Accept-Language", "en-US,en;q=0.9")
		r.Headers.Set("Referer", "https://otakudesu.blog/")
	})

	c.OnHTML(".rseries .detpost", func(h *colly.HTMLElement) {
		name := strings.TrimSpace(h.ChildText(".jdlflm"))
		imageURL := h.ChildAttr(".thumb img", "src")
		link := h.ChildAttr(".thumb a", "href")

		if name == "" || link == "" {
			return
		}

		parsedURL, err := url.Parse(link)
		if err != nil {
			return
		}

		data = append(data, Anime{
			Name:      name,
			Image_Url: imageURL,
			Link:      parsedURL.Path,
		})
	})

	if err := c.Visit("https://otakudesu.blog/"); err != nil {
		return nil, err
	}

	return data, nil
}

func GetAnimeDetailHelper(path string) (*AnimeDetail, error) {
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36"),
	)

	c.SetRequestTimeout(90 * time.Second)

	data := &AnimeDetail{
		Genre:    make([]string, 0),
		Episodes: make([]Episode, 0),
	}

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8")
		r.Headers.Set("Accept-Language", "en-US,en;q=0.9")
		r.Headers.Set("Referer", "https://otakudesu.blog/")
	})

	// FIX: id="venkonten" bukan class ".venkonten", dan ".venutama" tidak
	// pernah ada di halaman detail otakudesu. Selector lama tidak pernah
	// match sehingga callback ini tidak pernah jalan.
	c.OnHTML("#venkonten", func(h *colly.HTMLElement) {
		data.Name = strings.TrimSpace(h.ChildText(".jdlrx"))
		data.Image_Url = h.ChildAttr(".fotoanime img", "src")

		h.ForEach(".infozingle p", func(_ int, e *colly.HTMLElement) {
			label := strings.TrimSpace(e.ChildText("b"))
			value := strings.TrimSpace(e.Text)

			value = strings.TrimSpace(strings.TrimPrefix(value, label))
			value = strings.TrimSpace(strings.TrimPrefix(value, ":"))

			// FIX: label di otakudesu pakai Bahasa Indonesia (Tipe,
			// Tanggal Rilis, Durasi, dll), bukan Type/Aired/Duration.
			// Dua-duanya tetap ditangani biar aman kalau situs berubah.
			switch strings.ToLower(label) {
			case "judul":
				data.Name = value
			case "japanese":
				data.Japanese = value
			case "score", "skor":
				data.Score = value
			case "produser", "producer":
				data.Producer = value
			case "type", "tipe":
				data.Type = value
			case "status":
				data.Status = value
			case "total episode":
				data.TotalEpisode = value
			case "aired", "tanggal rilis":
				data.Aired = value
			case "duration", "durasi":
				data.Duration = value
			case "studio":
				data.Studio = value
			}
		})

		h.ForEach(".infozingle span", func(_ int, e *colly.HTMLElement) {
			text := strings.TrimSpace(e.Text)

			if !strings.HasPrefix(strings.ToLower(text), "genre:") {
				return
			}

			e.ForEach("a", func(_ int, genre *colly.HTMLElement) {
				name := strings.TrimSpace(genre.Text)

				if name != "" {
					data.Genre = append(data.Genre, name)
				}
			})
		})

		data.Description = strings.TrimSpace(h.ChildText(".sinopc"))

		h.ForEach(".episodelist ul li", func(_ int, e *colly.HTMLElement) {
			link := strings.TrimSpace(e.ChildAttr("a", "href"))
			name := strings.TrimSpace(e.ChildText("a"))
			date := strings.TrimSpace(e.ChildText(".zeebr"))

			if link == "" || name == "" {
				return
			}

			parsedURL, err := url.Parse(link)
			if err != nil {
				return
			}

			data.Episodes = append(data.Episodes, Episode{
				Name: name,
				Link: parsedURL.Path,
				Date: date,
			})
		})
	})

	if err := c.Visit("https://otakudesu.blog" + path); err != nil {
		return nil, err
	}

	return data, nil
}

// GetEpisodeHelper mengambil detail sebuah episode (judul, video stream
// default, navigasi prev/next, dan semua link download per resolusi)
// dari halaman episode otakudesu, mengikuti pola yang sama dengan
// GetAnimeDetailHelper.
func GetEpisodeHelper(path string) (*EpisodeDetail, error) {
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36"),
	)

	c.SetRequestTimeout(90 * time.Second)

	data := &EpisodeDetail{
		Streams:   make([]StreamMirror, 0),
		Downloads: make([]DownloadQuality, 0),
	}

	// Referensi mirror streaming (id/i/q terenkode base64 di data-content)
	// yang dikumpulkan dulu saat scraping, lalu di-resolve lewat AJAX
	// setelah halaman selesai di-visit.
	type mirrorRef struct {
		Quality  string
		Provider string
		Content  string
	}
	mirrorRefs := make([]mirrorRef, 0)

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8")
		r.Headers.Set("Accept-Language", "en-US,en;q=0.9")
		r.Headers.Set("Referer", "https://otakudesu.blog/")
	})

	c.OnHTML("#venkonten", func(h *colly.HTMLElement) {
		// Judul episode
		data.Title = strings.TrimSpace(h.ChildText(".venutama .posttl"))

		// Iframe streaming default (mirror yang otomatis dimuat saat halaman dibuka)
		data.Stream_Url = strings.TrimSpace(h.ChildAttr(".venutama .player-embed iframe", "src"))

		// Navigasi episode sebelumnya / selanjutnya + link balik ke halaman anime
		h.ForEach(".venutama .prevnext .flir a", func(_ int, e *colly.HTMLElement) {
			href := strings.TrimSpace(e.Attr("href"))
			title := strings.ToLower(strings.TrimSpace(e.Attr("title")))
			text := strings.ToLower(strings.TrimSpace(e.Text))

			if href == "" {
				return
			}

			parsedURL, err := url.Parse(href)
			if err != nil {
				return
			}

			switch {
			case strings.Contains(title, "sebelumnya") || strings.Contains(text, "previous"):
				data.Prev_Episode = parsedURL.Path
			case strings.Contains(title, "selanjutnya") || strings.Contains(text, "next"):
				data.Next_Episode = parsedURL.Path
			case strings.Contains(text, "all episode") || strings.Contains(text, "semua episode"):
				data.Anime_Link = parsedURL.Path
			}
		})

		// Kumpulkan semua mirror streaming (semua resolusi & provider) dari
		// .mirrorstream. Link aslinya belum ada di sini (masih data-content
		// terenkode), di-resolve lewat AJAX setelah c.Visit selesai.
		h.ForEach(".mirrorstream ul", func(_ int, ul *colly.HTMLElement) {
			quality := strings.TrimSpace(strings.TrimPrefix(ul.Attr("class"), "m"))

			ul.ForEach("li a", func(_ int, a *colly.HTMLElement) {
				content := strings.TrimSpace(a.Attr("data-content"))
				provider := strings.TrimSpace(a.Text)

				if content == "" || provider == "" {
					return
				}

				mirrorRefs = append(mirrorRefs, mirrorRef{
					Quality:  quality,
					Provider: provider,
					Content:  content,
				})
			})
		})

		// Semua link download, dikelompokkan per resolusi (Mp4 360p/480p/720p, dst)
		h.ForEach(".download ul li", func(_ int, e *colly.HTMLElement) {
			resolution := strings.TrimSpace(e.ChildText("strong"))
			size := strings.TrimSpace(e.ChildText("i"))

			if resolution == "" {
				return
			}

			mirrors := make([]VideoMirror, 0)
			e.ForEach("a", func(_ int, a *colly.HTMLElement) {
				provider := strings.TrimSpace(a.Text)
				link := strings.TrimSpace(a.Attr("href"))

				if provider == "" || link == "" {
					return
				}

				mirrors = append(mirrors, VideoMirror{
					Provider: provider,
					Url:      link,
				})
			})

			data.Downloads = append(data.Downloads, DownloadQuality{
				Resolution: resolution,
				Size:       size,
				Mirrors:    mirrors,
			})
		})
	})

	if err := c.Visit("https://otakudesu.blog" + path); err != nil {
		return nil, err
	}

	// Resolve semua mirror streaming (semua resolusi & provider) lewat AJAX.
	// Kalau gagal ambil nonce atau resolve satu mirror, tetap lanjut (skip
	// mirror yang gagal) supaya data lain (title, downloads, dll) tetap
	// dikembalikan.
	if len(mirrorRefs) > 0 {
		client := &http.Client{Timeout: 30 * time.Second}
		referer := "https://otakudesu.blog" + path

		nonce, err := fetchOtakudesuNonce(client, referer)
		if err == nil && nonce != "" {
			for _, ref := range mirrorRefs {
				embedURL, err := resolveOtakudesuStreamURL(client, referer, ref.Content, nonce)
				if err != nil || embedURL == "" {
					continue
				}

				data.Streams = append(data.Streams, StreamMirror{
					Quality:   ref.Quality,
					Provider:  ref.Provider,
					Embed_Url: embedURL,
				})
			}
		}
	}
	return data, nil
}
