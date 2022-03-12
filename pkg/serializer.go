package pkg

import (
	"encoding/json"
	"fmt"
)

func SerializeHints(hints []*Hint, resultFormat string) {
	switch resultFormat {
	case "json":
		hintsBytes, err := json.Marshal(hints)
		if err != nil {
			panic(err)
		}
		fmt.Println(string(hintsBytes))
	default:
		if len(hints) > 0 {
			for _, hint := range hints {
				fmt.Printf("Pos: %d, End: %d, Category: %s, Message: %s, Suggestion: %s \n", hint.Pos, hint.End, hint.Category, hint.Message, hint.Suggestion)
			}
		}
	}
}
