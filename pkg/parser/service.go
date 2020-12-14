package parser

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"golang.org/x/sync/errgroup"

	"parser/pkg/news"
)

type Service interface {
	ParsingPage(ctx context.Context, data Parser) error
}

type service struct {
	newsSvc news.Service
	cfg     *viper.Viper
	log     *logrus.Logger
}

func NewService(newsSvc news.Service, cfg *viper.Viper, log *logrus.Logger) Service {
	return &service{
		newsSvc: newsSvc,
		cfg:     cfg,
		log:     log,
	}
}

func (s *service) ParsingPage(ctx context.Context, data Parser) error {
	childLinks := make(map[string]struct{})
	//ch := make(chan string)

	doc, err := s.getDocument(data.Link)
	if err != nil {
		s.log.Errorf("error getting home page document: %s", err)
		return fmt.Errorf("get document: %w", err)
	}

	for _, attribute := range data.Attributes {
		doc.Find("." + attribute.DivClass).Each(func(i int, s *goquery.Selection) {
			link, ok := s.Find("a." + attribute.AClass).Attr("href")
			if ok {
				if data.URLIsNotFull && !strings.HasPrefix(link, "https://") {
					link = data.Link + link
				}
				childLinks[link] = struct{}{}
			}
		})
	}


	//go func() {
	//	for link, _ := range childLinks {
	//		ch <- link
	//	}
	//	defer close(ch)
	//}()

	//for link := range ch {
	//	//	delay := time.Duration(rand.Intn(s.cfg.GetInt("DELAY")))
	//	//	fmt.Println(delay)
	//	//	time.Sleep(time.Second * delay)
	//	//	if strings.Contains(link, data.Link) {
	//	//		if err := s.sendChildPages(ctx, link, data.ChildAttributes); err != nil {
	//	//			s.log.Errorf("error send child page: %s", err)
	//	//			return fmt.Errorf("send child pages: %w", err)
	//	//		}
	//	//	}
	//	//}

	g, ctx := errgroup.WithContext(ctx)

	for link, _ := range childLinks {
		if strings.Contains(link, data.Link) {
			g.Go(func() error {
				if err := s.sendChildPages(ctx, link, data.ChildAttributes); err != nil {
					s.log.Errorf("error send child page: %s", err)
					return fmt.Errorf("send child pages: %w", err)
				}
				return nil
			})
		}
	}

	if err := g.Wait(); err != nil {
		s.log.Errorf("error errorgroup wait: %s", err)
		return fmt.Errorf("g wait: %w", err)
	}

	return nil
}

func (s *service) getDocument(link string) (*goquery.Document, error) {
	res, err := http.Get(link)
	if err != nil {
		s.log.Errorf("error: %s, get link: %s", err, link)
		return nil, fmt.Errorf("error get home page parsing: %w", err)
	}

	defer res.Body.Close()
	if res.StatusCode != 200 {
		s.log.Infof("get documents code: %d, status: %s", res.StatusCode, res.Status)
		fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		s.log.Errorf("error new document from reader: %s", err)
		return nil, fmt.Errorf("error document from reader: %w", err)
	}

	return doc, nil
}

func (s *service) sendChildPages(ctx context.Context, childLink string, attributes ChildPagesAttribute) error {
	var (
		title       string
		description string
	)

	doc, err := s.getDocument(childLink)
	if err != nil {
		s.log.Errorf("error getting child page document: %s", err)
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
		s.log.Errorf("error create document: %s", err)
		return fmt.Errorf("error create news: %w", err)
	}

	return nil
}
