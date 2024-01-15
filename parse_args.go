package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

// Arg represents an argument passed to the program at initial execution.
// It has a name and a value where an Arg can only have one value
//
// An Arg can be a command, a switch, or a flag.
//
//	Commands generally have no dashes
//	Switches and flags generally have a dash (-) prepended to their names
type Arg struct {
	name  string
	value string
}

func IsValidCommand(name string) bool {
	return name == "root" ||
		name == "series" ||
		name == "movies"
}
func IsValidSwitch(name string) bool {
	return name == "--help" ||
		name == "-h" ||
		name == "--version" ||
		name == "-v"
}
func IsValidFlag(name string) bool {
	return name == "--keep-ep-nums" ||
		name == "-ken" ||
		name == "--starting-ep-num" ||
		name == "-sen" ||
		name == "--has-season-0" ||
		name == "-s0" ||
		name == "--naming-scheme" ||
		name == "-ns" ||
		name == "--logs" ||
		name == "-l" ||

		name == "--options" ||
		name == "-o"
}
func IsValidArgName(name string) bool {
	return IsValidCommand(name) || IsValidSwitch(name) || IsValidFlag(name)
}

func TokenizeArgs(args []string) ([]Arg, error) {
	defer timer("TokenizeArgs")()

	var tokenizedArgs []Arg
	var newArg Arg
	var value string
	readValues := false

	for i, arg := range args {
		if arg == "" {
			continue
		}
		if !readValues {
			if IsValidArgName(arg) {
				newArg.name = arg
				readValues = true
			} else {
				return nil, fmt.Errorf("start with invalid flag: '%s'", arg)
			}
			if i == len(args)-1 {
				newArg.value = value
				tokenizedArgs = append(tokenizedArgs, newArg)
			}
		} else {
			if IsValidArgName(arg) {
				newArg.value = value
				tokenizedArgs = append(tokenizedArgs, newArg)
				newArg = Arg{name: arg}
				value = ""
			} else {
				if value == "" {
					value = arg
				} else {
					return nil, fmt.Errorf("multiple values for flag: '%s' {%s, %s}", newArg.name, value, arg)
				}
			}
			if i == len(args)-1 {
				newArg.value = value
				tokenizedArgs = append(tokenizedArgs, newArg)
			}
		}
	}
	return tokenizedArgs, nil
}

// Args is a list of arguments needed by gorn to properly rename media files
type Args struct {
	root    []string
	series  []string
	movies  []string
	flags Flags
}

// Flags are options for modifying the behavior of renaming files
type Flags struct {
	keepEpNums    Option[bool]
	startingEpNum Option[int]
	hasSeason0    Option[bool]
	namingScheme  Option[string]
}
// Returns true if any of the flags are assigned
func (f *Flags) AnyAssigned() bool {
	return f.keepEpNums.IsSome() || f.startingEpNum.IsSome() || f.hasSeason0.IsSome() || f.namingScheme.IsSome()
}

// For text color on log headers
type LogHeader string
const (
	// for informational logs
	INFO LogHeader = "[INFO] " // default color (white)

	// can safely skip error, doesn't interrupt process
	WARN LogHeader = "\033[93m[WARN]\033[0m " // yellow

	// cannot safely skip error, must interrupt process
	FATAL LogHeader = "\033[91m[FATAL]\033[0m " // red

	// for timing purposes
	TIME LogHeader = "\033[94m[TIME]\033[0m " // blue
)

// LogLevel handles which logs to print based on level
type LogLevel int8 // can only be 1-4
const (
	NONE 		LogLevel = 0
	FATAL_LEVEL LogLevel = 1
	WARN_ONLY 	LogLevel = 2
	INFO_ONLY 	LogLevel = 3
	TIME_ONLY 	LogLevel = 4
	WARN_LEVEL 	LogLevel = 5
	INFO_LEVEL 	LogLevel = 6
	TIME_LEVEL 	LogLevel = 7
)

func (l *LogLevel) Headers() (string, error) {
	var headers string
	switch *l {
	case FATAL_LEVEL:
		headers += fmt.Sprint(FATAL)
	case WARN_LEVEL, WARN_ONLY:
		headers += fmt.Sprint(WARN)
	case INFO_LEVEL, INFO_ONLY:
		headers += fmt.Sprint(INFO)
	case TIME_LEVEL, TIME_ONLY:
		headers += fmt.Sprint(TIME)
	default:
		return "", fmt.Errorf("invalid log level: %d", *l)
	}
	return headers, nil
}

