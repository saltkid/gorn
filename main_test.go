package main

import (
	"path/filepath"
	"strings"
	"testing"
)

func Test_ParseArgs(t *testing.T) {
	t.Log("------------expects errors------------")
	
	cmd := "root -s0 yes"
	command := strings.Split(cmd, " ")
	_, err := ParseArgs(command)
	if err == nil {
		t.Errorf("expected error 'missing root dir'; got -s0 -s0")
	} else {
		t.Log(cmd, "\n\t", err)
	}

	cmd = "series -s0 yes"
	command = strings.Split(cmd, " ")
	_, err = ParseArgs(command)
	if err == nil {
		t.Errorf("expected error 'missing series dir'; got -s0")
	} else {
		t.Log(cmd, "\n\t", err)
	}

	cmd = "movies -s0 yes"
	command = strings.Split(cmd, " ")
	_, err = ParseArgs(command)
	if err == nil {
		t.Errorf("expected error 'missing movies dir'; got -s0")
	} else {
		t.Log(cmd, "\n\t", err)
	}

	cmd = "-s0 yes"
	command = strings.Split(cmd, " ")
	_, err = ParseArgs(command)
	if err == nil {
		t.Errorf("expected error 'missing dir'")
	} else {
		t.Log(cmd, "\n\t", err)
	}

	cmd = "root ./test_files series ./test_files"
	command = strings.Split(cmd, " ")
	_, err = ParseArgs(command)
	if err == nil {
		t.Errorf("expected error 'same dir: root test_files series test_files'")
	} else {
		t.Log(cmd, "\n\t", err)
	}

	cmd = "root ./test_files series ./test_files/series"
	command = strings.Split(cmd, " ")
	_, err = ParseArgs(command)
	if err == nil {
		t.Errorf("expected error 'series is a subdir of root: root test_files series test_files/series'")
	} else {
		t.Log(cmd, "\n\t", err)
	}
	
	cmd = "series ./test_files root ./test_files/series"
	command = strings.Split(cmd, " ")
	_, err = ParseArgs(command)
	if err == nil {
		t.Errorf("expected error 'root is a subdir of series: series test_files root test_files/series'")
	} else {
		t.Log(cmd, "\n\t", err)
	}

	cmd = "root ./test_files movies ./test_files"
	command = strings.Split(cmd, " ")
	_, err = ParseArgs(command)
	if err == nil {
		t.Errorf("expected error 'same dir: root test_files movies test_files'")
	} else {
		t.Log(cmd, "\n\t", err)
	}

	cmd = "root ./test_files movies ./test_files/movies"
	command = strings.Split(cmd, " ")
	_, err = ParseArgs(command)
	if err == nil {
		t.Errorf("expected error 'movies is a subdir of root: root test_files movies test_files/movies'")
	} else {
		t.Log(cmd, "\n\t", err)
	}
	
	cmd = "movies ./test_files root ./test_files/movies"
	command = strings.Split(cmd, " ")
	_, err = ParseArgs(command)
	if err == nil {
		t.Errorf("expected error 'root is a subdir of movies: movies test_files root test_files/movies'")
	} else {
		t.Log(cmd, "\n\t", err)
	}
	
	cmd = "root ./test_files -ken ye"
	command = strings.Split(cmd, " ")
	_, err = ParseArgs(command)
	if err == nil {
		t.Errorf("expected error 'invalid value for -ken: ye'")
	} else {
		t.Log(cmd, "\n\t", err)
	}
	
	cmd = "root ./test_files -sen ye"
	command = strings.Split(cmd, " ")
	_, err = ParseArgs(command)
	if err == nil {
		t.Errorf("expected error 'invalid value for -sen: ye'")
	} else {
		t.Log(cmd, "\n\t", err)
	}
		
	cmd = "root ./test_files -s0 ye"
	command = strings.Split(cmd, " ")
	_, err = ParseArgs(command)
	if err == nil {
		t.Errorf("expected error 'invalid value for -s0: ye'")
	} else {
		t.Log(cmd, "\n\t", err)
	}
	
	cmd = "root ./test_files -ns ye"
	command = strings.Split(cmd, " ")
	_, err = ParseArgs(command)
	if err == nil {
		t.Errorf("expected error 'invalid value for -ns: ye'")
	} else {
		t.Log(cmd, "\n\t", err)
	}

	t.Log("------------expects success------------")
	cmd = "root ./test_files -s0 yes"
	command = strings.Split(cmd, " ")
	_, err = ParseArgs(command)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	} else {
		t.Log(cmd)
	}
	
	cmd = "root ./test_files -s0 no"
	command = strings.Split(cmd, " ")
	_, err = ParseArgs(command)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	} else {
		t.Log(cmd)
	}

	cmd = "root ./test_files -s0 default"
	command = strings.Split(cmd, " ")
	_, err = ParseArgs(command)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	} else {
		t.Log(cmd)
	}

	cmd = "root ./test_files -s0 var"
	command = strings.Split(cmd, " ")
	_, err = ParseArgs(command)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	} else {
		t.Log(cmd)
	}

	cmd = `root ./test_files -ns "test"`
	command = strings.Split(cmd, " ")
	args, err := ParseArgs(command)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	} else {
		if args.options.namingScheme.IsNone() {
			t.Errorf("unexpected error: %s", err)
		} else {
			val, _ := args.options.namingScheme.Get()
			if val != "test" {
				t.Errorf("unexpected error: %s", err)
			}
			t.Log(cmd)
		}
	}

	cmd = "root ./test_files -ns default"
	command = strings.Split(cmd, " ")
	_, err = ParseArgs(command)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	} else {
		t.Log(cmd)
	}

	cmd = "root ./test_files -ns var"
	command = strings.Split(cmd, " ")
	_, err = ParseArgs(command)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	} else {
		t.Log(cmd)
	}

	cmd = "root ./test_files -ken yes"
	command = strings.Split(cmd, " ")
	_, err = ParseArgs(command)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	} else {
		t.Log(cmd)
	}
	
	cmd = "root ./test_files -ken no"
	command = strings.Split(cmd, " ")
	_, err = ParseArgs(command)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	} else {
		t.Log(cmd)
	}

	cmd = "root ./test_files -ken default"
	command = strings.Split(cmd, " ")
	_, err = ParseArgs(command)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	} else {
		t.Log(cmd)
	}

	cmd = "root ./test_files -ken var"
	command = strings.Split(cmd, " ")
	_, err = ParseArgs(command)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	} else {
		t.Log(cmd)
	}

	cmd = "root ./test_files -sen yes"
	command = strings.Split(cmd, " ")
	_, err = ParseArgs(command)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	} else {
		t.Log(cmd)
	}
	
	cmd = "root ./test_files -sen no"
	command = strings.Split(cmd, " ")
	_, err = ParseArgs(command)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	} else {
		t.Log(cmd)
	}

	cmd = "root ./test_files -sen default"
	command = strings.Split(cmd, " ")
	_, err = ParseArgs(command)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	} else {
		t.Log(cmd)
	}

	cmd = "root ./test_files -sen var"
	command = strings.Split(cmd, " ")
	_, err = ParseArgs(command)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	} else {
		t.Log(cmd)
	}

	cmd = `root ./test_files -ken -sen -s0 -ns "test"`
	command = strings.Split(cmd, " ")
	_, err = ParseArgs(command)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	} else {
		if args.options.namingScheme.IsNone() {
			t.Errorf("unexpected error: %s", err)
		} else {
			val, _ := args.options.namingScheme.Get()
			if val != "test" {
				t.Errorf("unexpected error: %s", err)
			}
			t.Log(cmd)
		}
	}
	
	cmd = "root ./test_files -o"
	command = strings.Split(cmd, " ")
	args, err = ParseArgs(command)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	} else {
		if args.options.keepEpNums.IsSome() || args.options.hasSeason0.IsSome() || args.options.startingEpNum.IsSome() || args.options.namingScheme.IsSome() {
			t.Errorf("unexpected error: %s", err)
		} else {
			t.Log(cmd)
		}
	}

	cmd = "-h"
	command = strings.Split(cmd, " ")
	_, err = ParseArgs(command)
	if err == nil {
		t.Errorf("supposed to exit safely")
	} else {
		t.Log(cmd, err)
	}

	cmd = "-v"
	command = strings.Split(cmd, " ")
	_, err = ParseArgs(command)
	if err == nil {
		t.Errorf("supposed to exit safely")
	} else {
		t.Log(cmd, err)
	}

	cmd = "-h -v"
	command = strings.Split(cmd, " ")
	_, err = ParseArgs(command)
	if err == nil {
		t.Errorf("supposed to exit safely")
	} else {
		t.Log(cmd, err)
	}

	cmd = "-h -o"
	command = strings.Split(cmd, " ")
	_, err = ParseArgs(command)
	if err == nil {
		t.Errorf("supposed to exit safely")
	} else {
		t.Log(cmd, err)
	}

	cmd = "-h -s0"
	command = strings.Split(cmd, " ")
	_, err = ParseArgs(command)
	if err == nil {
		t.Errorf("supposed to exit safely")
	} else {
		t.Log(cmd, err)
	}

	cmd = "-h -ken"
	command = strings.Split(cmd, " ")
	_, err = ParseArgs(command)
	if err == nil {
		t.Errorf("supposed to exit safely")
	} else {
		t.Log(cmd, err)
	}

	cmd = "-h -sen"
	command = strings.Split(cmd, " ")
	_, err = ParseArgs(command)
	if err == nil {
		t.Errorf("supposed to exit safely")
	} else {
		t.Log(cmd, err)
	}

	cmd = "-h -ns"
	command = strings.Split(cmd, " ")
	_, err = ParseArgs(command)
	if err == nil {
		t.Errorf("supposed to exit safely")
	} else {
		t.Log(cmd, err)
	}

	cmd = "series ./test_files/series -s0 yes"
	command = strings.Split(cmd, " ")
	_, err = ParseArgs(command)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	} else {
		t.Log(cmd)
	}
	
	cmd = "series ./test_files/series -s0 no"
	command = strings.Split(cmd, " ")
	_, err = ParseArgs(command)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	} else {
		t.Log(cmd)
	}

	cmd = "series ./test_files/series -s0 default"
	command = strings.Split(cmd, " ")
	_, err = ParseArgs(command)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	} else {
		t.Log(cmd)
	}

	cmd = "series ./test_files/series -s0 var"
	command = strings.Split(cmd, " ")
	_, err = ParseArgs(command)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	} else {
		t.Log(cmd)
	}

	cmd = `series ./test_files/series -ns "test"`
	command = strings.Split(cmd, " ")
	args, err = ParseArgs(command)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	} else {
		if args.options.namingScheme.IsNone() {
			t.Errorf("unexpected error: %s", err)
		} else {
			val, _ := args.options.namingScheme.Get()
			if val != "test" {
				t.Errorf("unexpected error: %s", err)
			}
			t.Log(cmd)
		}
	}

	cmd = "series ./test_files/series -ns default"
	command = strings.Split(cmd, " ")
	_, err = ParseArgs(command)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	} else {
		t.Log(cmd)
	}

	cmd = "series ./test_files/series -ns var"
	command = strings.Split(cmd, " ")
	_, err = ParseArgs(command)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	} else {
		t.Log(cmd)
	}

	cmd = "series ./test_files/series -ken yes"
	command = strings.Split(cmd, " ")
	_, err = ParseArgs(command)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	} else {
		t.Log(cmd)
	}
	
	cmd = "series ./test_files/series -ken no"
	command = strings.Split(cmd, " ")
	_, err = ParseArgs(command)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	} else {
		t.Log(cmd)
	}

	cmd = "series ./test_files/series -ken default"
	command = strings.Split(cmd, " ")
	_, err = ParseArgs(command)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	} else {
		t.Log(cmd)
	}

	cmd = "series ./test_files/series -ken var"
	command = strings.Split(cmd, " ")
	_, err = ParseArgs(command)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	} else {
		t.Log(cmd)
	}

	cmd = "series ./test_files/series -sen yes"
	command = strings.Split(cmd, " ")
	_, err = ParseArgs(command)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	} else {
		t.Log(cmd)
	}
	
	cmd = "series ./test_files/series -sen no"
	command = strings.Split(cmd, " ")
	_, err = ParseArgs(command)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	} else {
		t.Log(cmd)
	}

	cmd = "series ./test_files/series -sen default"
	command = strings.Split(cmd, " ")
	_, err = ParseArgs(command)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	} else {
		t.Log(cmd)
	}

	cmd = "series ./test_files/series -sen var"
	command = strings.Split(cmd, " ")
	_, err = ParseArgs(command)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	} else {
		t.Log(cmd)
	}

	cmd = `series ./test_files/series -ken -sen -s0 -ns "test"`
	command = strings.Split(cmd, " ")
	_, err = ParseArgs(command)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	} else {
		if args.options.namingScheme.IsNone() {
			t.Errorf("unexpected error: %s", err)
		} else {
			val, _ := args.options.namingScheme.Get()
			if val != "test" {
				t.Errorf("unexpected error: %s", err)
			}
			t.Log(cmd)
		}
	}
	
	cmd = "series ./test_files/series -o"
	command = strings.Split(cmd, " ")
	args, err = ParseArgs(command)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	} else {
		if args.options.keepEpNums.IsSome() || args.options.hasSeason0.IsSome() || args.options.startingEpNum.IsSome() || args.options.namingScheme.IsSome() {
			t.Errorf("unexpected error: %s", err)
		} else {
			t.Log(cmd)
		}
	}

	cmd = "movies ./test_files/movies -s0 yes"
	command = strings.Split(cmd, " ")
	_, err = ParseArgs(command)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	} else {
		t.Log(cmd)
	}
	
	cmd = "movies ./test_files/movies -s0 no"
	command = strings.Split(cmd, " ")
	_, err = ParseArgs(command)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	} else {
		t.Log(cmd)
	}

	cmd = "movies ./test_files/movies -s0 default"
	command = strings.Split(cmd, " ")
	_, err = ParseArgs(command)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	} else {
		t.Log(cmd)
	}

	cmd = "movies ./test_files/movies -s0 var"
	command = strings.Split(cmd, " ")
	_, err = ParseArgs(command)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	} else {
		t.Log(cmd)
	}

	cmd = `movies ./test_files/movies -ns "test"`
	command = strings.Split(cmd, " ")
	args, err = ParseArgs(command)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	} else {
		if args.options.namingScheme.IsNone() {
			t.Errorf("unexpected error: %s", err)
		} else {
			val, _ := args.options.namingScheme.Get()
			if val != "test" {
				t.Errorf("unexpected error: %s", err)
			}
			t.Log(cmd)
		}
	}

	cmd = "movies ./test_files/movies -ns default"
	command = strings.Split(cmd, " ")
	_, err = ParseArgs(command)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	} else {
		t.Log(cmd)
	}

	cmd = "movies ./test_files/movies -ns var"
	command = strings.Split(cmd, " ")
	_, err = ParseArgs(command)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	} else {
		t.Log(cmd)
	}

	cmd = "movies ./test_files/movies -ken yes"
	command = strings.Split(cmd, " ")
	_, err = ParseArgs(command)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	} else {
		t.Log(cmd)
	}
	
	cmd = "movies ./test_files/movies -ken no"
	command = strings.Split(cmd, " ")
	_, err = ParseArgs(command)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	} else {
		t.Log(cmd)
	}

	cmd = "movies ./test_files/movies -ken default"
	command = strings.Split(cmd, " ")
	_, err = ParseArgs(command)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	} else {
		t.Log(cmd)
	}

	cmd = "movies ./test_files/movies -ken var"
	command = strings.Split(cmd, " ")
	_, err = ParseArgs(command)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	} else {
		t.Log(cmd)
	}

	cmd = "movies ./test_files/movies -sen yes"
	command = strings.Split(cmd, " ")
	_, err = ParseArgs(command)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	} else {
		t.Log(cmd)
	}
	
	cmd = "movies ./test_files/movies -sen no"
	command = strings.Split(cmd, " ")
	_, err = ParseArgs(command)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	} else {
		t.Log(cmd)
	}

	cmd = "movies ./test_files/movies -sen default"
	command = strings.Split(cmd, " ")
	_, err = ParseArgs(command)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	} else {
		t.Log(cmd)
	}

	cmd = "movies ./test_files/movies -sen var"
	command = strings.Split(cmd, " ")
	_, err = ParseArgs(command)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	} else {
		t.Log(cmd)
	}

	cmd = `movies ./test_files/movies -ken -sen -s0 -ns "test"`
	command = strings.Split(cmd, " ")
	_, err = ParseArgs(command)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	} else {
		if args.options.namingScheme.IsNone() {
			t.Errorf("unexpected error: %s", err)
		} else {
			val, _ := args.options.namingScheme.Get()
			if val != "test" {
				t.Errorf("unexpected error: %s", err)
			}
			t.Log(cmd)
		}
	}
	
	cmd = "movies ./test_files/movies -o"
	command = strings.Split(cmd, " ")
	args, err = ParseArgs(command)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	} else {
		if args.options.keepEpNums.IsSome() || args.options.hasSeason0.IsSome() || args.options.startingEpNum.IsSome() || args.options.namingScheme.IsSome() {
			t.Errorf("unexpected error: %s", err)
		} else {
			t.Log(cmd)
		}
	}
}

