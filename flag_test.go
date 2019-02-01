package glarg

import (
	"flag"
	"fmt"
	"net/url"
	"testing"

	"github.com/google/uuid"
)

func TestUUIDFlag(t *testing.T) {
	expected1 := uuid.Nil
	flag1 := &UUIDFlag{}

	// a default UUIDFlag is Nil.
	if flag1.String() != expected1.String() {
		t.Errorf("Error. Expected: %s. Received: %s.", expected1.String(), flag1.String())
	}

	if tmp, ok := flag1.Get().(uuid.UUID); !ok {
		t.Errorf("Error. Expected a UUID. Got something else.")
	} else if tmp != expected1 {
		t.Errorf("Error. Expected: %s. Received: %s", expected1, tmp)
	}

	// Same tests as above, only with a provided pointer.
	expected2 := uuid.New()
	foo := expected2
	flag2 := NewUUIDFlag(&foo)

	if flag2.String() != expected2.String() {
		t.Errorf("Error. Expected: %s. Received: %s.", expected2.String(), flag2.String())
	}

	if tmp, ok := flag2.Get().(uuid.UUID); !ok {
		t.Errorf("Error. Expected a UUID. Got something else.")
	} else if tmp != expected2 {
		t.Errorf("Error. Expected: %s. Received: %s", expected2, tmp)
	}

	// Calling Set on a default UUID
	expected1 = uuid.New()
	if err := flag1.Set(expected1.String()); err != nil {
		t.Errorf("Error. Expected set to work. Received: %s", err)
	}

	if flag1.String() != expected1.String() {
		t.Errorf("Error. Expected: %s. Received: %s.", expected1.String(), flag1.String())
	}

	if tmp, ok := flag1.Get().(uuid.UUID); !ok {
		t.Errorf("Error. Expected a UUID. Got something else.")
	} else if tmp != expected1 {
		t.Errorf("Error. Expected: %s. Received: %s", expected1, tmp)
	}

	// Calling Set with an invalid UUID.
	if err := flag2.Set("Obviously not a UUID"); err == nil {
		t.Errorf("Error. Expected set to fail.")
	}

	// Instead of testing set with a pointer, using the flag package to invoke
	// set on both types.
	fs := flag.NewFlagSet("UUIDFlag", flag.ExitOnError)
	fs.Var(flag1, "flag1", "a UUID.")
	fs.Var(flag2, "flag2", "a UUID.")
	expected1, expected2 = uuid.New(), uuid.New()

	if foo == expected2 {
		t.Errorf("Error. foo and expected2 are already equal. ??")
	}

	args := []string{"-flag1", expected1.String(), "-flag2", expected2.String()}
	fs.Parse(args)

	if tmp, ok := flag1.Get().(uuid.UUID); !ok {
		t.Errorf("Error. Expected a UUID. Got something else.")
	} else if tmp != expected1 {
		t.Errorf("Error. Expected: %s. Received: %s", expected1, tmp)
	}

	if tmp, ok := flag2.Get().(uuid.UUID); !ok {
		t.Errorf("Error. Expected a UUID. Got something else.")
	} else if tmp != expected2 {
		t.Errorf("Error. Expected: %s. Received: %s", expected2, tmp)
	}

	if foo != expected2 {
		t.Errorf("Error. flag2 didn't update the pointer as expected.")
	}
}

func TestStringSliceFlag(t *testing.T) {
	result1 := make([]string, 0)
	expected1 := []string{"piece1", "piece2"}
	result2 := make([]string, 0)
	expected2 := []string{"piece3", "piece4", "piece5"}
	flag3 := NewSliceFlag(nil, "")
	expected3 := []string{"piece6", "piece7"}
	result4 := make([]string, 0)
	expected4 := []string{"piece8"}
	flag5 := NewSliceFlag(nil, "")
	expected5 := []string{}
	flag6 := NewSliceFlag(&StringSliceFlagTarget{}, "")
	expected6 := []string{"piece9", "piece10"}

	fs := flag.NewFlagSet("StringSlice", flag.ExitOnError)
	fs.Var(NewSliceFlag(&StringSliceFlagTarget{&result1}, ";"),
		"flag1", "comma seperated list of string.")
	fs.Var(NewSliceFlag(&StringSliceFlagTarget{&result2}, ""),
		"flag2", "comma seperated list of string.")
	fs.Var(flag3,
		"flag3", "comma seperated list of string.")
	fs.Var(NewSliceFlag(&StringSliceFlagTarget{&result4}, ""),
		"flag4", "comma seperated list of string.")
	fs.Var(flag6,
		"flag6", "comma seperated list of string.")

	args := []string{"-flag1", "piece1;piece2", "-flag2", "piece3,piece4,piece5",
		"-flag3", "piece6,piece7", "-flag4", "piece8",
		"-flag6", "piece9,piece10"}
	fs.Parse(args)

	if fmt.Sprintf("%v", result1) != fmt.Sprintf("%v", expected1) {
		t.Errorf("Error. Expected %v. Received %v", result1, expected1)
	}

	if fmt.Sprintf("%v", result2) != fmt.Sprintf("%v", expected2) {
		t.Errorf("Error. Expected %v. Received %v", result2, expected2)
	}

	if fmt.Sprintf("%v", flag3.Get()) != fmt.Sprintf("%v", expected3) {
		t.Errorf("Error. Expected %v. Received %v", flag3.Get(), expected3)
	}

	if fmt.Sprintf("%v", result4) != fmt.Sprintf("%v", expected4) {
		t.Errorf("Error. Expected %v. Received %v", result4, expected4)
	}

	if fmt.Sprintf("%v", flag5.Get()) != fmt.Sprintf("%v", expected5) {
		t.Errorf("Error. Expected %v. Received %v", flag5.Get(), expected5)
	}

	if fmt.Sprintf("%v", flag6.Get()) != fmt.Sprintf("%v", expected6) {
		t.Errorf("Error. Expected %v. Received %v", flag6.Get(), expected6)
	}
}

