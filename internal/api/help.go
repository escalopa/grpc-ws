package api

import (
	"fmt"
	"math/rand/v2"
)

var names = []string{
	"Flibber", "Snork", "Wizzle", "Zap", "Gloop", "Dingle", "Bork", "Splat",
	"Zoodle", "Quirk", "Bibble", "Twizzle", "Blurp", "Jibber", "Thunk", "Fizz",
	"Sprocket", "Whizzle", "Plonk", "Zizzle",
}

func GetUsername() string {
	var (
		first  = names[rand.IntN(len(names))]
		second = names[rand.IntN(len(names))]
	)
	return fmt.Sprintf("%s %s", first, second)
}
