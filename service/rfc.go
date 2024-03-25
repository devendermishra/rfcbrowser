package main

import (
	"database/sql"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// Structs for XML parsing
type RFCIndex struct {
	Entries []RFCEntry `xml:"rfc-entry"`
}

type RFCEntry struct {
	DocID       string   `xml:"doc-id"`
	Title       string   `xml:"title"`
	Authors     []Author `xml:"author"`
	Date        RFCDate  `xml:"date"`
	Status      string   `xml:"current-status"`
	Obsoletes   []DocID  `xml:"obsoletes>doc-id"`
	ObsoletedBy []DocID  `xml:"obsoleted-by>doc-id"`
	Updates     []DocID  `xml:"updates>doc-id"`
	UpdatedBy   []DocID  `xml:"updated-by>doc-id"`
}

type Author struct {
	Name string `xml:"name"`
}

type RFCDate struct {
	Year  string `xml:"year"`
	Month string `xml:"month"`
}

type DocID struct {
	Value string `xml:",chardata"`
}

// FetchRFCs downloads and parses the RFC index, then inserts entries into the database.
func fetchRFCs(db *sql.DB) error {
	resp, err := http.Get("https://www.rfc-editor.org/rfc-index.xml")
	if err != nil {
		return fmt.Errorf("error fetching RFC index: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading RFC index response: %v", err)
	}

	var index RFCIndex
	err = xml.Unmarshal(body, &index)
	if err != nil {
		return fmt.Errorf("error parsing RFC index XML: %v", err)
	}

	for _, entry := range index.Entries {
		authors := make([]string, len(entry.Authors))
		for i, author := range entry.Authors {
			authors[i] = author.Name
		}

		obsoletes := joinDocIDs(entry.Obsoletes)
		obsoletedBy := joinDocIDs(entry.ObsoletedBy)
		updates := joinDocIDs(entry.Updates)
		updatedBy := joinDocIDs(entry.UpdatedBy)

		_, err := db.Exec("INSERT INTO rfc (id, title, authors, year, month, obsoleted_by, obsoletes, updates, updated_by, status, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
			entry.DocID, entry.Title, strings.Join(authors, ", "), entry.Date.Year, entry.Date.Month, obsoletedBy, obsoletes, updates, updatedBy, entry.Status, time.Now(), time.Now())
		if err != nil {
			return fmt.Errorf("error inserting RFC entry into database: %v", err)
		}
	}

	return nil
}

// Helper function to join DocID slices into comma-separated strings
func joinDocIDs(docs []DocID) string {
	values := make([]string, len(docs))
	for i, doc := range docs {
		values[i] = doc.Value
	}
	return strings.Join(values, ", ")
}
