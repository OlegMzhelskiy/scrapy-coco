package models

import (
	"fmt"
	"strings"
)

type Event struct {
	Name        string
	Description string
	Meta        string
	Instructor  string
	Date        string
	Time        string
	URL         string
}

// Key form text for store
func (e Event) Key() string {
	return fmt.Sprintf("%s_%s_%s", e.Name, e.Date, e.Time)
}

func (e Event) String() string {
	b := strings.Builder{}
	b.WriteString(e.Name)
	if e.Meta != "" {
		b.WriteString("\n\n")
		b.WriteString("‼️ " + e.Meta)
	}
	b.WriteString("\n\n")
	b.WriteString("🗓 " + e.Date)
	b.WriteString("\n")
	b.WriteString(e.Time)
	b.WriteString("\n\n")
	b.WriteString("Instructor: " + e.Instructor)
	b.WriteString("\n\n")
	b.WriteString(e.Description)

	if e.URL != "" {
		b.WriteString("\n\n")
		b.WriteString(e.URL)
	}

	return b.String()
}
