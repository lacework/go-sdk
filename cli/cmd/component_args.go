//
// Copyright:: Copyright 2022, Lacework Inc.
// License:: Apache License, Version 2.0
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/pflag"
)

type componentArgParser struct {
	componentArgs []string
	cliArgs       []string
}

// Parsing component args is surprisingly tricky.  What we do here is mirror
// what pflags does itself, but instead of parsing arguments set aside (in componentArgs)
// those arguments that we should pass to the component, and keep track of the
// arguments that we should parse ourselves in cliArgs.

func (p *componentArgParser) parseArgs(globalFlags *pflag.FlagSet, args []string) {
	p.componentArgs = []string{}
	p.cliArgs = []string{}
	for len(args) > 0 {
		s := args[0]
		args = args[1:]
		if s == "help" || s == "--help" || s == "-h" {
			// always pass help down to the component
			p.componentArgs = append(p.componentArgs, s)
			continue
		}

		if len(s) == 0 || s[0] != '-' || len(s) == 1 {
			// not a flag, passthrough
			p.componentArgs = append(p.componentArgs, s)
			continue
		}
		if s[1] == '-' {
			if len(s) == 2 {
				// "--" terminates the flags, but we do want to pass along the --
				// to the compoent
				p.componentArgs = append(p.componentArgs, s)
				p.componentArgs = append(p.componentArgs, args...)
				break
			}
			args = p.parseLongArg(globalFlags, s, args)
		} else {
			args = p.parseShortArg(globalFlags, s, args)
		}
	}
}

func (p *componentArgParser) parseLongArg(flags *pflag.FlagSet, s string, args []string) []string {
	name := s[2:]
	if len(name) == 0 || name[0] == '-' || name[0] == '=' {
		// ---, or --= are not legal cobra flags, but we'll just pass
		// it through to the component and let that deal with it
		p.componentArgs = append(p.componentArgs, s)
		return args
	}
	split := strings.SplitN(name, "=", 2)
	name = split[0]

	flag := flags.Lookup(name)

	if flag == nil {
		p.componentArgs = append(p.componentArgs, s)
		// We're actually a bit stuck here as we don't know if this flag
		// takes an argument or not, so we don't know whether or not to consume
		// the next arg.  What we'll do is peek ahead, and if the next arg does
		// not start with - then we'll take it.
		if len(args) > 0 && len(args[0]) > 0 && args[0][0] == '-' {
			// This component flag does not take an arg
			return args
		}
		if len(args) > 0 {
			p.componentArgs = append(p.componentArgs, args[0])
			return args[1:]
		}
		return args
	}

	if len(split) == 2 {
		// '--flag=arg'
		p.cliArgs = append(p.cliArgs, s)
		return args
	}
	if flag.NoOptDefVal != "" {
		// '--flag' (arg was optional)
		p.cliArgs = append(p.cliArgs, s)
		return args
	}
	if len(args) > 0 {
		// '--flag arg'
		p.cliArgs = append(p.cliArgs, s, args[0])
		return args[1:]
	}
	// '--flag' (arg was required)
	p.cliArgs = append(p.cliArgs, s)
	return args
}

func (p *componentArgParser) parseShortArg(flags *pflag.FlagSet, s string, args []string) []string {
	shorthands := s[1:]

	// shorthands can be a repeated list, e.g. -vvv.
	for len(shorthands) > 0 {
		shorthand := shorthands[0:1]

		flag := flags.ShorthandLookup(shorthand)
		if flag == nil {
			// Not our flag, pass to the component.  Like the long form above we
			// don't know whether to consume an extra arg, so we'll do the same
			// thing: if the next arg does not start with - then pass it along
			p.componentArgs = append(p.componentArgs, fmt.Sprintf("-%s", shorthand))
			if len(shorthands) == 1 && (len(args) > 0 && len(args[0]) > 0 && args[0][0] == '-') {
				p.componentArgs = append(p.componentArgs, args[0])
				return args[2:]
			}
			shorthands = shorthands[1:]
			continue
		}

		if len(shorthands) > 2 && shorthands[1] == '=' {
			// '-f=arg'
			p.cliArgs = append(p.cliArgs, s)
			return args
		} else if flag.NoOptDefVal != "" {
			// '-f' (arg was optional)
			p.cliArgs = append(p.cliArgs, s)
		} else if len(shorthands) > 1 {
			// '-farg'
			p.componentArgs = append(p.componentArgs, s)
			return args
		} else if len(args) > 0 {
			// '-f arg'
			p.cliArgs = append(p.cliArgs, s, args[0])
			return args[1:]
		} else {
			// '-f' (arg was required)
			p.cliArgs = append(p.cliArgs, s)
			return args
		}
	}
	return args
}
