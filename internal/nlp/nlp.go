package nlp

import (
	"strings"
	"time"

	"github.com/tj/go-naturaldate"
)

// ExtractDate parses the natural language date from the title, removes it, and returns the cleaned title and parsed date.
// If no date is found, it returns the original title and nil.
func ExtractDate(title string) (string, *time.Time) {
	ref := time.Now()
	parsedTime, err := naturaldate.Parse(title, ref)
	if err != nil || parsedTime.Equal(ref) {
		return title, nil
	}

	words := strings.Fields(title)
	for i := len(words) - 1; i >= 0; i-- {
		prefix := strings.Join(words[:i], " ")
		suffix := strings.Join(words[i:], " ")

		tPrefix, _ := naturaldate.Parse(prefix, ref)
		tSuffix, _ := naturaldate.Parse(suffix, ref)

		if tPrefix.Equal(ref) && tSuffix.Equal(parsedTime) {
			prefix = strings.TrimSpace(prefix)
			stopWords := []string{"on", "at", "by", "for", "in", "due"}

			for {
				changed := false
				for _, sw := range stopWords {
					if strings.HasSuffix(strings.ToLower(prefix), " "+sw) {
						prefix = prefix[:len(prefix)-len(sw)-1]
						changed = true
					} else if strings.ToLower(prefix) == sw {
						prefix = ""
						changed = true
					}
				}
				if !changed {
					break
				}
			}

			if prefix == "" {
				return title, &parsedTime // If title is just a date, keep the original title
			}

			return prefix, &parsedTime
		}
	}

	return title, &parsedTime
}
