package filters

type Filter struct {
	Name  string
	Value int
}

var (
	Museum      = Filter{"museum", 1}
	Park        = Filter{"park", 2}
	Monument    = Filter{"monument", 4}
	Church      = Filter{"church", 8}
	Building    = Filter{"building", 16}
	FilterTypes = []Filter{Museum, Park, Monument, Church, Building}
)

type StringFilter []string

func (f StringFilter) Int() int {
	if len(f) == 0 {
		return emptyFilterValue()
	}

	result := 0
	for _, filter := range f {
		for _, type_ := range FilterTypes {
			if filter == type_.Name {
				result += type_.Value
			}
		}
	}
	return result
}

func emptyFilterValue() int {
	result := 0
	for _, type_ := range FilterTypes{
		result += type_.Value
	}
	return result
}