package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

type Args struct {
	root    []string
	series  []string
	movies  []string
	options AdditionalOptions
}
type AdditionalOptions struct {
	keepEpNums    Option[bool]
	startingEpNum Option[int]
	hasSeason0    Option[bool]
	namingScheme  Option[string]
}

func newArgs() Args {
	return Args{
		root:   make([]string, 0),
		series: make([]string, 0),
		movies: make([]string, 0),
		options: AdditionalOptions{
			hasSeason0:    none[bool](),
			keepEpNums:    none[bool](),
			startingEpNum: none[int](),
			namingScheme:  none[string](),
		},
	}
}

func ParseArgs(args []string) (Args, error) {
	if len(args) < 1 {
		return Args{}, fmt.Errorf("not enough arguments")
	}

	directoryArgs := map[string]bool{
		"root":   true,
		"series": true,
		"movies": true,
	}
	assignedVar := map[string]bool{
		"--options":         false,
		"--keep-ep-nums":    false,
		"--starting-ep-num": false,
		"--has-season-0":    false,
		"--naming-scheme":   false,
	}

	parsedArgs := newArgs()
	skipValue := 0
	for i, arg := range args {
		// valid values will be skipped
		if skipValue != 0 && i <= skipValue {
			continue

		} else if arg[0] != '-' && !directoryArgs[arg] {
			// catch invalid values acting as flags
			return Args{}, fmt.Errorf("invalid flag: '%s'", arg)

		} else if arg == "--help" || arg == "-h" {
			if len(args) <= i+1 {
				Help("")
			} else if len(args) > i+1 {
				Help(args[i+1])
			}
			return Args{}, fmt.Errorf("safe exit")

		} else if arg == "--version" || arg == "-v" {
			Version(version)
			return Args{}, fmt.Errorf("safe exit")

		} else if directoryArgs[arg] {
			// no value after flag / flag after flag
			if len(args) <= i+1 || (len(args) > i+1 && args[i+1][0] == '-') {
				return Args{}, fmt.Errorf("missing dir path value for flag '%s'", arg)
			}

			dir, err := filepath.Abs(args[i+1])
			if err != nil {
				return Args{}, err
			}
			_, err = os.Stat(args[i+1])
			if os.IsNotExist(err) {
				return Args{}, fmt.Errorf("'%s' is not a valid directory", args[i+1])
			}

			switch arg {
			case "root":
				parsedArgs.root = append(parsedArgs.root, dir)
			case "series":
				parsedArgs.series = append(parsedArgs.series, dir)
			case "movies":
				parsedArgs.movies = append(parsedArgs.movies, dir)
			}
			skipValue = i + 1

		} else if arg == "--has-season-0" || arg == "-s0" {
			if parsedArgs.options.hasSeason0.IsSome() {
				return Args{}, fmt.Errorf("only one --has-season-0 flag is allowed")
			}
			// use default value
			if len(args) <= i+1 || (len(args) > i+1 && args[i+1][0] == '-') {
				parsedArgs.options.hasSeason0 = some[bool](false)

			} else if args[i+1] != "yes" && args[i+1] != "var" && args[i+1] != "no" && args[i+1] != "default" {
				return Args{}, fmt.Errorf("invalid value '%s' for flag '%s'. Must be 'yes', 'no', 'var, or 'default", args[i+1], arg)

			} else {
				switch args[i+1] {
				case "yes":
					parsedArgs.options.hasSeason0 = some[bool](true)
				case "no", "default":
					parsedArgs.options.hasSeason0 = some[bool](false)
				case "var":
					parsedArgs.options.hasSeason0 = none[bool]()
					assignedVar["--has-season-0"] = true
				}
				skipValue = i + 1
			}

		} else if arg == "--keep-ep-nums" || arg == "-ken" {
			if parsedArgs.options.keepEpNums.IsSome() {
				return Args{}, fmt.Errorf("only one --keep-ep-nums flag is allowed")
			}

			// use default value
			if len(args) <= i+1 || (len(args) > i+1 && args[i+1][0] == '-') {
				parsedArgs.options.keepEpNums = some[bool](false)

			} else if args[i+1] != "yes" && args[i+1] != "var" && args[i+1] != "no" && args[i+1] != "default" {
				return Args{}, fmt.Errorf("invalid value '%s' for --keep-ep-nums. Must be 'yes', 'no', 'var', or 'default", args[i+1])

			} else {
				switch args[i+1] {
				case "yes":
					parsedArgs.options.keepEpNums = some[bool](true)
				case "no", "default":
					parsedArgs.options.keepEpNums = some[bool](false)
				case "var":
					parsedArgs.options.keepEpNums = none[bool]()
					assignedVar["--keep-ep-nums"] = true
				}
				skipValue = i + 1
			}

		} else if arg == "--starting-ep-num" || arg == "-sen" {
			if parsedArgs.options.startingEpNum.IsSome() {
				return Args{}, fmt.Errorf("only one --starting-ep-num flag is allowed")
			}

			// use default value
			if len(args) <= i+1 || (len(args) > i+1 && args[i+1][0] == '-') {
				parsedArgs.options.startingEpNum = some[int](1)
				assignedVar["--starting-ep-num"] = true

			} else if value, err := strconv.Atoi(args[i+1]); err != nil && value < 1 && args[i+1] != "var" && args[i+1] != "default" {
				return Args{}, fmt.Errorf("invalid value '%s' for --starting-ep-num. Must be a valid positive int or 'var", args[i+1])

			} else {
				switch args[i+1] {
				case "var":
					parsedArgs.options.startingEpNum = none[int]()
				case "default":
					parsedArgs.options.startingEpNum = some[int](1)
				default:
					parsedArgs.options.startingEpNum = some[int](value)
				}
				assignedVar["--starting-ep-num"] = true
				skipValue = i + 1
			}

		} else if arg == "--options" || arg == "-o" {
			if assignedVar["--options"] {
				return Args{}, fmt.Errorf("only one --options flag is allowed")
			}

			// all options are assigned `var` if no value is provided
			if len(args) <= i+1 || (len(args) > i+1 && (args[i+1][0] == '-' || args[i+1] == "var")) {
				assignedVar["--options"] = true
				parsedArgs.options.keepEpNums = none[bool]()
				parsedArgs.options.startingEpNum = none[int]()
				parsedArgs.options.hasSeason0 = none[bool]()
				parsedArgs.options.namingScheme = none[string]()

			} else if args[i+1] != "default" && args[i+1] != "var" {
				return Args{}, fmt.Errorf("invalid value '%s' for --options. Must be 'default' or 'var", args[i+1])

			} else if args[i+1] == "default" {
				// use default values
				assignedVar["--options"] = true
				parsedArgs.options.keepEpNums = some[bool](false)
				parsedArgs.options.startingEpNum = some[int](1)
				parsedArgs.options.hasSeason0 = some[bool](false)
				parsedArgs.options.namingScheme = some[string]("default")
			}

		} else if arg == "--naming-scheme" || arg == "-ns" {
			if parsedArgs.options.namingScheme.IsSome() {
				return Args{}, fmt.Errorf("only one --naming-scheme flag is allowed")
			}
			if len(args) <= i+1 || (len(args) > i+1 && args[i+1][0] == '-') {
				return Args{}, fmt.Errorf("missing value for --naming-scheme")
			}

			err := ValidateNamingScheme(args[i+1])
			if err != nil && args[i+1] != "default" && args[i+1] != "var" {
				return Args{}, fmt.Errorf("invalid value '%s' for --naming-scheme. Must be 'default', 'var', or a naming scheme enclosed in double quotes", args[i+1])

			} else {
				switch args[i+1] {
				case "default":
					parsedArgs.options.namingScheme = some[string]("default")
				case "var":
					parsedArgs.options.namingScheme = none[string]()
				default:
					namingScheme := strings.Trim(args[i+1], `"`)
					parsedArgs.options.namingScheme = some[string](namingScheme)
				}
				skipValue = i + 1
			}

		} else {
			return Args{}, fmt.Errorf("unknown flag: %s", arg)
		}
	}

	err := ValidateRoots(parsedArgs.root, parsedArgs.series, parsedArgs.movies)
	if err != nil {
		return Args{}, err
	}

	if !assignedVar["--options"] {
		// use default values for optional flags if not assigned var
		if parsedArgs.options.hasSeason0.IsNone() && !assignedVar["--has-season-0"] {
			parsedArgs.options.hasSeason0 = some[bool](false)
		}
		if parsedArgs.options.keepEpNums.IsNone() && !assignedVar["--keep-ep-nums"] {
			parsedArgs.options.keepEpNums = some[bool](false)
		}
		if parsedArgs.options.startingEpNum.IsNone() && !assignedVar["--starting-ep-num"] {
			parsedArgs.options.startingEpNum = some[int](1)
		}
		if parsedArgs.options.namingScheme.IsNone() && !assignedVar["--naming-scheme"] {
			parsedArgs.options.namingScheme = some[string]("default")
		}
	}
	return parsedArgs, nil
}

