package main

import (
	"fmt"
	"path/filepath"
	"strconv"
)

type Arg struct {
	flag string
	value string
}

type Args struct {
	root []Arg
	series []Arg
	movies []Arg
	has_season_0 Arg
	keep_ep_nums Arg
	starting_ep_num Arg
	naming_scheme Arg
}

func parse_args(args []string) (Args, error) {
	if len(args) < 2 {
		return Args{}, fmt.Errorf("not enough arguments")
	}

	directory_args := map[string]bool{
		"--root": true,
		"-r": true,
		"--series": true,
		"-s": true,
		"--movies": true,
		"-m": true,
	}

	var parsed_args Args
	var skip_iter int
	for i, arg := range args {
		// valid values will be skipped
		if i < skip_iter {
			continue

		// catch invalid values acting as flags
		} else if arg[0] != '-' {
			return Args{}, fmt.Errorf("invalid argument: '%s'", arg)

		} else if directory_args[arg] {
			// no value after flag / flag after flag
			if len(args) < i+1 || args[i+1][0] == '-' {
				return Args{}, fmt.Errorf("missing dir path value for flag '%s'", arg)

			// not a valid directory
			} else if _, err := filepath.Abs(args[i+1]); err != nil {
				return Args{}, err
			}

			if arg == "--root" || arg == "-r" {
				parsed_args.root = append(parsed_args.root, Arg{flag: arg, value: args[i+1]})
			} else if arg == "--series" || arg == "-s" {
				parsed_args.series = append(parsed_args.series, Arg{flag: arg, value: args[i+1]})
			} else if arg == "--movies" || arg == "-m" {
				parsed_args.movies = append(parsed_args.movies, Arg{flag: arg, value: args[i+1]})
			}

		} else if arg == "--season-0" || arg == "-s0" || 
				  arg == "--keep-ep-nums" || arg == "-ken" || 
				  arg == "--starting-ep-num" || arg == "-sen" {

			if len(args) < i+1 || args[i+1][0] == '-' {
				return Args{}, fmt.Errorf("missing value for flag '%s'", arg)

			} else if args[i+1] != "all" && args[i+1] != "var" {
				return Args{}, fmt.Errorf("invalid value '%s' for flag '%s'", args[i+1], arg)

			} else if args[i+1] == "all" {
				if len(args) < i+2 || args[i+2][0] == '-' {
					return Args{}, fmt.Errorf("missing value for flag '%s'", arg)

				} else if (arg == "--season-0" || arg == "-s0" || arg == "--keep-ep-nums" || arg == "-ken") &&
				          (args[i+2] != "yes" && args[i+2] != "no") {
					return Args{}, fmt.Errorf("invalid value '%s' for flag '%s'", args[i+2], arg)

				} else if (arg == "--starting-ep-num" || arg == "-sen") {
					_, err := strconv.Atoi(args[i+2])
					if err != nil {
						return Args{}, fmt.Errorf("invalid value '%s' for flag '%s'. Must be a positive integer", args[i+2], arg)
					}
				}

				parsed_args.has_season_0 = Arg{flag: arg, value: fmt.Sprintf("%s %s", args[i+1], args[i+2])}

			} else if args[i+1] == "var" {
				// if "var", must not be followed by anything; aka must be followed by a command
				if len(args) > i+2 && args[i+2][0] != '-' {
					return Args{}, fmt.Errorf("unexpected value '%s' for flag '%s'", args[i+1], arg)
				}

				// parsed_args = append(parsed_args, Arg{flag: arg})
				parsed_args.has_season_0 = Arg{flag: arg}
			}

		} else if arg == "--naming-scheme" || arg == "-ns" {
			if len(args) < i+1 || args[i+1][0] == '-' {
				return Args{}, fmt.Errorf("missing value for flag '%s'", arg)
			
			} else if args[i+1][0] != '"' && args[i+1][len(args[i+1])] != '"' {
				return Args{}, fmt.Errorf("invalid value '%s' for flag '%s'", args[i+1], arg)
			}

			parsed_args.naming_scheme = Arg{flag: arg, value: args[i+1][1:len(args[i+1])-1]}

		} else {
			return Args{}, fmt.Errorf("unknown flag: %s", arg)
		}
	}

	if len(parsed_args.root) == 0 && len(parsed_args.series) == 0 && len(parsed_args.movies) == 0 {
		return Args{}, fmt.Errorf("must specify at least one of --root, --series, or --movie")
	}

	return parsed_args, nil
}
