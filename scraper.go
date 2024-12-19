package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/gocolly/colly"
)

// entry is a template for a single calendar entry.
const entry = `
BEGIN:VEVENT
DTSTART{{.StartTime}}
DTEND{{.EndTime}}
DTSTAMP:{{.StampTime}}
SUMMARY:{{.Description}}
DESCRIPTION:{{.Title}}
END:VEVENT`

const stampTime = "20231111T072945Z"

const iCalHeaderBlock = `BEGIN:VCALENDAR
VERSION:2.0
PRODID:-//hacksw/handcal//NONSGML v1.0//EN
METHOD:PUBLISH`

const iCalFooterBlock = `
END:VCALENDAR`

func formatICalDateTime(date, time string) string {
	return ":" + date + "T" + time + "00Z"
}

// timetableEntry represents a single calendar entry.
type timetableEntry struct {
	Date        string
	StartTime   string
	EndTime     string
	StampTime   string
	Title       string
	Description string
}

func parseHeader(header string) string {
	parts := strings.Split(header, " – ")
	datePart := parts[1]
	parsedDate, err := time.Parse("Monday 2 January 2006", datePart)
	if err != nil {
		fmt.Println("bro")
	}
	iso8601Date := parsedDate.Format("20060102")
	return iso8601Date
}

func formatDate(dayOfMonth string, month string, year string) string {
	monthMap := map[string]string{
		"January":   "01",
		"February":  "02",
		"March":     "03",
		"April":     "04",
		"May":       "05",
		"June":      "06",
		"July":      "07",
		"August":    "08",
		"September": "09",
		"October":   "10",
		"November":  "11",
		"December":  "12",
	}

	monthNumber, exists := monthMap[month]
	if !exists {
		return "" // Invalid month
	}
	if len(dayOfMonth) == 1 {
		dayOfMonth = "0" + dayOfMonth
	}
	return year + monthNumber + dayOfMonth
}

func parseTableText(description string) string {
	desc_pattern := `(\d{1,2})\s*–\s*\w+\s*(\d{1,2})\s*(\w+)\s*(\d{4})`
	desc_re := regexp.MustCompile(desc_pattern)
	desc_matches := desc_re.FindStringSubmatch(description)

	if len(desc_matches) >= 5 {
		statement := parseHeader(description)
		return statement
	} else {
		statement := description
		return statement
	}

}

func writeEntries(entries []timetableEntry) {
	f, _ := os.Create("embl-conference-programme.ics")
	w := bufio.NewWriter(f)

	fmt.Fprintf(w, iCalHeaderBlock)

	// 🧠
	t := template.Must(template.New("entry").Parse(entry))
	for _, r := range entries {
		err := t.Execute(w, r)
		if err != nil {
			log.Println("wtf", err)
		}
	}

	fmt.Fprintf(w, iCalFooterBlock)
	w.Flush()
}

func main() {
	var allEntries []timetableEntry
	c := colly.NewCollector()

	fmt.Print("Enter conference programme URL: ")
	var website_url string
	fmt.Scanln(&website_url)

	url := website_url

	c.OnHTML("details.vf-details", func(section *colly.HTMLElement) {

		timetableentry := timetableEntry{
			StampTime: stampTime,
		}

		section.ForEach("summary.vf-details--summary", func(i int, row *colly.HTMLElement) {
			// Get date of event with regex
			text := row.Text
			pattern := `Day (\d+) – (\w+) (\d{1,2}) (\w+) (\d{4})`
			re := regexp.MustCompile(pattern)
			matches := re.FindAllStringSubmatch(text, -1)

			for _, match := range matches {
				dayOfMonth := match[3]
				month := match[4]
				year := match[5]
				formattedDate := formatDate(dayOfMonth, month, year)

				timetableentry.Date = formattedDate
			}
		})

		section.ForEach("tr", func(_ int, el *colly.HTMLElement) {
			descriptions := el.ChildText("td:nth-child(2)")
			statement := parseTableText(descriptions)

			if _, err := strconv.Atoi(statement); err == nil {
				timetableentry.Date = statement

			} else {
				description := el.ChildText("strong")
				if len(description) == 0 {
					description := el.ChildText("em")
					timetableentry.Description = description

					italics_title := strings.Replace(statement, description, "", -1)
					timetableentry.Title = italics_title
				} else {
					timetableentry.Description = description
					title := strings.Replace(statement, description, "", -1)
					timetableentry.Title = title
				}
			}

			// Get time of event with regex
			StampTime := el.ChildText("td:nth-child(1)")
			pattern := `(\d{2}:\d{2}) – (\d{2}:\d{2})`
			re := regexp.MustCompile(pattern)

			times := re.FindAllString(StampTime, -1)
			if len(times) > 0 {
				matches := re.FindStringSubmatch(times[0])
				startTime := matches[1]
				endTime := matches[2]

				startTimeStamp := strings.Replace(startTime, ":", "", -1)
				endTimeStamp := strings.Replace(endTime, ":", "", -1)

				startISO := formatICalDateTime(timetableentry.Date, startTimeStamp)
				endISO := formatICalDateTime(timetableentry.Date, endTimeStamp)
				timetableentry.StartTime = startISO
				timetableentry.EndTime = endISO
				allEntries = append(allEntries, timetableentry)
			}
		})
		writeEntries(allEntries)
	})
	c.OnError(func(_ *colly.Response, err error) {
		log.Println("Request URL:", url, "\nError:", err)
	})

	err := c.Visit(url)
	if err != nil {
		log.Fatal(err)
	}
}
