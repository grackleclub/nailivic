package main

import "fmt"

// example of composite types a template might receive page
// and parse page.nav.bar, page.head, page.footer, etc
type page struct {
	Nav    navbar
	Head   head
	Footer footer
}

type navbar struct {
	Items []string
}
type head struct {
	Stylesheets []string
}
type footer struct {
	Year string
}

// type stylesheet struct{}

func getNewPage(name string) (page, error) {
	switch name {
	case "login":
		return page{}, fmt.Errorf("not implemented")
	case "index":
		newPage := page{
			Nav: navbar{
				Items: []string{"home", "about", "contact"},
			},
			Head: head{
				Stylesheets: []string{
					"static/css/zero.css",
					"static/css/style.css",
				},
			},
			Footer: footer{
				Year: "2021",
			},
		}
		return newPage, nil
	default:
		return page{}, fmt.Errorf("invalid page name: %s", name)
	}
}

func doStuffWrapper() {
	page, err := getNewPage("index")
	if err != nil {
		fmt.Println("error getting new page:", err)
		return
	}
	fmt.Println("new page:", page)
}
