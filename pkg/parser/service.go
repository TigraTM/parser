package parser

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/sync/errgroup"

	"parser/pkg/news"
)

type Service interface {
	ParsingPage(ctx context.Context, data Parser) error
}

type service struct {
	newsSvc news.Service
}

func NewService(newsSvc news.Service) Service {
	return &service{
		newsSvc: newsSvc,
	}
}

func (s *service) ParsingPage(ctx context.Context, data Parser) error {
	childLinks := make(map[string]struct{})
	ch := make(chan string)

	doc, err := s.getDocument(data.Link)
	if err != nil {
		// TODO: Добавить логирование
		return fmt.Errorf("get document: %w", err)
	}

	for _, attribute := range data.Attributes {
		doc.Find("." + attribute.DivClass).Each(func(i int, s *goquery.Selection) {
			link, ok := s.Find("a." + attribute.AClass).Attr("href")
			if ok {
				if !data.ChildURLIsFull && !strings.HasPrefix(link, "https://") {
					link = data.Link + link
				}
				childLinks[link] = struct{}{}
			}
		})
	}

	go func() {
		for link, _ := range childLinks {
			ch <- link
		}
		defer close(ch)
	}()

	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		for link := range ch {
			//r := time.Duration(rand.Intn(3))
			//fmt.Println(r)
			//time.Sleep(time.Second * r)
			if strings.Contains(link, data.Link) {
				if err := s.sendChildPages(ctx, link, data.ChildAttributes); err != nil {
					return fmt.Errorf("send child pages: %w", err)
				}
			}
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		return fmt.Errorf("g wait: %w", err)
	}

	return nil
}

func (s *service) getDocument(link string) (*goquery.Document, error) {
	res, err := http.Get(link)
	if err != nil {
		//log.Fatal(err)
		return nil, fmt.Errorf("error get home page parsing: %w", err)
	}

	defer res.Body.Close()
	if res.StatusCode != 200 {
		fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		//log.Fatal(err)
		return nil, fmt.Errorf("error document from reader: %w", err)
	}

	return doc, nil
}

func (s *service) sendChildPages(ctx context.Context, childLink string, attributes ChildPagesAttribute) error {
	var (
		title string
		description string
	)

	doc, err := s.getDocument(childLink)
	if err != nil {
		return fmt.Errorf("get document: %w", err)
	}

	doc.Find("." + attributes.ChildDivClass).Each(func(i int, _s *goquery.Selection) {
		t := _s.Find("." + attributes.ClassTitle).Text()
		d := _s.Find("." + attributes.ClassDescription).Text()

		if t != "" && d != "" {
			title, description = t, d
		}
	})

	description = strings.ReplaceAll(description, "\t", "")
	description = strings.ReplaceAll(description, "\n", "")
	description = strings.TrimSpace(description)

	n := news.News{
		Link:        childLink,
		Title:       title,
		Description: description,
	}

	if err := s.newsSvc.CreateNews(ctx, n); err != nil {
		return fmt.Errorf("error create news: %w", err)
	}

	return nil
}
