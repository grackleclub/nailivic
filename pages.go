package main

import "fmt"

// example of composite types a template might receive page
// and parse page.nav.bar, page.head, page.footer, etc
type page struct {
	// Nav    navbar
	Head   head
	Footer footer
}

type head struct {
	Title       string
	Stylesheets []string
}

//	type navbar struct {
//		Items []string
//	}
type footer struct {
	Year string
}

func getNewPage(name string) (page, error) {
	switch name {
	case "login":
		return page{}, fmt.Errorf("not implemented")
	case "crazy":
		newPage := page{
			Head: head{
				Title: "Crazy Title",
				Stylesheets: []string{
					"/static/css/zero.css",
					"/static/css/style.css",
				},
			},
			Footer: footer{
				Year: "9876",
			},
		}
		return newPage, nil
	case "index":
		newPage := page{
			Head: head{
				Title: "Nailivic",
				Stylesheets: []string{
					"/static/css/zero.css",
					"/static/css/style.css",
				},
			},
			Footer: footer{
				Year: "2024",
			},
		}
		return newPage, nil
	case "dash":
		newPage := page{
			Head: head{
				Title: "Nailivic Dashboard",
				Stylesheets: []string{
					"/static/css/zero.css",
					"/static/css/style.css",
				},
			},
			Footer: footer{
				Year: "2024",
			},
		}
		return newPage, nil
	case "special":
		newPage := page{
			Head: head{
				Title: "Nailivic",
				Stylesheets: []string{
					"/static/css/zero.css",
					"/static/css/style.css",
				},
			},
			Footer: footer{
				Year: "2024",
			},
		}
		return newPage, nil
	default:
		return page{}, fmt.Errorf("invalid page name: %s", name)
	}
}

// func doStuffWrapper() {
// 	page, err := getNewPage("index")
// 	if err != nil {
// 		fmt.Println("error getting new page:", err)
// 		return
// 	}
// 	fmt.Println("new page:", page)
// }
