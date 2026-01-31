package tui

// ListItem represents a single selectable item in the enhanced list.
// For movies, Children is nil. For series, Children holds seasons.
type ListItem struct {
	Columns  []string   // Column values for the list row
	Label    string     // Display name for selection summary
	Detail   []KeyValue // Key-value pairs for the detail panel
	Children []ListItem // Sub-items (seasons); nil for leaf items (movies)
}

// KeyValue is a labeled field in the detail panel.
type KeyValue struct {
	Key   string
	Value string
}

// Selection identifies what was selected for deletion.
type Selection struct {
	ItemIndex  int // index into top-level items
	ChildIndex int // -1 = entire parent, >= 0 = specific child
}

// ColumnDef defines a column in the list display.
type ColumnDef struct {
	Title string
	Width int
}

// FetchItemsFunc loads data and returns rich ListItems.
type FetchItemsFunc func() ([]ListItem, error)

// DeleteItemsFunc processes structured selections and returns result messages.
type DeleteItemsFunc func(selections []Selection) []string

// visibleRow maps a rendered list position back to the data model.
type visibleRow struct {
	itemIdx  int // index into items[]
	childIdx int // -1 for parent row, >= 0 for child row
}
