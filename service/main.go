package main

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	dbPath := filepath.Join(os.Getenv("HOME"), ".rfc/sqlite3/rfc.db")
	err := initializeDB(dbPath)
	if err != nil {
		log.Fatalf("Could not initialize database: %v", err)
	}

	r := gin.Default()
	//c.Header("Cache-Control", "public, max-age=21600")
	addHandler(r, dbPath)

	//r.Use(cors.Default())
	r.Run() // listen and serve on 0.0.0.0:8080
}

func addHandler(r *gin.Engine, dbPath string) {
	r.Use(cors.Default())
	r.GET("/rfcs", func(c *gin.Context) {

		rfcs, err := getAllRFCs(dbPath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, rfcs)
	})

	r.POST("/refresh_rfcs", func(c *gin.Context) {
		rfcs, err := refreshRFCs(dbPath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, rfcs)
	})

	r.GET("/rfc/:rfc_id", func(c *gin.Context) {
		rfcID := c.Param("rfc_id")
		content, err := downloadRFC(rfcID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"id": rfcID, "content": content})
	})

	r.POST("/download/rfc/:rfc_id", func(c *gin.Context) {
		rfcID := c.Param("rfc_id")
		dbPath := filepath.Join(os.Getenv("HOME"), ".rfc/sqlite3/rfc.db")
		downloadsPath := filepath.Join(os.Getenv("HOME"), ".rfc/downloads/rfcs")

		db, err := sql.Open("sqlite3", dbPath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open database"})
			return
		}
		defer db.Close()

		if err := downloadAndSaveRFC(db, rfcID, downloadsPath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "RFC downloaded successfully", "rfc_id": rfcID})
	})

}

func initializeDB(dbPath string) error {
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		dir := filepath.Dir(dbPath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}

		db, err := sql.Open("sqlite3", dbPath)
		if err != nil {
			return err
		}
		defer db.Close()

		// Create tables
		sqlStmt := `
		CREATE TABLE IF NOT EXISTS rfc_update (update_date TIMESTAMP);
		CREATE TABLE IF NOT EXISTS rfc (
			id STRING PRIMARY KEY, 
			title STRING, 
			authors STRING, 
			year STRING, 
			month STRING, 
			obsoleted_by STRING, 
			obsoletes STRING, 
			updates STRING, 
			updated_by STRING, 
			also STRING, 
			status STRING, 
			is_downloaded BOOLEAN NOT NULL DEFAULT 0,
			updated_at TIMESTAMP, 
			created_at TIMESTAMP
		);
		`
		_, err = db.Exec(sqlStmt)
		if err != nil {
			return fmt.Errorf("Could not create tables: %v", err)
		}
	}
	return nil
}

// RFC represents a single RFC record.
type RFC struct {
	ID           string `json:"id"`
	Title        string `json:"title"`
	Authors      string `json:"authors"`
	Year         string `json:"year"`
	Month        string `json:"month"`
	ObsoletedBy  string `json:"obsoleted_by"`
	Obsoletes    string `json:"obsoletes"`
	Updates      string `json:"updates"`
	UpdatedBy    string `json:"updated_by"`
	Also         string `json:"also"`
	Status       string `json:"status"`
	IsDownloaded bool   `json:"is_downloaded"`
	UpdatedAt    string `json:"updated_at"`
	CreatedAt    string `json:"created_at"`
}

func refreshRFCs(dbPath string) ([]RFC, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	defer db.Close()
	_, err = db.Exec("DELETE FROM rfc")
	if err != nil {
		return nil, fmt.Errorf("error truncating rfc table: %v", err)
	}

	err = fetchRFCs(db)
	if err != nil {
		return nil, err
	}

	return getUpdateRFCs(db)
}

func getAllRFCs(dbPath string) ([]RFC, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	// Check if there are any rows
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM rfc").Scan(&count)
	if err != nil {
		return nil, err
	}

	// If no rows, fetch RFCs
	if count == 0 {
		err = fetchRFCs(db)
		if err != nil {
			return nil, err
		}
	}

	// Now query all rows
	// Check for errors after iterating through rows
	return getUpdateRFCs(db)
}

func getUpdateRFCs(db *sql.DB) ([]RFC, error) {
	rows, err := db.Query("SELECT id, title,authors, year, month, obsoleted_by, obsoletes, updates, updated_by, status, is_downloaded FROM rfc")
	if err != nil {
		return nil, fmt.Errorf("error querying RFCs: %v", err)
	}
	defer rows.Close()

	var rfcs []RFC
	for rows.Next() {
		var r RFC
		if err := rows.Scan(&r.ID, &r.Title, &r.Authors, &r.Year, &r.Month, &r.ObsoletedBy, &r.Obsoletes, &r.Updates, &r.UpdatedBy, &r.Status, &r.IsDownloaded); err != nil {
			return nil, fmt.Errorf("error scanning RFC row: %v", err)
		}
		rfcs = append(rfcs, r)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating through RFC rows: %v", err)
	}

	return rfcs, nil
}

func formatRFCID(rfcID string) string {
	// Ensure rfcID is in lowercase
	rfcID = strings.ToLower(rfcID)

	// Split "rfc" and the numeric part
	numericPart := strings.TrimPrefix(rfcID, "rfc")

	// Convert the numeric part to an integer
	num, err := strconv.Atoi(numericPart)
	if err != nil {
		// Handle the error if the conversion fails
		// This might happen if the rfcID format is not as expected
		// For simplicity, we'll just return the original rfcID in this case
		return rfcID
	}

	// Combine "rfc" with the numeric part without leading zeros
	formattedRFCID := fmt.Sprintf("rfc%d", num)

	return formattedRFCID
}

func downloadRFC(rfcID string) (string, error) {
	rfcID = formatRFCID(rfcID)
	downloadsPath := filepath.Join(os.Getenv("HOME"), ".rfc/downloads/rfcs")
	filePath := filepath.Join(downloadsPath, fmt.Sprintf("%s.html", rfcID))

	// Check if the RFC file exists
	if _, err := os.Stat(filePath); err == nil {
		// File exists, read and return its content
		content, err := os.ReadFile(filePath)
		if err != nil {
			return "", fmt.Errorf("failed to read the downloaded RFC file: %s", err)
		}
		return string(content), nil
	}

	resp, err := http.Get(fmt.Sprintf("https://www.rfc-editor.org/rfc/%s.html", rfcID))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

// downloadAndSaveRFC downloads the RFC content and saves it to the specified directory.
func downloadAndSaveRFC(db *sql.DB, rfcID, basePath string) error {
	// Format RFC ID to lower case and ensure correct directory structure
	formattedRFCID := formatRFCID(rfcID) // Assuming formatRFCID formats the RFC ID as previously described
	rfcPath := filepath.Join(basePath, fmt.Sprintf("%s.html", formattedRFCID))

	// Download RFC content
	resp, err := http.Get(fmt.Sprintf("https://www.rfc-editor.org/rfc/%s.html", formattedRFCID))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// Ensure the target directory exists
	if err := os.MkdirAll(filepath.Dir(rfcPath), 0755); err != nil {
		return err
	}

	// Save the content to a file
	if err := os.WriteFile(rfcPath, body, 0644); err != nil {
		return err
	}

	// Update the database to mark the RFC as downloaded
	_, err = db.Exec("UPDATE rfc SET is_downloaded = 1 WHERE id = ?", rfcID)
	if err != nil {
		return err
	}

	return nil
}