func Test_namingScheme_validation(t *testing.T) {
	t.Log("------------expects errors------------")
	err := ValidateNamingScheme("S<season_num:>E<episode_num:>")
	if err == nil {
		t.Errorf("expected error 'missing value for token: <season_num:>'")
	} else {
		t.Log("S<season_num:>E<episode_num:>", "\n\t", err, "\n")
	}

	err = ValidateNamingScheme(`S<season_num: 3l>`)
	if err == nil {
		t.Errorf("expected error '3l is not a valid arg. must be a valid positive integer'")
	} else {
		t.Log(`S<season_num: 3l>`, "\n\t", err, "\n")
	}

	err = ValidateNamingScheme(`S<season_num: 3`)
	if err == nil {
		t.Errorf("expected error 'reached end of string but still in an unclosed api: <season_num: 3'")
	} else {
		t.Log(`S<season_num: 3`, "\n\t", err, "\n")
	}

	err = ValidateNamingScheme(`<parent-parent:>`)
	if err == nil {
		t.Errorf("expected error 'missing value for token: <parent-parent:>'")
	} else {
		t.Log(`<parent-parent:>`, "\n\t", err, "\n")
	}

	err = ValidateNamingScheme(`E<episode_num: -2>`)
	if err == nil {
		t.Errorf("expected error '-2 is not a valid arg. must be a valid positive integer'")
	} else {
		t.Log(`E<episode_num: -2>`, "\n\t", err, "\n")
	}

	err = ValidateNamingScheme(`<parent-parent:1,>`)
	if err == nil {
		t.Errorf("expected error '1, is not a valid arg. must be two valid positive integers separated by a comma'")
	} else {
		t.Log(`<parent-parent:1,>`, "\n\t", err, "\n")
	}

	err = ValidateNamingScheme(`<p-3:1,2,3>`)
	if err == nil {
		t.Errorf("expected error '1,2,3 is not a valid arg. must be two valid positive integers separated by a comma'")
	} else {
		t.Log(`<p:1,2,3>`, "\n\t", err, "\n")
	}

	err = ValidateNamingScheme(`<parent: '\d+(.*)-.*>`)
	if err == nil {
		t.Errorf(`expected error ' "'\d+(.*)-.*" is not a valid arg. must be a valid regex expression enclosed by two single quotes '`)
	} else {
		t.Log(`<parent: '\d+(.*)-.*'>`, "\n\t", err, "\n")
	}

	err = ValidateNamingScheme(`<self: '[]'>`)
	if err == nil {
		t.Errorf(`expected error ' "[]" is not a valid regex '`)
	} else {
		t.Log(`<self: [>`, "\n\t", err, "\n")
	}

	err = ValidateNamingScheme(`<self: '>`)
	if err == nil {
		t.Errorf(`expected error ' ' is unclosed '`)
	} else {
		t.Log(`<self: '>`, "\n\t", err, "\n")
	}

	err = ValidateNamingScheme(`<self: '   '>`)
	if err == nil {
		t.Errorf(`expected error ' '   ' is empty '`)
	} else {
		t.Log(`<self: '   '>`, "\n\t", err, "\n")
	}

	err = ValidateNamingScheme(`<self: ' '  '>`)
	if err == nil {
		t.Errorf(`expected error ' ' '  ' is empty '`)
	} else {
		t.Log(`<self: ' '  '>`, "\n\t", err, "\n")
	}

	err = ValidateNamingScheme(`<self: 2.0,-3>`)
	if err == nil {
		t.Errorf(`expected error ' -2,-3 is not a valid arg. must be two valid positive integers separated by a comma'`)
	} else {
		t.Log(`<self: 2.0,-3>`, "\n\t", err, "\n")
	}

	err = ValidateNamingScheme(`<self: 'adhd' '>`)
	if err == nil {
		t.Errorf(`expected error ' "" is not a valid regex '`)
	} else {
		t.Log(`<self: 'adhd' '>`, "\n\t", err, "\n")
	}

	err = ValidateNamingScheme(`<self: 6,5>`)
	if err == nil {
		t.Errorf("expected error '6,5 is not a valid range. begin (6) must be less than or equal to end (5)'")
	} else {
		t.Log(`<self: 6,5>`, "\n\t", err, "\n")
	}

	t.Log("------------expects success------------")

	err = ValidateNamingScheme("S<season_num>E<episode_num> - <parent-parent> <parent> <p-3> <self>")
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	} else {
		t.Log("S<season_num>E<episode_num> - <parent-parent> <parent> <p-3> <self>")
	}

	err = ValidateNamingScheme(`S<season_num: 3>E<episode_num: 2> - <parent-parent: 2,3> <parent: '\d+(.*)-.*'> <p-3: '(\d+)'> <self>: 5,5`)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	} else {
		t.Log(`S<season_num: 3>E<episode_num: 2> - <parent-parent: 0,1> <parent: '\d+(.*)-.*'> <p-3: '(\d+)'> <self: 5,5>`)
	}

	err = ValidateNamingScheme(`<p> <p-2>`)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	} else {
		t.Log(`<p>`)
	}
}

