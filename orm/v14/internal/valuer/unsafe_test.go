//go:build v14

package valuer

import (
	"testing"
)

func Test_unsafeValue_SetColumns(t *testing.T) {
	testSetColumn(t, NewUnsafeValue)
}
