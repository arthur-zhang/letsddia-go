package file_utils

import "os"

type OpenOptions struct {
	read      bool
	write     bool
	append    bool
	truncate  bool
	create    bool
	createNew bool
	mode      os.FileMode
}

func NewOpenOptions() *OpenOptions {
	return &OpenOptions{
		read:      false,
		write:     false,
		append:    false,
		truncate:  false,
		create:    false,
		createNew: false,
		mode:      0666,
	}
}
func (o *OpenOptions) Read(read bool) *OpenOptions {
	o.read = read
	return o
}
func (o *OpenOptions) Write(write bool) *OpenOptions {
	o.write = write
	return o
}
func (o *OpenOptions) Append(append bool) *OpenOptions {
	o.append = append
	return o
}
func (o *OpenOptions) Truncate(truncate bool) *OpenOptions {
	o.truncate = truncate
	return o
}
func (o *OpenOptions) Create(create bool) *OpenOptions {
	o.create = create
	return o
}
func (o *OpenOptions) Mode(mode os.FileMode) *OpenOptions {
	o.mode = mode
	return o
}

func (o *OpenOptions) Open(path string) (*os.File, error) {
	var flag int
	if o.read {
		flag |= os.O_RDONLY
	}
	if o.write {
		flag |= os.O_WRONLY
	}
	if o.append {
		flag |= os.O_APPEND
	}
	if o.truncate {
		flag |= os.O_TRUNC
	}
	if o.create {
		flag |= os.O_CREATE
	}
	return os.OpenFile(path, flag, o.mode)
}