func Test_SplitRegexByPipe(t *testing.T) {
	t.Log("------------expects errors------------")
	parts := SplitRegexByPipe(``)
	if parts[0] != "" {
		t.Errorf("expected empty string")
	} else {
		t.Log(`''`, "\n\t", parts, "\n")
	}

	parts = SplitRegexByPipe(`|`)
	if len(parts) != 2 {
		if parts[0] != "" && parts[1] != "" {
			t.Errorf("expected empty string; got '%s', '%s'", parts[0], parts[1])
		}
	} else {
		t.Log(`'|'`, "\n\t", parts)
		for _, part := range parts {
			t.Log("\t", part, "has only one match group?", HasOnlyOneMatchGroup(part))
		}
	}

	parts = SplitRegexByPipe(`a|b|c`)
	if len(parts) != 3 {
		if parts[0] != "a" {
			t.Errorf("expected 'a'; got '%s'", parts[0])
		}
		if parts[1] != "b" {
			t.Errorf("expected 'b'; got '%s'", parts[1])
		}
		if parts[2] != "c" {
			t.Errorf("expected 'c'; got '%s'", parts[2])
		}
	} else {
		t.Log(`'a|b|c'`, "\n\t", parts)
		for _, part := range parts {
			t.Log("\t", part, "has only one match group?", HasOnlyOneMatchGroup(part))
		}
	}

	parts = SplitRegexByPipe(`(a|b|c)`)
	if len(parts) != 1 {
		if parts[0] != "(a|b|c)" {
			t.Errorf("expected '(a|b|c)'; got '%s'", parts[0])
		}
	} else {
		t.Log(`'(a|b|c)'`, "\n\t", parts)
		for _, part := range parts {
			t.Log("\t", part, "has only one match group?", HasOnlyOneMatchGroup(part))
		}
	}

	parts = SplitRegexByPipe(`(a)b(c)|(d|f)e`)
	if len(parts) != 2 {
		if parts[0] != "(a)b(c)" {
			t.Errorf("expected '(a)b(c)'; got '%s'", parts[0])
		}
		if parts[1] != "(d|f)e" {
			t.Errorf("expected '(d|f)e'; got '%s'", parts[1])
		}
	} else {
		t.Log(`'(a)b(c)|(d|f)e'`, "\n\t", parts)
		for _, part := range parts {
			t.Log("\t", part, "has only one match group?", HasOnlyOneMatchGroup(part))
		}
	}

	t.Log("------------expects success------------")

	parts = SplitRegexByPipe(`(a)bc|(d|f)e`)
	if len(parts) != 2 {
		if parts[0] != "(a)bc" {
			t.Errorf("expected '(a)bc'; got '%s'", parts[0])
		}
		if parts[1] != "(d|f)e" {
			t.Errorf("expected '(d|f)e'; got '%s'", parts[1])
		}
	} else {
		t.Log(`'(a)bc|(d|f)e'`, "\n\t", parts)
		for _, part := range parts {
			t.Log("\t", part, "has only one match group?", HasOnlyOneMatchGroup(part))
		}
	}
}

