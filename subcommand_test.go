package glarg

import (
	//"flag"
	//"fmt"
	//"net/url"
	"context"
	"testing"
	//"github.com/google/uuid"
)

func TestInvoke(t *testing.T) {
	EmptyCommand := SubcommandNoOp{Name: "empty"}
	rc := Invoke(context.Background(), &EmptyCommand, []string{"cmd"})
	if rc != 0 {
		t.Errorf("Error. Expected: 0. Received: %d.", rc)
	}

	ErrorCommand := SubcommandNoOp{Name: "error", ExecuteInt: 1}
	rc = Invoke(context.Background(), &ErrorCommand, []string{"cmd"})
	if rc != 1 {
		t.Errorf("Error. Expected: 1. Received: %d.", rc)
	}

	// Create some test cases for the subcommands.
	RootCommand := Subcommands{
		Name:     "root",
		Children: []Subcommand{&EmptyCommand, &ErrorCommand},
	}
	rc = Invoke(context.Background(), &RootCommand, []string{"cmd"})
	if rc != 1 {
		t.Errorf("Error. Expected: 1. Received: %d.", rc)
	}
	rc = Invoke(context.Background(), &RootCommand, []string{"cmd", "empty"})
	if rc != 0 {
		t.Errorf("Error. Expected: 0. Received: %d.", rc)
	}

	// make a sub command of a subcommand.
	RealRootCommand := Subcommands{
		Name:     "realroot",
		Children: []Subcommand{&RootCommand, &EmptyCommand},
	}
	rc = Invoke(context.Background(), &RealRootCommand, []string{"cmd"})
	if rc != 1 {
		t.Errorf("Error. Expected: 1. Received: %d.", rc)
	}
	rc = Invoke(context.Background(), &RealRootCommand, []string{"cmd", "root"})
	if rc != 1 {
		t.Errorf("Error. Expected: 1. Received: %d.", rc)
	}
	rc = Invoke(context.Background(), &RealRootCommand, []string{"cmd", "empty"})
	if rc != 0 {
		t.Errorf("Error. Expected: 0. Received: %d.", rc)
	}
	rc = Invoke(context.Background(), &RealRootCommand, []string{"cmd", "root", "unknown"})
	if rc != 1 {
		t.Errorf("Error. Expected: 0. Received: %d.", rc)
	}
	rc = Invoke(context.Background(), &RealRootCommand, []string{"cmd", "root", "error"})
	if rc != 1 {
		t.Errorf("Error. Expected: 0. Received: %d.", rc)
	}
	rc = Invoke(context.Background(), &RealRootCommand, []string{"cmd", "root", "empty"})
	if rc != 0 {
		t.Errorf("Error. Expected: 0. Received: %d.", rc)
	}
}