func TestURLSliceFlag(t *testing.T) {
	result1 := make([]*url.URL, 0)
	expected1 := []string{"http://www.piece1.com/", "http://www.piece2.com/"}
	result2 := make([]*url.URL, 0)
	expected2 := []string{"file:///piece3.com", "/piece4", "https://secure.piece5.com/test?foo=bar"}
	result3 := make([]*url.URL, 0)
	expected3 := []string{"http://piece6.net"}
	flag4 := NewSliceFlag(&URLSliceFlagTarget{}, "")
	expected4 := []string{"http://piece7.org", "http://piece8.fr"}
	expected4_2 := "http://piece7.org,http://piece8.fr"

	fs := flag.NewFlagSet("URLSlice", flag.ExitOnError)
	fs.Var(NewSliceFlag(&URLSliceFlagTarget{&result1}, ";"),
		"flag1", "comma seperated list of URLs.")
	fs.Var(NewSliceFlag(&URLSliceFlagTarget{&result2}, ""),
		"flag2", "comma seperated list of URLs.")
	fs.Var(NewSliceFlag(&URLSliceFlagTarget{&result3}, ""),
		"flag3", "comma seperated list of URLs.")
	fs.Var(flag4,
		"flag4", "comma seperated list of URLs.")

	args := []string{"-flag1", "http://www.piece1.com/;http://www.piece2.com/",
		"-flag2", "file:///piece3.com,/piece4,https://secure.piece5.com/test?foo=bar",
		"-flag3", "http://piece6.net",
		"-flag4", "http://piece7.org,http://piece8.fr",
	}
	fs.Parse(args)

	if fmt.Sprintf("%v", result1) != fmt.Sprintf("%v", expected1) {
		t.Errorf("Error. Received %v. Expected %v", result1, expected1)
	}

	if fmt.Sprintf("%v", result2) != fmt.Sprintf("%v", expected2) {
		t.Errorf("Error. Received %v. Expected %v", result2, expected2)
	}

	if fmt.Sprintf("%v", result3) != fmt.Sprintf("%v", expected3) {
		t.Errorf("Error. Received %v. Expected %v", result3, expected3)
	}

	if fmt.Sprintf("%v", flag4.Get()) != fmt.Sprintf("%v", expected4) {
		t.Errorf("Error. Received %v. Expected %v", flag4.Get(), expected4)
	}

	if fmt.Sprintf("%v", flag4.String()) != fmt.Sprintf("%v", expected4_2) {
		t.Errorf("Error. Received %v. Expected %v", flag4.String(), expected4_2)
	}
}

