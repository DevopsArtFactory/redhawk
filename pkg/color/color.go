/*
Copyright 2020 The redhawk Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package color

import (
	"fmt"
	"io"
	"strings"

	"github.com/fatih/color"
)

type Color struct {
	color *color.Color
}

var (
	Red    = Color{color: color.New(color.FgRed)}
	Blue   = Color{color: color.New(color.FgBlue)}
	Green  = Color{color: color.New(color.FgGreen)}
	Yellow = Color{color: color.New(color.FgYellow)}
	Cyan   = Color{color: color.New(color.FgCyan)}
)

// Fprintln outputs the result to out, followed by a newline.
func (c Color) Fprintln(out io.Writer, a ...interface{}) {
	if c.color == nil {
		fmt.Fprintln(out, a...)
		return
	}

	fmt.Fprintln(out, c.color.Sprint(strings.TrimSuffix(fmt.Sprintln(a...), "\n")))
}

// Fprintf outputs the result to out.
func (c Color) Fprintf(out io.Writer, format string, a ...interface{}) {
	if c.color == nil {
		fmt.Fprintf(out, format+"\n", a...)
		return
	}

	fmt.Fprint(out, c.color.Sprintf(format+"\n", a...))
}

// DecorateAttr decorate strings with a color or an emoji, respecting the user
// preference if no colour needed.
func DecorateAttr(attrString, message string) string {
	if color.NoColor {
		return message
	}

	switch attrString {
	case "bullet":
		return fmt.Sprintf("∙ %s", message)
	case "check":
		return "✔ "
	case "capacity":
		return "📦 "
	case "tags":
		return "⚓ "
	case "instance_statistics":
		return "🖥 "
	case "security groups":
		return "🚥 "
	case "message":
		return "💌 "
	}

	attr := color.Reset
	switch attrString {
	case "underline":
		attr = color.Underline
	case "underline bold":
		return color.New(color.Underline).Add(color.Bold).Sprintf(message)
	case "bold":
		attr = color.Bold
	case "yellow":
		attr = color.FgHiYellow
	case "green":
		attr = color.FgHiGreen
	case "red":
		attr = color.FgHiRed
	case "blue":
		attr = color.FgHiBlue
	case "magenta":
		attr = color.FgHiMagenta
	case "cyan":
		attr = color.FgHiCyan
	case "black":
		attr = color.FgHiBlack
	case "white":
		attr = color.FgHiWhite
	}

	return color.New(attr).Sprintf(message)
}
