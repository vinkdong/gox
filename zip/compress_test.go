package zip

import "testing"

func TestCompress(t *testing.T) {
	Compress("../ssh", "test.zip", false)
}
