package main

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

type Arg struct {
	flag  string
	value string
}

type Args struct {
	root            []Arg
	series          []Arg
	movies          []Arg
	has_season_0    Arg
	keep_ep_nums    Arg
	starting_ep_num Arg
	naming_scheme   Arg
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

	var parsed_args Args
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
				parsed_args.root = append(parsed_args.root, Arg{flag: arg, value: dir})
			} else if arg == "--series" || arg == "-s" {
				parsed_args.series = append(parsed_args.series, Arg{flag: arg, value: dir})
			} else if arg == "--movies" || arg == "-m" {
				parsed_args.movies = append(parsed_args.movies, Arg{flag: arg, value: dir})
			}
			skip_iter = i + 1

		} else if arg == "--season-0" || arg == "-s0" {
			if parsed_args.has_season_0.flag != "" {
				return Args{}, fmt.Errorf("only one --season-0 flag is allowed")
			}

			// default value
			if len(args) <= i+1 || (len(args) > i+1 && args[i+1][0] == '-') {
				parsed_args.has_season_0 = Arg{flag: arg, value: "all yes"}

			} else if args[i+1] != "all" && args[i+1] != "var" {
				return Args{}, fmt.Errorf("invalid value '%s' for flag '%s'. Must be 'all' or 'var", args[i+1], arg)

			} else if args[i+1] == "all" {
				if len(args) < i+2 || args[i+2][0] == '-' {
					return Args{}, fmt.Errorf("all must be followed by 'yes' or 'no' for --season-0. %s is invalid", args[i+2])

				} else if args[i+2] != "yes" && args[i+2] != "no" {
					return Args{}, fmt.Errorf("all must be followed by 'yes' or 'no' for --season-0. %s is invalid", args[i+2])
				}

				parsed_args.has_season_0 = Arg{flag: arg, value: args[i+1] + " " + args[i+2]}
				skip_iter = i + 2

			} else if args[i+1] == "var" {
				parsed_args.has_season_0 = Arg{flag: arg}
			}

		} else if arg == "--keep-ep-nums" || arg == "-ken" {
			if parsed_args.keep_ep_nums.flag != "" {
				return Args{}, fmt.Errorf("only one --keep-ep-nums flag is allowed")
			}

			// default value
			if len(args) <= i+1 || (len(args) > i+1 && args[i+1][0] == '-') {
				parsed_args.keep_ep_nums = Arg{flag: arg, value: "all yes"}

			} else if args[i+1] != "all" && args[i+1] != "var" {
				return Args{}, fmt.Errorf("invalid value '%s' for --keep-ep-nums. Must be 'all' or 'var", args[i+1])

			} else if args[i+1] == "all" {
				if len(args) < i+2 || args[i+2][0] == '-' {
					return Args{}, fmt.Errorf("all must be followed by 'yes' or 'no' for --keep-ep-nums. %s is invalid", args[i+2])

				} else if args[i+2] != "yes" && args[i+2] != "no" {
					return Args{}, fmt.Errorf("all must be followed by 'yes' or 'no' for --keep-ep-nums. %s is invalid", args[i+2])
				}

				parsed_args.keep_ep_nums = Arg{flag: arg, value: args[i+1] + " " + args[i+2]}
				skip_iter = i + 2

			} else if args[i+1] == "var" {
				parsed_args.has_season_0 = Arg{flag: arg}
			}

		} else if arg == "--starting-ep-num" || arg == "-sen" {
			if parsed_args.starting_ep_num.flag != "" {
				return Args{}, fmt.Errorf("only one --starting-ep-num flag is allowed")
			}

			// default value
			if len(args) <= i+1 || (len(args) > i+1 && args[i+1][0] == '-') {
				parsed_args.starting_ep_num = Arg{flag: arg, value: "all yes"}

			} else if args[i+1] != "all" && args[i+1] != "var" {
				return Args{}, fmt.Errorf("invalid value '%s' for --starting-ep-num. Must be 'all' or 'var", args[i+1])

			} else if args[i+1] == "all" {
				if len(args) < i+2 || args[i+2][0] == '-' {
					return Args{}, fmt.Errorf("all must be followed by a positive int for --starting-ep-num. %s is not a valid positive int", args[i+2])
				}

				_, err := strconv.Atoi(args[i+2])
				if err != nil {
					return Args{}, fmt.Errorf("all must be followed by a positive int for --starting-ep-num. %s is not a valid positive int", args[i+2])
				}

				parsed_args.starting_ep_num = Arg{flag: arg, value: args[i+1] + " " + args[i+2]}
				skip_iter = i + 2

			} else if args[i+1] == "var" {
				parsed_args.starting_ep_num = Arg{flag: arg}
			}

		} else if arg == "--naming-scheme" || arg == "-ns" {
			if parsed_args.naming_scheme.flag != "" {
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

				parsed_args.naming_scheme = Arg{flag: arg, value: args[i+1] + " " + args[i+2]}
				skip_iter = i + 2

			} else if args[i+1] == "var" {
				parsed_args.naming_scheme = Arg{flag: arg}
			}

		} else {
			return Args{}, fmt.Errorf("unknown flag: %s", arg)
		}
	}

	if len(parsed_args.root) == 0 && len(parsed_args.series) == 0 && len(parsed_args.movies) == 0 {
		return Args{}, fmt.Errorf("must specify at least one of --root, --series, or --movie")
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
	valid_range := regexp.MustCompile(`^\d+,\s*\d+$`)

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
			res := strings.SplitN(val, ",", 2)
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