// ToLogLevel converts a string to a LogLevel if the string passed is valid; otherwise it returns an error.
//
// Valid log levels:
//
//	none, all, fatal, warn, info, time
func ToLogLevel(s string) (LogLevel, error) {
	switch s {
	case "fatal", "fatal-only":
		return FATAL_LEVEL, nil
	case "warn-only":
		return WARN_ONLY, nil
	case "info-only":
		return INFO_ONLY, nil
	case "time-only":
		return TIME_ONLY, nil
	case "warn":
		return WARN_LEVEL, nil
	case "info", "", "all":
		return INFO_LEVEL, nil
	case "time":
		return TIME_LEVEL, nil
	case "none":
		return NONE, nil
	}
	return -1, fmt.Errorf("invalid value '%s' for --logs. Must be 'all', 'none', or a valid log header", s)
}

// newArgs returns a new Args struct with default values for the Flags
//
// Commands though are empty by default and need to be populated
func newArgs() Args {
	return Args{
		root:   make([]string, 0),
		series: make([]string, 0),
		movies: make([]string, 0),
		flags: Flags{
			hasSeason0:    some[bool](false),
			keepEpNums:    some[bool](false),
			startingEpNum: some[int](1),
			namingScheme:  some[string]("default"),
		},
	}
}

func (args *Args) Log() {
	defer timer("Args.Log")()

	if len(args.root) > 0 {
		gornLog(INFO, "root directories: ")
		for _, root := range args.root {
			gornLog(INFO, "\t", root)
		}
	}
	if len(args.series) > 0 {
		gornLog(INFO, "series sources:")
		for _, series := range args.series {
			gornLog(INFO, "\t", series)
		}
	}
	if len(args.movies) > 0 {
		gornLog(INFO, "movies sources:")
		for _, movie := range args.movies {
			gornLog(INFO, "\t", movie)
		}
	}
	ken, err := args.flags.keepEpNums.Get()
	if err == nil {
		gornLog(INFO, "keep episode numbers: ", ken)
	}
	sen, err := args.flags.startingEpNum.Get()
	if err == nil {
		gornLog(INFO, "starting episode number: ", sen)
	}
	s0, err := args.flags.hasSeason0.Get()
	if err == nil {
		gornLog(INFO, "has season 0: ", s0)
	}
	ns, err := args.flags.namingScheme.Get()
	if err == nil {
		gornLog(INFO, "naming scheme: ", ns)
	}
	headers, err := logLevel.Headers()
	if err == nil {
		gornLog(INFO, "Showing the following logs:", headers)
	}
}

