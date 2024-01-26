// Package zbaction defines the schema, executor and procedures for Zeabur Action.
//
// Zeabur Action is a workflow language like GitHub Actions, and try to be Go-native.
//
// To use it:
//
//	import (
//	 	"github.com/zeabur/action"
//		_ "github/zeabur/builder/zbaction/procedures" /* builtin procedures */
//	 	/* your procedures plugins */
//	 )
//
//	err := zbaction.RunAction(context.TODO() /* your context */ , zbaction.Action{} /* your action */)
package zbaction
