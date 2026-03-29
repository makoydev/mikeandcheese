package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/makoydev/mikeandcheese/backend/internal/models"
)

type PostHandler struct {
	DB *sql.DB
}

func NewPostHandler(db *sql.DB) *PostHandler {
	return &PostHandler{DB: db}
}

func (h *PostHandler) GetPosts(w http.ResponseWriter, r *http.Request) {
	rows, err := h.DB.Query(`
		SELECT id, title, slug, excerpt, content, author, category, tags, featured, created_at, updated_at
		FROM posts
		ORDER BY created_at DESC
	`)
	if err != nil {
		log.Printf("Error querying posts: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	posts, err := scanPosts(rows)
	if err != nil {
		log.Printf("Error scanning posts: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, posts)
}

func (h *PostHandler) GetPost(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	if slug == "" {
		http.Error(w, "Slug is required", http.StatusBadRequest)
		return
	}

	row := h.DB.QueryRow(`
		SELECT id, title, slug, excerpt, content, author, category, tags, featured, created_at, updated_at
		FROM posts
		WHERE slug = ?
	`, slug)

	var post models.Post
	err := row.Scan(
		&post.ID, &post.Title, &post.Slug, &post.Excerpt, &post.Content,
		&post.Author, &post.Category, &post.Tags, &post.Featured,
		&post.CreatedAt, &post.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}
	if err != nil {
		log.Printf("Error scanning post: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, post)
}

func (h *PostHandler) GetFeaturedPosts(w http.ResponseWriter, r *http.Request) {
	rows, err := h.DB.Query(`
		SELECT id, title, slug, excerpt, content, author, category, tags, featured, created_at, updated_at
		FROM posts
		WHERE featured = 1
		ORDER BY created_at DESC
	`)
	if err != nil {
		log.Printf("Error querying featured posts: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	posts, err := scanPosts(rows)
	if err != nil {
		log.Printf("Error scanning featured posts: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, posts)
}

func generateSlug(title string) string {
	slug := strings.ToLower(title)
	slug = strings.ReplaceAll(slug, " ", "-")
	reg := regexp.MustCompile(`[^a-z0-9-]`)
	slug = reg.ReplaceAllString(slug, "")
	return slug
}

func (h *PostHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
	var req models.CreatePostRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(req.Title) == "" || strings.TrimSpace(req.Content) == "" {
		http.Error(w, "Title and content are required", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(req.Slug) == "" {
		req.Slug = generateSlug(req.Title)
	}

	now := time.Now()
	result, err := h.DB.Exec(`
		INSERT INTO posts (title, slug, excerpt, content, author, category, tags, featured, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, req.Title, req.Slug, req.Excerpt, req.Content, req.Author, req.Category, req.Tags, req.Featured, now, now)
	if err != nil {
		log.Printf("Error creating post: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	id, err := result.LastInsertId()
	if err != nil {
		log.Printf("Error getting last insert ID: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	post := models.Post{
		ID:        id,
		Title:     req.Title,
		Slug:      req.Slug,
		Excerpt:   req.Excerpt,
		Content:   req.Content,
		Author:    req.Author,
		Category:  req.Category,
		Tags:      req.Tags,
		Featured:  req.Featured,
		CreatedAt: now,
		UpdatedAt: now,
	}

	writeJSON(w, http.StatusCreated, post)
}

func (h *PostHandler) UpdatePost(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	if slug == "" {
		http.Error(w, "Slug is required", http.StatusBadRequest)
		return
	}

	var req models.CreatePostRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	now := time.Now()
	newSlug := req.Slug
	if strings.TrimSpace(newSlug) == "" {
		newSlug = slug
	}

	result, err := h.DB.Exec(`
		UPDATE posts
		SET title = ?, slug = ?, excerpt = ?, content = ?, author = ?, category = ?, tags = ?, featured = ?, updated_at = ?
		WHERE slug = ?
	`, req.Title, newSlug, req.Excerpt, req.Content, req.Author, req.Category, req.Tags, req.Featured, now, slug)
	if err != nil {
		log.Printf("Error updating post: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Error getting rows affected: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if rowsAffected == 0 {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	// Fetch the updated post
	row := h.DB.QueryRow(`
		SELECT id, title, slug, excerpt, content, author, category, tags, featured, created_at, updated_at
		FROM posts
		WHERE slug = ?
	`, newSlug)

	var post models.Post
	err = row.Scan(
		&post.ID, &post.Title, &post.Slug, &post.Excerpt, &post.Content,
		&post.Author, &post.Category, &post.Tags, &post.Featured,
		&post.CreatedAt, &post.UpdatedAt,
	)
	if err != nil {
		log.Printf("Error fetching updated post: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, post)
}

func (h *PostHandler) DeletePost(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	if slug == "" {
		http.Error(w, "Slug is required", http.StatusBadRequest)
		return
	}

	result, err := h.DB.Exec(`DELETE FROM posts WHERE slug = ?`, slug)
	if err != nil {
		log.Printf("Error deleting post: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Error getting rows affected: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if rowsAffected == 0 {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func scanPosts(rows *sql.Rows) ([]models.Post, error) {
	var posts []models.Post
	for rows.Next() {
		var post models.Post
		err := rows.Scan(
			&post.ID, &post.Title, &post.Slug, &post.Excerpt, &post.Content,
			&post.Author, &post.Category, &post.Tags, &post.Featured,
			&post.CreatedAt, &post.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	// Return empty array instead of null in JSON
	if posts == nil {
		posts = []models.Post{}
	}
	return posts, nil
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Error encoding JSON response: %v", err)
	}
}
