package shell

import (
	"fmt"
	"os"

	s "github.com/PoC-Consortium/Aspera/pkg/store"

	"github.com/c-bata/go-prompt"
)

func completer(d prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{
		{Text: "rebuild", Description: "rebuilds the database from the raw storage"},
		{Text: "shutdown", Description: "clean wallet shutdown"},
	}
	return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
}

func Prompt(store *s.Store) {
	go func() {
	PROMPT:
		t := prompt.Input("> ", completer)
		switch t {
		case "rebuild":
			fmt.Println("rebulding ... ")
			store.ChainStore.Rebuild()
			fmt.Println("done.")
		case "shutdown":
			os.Exit(0)
		default:
			fmt.Println("unknown command " + t)
		}
		goto PROMPT
	}()
}
