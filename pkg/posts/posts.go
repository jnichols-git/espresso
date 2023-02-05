package posts

import "github.com/google/uuid"

// Posts are Markdown files.
type PostVersion struct {
	PostID    string
	VersionID string
	Content   string
}

func NewPost(content string) *PostVersion {
	return &PostVersion{
		PostID:    uuid.NewString(),
		VersionID: uuid.NewString(),
		Content:   content,
	}
}

func NewVersion(post_id string, content string) *PostVersion {
	return &PostVersion{
		PostID:    post_id,
		VersionID: uuid.NewString(),
		Content:   content,
	}
}
