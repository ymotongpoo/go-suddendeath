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

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/ymotongpoo/go-suddendeath"
)

var (
	slack   = flag.Bool("slack", false, "enable escape for slack")
	message = flag.String("message", "", "message to be formatted in 突然の死")
)

func init() {
	flag.BoolVar(slack, "s", false, "enable escape for slack")
	flag.StringVar(message, "m", "", "message to be formatted in 突然の死")
}

func main() {
	flag.Parse()

	if *message != "" {
		PrintSlackEscaped(*message, *slack)
		return
	}
	bytes, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal("unable to read stdin: %v", err)
	}
	PrintSlackEscaped(string(bytes), *slack)
}

func PrintSlackEscaped(msg string, slack bool) {
	formatted := suddendeath.Format(msg)
	if !slack {
		fmt.Println(formatted)
		return
	}
	lines := strings.Split(formatted, "\n")
	newLines := make([]string, len(lines), len(lines))
	for i, l := range lines {
		newLines[i] = "." + l
	}
	fmt.Println(strings.Join(newLines, "\n"))
}
