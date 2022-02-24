// Intset provides a set data structure for small non-negative integers.
package intset

import (
	"bytes"
	"fmt"
)

// N is 32 or 64 on 32-bit or 64-bit platform respectively.
const N = 32 << (^uint(0) >> 63)

// An IntSet is a set of small non-negative integers.
// Its zero value represents the empty set.
type IntSet struct {
	words []uint
}

// Has reports whether the set contains the non-negative value x.
func (s *IntSet) Has(x int) bool {
	word, bit := x/N, uint(x%N)
	return word < len(s.words) && s.words[word]&(1<<bit) != 0
}

// Add adds the non-negative value x to the set.
func (s *IntSet) Add(x int) {
	word, bit := x/N, uint(x%N)
	for word >= len(s.words) {
		s.words = append(s.words, 0)
	}
	s.words[word] |= 1 << bit
}

// AddAll adds each non-negative value in vals to the set.
func (s *IntSet) AddAll(vals ...int) {
	for _, x := range vals {
		s.Add(x)
	}
}

// UnionWith sets s to the union of s and t.
func (s *IntSet) UnionWith(t *IntSet) {
	for i, tword := range t.words {
		if i < len(s.words) {
			s.words[i] |= tword
		} else {
			s.words = append(s.words, tword)
		}
	}
}

// IntersectWith sets s to the intersect of s and t.
func (s *IntSet) IntersectWith(t *IntSet) {
	for i := range t.words {
		if i >= len(s.words) {
			break
		}
		s.words[i] &= t.words[i]
	}
}

// DifferenceWith sets s to the difference of s and t.
func (s *IntSet) DifferenceWith(t *IntSet) {
	for i := range t.words {
		if i >= len(s.words) {
			break
		}
		s.words[i] &^= t.words[i]
	}
}

// SymmetricDifference sets s to the symmetric difference of s and t.
func (s *IntSet) SymmetricDifference(t *IntSet) {
	for i, tword := range t.words {
		if i < len(s.words) {
			s.words[i] ^= tword
		} else {
			s.words = append(s.words, tword)
		}
	}
}

// String returns the set as a string of the form "{1 2 3}".
func (s IntSet) String() string {
	var buf bytes.Buffer
	buf.WriteByte('{')
	for i, word := range s.words {
		if word == 0 {
			continue
		}
		for j := 0; j < N; j++ {
			if word&(1<<j) == 0 {
				continue
			}
			if buf.Len() > len("{") {
				buf.WriteByte(' ')
			}
			fmt.Fprintf(&buf, "%d", i*N+j)
		}
	}
	buf.WriteByte('}')
	return buf.String()
}

// Len returns the number of elements in the set.
func (s *IntSet) Len() int {
	n := 0
	for _, word := range s.words {
		if word == 0 {
			continue
		}
		for j := 0; j < N; j++ {
			if word&(1<<j) != 0 {
				n++
			}
		}
	}
	return n
}

// Remove removes x from the set.
func (s *IntSet) Remove(x int) {
	word, bit := x/N, uint(x%N)
	if word >= len(s.words) {
		return
	}
	s.words[word] &^= 1 << bit
}

// Clear removes all elements from the set.
func (s *IntSet) Clear() {
	s.words = nil
}

// Copy returns a copy of the set.
func (s *IntSet) Copy() *IntSet {
	t := &IntSet{}
	t.words = make([]uint, len(s.words))
	copy(t.words, s.words)
	return t
}

// Elems return the set as a slice of int.
func (s *IntSet) Elems() []int {
	var res []int
	for i, word := range s.words {
		for j := 0; j < N; j++ {
			if word&(1<<j) != 0 {
				res = append(res, i*N+j)
			}
		}
	}
	return res
}
