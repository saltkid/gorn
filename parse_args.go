package main

import (
	"fmt"
	"path/filepath"
	"strconv"
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
				// todo validate naming scheme

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
