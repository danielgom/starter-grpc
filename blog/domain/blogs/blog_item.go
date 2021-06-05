package blogs

import "go.mongodb.org/mongo-driver/bson/primitive"

type BlogItem struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	AuthorID string             `bson:"author_id"`
	Content  string             `bson:"content"`
	Title    string             `bson:"title"`
}

func (b *BlogItem) Id() primitive.ObjectID {
	if b != nil {
		return b.ID
	}
	return primitive.ObjectID{0}
}

func (b *BlogItem) GetAuthorID() string {
	if b != nil {
		return b.AuthorID
	}
	return ""
}

func (b *BlogItem) GetContent() string {
	if b != nil {
		return b.Content
	}
	return ""
}

func (b *BlogItem) GetTitle() string {
	if b != nil {
		return b.Title
	}
	return ""
}

func (b *BlogItem) SetId(id primitive.ObjectID) {
	b.ID = id
}

func (b *BlogItem) SetAuthorID(authorID string) {
	b.AuthorID = authorID
}

func (b *BlogItem) SetContent(content string) {
	b.Content = content
}

func (b *BlogItem) SetTitle(title string) {
	b.Title = title
}