func ParseArgs(args []Arg) (Args, error) {
	defer timer("ParseArgs")()

	if len(args) < 1 {
		return Args{}, fmt.Errorf("not enough arguments: '%v'", args)
	}
	isAssigned := map[string]bool{
		"--options":         false,
		"--keep-ep-nums":    false,
		"--starting-ep-num": false,
		"--has-season-0":    false,
		"--naming-scheme":   false,
		"--logs":            false,
	}

	parsedArgs := newArgs()
	for i, arg := range args {
		if arg.name == "--help" || arg.name == "-h" {
			if len(args) <= i+1 {
				Help(arg.value)
			} else if len(args) > i+1 {
				Help(args[i+1].name)
			}
			return Args{}, SafeErrorF("safe exit")

		} else if arg.name == "--version" || arg.name == "-v" {
			Version(version)
			if len(args) > i+1 {
				gornLog(WARN, "There are no arguments for --version/-v. got:", args[i+1].name)
			}
			return Args{}, SafeErrorF("safe exit")

		} else if IsValidCommand(arg.name) {
			// no value after command
			if arg.value == "" {
				return Args{}, fmt.Errorf("missing dir path value for flag '%s'", arg)
			}

			dir, err := filepath.Abs(arg.value)
			if err != nil {
				return Args{}, err
			}
			_, err = os.Stat(dir)
			if os.IsNotExist(err) {
				return Args{}, fmt.Errorf("'%s' is not a valid directory", dir)
			}

			switch arg.name {
			case "root":
				parsedArgs.root = append(parsedArgs.root, dir)
			case "series":
				parsedArgs.series = append(parsedArgs.series, dir)
			case "movies":
				parsedArgs.movies = append(parsedArgs.movies, dir)
			}

		} else if arg.name == "--has-season-0" || arg.name == "-s0" {
			if parsedArgs.flags.hasSeason0.IsSome() && isAssigned["--has-season-0"] {
				return Args{}, fmt.Errorf("only one --has-season-0 flag is allowed")
			}
			// use default value
			if arg.value == "" {
				parsedArgs.flags.hasSeason0 = some[bool](true)
				isAssigned["--has-season-0"] = true

			} else if arg.value != "yes" && arg.value != "var" && arg.value != "no" && arg.value != "default" {
				return Args{}, fmt.Errorf("invalid value '%s' for flag --has-season-0. Must be 'yes', 'no', 'var, or 'default", arg.value)

			} else {
				switch arg.value {
				case "yes":
					parsedArgs.flags.hasSeason0 = some[bool](true)
				case "no", "default":
					parsedArgs.flags.hasSeason0 = some[bool](false)
				case "var":
					parsedArgs.flags.hasSeason0 = none[bool]()
				}
				isAssigned["--has-season-0"] = true
			}

		} else if arg.name == "--keep-ep-nums" || arg.name == "-ken" {
			if parsedArgs.flags.keepEpNums.IsSome() && isAssigned["--keep-ep-nums"] {
				return Args{}, fmt.Errorf("only one --keep-ep-nums flag is allowed")
			}

			// use default value
			if arg.value == "" {
				parsedArgs.flags.keepEpNums = some[bool](true)
				isAssigned["--keep-ep-nums"] = true

			} else if arg.value != "yes" && arg.value != "var" && arg.value != "no" && arg.value != "default" {
				return Args{}, fmt.Errorf("invalid value '%s' for --keep-ep-nums. Must be 'yes', 'no', 'var', or 'default", arg.value)

			} else {
				switch arg.value {
				case "yes":
					parsedArgs.flags.keepEpNums = some[bool](true)
				case "no", "default":
					parsedArgs.flags.keepEpNums = some[bool](false)
				case "var":
					parsedArgs.flags.keepEpNums = none[bool]()
				}
				isAssigned["--keep-ep-nums"] = true
			}

		} else if arg.name == "--starting-ep-num" || arg.name == "-sen" {
			if parsedArgs.flags.startingEpNum.IsSome() && isAssigned["--starting-ep-num"] {
				return Args{}, fmt.Errorf("only one --starting-ep-num flag is allowed")
			}

			// use default value
			if arg.value == "" {
				parsedArgs.flags.startingEpNum = some[int](1)
				isAssigned["--starting-ep-num"] = true

			} else if value, err := strconv.Atoi(arg.value); err != nil && value < 1 && arg.value != "var" && arg.value != "default" {
				return Args{}, fmt.Errorf("invalid value '%s' for --starting-ep-num. Must be a valid positive int or 'var", arg.value)

			} else {
				switch arg.value {
				case "var":
					parsedArgs.flags.startingEpNum = none[int]()
				case "default":
					parsedArgs.flags.startingEpNum = some[int](1)
				default:
					parsedArgs.flags.startingEpNum = some[int](value)
				}
				isAssigned["--starting-ep-num"] = true
			}

		} else if arg.name == "--options" || arg.name == "-o" {
			if isAssigned["--options"] {
				return Args{}, fmt.Errorf("only one --options flag is allowed")
			}
			isAssigned["--options"] = true

			if !isAssigned["--keep-ep-nums"] {
				parsedArgs.flags.keepEpNums = none[bool]()
			}
			if !isAssigned["--has-season-0"] {
				parsedArgs.flags.hasSeason0 = none[bool]()
			}
			if !isAssigned["--starting-ep-num"] {
				parsedArgs.flags.startingEpNum = none[int]()
			}
			if !isAssigned["--naming-scheme"] {
				parsedArgs.flags.namingScheme = none[string]()
			}

		} else if arg.name == "--naming-scheme" || arg.name == "-ns" {
			if parsedArgs.flags.namingScheme.IsSome() && isAssigned["--naming-scheme"] {
				return Args{}, fmt.Errorf("only one --naming-scheme flag is allowed")
			}
			if arg.value == "" {
				return Args{}, fmt.Errorf("missing value for --naming-scheme")
			}

			err := ValidateNamingScheme(arg.value)
			if err != nil && arg.value != "default" && arg.value != "var" {
				return Args{}, fmt.Errorf("invalid value '%s' for --naming-scheme. Must be 'default', 'var', or a naming scheme enclosed in double quotes", arg.value)

			} else {
				switch arg.value {
				case "default":
					parsedArgs.flags.namingScheme = some[string]("default")
				case "var":
					parsedArgs.flags.namingScheme = none[string]()
				default:
					namingScheme := strings.Trim(arg.value, `"`)
					parsedArgs.flags.namingScheme = some[string](namingScheme)
				}
			}
		} else if arg.name == "--logs" || arg.name == "-l" {
			if isAssigned["--logs"] {
				return Args{}, fmt.Errorf("only one --logs flag is allowed")
			}
			tmp, err := ToLogLevel(strings.ToLower(arg.value))
			if err != nil {
				return Args{}, err
			}
			logLevel = tmp
			isAssigned["--logs"] = true

		} else {
			return Args{}, fmt.Errorf("unknown flag: %s", arg.name)
		}
	}

	err := ValidateRoots(parsedArgs.root, parsedArgs.series, parsedArgs.movies)
	if err != nil {
		return Args{}, err
	}

	return parsedArgs, nil
}

// ValidateNamingScheme checks if a naming scheme is valid by:
//   - tokenizing APIs
//   - checking if each API is valid
//   - validating each API's value if user provided any
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

// Splits the naming scheme into a list of APIs
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

// ValidateRoots checks if:
//   - there is at least one root/source directory
//   - each root/source directory exists
//   - each root/source directory is not a subdirectory of another root/source directory
//   - each root/source directory is not a duplicate of another root/source directory
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
