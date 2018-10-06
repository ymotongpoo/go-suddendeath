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
	HeaderChars = []string{"人"}
	FooterChars = []string{"Ｙ"}
	LeftChar    = "＞"
	RightChar   = "＜"
)

// cycle is a helper function to repeatedly yeild strings
// specified in item variadic parameters.
func cycle(item ...string) <-chan string {
	cycleCh := make(chan string)
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
func repeatPart(runes []string, width int) string {
	if len(runes) == 1 {
		w := MessageWidth(runes[0])
		return strings.Repeat(runes[0], width/w+width%1)
	}
	iter := cycle(runes...)
	part := ""
	for i, w := 0, 0; i < width; i += w {
		r := <-iter
		w = MessageWidth(r)
		part += r
	}
	return part
}

// MessageWidth is the display width of assined string msg.
// It is based on the information from Unicode. So-called
// Multibyte characters' width is counted as 2. Other characters'
// are 1. See https://godoc.org/golang.org/x/text/width for the details.
func MessageWidth(msg string) int {
	w := 0
	for _, c := range msg {
		p := width.LookupRune(c)
		switch p.Kind() {
		case width.EastAsianWide, width.EastAsianFullwidth:
			w += 2
		case width.EastAsianNarrow, width.EastAsianHalfwidth:
			w += 1
		default:
			w += 1
		}
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
	header := "＿" + repeatPart(HeaderChars, maxWidth+4) + "＿"
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
	footer := "￣" + repeatPart(FooterChars, maxWidth+4) + "￣"
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
// Body adds side ornamengs
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
