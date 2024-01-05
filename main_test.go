package main

import (
	"path/filepath"
	"strings"
	"testing"
)

func Test_ParseArgs(t *testing.T) {
	t.Log("------------expects errors------------")

	command := []string{"--root", "--series", "--movies"}
	_, err := ParseArgs(command)
	if err == nil {
		t.Errorf("expected error 'missing root dir path'")
	} else {
		t.Log("--root", "--series", "--movies", "\n\t", err, "\n")
	}

	command = []string{"-s", "--series", "-r", "--root", "-m", "--movies"}
	_, err = ParseArgs(command)
	if err == nil {
		t.Errorf("expected error 'missing series dir path'")
	} else {
		t.Log("-s", "--series", "-r", "--root", "-m", "--movies", "\n\t", err, "\n")
	}

	command = []string{"-r", "-s", "-m"}
	_, err = ParseArgs(command)
	if err == nil {
		t.Errorf("expected error 'missing root dir path'")
	} else {
		t.Log("-r", "-s", "-m", "\n\t", err, "\n")
	}

	command = []string{"-m", "--movies", "-r", "--root", "-s", "--series"}
	_, err = ParseArgs(command)
	if err == nil {
		t.Errorf("expected error 'missing movies dir path'")
	} else {
		t.Log("-m", "--movies", "-r", "--root", "-s", "--series", "\n\t", err, "\n")
	}

	command = []string{"--root", "./test_files", "./test_files"}
	_, err = ParseArgs(command)
	if err == nil {
		t.Errorf("expected error 'multiple values for one flag is not allowed'")
	} else {
		t.Log("--root", "./test_files", "./test_files", "\n\t", err, "\n")
	}

	command = []string{"-m", "./test_files", "./test_files"}
	_, err = ParseArgs(command)
	if err == nil {
		t.Errorf("expected error 'multiple values for one flag is not allowed'")
	} else {
		t.Log("-m", "./test_files", "./test_files", "\n\t", err, "\n")
	}

	command = []string{"-s", "./test_files", "./test_files"}
	_, err = ParseArgs(command)
	if err == nil {
		t.Errorf("expected error 'multiple values for one flag is not allowed'")
	} else {
		t.Log("-s", "./test_files", "./test_files", "\n\t", err, "\n")
	}

	command = []string{"./test_files"}
	_, err = ParseArgs(command)
	if err == nil {
		t.Errorf("expected error 'not valid flag'")
	} else {
		t.Log("./test_files", "\n\t", err, "\n")
	}

	command = []string{"-r", "./test_files", "--season-0", "all", "1"}
	_, err = ParseArgs(command)
	if err == nil {
		t.Errorf("expected error '1 is not a valid arg. must be yes or no'")
	} else {
		t.Log("-r", "./test_files", "--season-0", "all", "1", "\n\t", err, "\n")
	}

	command = []string{"-s0", "all", "yes"}
	_, err = ParseArgs(command)
	if err == nil {
		t.Errorf("expected error 'missing root/series/movies dir path'")
	} else {
		t.Log("-s0", "all", "yes", "\n\t", err, "\n")
	}

	command = []string{"-r", "./test_files", "--season-0", "yes"}
	_, err = ParseArgs(command)
	if err == nil {
		t.Errorf("expected error 'yes is not a valid arg. must be all or var'")
	} else {
		t.Log("-r", "./test_files", "--season-0", "yes", "\n\t", err, "\n")
	}

	command = []string{"-ken", "all", "yes"}
	_, err = ParseArgs(command)
	if err == nil {
		t.Errorf("expected error 'missing root/series/movies dir path'")
	} else {
		t.Log("-ken", "all", "yes", "\n\t", err, "\n")
	}

	command = []string{"-r", "./test_files", "-ken", "all", "1"}
	_, err = ParseArgs(command)
	if err == nil {
		t.Errorf("expected error '1 is not a valid arg. must be yes or no'")
	} else {
		t.Log("-r", "./test_files", "-ken", "all", "1", "\n\t", err, "\n")
	}

	command = []string{"-r", "./test_files", "--season-0", "-s0"}
	_, err = ParseArgs(command)
	if err == nil {
		t.Errorf("expected error 'only one of --season-0 and -s0 is allowed'")
	} else {
		t.Log("-r", "./test_files", "--season-0", "-s0", "\n\t", err, "\n")
	}

	command = []string{"-r", "./test_files", "-ken", "all", "yes"}
	_, err = ParseArgs(command)
	if err == nil {
		t.Errorf("expected error 'all is not a valid arg. must be yes no default or var'")
	} else {
		t.Log("-r", "./test_files", "-ken", "yes", "\n\t", err, "\n")
	}

	command = []string{"-r", "./test_files", "--keep-ep-nums", "-ken"}
	_, err = ParseArgs(command)
	if err == nil {
		t.Errorf("expected error 'multiple keep-ep-nums flags'")
	} else {
		t.Log("-r", "./test_files", "--keep-ep-nums", "-ken", "\n\t", err, "\n")
	}

	command = []string{"-sen", "all", "yes"}
	_, err = ParseArgs(command)
	if err == nil {
		t.Errorf("expected error 'missing root/series/movies dir path'")
	} else {
		t.Log("-sen", "all", "yes", "\n\t", err, "\n")
	}

	command = []string{"-r", "./test_files", "-sen", "all", "yes"}
	_, err = ParseArgs(command)
	if err == nil {
		t.Errorf("expected error '1 is not a valid arg. must be a valid positive integer'")
	} else {
		t.Log("-r", "./test_files", "-sen", "all", "yes", "\n\t", err, "\n")
	}

	command = []string{"-r", "./test_files", "-sen", "all", "1"}
	_, err = ParseArgs(command)
	if err == nil {
		t.Errorf("expected error 'all is not a valid arg. must be int default or var'")
	} else {
		t.Log("-r", "./test_files", "-sen", "yes", "\n\t", err, "\n")
	}

	command = []string{"-r", "./test_files", "--starting-ep-num", "-sen"}
	_, err = ParseArgs(command)
	if err == nil {
		t.Errorf("expected error 'multiple starting-ep-num flags'")
	} else {
		t.Log("-r", "./test_files", "--starting-ep-num", "-sen", "\n\t", err, "\n")
	}

	command = []string{"--root", "./test_files", "-r", "./test_files"}
	_, err = ParseArgs(command)
	if err == nil {
		t.Errorf("expected error 'root directory ./test_files is a duplicate'")
	} else {
		t.Log("--root", "./test_files", "-r", "./test_files", "\n\t", err, "\n")
	}

	command = []string{"--root", `.\test_files`, "-s", `.\test_files\Series`, "-m", "./test_files/Movies"}
	_, err = ParseArgs(command)
	if err == nil {
		t.Errorf("expected error 'series directory ./test_files/Series is a subdirectory of root directory ./test_files'")
	} else {
		t.Log("--root", "./test_files", "-s", "./test_files/Series", "-m", "./test_files/Movies", "\n\t", err, "\n")
	}

	command = []string{"-r", "./test_files", "-sen", "1", "--starting-ep-num", "2"}
	_, err = ParseArgs(command)
	if err == nil {
		t.Errorf("expected error 'multiple starting-ep-num flags'")
	} else {
		t.Log("-r", "./test_files", "--naming-scheme", "S01E01", "\n\t", err, "\n")
	}
	t.Log("------------expects success------------")

	command = []string{"--root", "./test_files", "-s0", "yes"}
	_, err = ParseArgs(command)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	} else {
		t.Log("--root", "./test_files", "-s0", "yes")
	}

	command = []string{"--root", "./test_files"}
	_, err = ParseArgs(command)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	} else {
		t.Log("--root", "./test_files")
	}

	command = []string{"--root", "./test_files", "-s0"}
	_, err = ParseArgs(command)
	if err != nil {
		t.Errorf("--season-0 can have no value: %s", err)
	} else {
		t.Log("--root", "./test_files", "-s0")
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
