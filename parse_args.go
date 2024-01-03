package main

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"os"
)

type Args struct {
	root            	[]string
	series          	[]string
	movies          	[]string
	options 	AdditionalOptions
}
type AdditionalOptions struct {
	keep_ep_nums    Option[bool]
	starting_ep_num Option[int]
	has_season_0    Option[bool]
	naming_scheme   Option[string]
}
func new_Args() Args {
	return Args{
		root:            make([]string, 0),
		series:          make([]string, 0),
		movies:          make([]string, 0),
		options: AdditionalOptions{
			has_season_0:    none[bool](),
			keep_ep_nums:    none[bool](),
			starting_ep_num: none[int](),
			naming_scheme:   none[string](),
		},
	}
}

func parse_args(args []string) (Args, error) {
	if len(args) < 1 {
		return Args{}, fmt.Errorf("not enough arguments")
	}

	directory_args := map[string]bool{
		"--root":   true,
		"-r":       true,
		"--series": true,
		"-s":       true,
		"--movies": true,
		"-m":       true,
	}

	parsed_args := new_Args()
	skip_iter := 0
	for i, arg := range args {
		// valid values will be skipped
		if skip_iter != 0 && i <= skip_iter {
			continue

		// catch invalid values acting as flags
		} else if arg[0] != '-' {
			return Args{}, fmt.Errorf("invalid flag: '%s'", arg)

		} else if directory_args[arg] {
			// no value after flag / flag after flag
			if len(args) <= i+1 || (len(args) > i+1 && args[i+1][0] == '-') {
				return Args{}, fmt.Errorf("missing dir path value for flag '%s'", arg)

			// not a valid directory
			} else if _, err := filepath.Abs(args[i+1]); err != nil {
				return Args{}, err
			}

			dir, err := filepath.Abs(args[i+1])
			if err != nil {
				return Args{}, err
			}

			if arg == "--root" || arg == "-r" {
				parsed_args.root = append(parsed_args.root, dir)
			} else if arg == "--series" || arg == "-s" {
				parsed_args.series = append(parsed_args.series, dir)
			} else if arg == "--movies" || arg == "-m" {
				parsed_args.movies = append(parsed_args.movies, dir)
			}
			skip_iter = i + 1

		} else if arg == "--season-0" || arg == "-s0" {
			if parsed_args.options.has_season_0.is_some() {
				return Args{}, fmt.Errorf("only one --season-0 flag is allowed")
			}

			// use default value
			if len(args) <= i+1 || (len(args) > i+1 && args[i+1][0] == '-') {
				parsed_args.options.has_season_0 = some[bool](false)

			} else if args[i+1] != "all" && args[i+1] != "var" {
				return Args{}, fmt.Errorf("invalid value '%s' for flag '%s'. Must be 'all' or 'var", args[i+1], arg)

			} else if args[i+1] == "all" {
				if len(args) < i+2 || args[i+2][0] == '-' {
					return Args{}, fmt.Errorf("all must be followed by 'yes' or 'no' for --season-0. %s is invalid", args[i+2])

				} else if args[i+2] != "yes" && args[i+2] != "no" {
					return Args{}, fmt.Errorf("all must be followed by 'yes' or 'no' for --season-0. %s is invalid", args[i+2])
				}
				var value bool
				if args[i+2] == "yes" {
					value = true
				} else {
					value = false
				}
				parsed_args.options.has_season_0 = some[bool](value)
				skip_iter = i + 2

			} else if args[i+1] == "var" {
				parsed_args.options.has_season_0 = none[bool]()
			}

		} else if arg == "--keep-ep-nums" || arg == "-ken" {
			if parsed_args.options.keep_ep_nums.is_some() {
				return Args{}, fmt.Errorf("only one --keep-ep-nums flag is allowed")
			}

			// use default value
			if len(args) <= i+1 || (len(args) > i+1 && args[i+1][0] == '-') {
				parsed_args.options.keep_ep_nums = some[bool](false)

			} else if args[i+1] != "all" && args[i+1] != "var" {
				return Args{}, fmt.Errorf("invalid value '%s' for --keep-ep-nums. Must be 'all' or 'var", args[i+1])

			} else if args[i+1] == "all" {
				if len(args) < i+2 || args[i+2][0] == '-' {
					return Args{}, fmt.Errorf("all must be followed by 'yes' or 'no' for --keep-ep-nums. %s is invalid", args[i+2])

				} else if args[i+2] != "yes" && args[i+2] != "no" {
					return Args{}, fmt.Errorf("all must be followed by 'yes' or 'no' for --keep-ep-nums. %s is invalid", args[i+2])
				}

				var value bool
				if args[i+2] == "yes" {
					value = true
				} else {
					value = false
				}
				parsed_args.options.keep_ep_nums = some[bool](value)
				skip_iter = i + 2

			} else if args[i+1] == "var" {
				parsed_args.options.has_season_0 = none[bool]()
			}

		} else if arg == "--starting-ep-num" || arg == "-sen" {
			if parsed_args.options.starting_ep_num.is_some() {
				return Args{}, fmt.Errorf("only one --starting-ep-num flag is allowed")
			}

			// use default value
			if len(args) <= i+1 || (len(args) > i+1 && args[i+1][0] == '-') {
				parsed_args.options.starting_ep_num = some[int](1)

			} else if args[i+1] != "all" && args[i+1] != "var" {
				return Args{}, fmt.Errorf("invalid value '%s' for --starting-ep-num. Must be 'all' or 'var", args[i+1])

			} else if args[i+1] == "all" {
				if len(args) < i+2 || args[i+2][0] == '-' {
					return Args{}, fmt.Errorf("all must be followed by a positive int for --starting-ep-num. %s is not a valid positive int", args[i+2])
				}

				value, err := strconv.Atoi(args[i+2])
				if err != nil {
					return Args{}, fmt.Errorf("all must be followed by a positive int for --starting-ep-num. %s is not a valid positive int", args[i+2])
				}

				parsed_args.options.starting_ep_num = some[int](value)
				skip_iter = i + 2

			} else if args[i+1] == "var" {
				parsed_args.options.starting_ep_num = none[int]()
			}

		} else if arg == "--naming-scheme" || arg == "-ns" {
			if parsed_args.options.naming_scheme.is_some() {
				return Args{}, fmt.Errorf("only one --naming-scheme flag is allowed")
			}

			if len(args) <= i+1 || (len(args) > i+1 && args[i+1][0] == '-') {
				return Args{}, fmt.Errorf("missing value for --naming-scheme")

			} else if args[i+1] != "all" && args[i+1] != "var" {
				return Args{}, fmt.Errorf("invalid value '%s' for --naming-scheme. Must be 'all' or 'var", args[i+1])
			
			} else if args[i+1] == "all" {
				if len(args) < i+2 || args[i+2][0] == '-' {
					return Args{}, fmt.Errorf("'all' must be followed by a naming scheme string enclosed in double quotes")
				}
				err := validate_naming_scheme(args[i+2])
				if err != nil {
					return Args{}, err
				}

				parsed_args.options.naming_scheme = some[string](args[i+2])
				skip_iter = i + 2

			} else if args[i+1] == "var" {
				parsed_args.options.naming_scheme = none[string]()
			}

		} else {
			return Args{}, fmt.Errorf("unknown flag: %s", arg)
		}
	}

	err := validate_roots(parsed_args.root, parsed_args.series, parsed_args.movies)
	if err != nil {
		return Args{}, err
	}

	// use default values for additional options
	if parsed_args.options.has_season_0.is_none() {
		parsed_args.options.has_season_0 = some[bool](false)
	}
	if parsed_args.options.keep_ep_nums.is_none() {
		parsed_args.options.keep_ep_nums = some[bool](false)
	}
	if parsed_args.options.starting_ep_num.is_none() {
		parsed_args.options.starting_ep_num = some[int](1)
	}
	if parsed_args.options.naming_scheme.is_none() {
		parsed_args.options.naming_scheme = none[string]()
	}

	return parsed_args, nil
}

