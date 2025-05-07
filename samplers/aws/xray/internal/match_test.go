// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0
// Modifications Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.

package internal

import (
	"bytes"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// 2025-05-07: Begin Amazon modification.
// assert wildcard match is positive.
func TestWildCardMatchPositive(t *testing.T) {
	tests := []struct {
		pattern string
		text    string
	}{
		// wildcard positive test set
		{"*", ""},
		{"foo", "foo"},
		{"foo*bar*?", "foodbaris"},
		{"?o?", "foo"},
		{"*oo", "foo"},
		{"foo*", "foo"},
		{"*o?", "foo"},
		{"*", "boo"},
		{"", ""},
		{"a", "a"},
		{"*a", "a"},
		{"*a", "ba"},
		{"a*", "a"},
		{"a*", "ab"},
		{"a*a", "aa"},
		{"a*a", "aba"},
		{"a*a*", "aaaaaaaaaaaaaaaaaaaaaaa"},
		{
			"a*b*a*b*a*b*a*b*a*",
			"akljd9gsdfbkjhaabajkhbbyiaahkjbjhbuykjakjhabkjhbabjhkaabbabbaaakljdfsjklababkjbsdabab",
		},
		{"a*na*ha", "anananahahanahanaha"},
		{"***a", "a"},
		{"**a**", "a"},
		{"a**b", "ab"},
		{"*?", "a"},
		{"*??", "aa"},
		{"*?", "a"},
		{"*?*a*", "ba"},
		{"?at", "bat"},
		{"?at", "cat"},
		{"?o?se", "horse"},
		{"?o?se", "mouse"},
		{"*s", "horses"},
		{"J*", "Jeep"},
		{"J*", "jeep"},
		{"*/foo", "/bar/foo"},
		{"*", "HelloWorld"},
		{"HelloWorld", "HelloWorld"},
		{"Hello*", "HelloWorld"},
		{"*World", "HelloWorld"},
		{"?ello*", "HelloWorld"},
		{"Hell?W*d", "HelloWorld"},
		{"*.World", "Hello.World"},
		{"*.World", "Bye.World"},
	}

	for _, test := range tests {
		match, err := wildcardMatch(test.pattern, test.text)
		require.NoError(t, err)
		assert.True(t, match, test.text)
	}
}

// assert wildcard match is negative.
func TestWildCardMatchNegative(t *testing.T) {
	tests := []struct {
		pattern string
		text    string
	}{
		// wildcard negative test set
		{"", "whatever"},
		{"foo", "bar"},
		{"f?o", "boo"},
		{"f??", "boo"},
		{"fo*", "boo"},
		{"f?*", "boo"},
		{"abcd", "abc"},
		{"??", "a"},
		{"??", "a"},
		{"*?*a", "a"},
		{"/", "target"},
		{"/", "/target"},
		{"a*na*ha", "anananahahanahana"},
		{"*s", "horse"},
	}

	for _, test := range tests {
		match, err := wildcardMatch(test.pattern, test.text)
		require.NoError(t, err)
		assert.False(t, match)
	}
}

// End of Amazon modification.

func TestLongStrings(t *testing.T) {
	chars := []byte{'a', 'b', 'c', 'd'}
	text := bytes.NewBufferString("a")
	for i := 0; i < 8192; i++ {
		_, _ = text.WriteString(string(chars[rand.Intn(len(chars))]))
	}
	_, _ = text.WriteString("b")

	match, err := wildcardMatch("a*b", text.String())
	require.NoError(t, err)
	assert.True(t, match)
}
