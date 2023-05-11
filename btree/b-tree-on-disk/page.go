package b_tree_on_disk

import (
	"encoding/binary"
	"io"
	"os"
)

const PageSize = 4096

type Pager struct {
	file *os.File
	pos  int64
}
type Page struct {
	data [PageSize]byte
}

func (page *Page) ToNode() *Node {
	r := NewBinaryReader(page.data[:])
	isLeafB, _ := r.ReadByte()
	isLeaf := isLeafB == 1
	numOfKey, _ := r.ReadUint64At(1, binary.BigEndian)
	offset := uint64(1 + 8)
	node := NewNode(isLeaf)
	if !isLeaf {
		for i := uint64(0); i < numOfKey+1; i++ {
			childOffset, _ := r.ReadInt64At(offset, binary.BigEndian)
			offset += 8
			node.children = append(node.children, childOffset)
		}
	}
	for i := uint64(0); i < numOfKey; i++ {
		data, _ := r.ReadBytesAt(2*N, offset)
		offset += 2 * N
		key := data[:N]
		value := data[N:]
		node.items = append(node.items, &Item{key: [N]byte(key), value: [N]byte(value)})
	}
	return node
}
func (node *Node) ToPage() *Page {
	page := &Page{}
	w := NewBinaryWriter(page.data)
	w.WriteBool(node.isLeaf)
	numOfKeys := len(node.items)
	w.WriteUint64(binary.BigEndian, uint64(numOfKeys))
	if !node.isLeaf {
		for i := 0; i < numOfKeys+1; i++ {
			w.WriteInt64(binary.BigEndian, node.children[i])
		}
	}
	for i := 0; i < numOfKeys; i++ {
		w.WriteBytes(node.items[i].key[:])
		w.WriteBytes(node.items[i].value[:])
	}
	page.data = w.buf
	return page
}
func NewPager(path string) *Pager {
	file, _ := os.OpenFile(path, os.O_TRUNC|os.O_RDWR|os.O_CREATE, 0755)
	return &Pager{
		file: file,
		pos:  0,
	}
}

func (p *Pager) ReadPage(offset int64) *Page {
	_, err := p.file.Seek(offset, 0)
	if err != nil {
		return nil
	}
	buf := make([]byte, PageSize)
	_, err = io.ReadFull(p.file, buf)
	if err != nil {
		return nil
	}
	return &Page{
		data: [PageSize]byte(buf),
	}
}

func (p *Pager) WritePage(page *Page) (int64, error) {
	_, err := p.file.Seek(p.pos, 0)
	if err != nil {
		return 0, err
	}
	_, err = p.file.Write(page.data[:])
	if err != nil {
		return 0, err
	}
	pos := p.pos
	p.pos += PageSize
	return pos, nil
}

func (p *Pager) WritePageAtOffset(page *Page, offset int64) {
	_, err := p.file.Seek(offset, 0)
	if err != nil {
		panic(err)
	}
	_, err = p.file.Write(page.data[:])
	if err != nil {
		panic(err)
	}
}
