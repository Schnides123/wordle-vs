package util

import "encoding/json"

type Set struct {
	values map[string]bool
}

func NewSet() *Set {
	return &Set{
		values: make(map[string]bool),
	}
}

func (s *Set) Add(el string) {
	s.values[el] = true
}

func (s *Set) Remove(el string) {
	delete(s.values, el)
}

func (s Set) Contains(el string) bool {
	_, ok := s.values[el]
	return ok
}

func (s Set) Values() []string {
	vals := []string{}
	for k := range s.values {
		vals = append(vals, k)
	}
	return vals
}

func (s Set) Length() int {
	return len(s.values)
}

func (s Set) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Values())
}

func (s *Set) UnmarshalJSON(str []byte) error {
	var vals []string
	if err := json.Unmarshal(str, &vals); err != nil {
		return err
	}

	s.values = map[string]bool{}

	for _, v := range vals {
		s.Add(v)
	}

	return nil
}
