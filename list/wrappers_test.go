package list

import "testing"

func TestWrappedList(t *testing.T) {
	t.Run("SafeList", testSafeListWrapper)
	t.Run("SetList", testCompileSetList)
}
