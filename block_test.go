package block

import "testing"

func TestBlock(t *testing.T) {
	t.Logf("timeout 1s, block 5s, blockWrap returns %v", blockWrap(block, 1, 5))
	t.Logf("timeout 5s, block 1s, blockWrap returns %v", blockWrap(block, 5, 1))
}
