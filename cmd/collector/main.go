package main

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/go-resty/resty/v2"
	"golang.org/x/net/html"

	"github.com/dcwk/hhparser/internal/models"
)

func main() {
	client := resty.New()
	resp, err := client.R().Get("https://hh.ru/search/vacancy?text=php&salary=&ored_clusters=true&area=1&hhtmFrom=vacancy_search_list&hhtmFromLabel=vacancy_search_line&customDomain=1")
	if err != nil {
		log.Fatal(err)
	}

	page, err := html.Parse(strings.NewReader(resp.String()))
	if err != nil {
		log.Fatal(err)
	}

	cardNode := getNodeByTagAndClass(page, "div", "vacancy-search-item__card")
	card := convertNodeToStruct(cardNode)
	fmt.Println(card)
}

func getNodeByTagAndClass(node *html.Node, tagName string, className string) *html.Node {
	if node.Type == html.ElementNode && node.Data == tagName {
		for _, div := range node.Attr {
			matchedClassName, _ := regexp.MatchString(`(`+className+`$)|(`+className+`\s.+)`, div.Val)
			if div.Key == "class" && matchedClassName {
				return node
			}
		}
	}

	for c := node.FirstChild; c != nil; c = c.NextSibling {
		targetNode := getNodeByTagAndClass(c, tagName, className)
		if targetNode != nil {
			return targetNode
		}
	}

	return nil
}

func convertNodeToStruct(cardNode *html.Node) *models.Card {
	card := &models.Card{}
	nameNode := getNodeByTagAndClass(cardNode, "span", "serp-item__title-link")
	if nameNode != nil {
		card.Name = nameNode.FirstChild.Data
	}

	return card
}
