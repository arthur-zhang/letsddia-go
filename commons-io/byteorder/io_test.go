package byteorder

import (
	"os"
	"testing"
)

func TestBinaryIO(t *testing.T) {
	f, _ := os.OpenFile("/tmp/test", os.O_CREATE|os.O_RDWR, 0666)
	defer f.Close()
	//file := NewBinaryIO(f)

}
