// Copyright (c) 2016 Daniel Oaks <daniel@danieloaks.net>
// released under the MIT license

package isupport

import "fmt"
import "sort"

const (
	maxLastArgLength = 400
)

// List holds a list of ISUPPORT tokens
type List struct {
	Tokens      map[string]*string
	CachedReply [][]string
}

// NewList returns a new List
func NewList() *List {
	var il List
	il.Tokens = make(map[string]*string)
	il.CachedReply = make([][]string, 0)
	return &il
}

// Add adds an RPL_ISUPPORT token to our internal list
func (il *List) Add(name string, value string) {
	il.Tokens[name] = &value
}

// AddNoValue adds an RPL_ISUPPORT token that does not have a value
func (il *List) AddNoValue(name string) {
	il.Tokens[name] = nil
}

// getTokenString gets the appropriate string for a token+value.
func getTokenString(name string, value *string) string {
	if value == nil {
		return name
	}
	return fmt.Sprintf("%s=%s", name, *value)
}

// GetDifference returns the difference between two token lists.
func (il *List) GetDifference(newil *List) [][]string {
	var outTokens sort.StringSlice

	// append removed tokens
	for name := range il.Tokens {
		_, exists := newil.Tokens[name]
		if exists {
			continue
		}

		token := fmt.Sprintf("-%s", name)

		outTokens = append(outTokens, token)
	}

	// append added tokens
	for name, value := range newil.Tokens {
		newval, exists := il.Tokens[name]
		if exists && ((value == nil && newval == nil) || (value != nil && newval != nil && *value == *newval)) {
			continue
		}

		token := getTokenString(name, value)

		outTokens = append(outTokens, token)
	}

	sort.Sort(outTokens)

	// create output list
	replies := make([][]string, 0)
	var length int     // Length of the current cache
	var cache []string // Token list cache

	for _, token := range outTokens {
		if len(token)+length <= maxLastArgLength {
			// account for the space separating tokens
			if len(cache) > 0 {
				length++
			}
			cache = append(cache, token)
			length += len(token)
		}

		if len(cache) == 13 || len(token)+length >= maxLastArgLength {
			replies = append(replies, cache)
			cache = make([]string, 0)
			length = 0
		}
	}

	if len(cache) > 0 {
		replies = append(replies, cache)
	}

	return replies
}

// RegenerateCachedReply regenerates the cached RPL_ISUPPORT reply
func (il *List) RegenerateCachedReply() {
	il.CachedReply = make([][]string, 0)
	var length int     // Length of the current cache
	var cache []string // Token list cache

	// make sure we get a sorted list of tokens, needed for tests and looks nice
	var tokens sort.StringSlice
	for name := range il.Tokens {
		tokens = append(tokens, name)
	}
	sort.Sort(tokens)

	for _, name := range tokens {
		token := getTokenString(name, il.Tokens[name])

		if len(token)+length <= maxLastArgLength {
			// account for the space separating tokens
			if len(cache) > 0 {
				length++
			}
			cache = append(cache, token)
			length += len(token)
		}

		if len(cache) == 13 || len(token)+length >= maxLastArgLength {
			il.CachedReply = append(il.CachedReply, cache)
			cache = make([]string, 0)
			length = 0
		}
	}

	if len(cache) > 0 {
		il.CachedReply = append(il.CachedReply, cache)
	}
}