func validate_naming_scheme(s string) error {
	tokens, err := tokenize_naming_scheme(s)
	if err != nil {
		return err
	}

	valid_api := regexp.MustCompile(`^season_num$|^episode_num$|^self$`)
	valid_parent_api := regexp.MustCompile(`^parent(-parent)*$|^p(-\d+)?$`)
	valid_range := regexp.MustCompile(`^\d+(\s*,\s*\d+)?$`)

	for _,token := range tokens {
		var api, val string

		// for api with values
		if strings.Contains(token, ":") {
			res := strings.SplitN(token, ":", 2)
			api, val = strings.TrimSpace(res[0]), strings.TrimSpace(res[1])
		} else {
			api, val = strings.TrimSpace(token), "none"
		}

		if !valid_api.MatchString(api) && !valid_parent_api.MatchString(api) {
			return fmt.Errorf("invalid api: %s", api)
		}

		if api == "season_num" || api ==  "episode_num" {
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

		} else if api == "self" || valid_parent_api.MatchString(api) {
			if val == "none" {
				continue
			}
			if val == "" {
				return fmt.Errorf("%s's value cannot be empty", api)
			}
			if val[0] == '\'' {
				single_closed_string := regexp.MustCompile(`^'[^']*'$`)
				if !single_closed_string.MatchString(val) {
					return fmt.Errorf("%s is either an unclosed string or contains more than 2 single quotes", val)
				}
				
				val_trimmed := strings.Trim(val, "'")
				just_whitespace := regexp.MustCompile(`^\s*$`)
				if just_whitespace.MatchString(val_trimmed) {
					return fmt.Errorf("%s is an empty string (just whitespace/s)", val)
				}

				if _, err := regexp.Compile(val_trimmed); err != nil {
					return fmt.Errorf("invalid regex: %s", val)

				} else {
					parts := split_regex_by_pipe(val_trimmed)
					for _, part := range parts {
						if !has_only_one_match_group(part) {
							return fmt.Errorf("regex should have only one match group per part (parts are separated by outermost pipes |): %s", val)
						}
					}
					// valid regex
					continue
				}

			} else if !valid_range.MatchString(val) {
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

func tokenize_naming_scheme(s string) ([]string, error) {
	is_token := false
	builder := strings.Builder{}
	naming_scheme := make([]string, 0)

	for i, c := range s {
		if is_token && i+1 == len(s) && c != '>' {
		    return nil, fmt.Errorf("reached end of string but still in an unclosed api: '<%s%s'", builder.String(), string(c))
		}

		// start of token
		if c == '<' && !is_token {
			is_token = true
			builder.Reset()
			continue

		// end of token
		} else if c == '>' && is_token {
			is_token = false
			naming_scheme = append(naming_scheme, builder.String())
			builder.Reset()
			continue
		}

		if is_token {
			_, err := builder.WriteRune(c)
			if err != nil {
				return nil, err
			}
		}
	}

	return naming_scheme, nil
}

func validate_roots(root []string, series []string, movies []string) error {
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