func TestUUIDSliceFlag(t *testing.T) {
	result1 := make([]uuid.UUID, 0)
	expected1 := []string{"63a36905-a4ea-42f4-8133-91951057c10d",
		"bc938938-be7e-4ecc-acb5-b111ef6275f7"}
	result2 := make([]uuid.UUID, 0)
	expected2 := []string{"e8085a3d-5ed5-4789-b048-566dd249a431",
		"732c6e88-4046-4426-88c5-49b77cfcdf71",
		"a420d00e-40cc-40eb-b4ac-51aeda15b890"}
	result3 := make([]uuid.UUID, 0)
	expected3 := []string{"3a884268-0b31-458b-888b-00532c8f0e17"}
	flag4 := NewSliceFlag(&UUIDSliceFlagTarget{}, "")
	expected4 := []string{"b98718d2-d4ef-4e32-8c88-527bcd3ba21c",
		"15f397b2-4209-428a-a207-941285fd85e7"}
	expected4_2 := "b98718d2-d4ef-4e32-8c88-527bcd3ba21c,15f397b2-4209-428a-a207-941285fd85e7"

	fs := flag.NewFlagSet("UUIDSlice", flag.ExitOnError)
	fs.Var(NewSliceFlag(&UUIDSliceFlagTarget{&result1}, ";"),
		"flag1", "comma seperated list of UUIDs.")
	fs.Var(NewSliceFlag(&UUIDSliceFlagTarget{&result2}, ""),
		"flag2", "comma seperated list of UUIDs.")
	fs.Var(NewSliceFlag(&UUIDSliceFlagTarget{&result3}, ""),
		"flag3", "comma seperated list of UUIDs.")
	fs.Var(flag4,
		"flag4", "comma seperated list of UUIDs.")

	args := []string{"-flag1", "63a36905-a4ea-42f4-8133-91951057c10d;bc938938-be7e-4ecc-acb5-b111ef6275f7",
		"-flag2", "e8085a3d-5ed5-4789-b048-566dd249a431,732c6e88-4046-4426-88c5-49b77cfcdf71,a420d00e-40cc-40eb-b4ac-51aeda15b890",
		"-flag3", "3a884268-0b31-458b-888b-00532c8f0e17",
		"-flag4", "b98718d2-d4ef-4e32-8c88-527bcd3ba21c,15f397b2-4209-428a-a207-941285fd85e7",
	}
	fs.Parse(args)

	if fmt.Sprintf("%v", result1) != fmt.Sprintf("%v", expected1) {
		t.Errorf("Error. Received %v. Expected %v", result1, expected1)
	}

	if fmt.Sprintf("%v", result2) != fmt.Sprintf("%v", expected2) {
		t.Errorf("Error. Received %v. Expected %v", result2, expected2)
	}

	if fmt.Sprintf("%v", result3) != fmt.Sprintf("%v", expected3) {
		t.Errorf("Error. Received %v. Expected %v", result3, expected3)
	}

	if fmt.Sprintf("%v", flag4.Get()) != fmt.Sprintf("%v", expected4) {
		t.Errorf("Error. Received %v. Expected %v", flag4.Get(), expected4)
	}

	if fmt.Sprintf("%v", flag4.String()) != fmt.Sprintf("%v", expected4_2) {
		t.Errorf("Error. Received %v. Expected %v", flag4.String(), expected4_2)
	}
}

func TestURLFlag(t *testing.T) {
	expected1 := &url.URL{}
	flag1 := &URLFlag{}

	// a default URLFlag is the same as a zero URL.
	if flag1.String() != expected1.String() {
		t.Errorf("Error. Expected: %s. Received: %s.", expected1.String(), flag1.String())
	}

	if tmp, ok := flag1.Get().(*url.URL); !ok {
		t.Errorf("Error. Expected a URL. Got something else.")
	} else if *tmp != *expected1 {
		t.Errorf("Error. Expected: %s. Received: %s", expected1, tmp)
	}

	// Same tests as above, only with a provided pointer.
	expected2, _ := url.Parse("file:///a")
	foo := *expected2
	flag2 := NewURLFlag(&foo)

	if flag2.String() != expected2.String() {
		t.Errorf("Error. Expected: %s. Received: %s.", expected2.String(), flag2.String())
	}

	if tmp, ok := flag2.Get().(*url.URL); !ok {
		t.Errorf("Error. Expected a URL. Got something else.")
	} else if *tmp != *expected2 {
		t.Errorf("Error. Expected: %s. Received: %s", expected2, tmp)
	}

	// Calling Set on a default URL
	expected1, _ = url.Parse("file:///b")
	if err := flag1.Set(expected1.String()); err != nil {
		t.Errorf("Error. Expected set to work. Received: %s", err)
	}

	if flag1.String() != expected1.String() {
		t.Errorf("Error. Expected: %s. Received: %s.", expected1.String(), flag1.String())
	}

	if tmp, ok := flag1.Get().(*url.URL); !ok {
		t.Errorf("Error. Expected a URL. Got something else.")
	} else if *tmp != *expected1 {
		t.Errorf("Error. Expected: %s. Received: %s", expected1, tmp)
	}

	// Calling Set with an invalid URL.
	if err := flag2.Set("http://bad host/"); err == nil {
		t.Errorf("Error. Expected set to fail.")
	}

	// Instead of testing set with a pointer, using the flag package to invoke
	// set on both types.
	fs := flag.NewFlagSet("URLFlag", flag.ExitOnError)
	fs.Var(flag1, "flag1", "a URL.")
	fs.Var(flag2, "flag2", "a URL.")
	expected1, _ = url.Parse("file:///c")
	expected2, _ = url.Parse("file:///d")

	if foo == *expected2 {
		t.Errorf("Error. foo and expected2 are already equal. ??")
	}

	args := []string{"-flag1", expected1.String(), "-flag2", expected2.String()}
	fs.Parse(args)

	if tmp, ok := flag1.Get().(*url.URL); !ok {
		t.Errorf("Error. Expected a URL. Got something else.")
	} else if *tmp != *expected1 {
		t.Errorf("Error. Expected: %s. Received: %s", expected1, tmp)
	}

	if tmp, ok := flag2.Get().(*url.URL); !ok {
		t.Errorf("Error. Expected a URL. Got something else.")
	} else if *tmp != *expected2 {
		t.Errorf("Error. Expected: %s. Received: %s", expected2, tmp)
	}

	if foo != *expected2 {
		t.Errorf("Error. flag2 didn't update the pointer as expected.")
	}
}
