//    Copyright 2018 Yoshi Yamaguchi
//
//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.

package suddendeath

import (
	"fmt"
	"strings"

	"golang.org/x/text/width"
)

var (
	HeaderRunes = []rune{'人'}
	FooterRunes = []rune{'Ｙ'}
	LeftChar    = "＞"
	RightChar   = "＜"
)

// cycle is a helper function to repeatedly yeild strings
// specified in item variadic parameters.
func cycle(item ...rune) <-chan rune {
	cycleCh := make(chan rune)
	go func() {
		for {
			for _, i := range item {
				cycleCh <- i
			}
		}
	}()
	return cycleCh
}

// repeatPart returns repeated part in the header and footers
// to fit the width with specified header and footer chars.
func repeatPart(runes []rune, width int) string {
	if len(runes) == 1 {
		w := RuneWidth(runes[0])
		return strings.Repeat(string(runes[0]), width/w+width%1)
	}
	iter := cycle(runes...)
	part := ""
	for i, w := 0, 0; i < width; i += w {
		r := <-iter
		w = RuneWidth(r)
		part += string(r)
	}
	return part
}

// RuneWidth returns the display width of r.
func RuneWidth(r rune) int {
	p := width.LookupRune(r)
	switch p.Kind() {
	case width.EastAsianWide, width.EastAsianFullwidth:
		return 2
	case width.EastAsianNarrow, width.EastAsianHalfwidth:
		return 1
	default:
		return 1
	}
}

// MessageWidth is the display width of assined string msg.
// msg is expected to be single line; i.e. not including \r and \n.
// It is based on the information from Unicode. So-called
// Multibyte characters' width is counted as 2. Other characters'
// are 1. See https://godoc.org/golang.org/x/text/width for the details.
func MessageWidth(msg string) int {
	w := 0
	for _, r := range msg {
		w += RuneWidth(r)
	}
	return w
}

// Header returns the header string of "突然の死" like message
// based on the display width of msg.
func Header(msg string) string {
	lines := strings.Split(msg, "\n")
	maxWidth := 0
	for _, l := range lines {
		if w := MessageWidth(l); w > maxWidth {
			maxWidth = w
		}
	}
	header := "＿" + repeatPart(HeaderRunes, maxWidth+4) + "＿"
	return header
}

// Footer returns the footer string of "突然の死" like message
// based on the display width of msg.
func Footer(msg string) string {
	lines := strings.Split(msg, "\n")
	maxWidth := 0
	for _, l := range lines {
		if w := MessageWidth(l); w > maxWidth {
			maxWidth = w
		}
	}
	footer := "￣" + repeatPart(FooterRunes, maxWidth+4) + "￣"
	return footer
}

// centerAlign returns string for line with padding in both side in single white space
// to achieve center alignment to fit width space.
func centerAlign(line string, width int) string {
	w := MessageWidth(line)
	if width <= w {
		return line
	}
	gap := width - w
	leftPad := strings.Repeat(" ", gap/2)
	rightPad := strings.Repeat(" ", (gap+1)/2)
	return leftPad + line + rightPad
}

// Body returns the message body part of "突然の死" like message.
// Body adds side ornament characters and spaces in each lines.
func Body(msg string) string {
	lines := strings.Split(msg, "\n")
	maxWidth := 0
	for _, l := range lines {
		if w := MessageWidth(l); w > maxWidth {
			maxWidth = w
		}
	}
	newLines := make([]string, len(lines), len(lines))
	for i, l := range lines {
		newLines[i] = LeftChar + "　" + centerAlign(l, maxWidth) + "　" + RightChar
	}
	return strings.Join(newLines, "\n")
}

// Format returns "突然の死"-nized version of msg.
func Format(msg string) string {
	return strings.Join([]string{Header(msg), Body(msg), Footer(msg)}, "\n")
}

// Print simply prints the string using fmt.Println.
func Print(msg string) {
	fmt.Println(Format(msg))
}
