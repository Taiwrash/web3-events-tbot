package scrape

import (
	"fmt"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
)

type Event struct {
	Title    string `json:"title"`
	Location string `json:"location"`
	Date     string `json:"date"`
	URL      string `json:"url"`
}

// var events []Event

func Scrape() ([]Event, error) {
	c := colly.NewCollector()
	eventsMap := make(map[string]Event)

	c.OnHTML("section.event-card-details", func(e *colly.HTMLElement) {
		title := e.ChildText("h2")
		location := e.ChildAttr("a", "data-event-location")
		date := e.ChildText("p")
		url := e.ChildAttr("a", "href")

		event := Event{
			Title:    title,
			Location: location,
			Date:     date,
			URL:      url,
		}

		key := title + date
		if _, exists := eventsMap[key]; !exists {
			eventsMap[key] = event
		}
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	err := c.Visit("https://www.eventbrite.com/d/nigeria/blockchain-events/")
	// err := c.Visit("https://www.eventbrite.com/d/nigeria/devops/")
	if err != nil {
		return nil, err
	}
	fmt.Println("visited")

	events := make([]Event, 0, len(eventsMap))
	for _, event := range eventsMap {
		events = append(events, event)
	}

	return events, nil
}

func DecodeScrape(events []Event, search string) ([]Event, error) {
	switch strings.ToLower(search) {
	case "January":
		search = "jan"
	case "February":
		search = "feb"
	case "March":
		search = "mar"
	case "April":
		search = "apr"
	case "May":
		search = "may"
	case "June":
		search = "jun"
	case "July":
		search = "jul"
	case "August":
		search = "aug"
	case "September":
		search = "sep"
	case "October":
		search = "oct"
	case "November":
		search = "nov"
	case "December":
		search = "dec"
	default:
		search = fmt.Sprintf("%v", time.Now().Month())[:3]
	}

	searchMonth, err := time.Parse("Jan", search)
	if err != nil {
		return nil, err
	}

	var matchingEvents []Event

	for _, event := range events {
		parts := strings.SplitN(event.Date, ",", 3)
		if len(parts) < 3 {
			fmt.Printf("Invalid date format '%s':\n", event.Date)
			continue
		}

		datePart := strings.TrimSpace(parts[0] + "," + parts[1])

		eventDate, err := time.Parse("Mon, Jan 2", datePart)
		if err != nil {
			fmt.Printf("Error parsing date '%s':\n", datePart)
			continue
		}
		if eventDate.Month() == searchMonth.Month() {
			matchingEvents = append(matchingEvents, event)
		}
	}

	if len(matchingEvents) == 0 {
		return nil, fmt.Errorf("no events found for month %s", search)
	}

	return matchingEvents, nil
}
