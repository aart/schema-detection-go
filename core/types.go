package core

// LineInfo holds a line of content along with its origin.
type LineInfo struct {
	Content    string
	Filename   string
	LineNumber int
}