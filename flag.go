package glarg

import (
	"net/url"
	"strings"

	"github.com/google/uuid"
)

const (
	DEFAULT_DELIMITER = ","
)

// This is an interface to enable the SliceFlag below
// to be able to operate on Slices without understanding
// the underlying slice type. Examples of the SliceFlagTarget
// are in this file.
type SliceFlagTarget interface {
	Clear()
	Append(item string) (SliceFlagTarget, error)
	Join(del string) string
	Get() interface{}
}

// Deal with getting multiple string values on the command line.
// By default it slices on comma, but you can change that
// during the subcommand setup.
// In theory, this could be expanded to work for anything, not
// just strings.
type SliceFlag struct {
	delimiter string
	target    SliceFlagTarget
}

func NewSliceFlag(v SliceFlagTarget, sep string) *SliceFlag {
	return &SliceFlag{
		delimiter: sep,
		target:    v,
	}
}

func (self SliceFlag) String() string {
	if self.target == nil {
		return ""
	} else {
		if self.delimiter == "" {
			return self.target.Join(DEFAULT_DELIMITER)
		} else {
			return self.target.Join(self.delimiter)
		}
	}
}

func (self *SliceFlag) Set(s string) error {
	var pieces []string
	if self.delimiter == "" {
		pieces = strings.Split(s, DEFAULT_DELIMITER)
	} else {
		pieces = strings.Split(s, self.delimiter)
	}

	if self.target == nil {
		self.target = &StringSliceFlagTarget{&[]string{}}
	}
	self.target.Clear()
	for _, v := range pieces {
		var err error
		self.target, err = self.target.Append(v)
		if err != nil {
			return err
		}
	}
	return nil
}

func (self SliceFlag) Get() interface{} {
	if self.target == nil {
		self.target = &StringSliceFlagTarget{&[]string{}}
	}
	return self.target.Get()
}

// StringSliceFlagTarget is a String Target for a
// SliceFlag.
type StringSliceFlagTarget struct {
	Target *[]string
}

func (self *StringSliceFlagTarget) makeSafe() {
	if self.Target == nil {
		self.Target = &[]string{}
	}
}

func (self *StringSliceFlagTarget) Clear() {
	self.makeSafe()
	*self.Target = (*self.Target)[:0]
}

func (self *StringSliceFlagTarget) Append(item string) (SliceFlagTarget, error) {
	self.makeSafe()
	*self.Target = append(*self.Target, item)
	return self, nil
}

func (self *StringSliceFlagTarget) Join(del string) string {
	self.makeSafe()
	return strings.Join(*self.Target, del)
}

func (self *StringSliceFlagTarget) Get() interface{} {
	self.makeSafe()
	return *self.Target
}

// UUIDSliceFlagTarget is a UUID Target for a
// SliceFlag.
type UUIDSliceFlagTarget struct {
	Target *[]uuid.UUID
}

func (self *UUIDSliceFlagTarget) makeSafe() {
	if self.Target == nil {
		self.Target = &[]uuid.UUID{}
	}
}

func (self *UUIDSliceFlagTarget) Clear() {
	self.makeSafe()
	*self.Target = (*self.Target)[:0]
}

func (self *UUIDSliceFlagTarget) Append(item string) (SliceFlagTarget, error) {
	self.makeSafe()
	if v, err := uuid.Parse(item); err != nil {
		return nil, err
	} else {
		*self.Target = append(*self.Target, v)
	}
	return self, nil
}

func (self *UUIDSliceFlagTarget) Join(del string) string {
	self.makeSafe()
	pieces := make([]string, len(*self.Target))
	for k, v := range *self.Target {
		pieces[k] = v.String()
	}
	return strings.Join(pieces, del)
}

func (self *UUIDSliceFlagTarget) Get() interface{} {
	self.makeSafe()
	return *self.Target
}

// UUID flag getter. Deals with parsing UUID inputs. This
// type implements the flag.Getter interface.
type UUIDFlag struct {
	ptr *uuid.UUID
}

func NewUUIDFlag(v *uuid.UUID) *UUIDFlag {
	return &UUIDFlag{
		ptr: v,
	}
}

func (self UUIDFlag) String() string {
	if self.ptr == nil {
		return uuid.Nil.String()
	} else {
		return self.ptr.String()
	}
}

func (self *UUIDFlag) Set(s string) error {
	if self.ptr == nil {
		self.ptr = &uuid.UUID{}
	}

	if v, err := uuid.Parse(s); err != nil {
		return err
	} else {
		*self.ptr = v
	}
	return nil
}

func (self UUIDFlag) Get() interface{} {
	if self.ptr == nil {
		return uuid.Nil
	} else {
		return *self.ptr
	}
}

// URLSliceFlagTarget is a URL Target for a
// SliceFlag.
type URLSliceFlagTarget struct {
	Target *[]*url.URL
}

func (self *URLSliceFlagTarget) makeSafe() {
	if self.Target == nil {
		self.Target = &[]*url.URL{}
	}
}

func (self *URLSliceFlagTarget) Clear() {
	self.makeSafe()
	*self.Target = (*self.Target)[:0]
}

func (self *URLSliceFlagTarget) Append(item string) (SliceFlagTarget, error) {
	self.makeSafe()
	if v, err := url.Parse(item); err != nil {
		return nil, err
	} else {
		*self.Target = append(*self.Target, v)
	}
	return self, nil
}

func (self *URLSliceFlagTarget) Join(del string) string {
	self.makeSafe()
	pieces := make([]string, len(*self.Target))
	for k, v := range *self.Target {
		pieces[k] = v.String()
	}
	return strings.Join(pieces, del)
}

func (self *URLSliceFlagTarget) Get() interface{} {
	self.makeSafe()
	return *self.Target
}

// URL flag getter. Deals with parsing URL inputs
// This type implements the flag.Getter interface.
type URLFlag struct {
	ptr *url.URL
}

func NewURLFlag(v *url.URL) *URLFlag {
	return &URLFlag{
		ptr: v,
	}
}

func (self URLFlag) String() string {
	if self.ptr == nil {
		return (&url.URL{}).String()
	} else {
		return self.ptr.String()
	}
}

func (self *URLFlag) Set(s string) error {
	if self.ptr == nil {
		self.ptr = &url.URL{}
	}

	if v, err := url.Parse(s); err != nil {
		return err
	} else {
		*self.ptr = *v
	}
	return nil
}

func (self URLFlag) Get() interface{} {
	if self.ptr == nil {
		return &url.URL{}
	} else {
		return self.ptr
	}
}
