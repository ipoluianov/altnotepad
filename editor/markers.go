package editor

import "sort"

// LineMarkers is a set of lines, e.g. the bookmarks, that follow the edits:
// the marked lines move with their text
type LineMarkers struct {
	lines []int // sorted
}

func (m *LineMarkers) Has(line int) bool {
	i := sort.SearchInts(m.lines, line)
	return i < len(m.lines) && m.lines[i] == line
}

func (m *LineMarkers) Add(line int) {
	i := sort.SearchInts(m.lines, line)
	if i < len(m.lines) && m.lines[i] == line {
		return
	}
	m.lines = append(m.lines, 0)
	copy(m.lines[i+1:], m.lines[i:])
	m.lines[i] = line
}

func (m *LineMarkers) Remove(line int) {
	i := sort.SearchInts(m.lines, line)
	if i < len(m.lines) && m.lines[i] == line {
		m.lines = append(m.lines[:i], m.lines[i+1:]...)
	}
}

// Toggle adds the line or removes it; returns whether it is marked now
func (m *LineMarkers) Toggle(line int) bool {
	if m.Has(line) {
		m.Remove(line)
		return false
	}
	m.Add(line)
	return true
}

func (m *LineMarkers) Clear()       { m.lines = nil }
func (m *LineMarkers) Len() int     { return len(m.lines) }
func (m *LineMarkers) Lines() []int { return append([]int(nil), m.lines...) }

// Set replaces the marked lines
func (m *LineMarkers) Set(lines []int) {
	m.lines = append([]int(nil), lines...)
	sort.Ints(m.lines)
	out := m.lines[:0]
	for i, l := range m.lines {
		if i == 0 || l != m.lines[i-1] {
			out = append(out, l)
		}
	}
	m.lines = out
}

// Next returns the first marked line after line, wrapping around; -1 if none
func (m *LineMarkers) Next(line int) int {
	if len(m.lines) == 0 {
		return -1
	}
	i := sort.SearchInts(m.lines, line+1)
	if i < len(m.lines) {
		return m.lines[i]
	}
	return m.lines[0]
}

// Prev returns the last marked line before line, wrapping around; -1 if none
func (m *LineMarkers) Prev(line int) int {
	if len(m.lines) == 0 {
		return -1
	}
	i := sort.SearchInts(m.lines, line) - 1
	if i >= 0 {
		return m.lines[i]
	}
	return m.lines[len(m.lines)-1]
}

// Adjust moves the lines after a change of the text
func (m *LineMarkers) Adjust(c *Change) {
	if len(m.lines) == 0 || (c.LinesAdded == 0 && c.LinesRemoved == 0 && c.multi == nil) {
		return
	}
	changed := false
	for i, l := range m.lines {
		nl := c.MapLine(l)
		if nl != l {
			m.lines[i] = nl
			changed = true
		}
	}
	if changed {
		m.Set(m.lines)
	}
}

// Range is a part of the text from Start to End
type Range struct {
	Start, End int
}

// RangeSet is a set of ranges of the text that follow the edits, e.g. the
// ranges marked by the search
type RangeSet struct {
	ranges []Range // sorted by Start, not overlapping
}

func (s *RangeSet) Len() int        { return len(s.ranges) }
func (s *RangeSet) Clear()          { s.ranges = nil }
func (s *RangeSet) Ranges() []Range { return s.ranges }

// Add adds a range; ranges overlapping it are joined with it
func (s *RangeSet) Add(r Range) {
	if r.End <= r.Start {
		return
	}
	// Adding in order is the usual case
	if n := len(s.ranges); n == 0 || s.ranges[n-1].End < r.Start {
		s.ranges = append(s.ranges, r)
		return
	}
	i := sort.Search(len(s.ranges), func(i int) bool { return s.ranges[i].End >= r.Start })
	j := i
	for j < len(s.ranges) && s.ranges[j].Start <= r.End {
		r.Start = min(r.Start, s.ranges[j].Start)
		r.End = max(r.End, s.ranges[j].End)
		j++
	}
	out := make([]Range, 0, len(s.ranges)-(j-i)+1)
	out = append(out, s.ranges[:i]...)
	out = append(out, r)
	out = append(out, s.ranges[j:]...)
	s.ranges = out
}

// InRange returns the ranges overlapping start..end
func (s *RangeSet) InRange(start, end int) []Range {
	i := sort.Search(len(s.ranges), func(i int) bool { return s.ranges[i].End > start })
	j := i
	for j < len(s.ranges) && s.ranges[j].Start < end {
		j++
	}
	return s.ranges[i:j]
}

// At returns the range holding the position
func (s *RangeSet) At(pos int) (Range, bool) {
	i := sort.Search(len(s.ranges), func(i int) bool { return s.ranges[i].End > pos })
	if i < len(s.ranges) && s.ranges[i].Start <= pos {
		return s.ranges[i], true
	}
	return Range{}, false
}

// Adjust moves the ranges after a change; removed text shrinks them
func (s *RangeSet) Adjust(c *Change) {
	if len(s.ranges) == 0 {
		return
	}
	i := sort.Search(len(s.ranges), func(i int) bool { return s.ranges[i].End >= c.Pos })
	if i == len(s.ranges) {
		return
	}
	delta := c.Inserted - c.Removed
	delEnd := c.Pos + c.Removed
	out := s.ranges[:i]
	for _, r := range s.ranges[i:] {
		if r.Start >= delEnd && (r.Start > c.Pos || c.Removed == 0) {
			r.Start += delta
			r.End += delta
		} else {
			// Overlaps the change: keep what is left of it
			start := r.Start
			if start >= c.Pos {
				start = c.Pos + c.Inserted
			}
			end := r.End
			if end > c.Pos {
				if end <= delEnd {
					end = c.Pos
				} else {
					end += delta
				}
			}
			r = Range{start, end}
		}
		if r.End > r.Start {
			out = append(out, r)
		}
	}
	s.ranges = out
}

// AdjustMulti moves the ranges after a change made of many edits
func (s *RangeSet) AdjustMulti(c *Change) {
	if c.multi == nil {
		s.Adjust(c)
		return
	}
	if len(s.ranges) == 0 {
		return
	}
	out := s.ranges[:0]
	for _, r := range s.ranges {
		r.Start, r.End = c.MapPos(r.Start), c.MapPos(r.End)
		if r.End > r.Start {
			out = append(out, r)
		}
	}
	s.ranges = out
}
