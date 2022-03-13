package tools

// NewSet constructor
func NewSet() *Set {
	var set Set
	set.Elements = make(map[Any]struct{})
	return &set
}

// Add - adds an element to our set
func (s *Set) Add(elm Any) {
	s.Elements[elm] = struct{}{}
}

func (s *Set) Delete(elm Any) error {
	if _, exists := s.Elements[elm]; !exists {
		return ErrElementNotFound
	}
	delete(s.Elements, elm)
	return nil
}

func (s *Set) Contains(elm Any) bool {
	_, exists := s.Elements[elm]
	return exists
}

// ToSliceOfAny - return a slice of Any
func (s *Set) ToSliceOfAny() []Any {
	var result []Any
	for k := range s.Elements {
		result = append(result, k)
	}
	return result
}
