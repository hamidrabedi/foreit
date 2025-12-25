# Utils Module

Shared utility functions for Gogo framework.

## Functions

### ID Generation
- `GenerateID()` - Generates a unique ID
- `GenerateToken(length)` - Generates a random token

### String Utilities
- `Slugify(s)` - Converts string to URL-friendly slug
- `Truncate(s, maxLen)` - Truncates string
- `Pluralize(word, count)` - Returns plural form

### Formatting
- `FormatDuration(d)` - Formats duration
- `FormatBytes(bytes)` - Formats bytes

### Slice Utilities
- `Contains(slice, value)` - Checks if slice contains value
- `Unique(slice)` - Removes duplicates
- `Map(slice, fn)` - Maps function over slice
- `Filter(slice, fn)` - Filters slice

## Usage

```go
id := utils.GenerateID()
slug := utils.Slugify("Hello World") // "hello-world"
unique := utils.Unique([]string{"a", "b", "a"}) // ["a", "b"]
```

