package console

// Filter helpers for common filter types

// TextFilter creates a text filter
func TextFilter(field, label string) Filter {
	return Filter{
		Field: field,
		Type:  FilterTypeText,
		Label: label,
	}
}

// NumberFilter creates a number filter
func NumberFilter(field, label string) Filter {
	return Filter{
		Field: field,
		Type:  FilterTypeNumber,
		Label: label,
	}
}

// DateFilter creates a date filter
func DateFilter(field, label string) Filter {
	return Filter{
		Field: field,
		Type:  FilterTypeDate,
		Label: label,
	}
}

// DateTimeFilter creates a datetime filter
func DateTimeFilter(field, label string) Filter {
	return Filter{
		Field: field,
		Type:  FilterTypeDateTime,
		Label: label,
	}
}

// BooleanFilter creates a boolean filter
func BooleanFilter(field, label string) Filter {
	return Filter{
		Field: field,
		Type:  FilterTypeBoolean,
		Label: label,
	}
}

// ChoiceFilter creates a choice filter
func ChoiceFilter(field, label string, choices []Choice) Filter {
	return Filter{
		Field:   field,
		Type:    FilterTypeChoice,
		Label:   label,
		Choices: choices,
	}
}

