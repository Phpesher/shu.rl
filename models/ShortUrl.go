package models

type Url struct {
	Id        string
	SourceUrl string
	NewUrl    string
}

// Constructor
func NewUrl(id, sourceUrl, newUrl string) *Url{
	return &Url{id, sourceUrl, newUrl}
}