func Test_GenerateNewName(t *testing.T) {
	path := filepath.Clean(`.test_files\Series\Series_seasonal\Season 1\1234567890.mp4`)
	t.Log("------------expects success------------")
	name, err := GenerateNewName(some[string](`S<season_num: 3>E<episode_num: 2> - <parent-parent: '([^_]+)_.*$'> <parent: '([^ ]+) \d+'> <p-3: 'r(.*)$'> <self: 5,6>`),
		2, 1, 3, 2,
		"title", path)
	if err != nil {
		t.Error("expected no error; got", err)
	} else {
		if strings.ReplaceAll(name, `.test_files\Series\Series_seasonal\Season 1\S001E02 - Series Season ies 67.mp4`, "") != "" {
			t.Errorf(`expected '.test_files\Series\Series_seasonal\Season 1\S001E02 - Series Season ies 67.mp4' got '%s'`, name)
		} else {
			t.Log("\n\told:\t\t", filepath.Base(path), "\n\tnaming scheme:\t", `S<season_num: 3>E<episode_num: 2> - <parent-parent: '([^_]+)_.*$'> <parent: '([^ ]+) \d+'> <p-3: 'r(.*)$'> <self: 5,6>`, "\n\tnew:\t\t", name)
		}
	}

	name, err = GenerateNewName(some[string](`<p>`),
		2, 1, 3, 2,
		"title", path)
	if err != nil {
		t.Error("expected no error; got", err)
	} else {
		if strings.ReplaceAll(name, `.test_files\Series\Series_seasonal\Season 1\Season 1.mp4`, "") != "" {
			t.Errorf(`expected '.test_files\Series\Series_seasonal\Season 1\Season 1.mp4' got '%s'`, name)
		} else {
			t.Log("\n\told:\t\t", filepath.Base(path), "\n\tnaming scheme:\t", `<p> <p-2>`, "\n\tnew:\t\t", name)
		}
	}

}
