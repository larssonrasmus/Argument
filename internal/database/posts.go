package database

import "fmt"

// Post is a new thread on the discussion board
type Post struct {
	ID          int
	Title       string
	DateCreated string
}

// GetPosts loads all the available posts from the database
func GetPosts() []Post {
	results, err := db.Query("SELECT id, title, created_at FROM posts")

	if err != nil {
		panic(fmt.Errorf("query posts: %w", err))
	}
	defer results.Close()

	var posts []Post

	for results.Next() {
		var post Post
		results.Scan(&post.ID, &post.Title, &post.DateCreated)
		posts = append(posts, post)
	}

	return posts
}

// SavePost persists a given post to database
func SavePost(post Post) (int64, error) {
	query, _ := db.Prepare("INSERT INTO posts(title, created_at, created_by) values (?, NOW(), 1)")

	result, err := query.Exec(post.Title)

	if err != nil {
		id, _ := result.LastInsertId()

		return id, nil
	}

	return 0, err
}