func ValidateNamingScheme(s string) error {
	if s[0] != '"' || s[len(s)-1] != '"' {
		return fmt.Errorf("naming scheme must be enclosed in double quotes: %s", s)
	}

	tokens, err := TokenizeNamingScheme(s)
	if err != nil {
		return err
	}

	validAPI := regexp.MustCompile(`^season_num$|^episode_num$|^self$`)
	validParentAPI := regexp.MustCompile(`^parent(-parent)*$|^p(-\d+)?$`)
	validRange := regexp.MustCompile(`^\d+(\s*,\s*\d+)?$`)

	for _, token := range tokens {
		var api, val string

		// for api with values
		if strings.Contains(token, ":") {
			res := strings.SplitN(token, ":", 2)
			api, val = strings.TrimSpace(res[0]), strings.TrimSpace(res[1])
		} else {
			api, val = strings.TrimSpace(token), "none"
		}

		if !validAPI.MatchString(api) && !validParentAPI.MatchString(api) {
			return fmt.Errorf("invalid api: %s", api)
		}

		if api == "season_num" || api == "episode_num" {
			if val == "none" {
				continue
			}
			if val == "" {
				return fmt.Errorf("season_num's value cannot be empty and must be a positive integer")
			}

			int_val, err := strconv.Atoi(val)
			if err != nil {
				return fmt.Errorf("%s's value must be a positive integer. '%s' is not a valid integer or 0", api, val)
			}
			if int_val < 0 {
				return fmt.Errorf("%s's value must be a positive integer. '%s' is not a positive integer or 0", api, val)
			}

		} else if api == "self" || validParentAPI.MatchString(api) {
			if val == "none" {
				continue
			}
			if val == "" {
				return fmt.Errorf("%s's value cannot be empty", api)
			}
			if val[0] == '\'' {
				singleClosedString := regexp.MustCompile(`^'[^']*'$`)
				if !singleClosedString.MatchString(val) {
					return fmt.Errorf("%s is either an unclosed string or contains more than 2 single quotes", val)
				}

				valTrimmed := strings.Trim(val, "'")
				justWhitespace := regexp.MustCompile(`^\s*$`)
				if justWhitespace.MatchString(valTrimmed) {
					return fmt.Errorf("%s is an empty string (just whitespace/s)", val)
				}

				if _, err := regexp.Compile(valTrimmed); err != nil {
					return fmt.Errorf("invalid regex: %s", val)

				} else {
					parts := SplitRegexByPipe(valTrimmed)
					for _, part := range parts {
						if !HasOnlyOneMatchGroup(part) {
							return fmt.Errorf("regex should have only one match group per part (parts are separated by outermost pipes |): %s", val)
						}
					}
					// valid regex
					continue
				}

			} else if !validRange.MatchString(val) {
				return fmt.Errorf("%s's value must be in the format <start>,<end> where <start> and <end> are positive integers and 0. '%s' is not a valid range", api, val)
			}

			// valid range
			var res []string
			if strings.Contains(val, ",") {
				res = strings.SplitN(val, ",", 2)
			} else {
				res = []string{val, val}
			}
			begin, end := strings.TrimSpace(res[0]), strings.TrimSpace(res[1])
			if begin > end {
				return fmt.Errorf("%s is an invalid range. begin (%s) must be less than or equal to end (%s)", val, begin, end)
			}
		}
	}

	return nil
}

