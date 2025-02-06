package wheel

import "testing"

func TestArraySortDirection(t *testing.T) {
	if !Ascending.IsValid() {
		t.Error("direction Ascending is invalid")
	}
	if !Descending.IsValid() {
		t.Error("direction Descending is invalid")
	}

	d := SortOrder(69)
	if d.IsValid() {
		t.Errorf("invalid direction %d is valid", d)
	}
}
