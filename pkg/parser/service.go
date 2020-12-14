package parser

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"

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

	doc, err := s.getDocument(data.Link)
	if err != nil {
		// TODO: Добавить логирование
		return fmt.Errorf("get document: %w", err)
	}

	for _, attribute := range data.Attributes {
		doc.Find("." + attribute.DivClass).Each(func(i int, s *goquery.Selection) {
			link, ok := s.Find("a." + attribute.AClass).Attr("href")
			if ok {
				childLinks[link] = struct{}{}
			}
		})
	}

	for childLink, _ := range childLinks {
		if err := s.sendChildPages(ctx, childLink, data.ChildAttributes); err != nil {
			return fmt.Errorf("error send child page: %w", err)
		}
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
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		//log.Fatal(err)
		return nil, fmt.Errorf("error document from reader: %w", err)
	}

	return doc, nil
}

func (s *service) sendChildPages(ctx context.Context, childLink string, attributes ChildPagesAttribute) error {
	doc, err := s.getDocument(childLink)
	if err != nil {
		return fmt.Errorf("get document: %w", err)
	}


	doc.Find("." + attributes.ChildDivClass).Each(func(i int, s *goquery.Selection) {
		title := s.Find("." + attributes.ClassTitle).Text()
		description := s.Find("." + attributes.ClassDescription).Text()

		n := news.News{
			Link:        childLink,
			Title:       title,
			Description: description,
		}

		fmt.Println(n)
	})

	//if err := s.newsSvc.CreateNews(ctx, n); err != nil {
	//	return fmt.Errorf("error create news: %w", err)
	//}

	return nil
}