func TokenizeNamingScheme(s string) ([]string, error) {
	isToken := false
	builder := strings.Builder{}
	namingScheme := make([]string, 0)

	for i, c := range s {
		if isToken && i+1 == len(s) && c != '>' {
			return nil, fmt.Errorf("reached end of string but still in an unclosed api: '<%s%s'", builder.String(), string(c))
		}

		// start of token
		if c == '<' && !isToken {
			isToken = true
			builder.Reset()
			continue

			// end of token
		} else if c == '>' && isToken {
			isToken = false
			namingScheme = append(namingScheme, builder.String())
			builder.Reset()
			continue
		}

		if isToken {
			_, err := builder.WriteRune(c)
			if err != nil {
				return nil, err
			}
		}
	}

	return namingScheme, nil
}

func ValidateRoots(root []string, series []string, movies []string) error {
	// must at least have one of any
	if len(root) == 0 && len(series) == 0 && len(movies) == 0 {
		return fmt.Errorf("must specify at least one root directory")
	}

	// check if exists
	for _, r := range root {
		if _, err := os.Stat(r); err != nil {
			return fmt.Errorf("root directory %s does not exist", r)
		}
	}
	for _, r := range series {
		if _, err := os.Stat(r); err != nil {
			return fmt.Errorf("series directory %s does not exist", r)
		}
	}
	for _, r := range movies {
		if _, err := os.Stat(r); err != nil {
			return fmt.Errorf("movies directory %s does not exist", r)
		}
	}

	// check if any of the series and movies directories are subdirectories of a root directory OR vice versa
	// and in turn checking if any of the series and movies directories are duplicates of root directories
	for _, r := range root {
		for _, s := range series {
			if strings.EqualFold(filepath.Dir(s), r) {
				return fmt.Errorf("series directory %s is a subdirectory of root directory %s", s, r)
			} else if strings.EqualFold(filepath.Dir(r), s) {
				return fmt.Errorf("root directory %s is a subdirectory of series directory %s", r, s)
			} else if strings.EqualFold(s, r) {
				return fmt.Errorf("series directory %s is a duplicate of root directory %s", s, r)
			}
		}
		for _, m := range movies {
			if strings.EqualFold(filepath.Dir(m), r) {
				return fmt.Errorf("movies directory %s is a subdirectory of root directory %s", m, r)
			} else if strings.EqualFold(filepath.Dir(r), m) {
				return fmt.Errorf("root directory %s is a subdirectory of movies directory %s", r, m)
			} else if strings.EqualFold(m, r) {
				return fmt.Errorf("movies directory %s is a duplicate of root directory %s", m, r)
			}
		}
	}

	// check if any of the series and movies directories are subdirectories of each other
	// and in turn checking if any of the series and movies directories are duplicates of each other
	for _, s := range series {
		for _, m := range movies {
			if strings.EqualFold(filepath.Dir(m), s) {
				return fmt.Errorf("series directory %s is a subdirectory of movies directory %s", s, m)

			} else if strings.EqualFold(filepath.Dir(s), m) {
				return fmt.Errorf("movies directory %s is a subdirectory of series directory %s", m, s)

			} else if strings.EqualFold(s, m) {
				return fmt.Errorf("series directory %s is a duplicate of movies directory %s", s, m)
			}
		}
	}

	// check if there are any duplicates in each individually
	for i, r := range root {
		for _, r1 := range root[i+1:] {
			if strings.EqualFold(r, r1) {
				return fmt.Errorf("there are multiple root directories that share the same path: %s", r)
			}
		}
	}
	for i, s := range series {
		for _, s1 := range series[i+1:] {
			if strings.EqualFold(s, s1) {
				return fmt.Errorf("there are multiple series directories that share the same path: %s", s)
			}
		}
	}
	for i, m := range movies {
		for _, m1 := range movies[i+1:] {
			if strings.EqualFold(m, m1) {
				return fmt.Errorf("there are multiple movies directories that share the same path: %s", m)
			}
		}
	}

	// all done
	return nil
}
