# EMBL Conference ICS Generator

A Go-based web scraper that automatically converts EMBL conference program schedules into iCalendar (.ics) files for easy import into your calendar application.

## ⚠️ Disclaimer

This tool is designed specifically for EMBL conference program webpages. It may not work with other conference websites due to different HTML structures. This is not an official EMBL tool.

## 🎯 Features

- **Web Scraping**: Automatically extracts conference schedules from EMBL conference program pages
- **iCalendar Generation**: Creates standard .ics files compatible with all major calendar applications
- **Event Details**: Captures event titles, descriptions, dates, and times
- **Cross-Platform**: Works on Windows, macOS, and Linux
- **Easy to Use**: Simple command-line interface with minimal setup

## 📋 Requirements

- **Go 1.21.3** or higher
- Internet connection to access EMBL conference pages

## 🚀 Installation

### Prerequisites
1. Install Go from [golang.org](https://golang.org/dl/)
2. Verify installation: `go version`

### Setup
1. Clone or download this repository:
   ```bash
   git clone <repository-url>
   cd EMBLConferenceICS
   ```

2. Install dependencies:
   ```bash
   go mod tidy
   ```

## 💻 Usage

1. **Run the application**:
   ```bash
   go run scraper.go
   ```

2. **Enter the conference program URL** when prompted:
   ```
   Enter Conference Programme URL: https://www.embl.org/events/your-conference/programme/
   ```

3. **Find your generated file**: The script will create `embl-conference-programme.ics` in the current directory

4. **Import to your calendar**:
   - **Google Calendar**: Settings → Import & Export → Import
   - **Apple Calendar**: File → Import → Select your .ics file
   - **Outlook**: File → Open & Export → Import/Export → Import an iCalendar file
   - **Other apps**: Most calendar applications support .ics file imports

## 📁 Project Structure

```
EMBLConferenceICS/
├── scraper.go                    # Main application code
├── go.mod                        # Go module definition
├── go.sum                        # Dependency checksums
├── README.md                     # This file
└── embl-conference-programme.ics # Generated calendar file (after running)
```

## 🔧 How It Works

The application uses the [Colly](https://github.com/gocolly/colly) web scraping framework to:

1. **Parse HTML Structure**: Targets specific CSS selectors (`details.vf-details`, `summary.vf-details--summary`)
2. **Extract Event Data**: Captures dates, times, titles, and descriptions from table rows
3. **Format Data**: Converts dates and times to iCalendar format
4. **Generate .ics File**: Creates a properly formatted iCalendar file with VEVENT entries

## 📝 Example Output

The generated .ics file contains events in this format:
```
BEGIN:VEVENT
DTSTART:20240115T090000Z
DTEND:20240115T103000Z
DTSTAMP:20231111T072945Z
SUMMARY:Opening Session
DESCRIPTION:Welcome and Introduction
END:VEVENT
```





## 🐛 Troubleshooting

**Common Issues:**

1. **"No such file or directory"**: Ensure you're in the correct directory and Go is installed
2. **"Cannot find module"**: Run `go mod tidy` to install dependencies
3. **Empty .ics file**: Check that the URL is correct and the page is accessible
4. **Import errors**: Verify the .ics file was generated successfully

**Getting Help:**
- Check that the conference program URL is correct and accessible
- Ensure you have an active internet connection
- Verify the page structure hasn't changed (the scraper targets specific HTML elements)

---

**Note**: This tool is a work in progress and may need updates if EMBL changes their website structure.

