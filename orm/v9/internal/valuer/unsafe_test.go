//go:build v9

package valuer

import (
	"testing"
)

func Test_unsafeValue_SetColumns(t *testing.T) {
	testSetColumn(t, NewUnsafeValue)
}
