package main

import (
	"math/rand"
	"testing"
	"time"
	"unicode"
)

func TestRandomNonPalindrome(t *testing.T) {
	seed := time.Now().UnixNano()
	t.Logf("Random seed: %d", seed)

	source := genRandomPalindrome(seed)

	for i := 0; i < 1000; i++ {
		input := <-source
		if !IsPalindrome(input) {
			t.Errorf("cacl [%s] as nonpalindrome, expect true", input)
		}
	}
}

func genRandomPalindrome(seed int64) <-chan string {
	const maxLength = 16
	matter, ignore := genValidRunes()

	ch := make(chan string)
	random := rand.New(rand.NewSource(seed))
	go func() {
		ch <- "" // use empty string as first output
		for {
			var str []rune

			length := random.Intn(maxLength) + 1 // [1,maxLength] smaller length cause palindrome

			// fulfill [0,med] with randomRune
			for i := 0; i < (length+1)/2; i++ {
				str = append(str, matter[random.Intn(len(matter))])
			}

			// copy to (med,end)
			for i := (length + 1) / 2; i < length; i++ {
				r := str[length-i-1]

				// randomly flop case
				if random.Float64() < 0.1 {
					if unicode.IsLower(r) {
						r = unicode.ToUpper(r)
					} else {
						r = unicode.ToLower(r)
					}
				}

				str = append(str, r)
			}

			// randomly add rune that shall ignore at random place
			for random.Float64() < 0.8 {
				index := random.Intn(len(str))
				str = []rune(string(str[:index]) + string(ignore[random.Intn(len(ignore))]) + string(str[index:]))
			}

			ch <- string(str)
		}
	}()
	return ch
}

func genValidRunes() (matter, ignore []rune) {
	ignore = []rune{'.', ',', '，', '。', '/', '?', '!', '-', ' '}

	for i := 0; i < 26; i++ {
		matter = append(matter, rune('a'+i))
		matter = append(matter, rune('A'+i))
	}

	return
}
