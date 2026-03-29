package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

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
