package main

import (
	"math/rand"
	"sort"
	"strings"
	"testing"
	"time"
)

func TestRandomNonPalindrome(t *testing.T) {
	seed := time.Now().UnixNano()
	t.Logf("Random seed: %d", seed)

	source := genRandomNonPalindrome(seed)

	for i := 0; i < 1000; i++ {
		input := <-source
		if IsPalindrome(sort.StringSlice(strings.Split(input, ""))) {
			t.Errorf("cacl [%s] as palindrome, expect not", input)
		}
	}
}

func genRandomNonPalindrome(seed int64) <-chan string {
	const maxLength = 16
	alphabet := genValidRunes()

	ch := make(chan string, 128)
	random := rand.New(rand.NewSource(seed))
	go func() {
		for {
			var str []rune

			// fulfill with randomRune
			length := random.Intn(maxLength-1) + 2 // [2,maxLength] smaller length cause palindrome
			for i := 0; i < length; i++ {
				randomRune := alphabet[random.Intn(len(alphabet))]
				str = append(str, randomRune)
			}

			// force modify to break potential palindrome on longer string
			index := random.Intn(length / 2) // before med
			for str[index] == str[length-index-1] {
				// If dice always return a same number, would block.
				str[index] = alphabet[random.Intn(len(alphabet))]
			}
			ch <- string(str)
		}
	}()
	return ch
}

func genValidRunes() []rune {
	alphabet := []rune{'`', '.', ',', '，', '。', '/', '?', '!', '-'}
	for i := 0; i <= 9; i++ {
		alphabet = append(alphabet, rune('0'+i))
	}
	for i := 0; i < 26; i++ {
		alphabet = append(alphabet, rune('a'+i))
		alphabet = append(alphabet, rune('A'+i))
	}
	return alphabet
}
