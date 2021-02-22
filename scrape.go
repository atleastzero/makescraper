// EXAMPLE: chromedp example demonstrating how to use a selector to click on an element.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
)

type Major struct {
	Name               string `json:"name"`
	DegreeRequirements string `json:"degree requirements"`
}

type Course struct {
	Department    string `json:"department"`
	Number        string `json:"number"`
	Name          string `json:"name"`
	Prerequisites string `json:"prerequisites"`
	Type          string `json:"type"`
	// Syllabus      string `json:"syllabus"`
	Description   string `json:"description"`
}

type University struct {
	Name    string   `json:"name"`
	Majors  []Major  `json:"majors"`
	Courses []Course `json:"courses"`
}

func main() {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.DisableGPU,
		chromedp.Flag("headless", false),
	)

	ctx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	major := Major {
		Name: "B.S. in Applied Computer Science",
		DegreeRequirements: "9 Science & Letters Courses (S&L)" +
			"9 Software Product Development Courses (SPD)" +
			"6 Intensives (INT)" +
			"5 Foundation Courses" +
			"1 Supervised Work Experience (SWE) Course" +
			"6 Concentration Courses (Note: Data Science has 5)" +
			"4+ Technical Electives" +
			"1 Foreign Language",
	}

	var majors []Major
	majors = append(majors, major)

	var courses []Course

	var rows []*cdp.Node

	university := University{
		Name: "Make School",
		Majors: majors,
		Courses: courses,
	}

	// create chrome instance
	ctx, cancel = chromedp.NewContext(ctx)
	defer cancel()

	err := chromedp.Run(ctx,
		chromedp.Navigate(`https://makeschool.com/course-offerings`),
		chromedp.WaitVisible(`//a[@class='index-Anchor--link-1WshRApLNQqrEwrQmzFIhE']`),
	)
	if err != nil {
		log.Fatal(err)
	}

	err = chromedp.Run(ctx,
		chromedp.Nodes(
			`.CatalogTable-TableRow-1tCv9xxpmFyhZVtK5_DGzu`,
			&rows,
			chromedp.ByQueryAll,
		),
	)
	if err != nil {
		log.Fatal(err)
	}

	for _, row := range rows {
		var course Course

		err := chromedp.Run(ctx,
			chromedp.Text(row.FullXPathByID() + `/div[2]`, &course.Department),
			chromedp.Text(row.FullXPathByID() + `/div[3]`, &course.Number),
			chromedp.Text(row.FullXPathByID() + `/div[4]`, &course.Name),
		)
		if err != nil {
			log.Fatal(err)
		}

		err = chromedp.Run(ctx,
			chromedp.Click(row.FullXPathByID(), chromedp.NodeVisible),
			chromedp.WaitVisible(`//div[@class='ReactModalPortal'][1]`),
		)
		if err != nil {
			log.Fatal(err)
		}

		portalPath := `//div[@class='ReactModalPortal'][1]/div[1]/div[1]/div[1]/div[2]/div[1]/div[1]/`
		err = chromedp.Run(ctx,
			chromedp.Text(portalPath + `p[1]`, &course.Prerequisites),
			chromedp.Text(portalPath + `div[1]/div[1]/span[1]`, &course.Type),
			// chromedp.Text(portalPath + `div[2]/a[1]`, &course.Syllabus),
			chromedp.Text(portalPath + `p[2]`, &course.Description),
		)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("hello")

		err = chromedp.Run(ctx,
			chromedp.Click(`ReactModal__Overlay`, chromedp.NodeNotVisible),
		)
		if err != nil {
			log.Fatal(err)
		} else {
			fmt.Printf(course.Name)
			university.Courses = append(university.Courses, course)
		}
	}

	makeSchoolJson, err := json.Marshal(university)
	if err != nil {
		log.Fatal(err)
	}

	ioutil.WriteFile("MakeSchool.json", makeSchoolJson, 0644)
}
