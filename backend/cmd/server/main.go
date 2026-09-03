package main

import (
	"database/sql"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func initDB() {
	dbPath := os.Getenv("DATABASE_PATH")
	if dbPath == "" {
		dbPath = "/app/data/audits.db"
	}
	os.MkdirAll(filepath.Dir(dbPath), 0755)

	var err error
	db, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		panic(err)
	}

	db.Exec(`CREATE TABLE IF NOT EXISTS audits (
		id TEXT PRIMARY KEY, url TEXT, status TEXT DEFAULT 'pending', score INTEGER, created_at TEXT, completed_at TEXT
	)`)
	db.Exec(`CREATE TABLE IF NOT EXISTS audit_categories (
		id TEXT PRIMARY KEY, audit_id TEXT, category TEXT, score INTEGER, findings TEXT
	)`)
	db.Exec(`CREATE TABLE IF NOT EXISTS audit_recommendations (
		id TEXT PRIMARY KEY, audit_id TEXT, priority TEXT, title TEXT, description TEXT
	)`)
}

func health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "healthy"})
}

func root(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"service": "Web Audits API", "version": "1.0.0"})
}

type AuditReq struct {
	URL string `json:"url" binding:"required"`
}

func createAudit(c *gin.Context) {
	var req AuditReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	aid := uuid.New().String()
	_, err := db.Exec("INSERT INTO audits (id,url,status,created_at) VALUES (?,?,?,?)",
		aid, req.URL, "running", time.Now().Format(time.RFC3339))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Simulate audit
	score := 45 + len(req.URL)%30
	db.Exec("UPDATE audits SET status=?,score=?,completed_at=? WHERE id=?",
		"completed", score, time.Now().Format(time.RFC3339), aid)

	// Insert categories
	categories := []struct {
		name    string
		score   int
		findings string
	}{
		{"performance", 70, "Found 5 scripts, 3 stylesheets, 10 images"},
		{"mobile", 100, "Viewport meta tag found"},
		{"seo", 80, "Title: yes, Meta desc: yes, H1: yes"},
		{"cta", 60, "Found 3 buttons, 1 forms"},
		{"accessibility", 90, "Images with alt: 8, without: 2"},
	}

	for _, cat := range categories {
		cid := uuid.New().String()
		db.Exec("INSERT INTO audit_categories (id,audit_id,category,score,findings) VALUES (?,?,?,?,?)",
			cid, aid, cat.name, cat.score, cat.findings)
	}

	// Insert recommendations
	recs := []struct {
		priority string
		title    string
		desc     string
	}{
		{"high", "Add viewport meta tag", "Mobile usability is critical for SEO"},
		{"medium", "Optimize performance", "Reduce scripts and images"},
		{"low", "Add alt text to images", "Improve accessibility"},
	}

	for _, rec := range recs {
		rid := uuid.New().String()
		db.Exec("INSERT INTO audit_recommendations (id,audit_id,priority,title,description) VALUES (?,?,?,?,?)",
			rid, aid, rec.priority, rec.title, rec.desc)
	}

	c.JSON(http.StatusOK, gin.H{"audit_id": aid, "status": "completed"})
}

func listAudits(c *gin.Context) {
	rows, err := db.Query("SELECT * FROM audits ORDER BY created_at DESC")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var audits []gin.H
	for rows.Next() {
		var id, url, status, createdAt, completedAt string
		var score int
		rows.Scan(&id, &url, &status, &score, &createdAt, &completedAt)
		audits = append(audits, gin.H{
			"id": id, "url": url, "status": status, "score": score,
			"created_at": createdAt, "completed_at": completedAt,
		})
	}
	c.JSON(http.StatusOK, gin.H{"audits": audits})
}

func getAudit(c *gin.Context) {
	aid := c.Param("aid")
	var id, url, status, createdAt, completedAt string
	var score int
	err := db.QueryRow("SELECT * FROM audits WHERE id=?", aid).Scan(&id, &url, &status, &score, &createdAt, &completedAt)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Audit not found"})
		return
	}

	rows, err := db.Query("SELECT * FROM audit_categories WHERE audit_id=?", aid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var categories []gin.H
	for rows.Next() {
		var id, auditID, category, findings string
		var score int
		rows.Scan(&id, &auditID, &category, &score, &findings)
		categories = append(categories, gin.H{
			"id": id, "audit_id": auditID, "category": category, "score": score, "findings": findings,
		})
	}

	rows2, err := db.Query("SELECT * FROM audit_recommendations WHERE audit_id=? ORDER BY priority", aid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows2.Close()

	var recs []gin.H
	for rows2.Next() {
		var id, auditID, priority, title, description string
		rows2.Scan(&id, &auditID, &priority, &title, &description)
		recs = append(recs, gin.H{
			"id": id, "audit_id": auditID, "priority": priority, "title": title, "description": description,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"id": id, "url": url, "status": status, "score": score,
		"created_at": createdAt, "completed_at": completedAt,
		"categories": categories, "recommendations": recs,
	})
}

func main() {
	initDB()
	defer db.Close()

	r := gin.Default()
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "*")
		c.Header("Access-Control-Allow-Headers", "*")
		c.Next()
	})

	r.GET("/health", health)
	r.GET("/", root)
	r.POST("/audits", createAudit)
	r.GET("/audits", listAudits)
	r.GET("/audits/:aid", getAudit)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	r.Run(":" + port)
}
