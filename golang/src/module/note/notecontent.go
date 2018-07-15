package note

import (
	"html/template"
)

type NoteContent struct {
	Path		string
	Hash		string
	EditTime	string
	Title		string
	Format		string
	Content		string
	ToHTML		interface{}
}

func NewNote(path string) *NoteContent {
	if path[0] == 47 {
		path = path[1:]
	}
	return &NoteContent{
		Path:	path,
		Hash:	PathHash(path),
	}
}

func (n *NoteContent) LoadData() error {
	return stmtQueryNoteData.QueryRow(n.Hash).Scan(&n.EditTime,&n.Title,&n.Format,&n.Content)
}

func (n *NoteContent) Show() {
	n.EditTime = n.EditTime[:10]
	switch n.Format {
	case "rich":
		n.ToHTML = template.HTML(n.Content)
	default:
		n.ToHTML = n.Content
	}
}