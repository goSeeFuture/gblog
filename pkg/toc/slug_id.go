package toc

import (
	"github.com/gosimple/slug"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/util"
)

type Head struct {
	Level int
	ID    string
	Text  string
}

type SlugID struct {
	values   map[string]bool
	sequence []Head
}

func NewSlugID() parser.IDs {
	return &SlugID{
		values: map[string]bool{},
	}
}

func (s *SlugID) Generate(value []byte, kind ast.NodeKind, level int) []byte {
	id := slug.Make(util.BytesToReadOnlyString(value))
	s.sequence = append(s.sequence, Head{Level: level, ID: id, Text: util.BytesToReadOnlyString(value)})
	return util.StringToReadOnlyBytes(id)
}

func (s *SlugID) Put(value []byte) {
	s.values[util.BytesToReadOnlyString(value)] = true
}

func (s *SlugID) Heads() []Head {
	return s.sequence
}
