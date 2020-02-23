package render

import (
	"reflect"

	"github.com/rs/zerolog/log"
)

const ageCol = "AGE"

// HeaderColumn represent a table header
type HeaderColumn struct {
	Name      string
	Align     int
	Decorator DecoratorFunc
	Hide      bool
	Wide      bool
	MX        bool
	Time      bool
}

// Clone copies a header.
func (h HeaderColumn) Clone() HeaderColumn {
	return h
}

// ----------------------------------------------------------------------------

// Header represents a table header.
type Header []HeaderColumn

// Clone duplicates a header.
func (h Header) Clone() Header {
	header := make(Header, len(h))
	for i, c := range h {
		header[i] = c.Clone()
	}

	return header
}

// MapIndices returns a collection of mapped column indices based of the requested columns.
func (h Header) MapIndices(cols []string, wide bool, ii []int) {
	cc := make(map[int]struct{}, len(cols))
	var lastIndex int
	log.Debug().Msgf("MAP %d -- %d ", len(cols), len(ii))
	for i, col := range cols {
		idx := h.IndexOf(col, true)
		ii[i], cc[idx] = idx, struct{}{}
		lastIndex = i
	}
	if !wide {
		return
	}

	for i := range h {
		if _, ok := cc[i]; ok {
			continue
		}
		lastIndex++
		if lastIndex < len(ii) {
			ii[lastIndex] = i
		}
	}
}

// Customize builds a header from custom col definitions.
func (h Header) Customize(cols []string, wide bool) Header {
	if len(cols) == 0 {
		return h
	}
	cc := make(Header, 0, len(h))
	xx := make(map[int]struct{}, len(h))
	for _, c := range cols {
		idx := h.IndexOf(c, true)
		if idx == -1 {
			log.Warn().Msgf("Column %s is not available on this resource", c)
			continue
		}
		xx[idx] = struct{}{}
		col := h[idx].Clone()
		col.Wide = false
		cc = append(cc, col)
	}

	if !wide {
		return cc
	}

	for i, c := range h {
		if _, ok := xx[i]; ok {
			continue
		}
		col := c.Clone()
		col.Wide = true
		cc = append(cc, col)
	}

	return cc
}

// Diff returns true if the header changed.
func (h Header) Diff(header Header) bool {
	if len(h) != len(header) {
		return true
	}
	return !reflect.DeepEqual(h, header)
}

// Columns return header as a collection of strings.
func (h Header) Columns(wide bool) []string {
	if len(h) == 0 {
		return nil
	}
	var cc []string
	for _, c := range h {
		if !wide && c.Wide {
			continue
		}
		cc = append(cc, c.Name)
	}

	return cc
}

// HasAge returns true if table has an age column.
func (h Header) HasAge() bool {
	return h.IndexOf(ageCol, true) != -1
}

// AgeCol checks if given column index is the age column.
func (h Header) IsAgeCol(col int) bool {
	if !h.HasAge() || col >= len(h) {
		return false
	}
	return h[col].Time
}

// ValidColIndex returns the valid col index or -1 if none.
func (h Header) ValidColIndex() int {
	return h.IndexOf("VALID", true)
}

// IndexOf returns the col index or -1 if none.
func (h Header) IndexOf(colName string, includeWide bool) int {
	for i, c := range h {
		if c.Wide && !includeWide {
			continue
		}
		if c.Name == colName {
			return i
		}
	}
	return -1
}
