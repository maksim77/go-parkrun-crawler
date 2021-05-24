package goparkruncrawler

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const baseURL = "https://www.parkrun.ru/results/athleteresultshistory"

func prepareURL(parkrunID string) (*url.URL, error) {
	parkrunURL, err := url.Parse(baseURL)
	if err != nil {
		return &url.URL{}, err
	}

	v := url.Values{}
	v.Add("athleteNumber", parkrunID)
	parkrunURL.RawQuery = v.Encode()

	return parkrunURL, nil
}

func getPage(ctx context.Context, parkrunID string) (*goquery.Document, error) {
	parkrunURL, err := prepareURL(parkrunID)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "GET", parkrunURL.String(), nil)
	if err != nil {
		return nil, err
	}

	client := http.Client{}

	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	rootNode, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		return nil, err
	}

	return rootNode, nil
}

// GetRecentRuns returns a list of the last races on the parkrun
func GetRecentRuns(ctx context.Context, parkrunID string) ([]ParkrunRun, error) {
	var parkruns []ParkrunRun

	parkrunMaps := make(map[string]*Parkrun)

	doc, err := getPage(ctx, parkrunID)
	if err != nil {
		return nil, err
	}

	doc.Find("#most-recent").Next().Find("tbody > tr").
		Each(func(i int, s *goquery.Selection) {
			var parkrunRun ParkrunRun
			s.Find("td").Each(func(i int, s *goquery.Selection) {
				switch i {
				case 0:
					if parkrun, ok := parkrunMaps[s.Text()]; ok {
						parkrunRun.Parkrun = parkrun
					} else {
						var pUrl url.URL
						href, ok := s.Children().Attr("href")
						if ok {
							u, _ := url.Parse(href)
							pUrl = *u
						}
						parkrunMaps[s.Text()] = &Parkrun{
							Parkrun:     s.Text(),
							ParkrunLink: pUrl,
						}
						parkrunRun.Parkrun = parkrunMaps[s.Text()]
					}
				case 1:
					date, err := time.Parse("02/01/2006", s.Text())
					if err == nil {
						parkrunRun.Date = date
					}

				case 2:
					i, err := strconv.ParseInt(s.Text(), 0, 16)
					if err == nil {
						parkrunRun.GenderPosition = i
					}
				case 3:
					i, err := strconv.ParseInt(s.Text(), 0, 16)
					if err == nil {
						parkrunRun.OverallPosition = i
					}
				case 4:
					textSplits := strings.Split(s.Text(), ":")
					d, err := time.ParseDuration(fmt.Sprintf("%sm%ss", textSplits[0], textSplits[1]))
					if err == nil {
						parkrunRun.Time = d
					}

				case 5:

					textSplits := strings.Split(s.Text(), "%")
					f, err := strconv.ParseFloat(textSplits[0], 64)
					if err == nil {
						parkrunRun.AgeGrade = f
					}
				}
			})
			parkruns = append(parkruns, parkrunRun)
		})

	return parkruns, nil
}
