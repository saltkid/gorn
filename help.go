package main

import "fmt"

func WelcomeMsg(version string) {
	fmt.Println("gorn - go rename tool")
	fmt.Println("version:", version)
	fmt.Println("renames series/movies media files based on directory naming and structure")
	fmt.Println("for more usage information, run 'gorn -h'")
	fmt.Println("https://github.com/saltkid/gorn")
}

func Help(flag string) {
	switch flag {
	case "":
		fmt.Println("Basic usage: gorn -r path/to/root")
		fmt.Println("at least one of the three should be present: --root, --series, --movies")
		fmt.Println("\nOptions:")
		HelpHelp(false)
		HelpVersion(false)
		HelpRoot(false)
		HelpSeries(false)
		HelpMovies(false)
		HelpKEN(false)
		HelpSEN(false)
		HelpS0(false)
		HelpNS(false)
	case "-h", "--help":
		HelpHelp(true)
	case "-v", "--version":
		HelpVersion(true)
	case "-r", "--root":
		HelpRoot(true)
	case "-s", "--series":
		HelpSeries(true)
	case "-m", "--movies":
		HelpMovies(true)
	case "-ken", "--keep-ep-nums":
		HelpKEN(true)
	case "-sen", "--starting-ep-num":
		HelpSEN(true)
	case "-s0", "--has-season-0":
		HelpS0(true)
	case "-ns", "--naming-scheme":
		HelpNS(true)
	default:
		fmt.Printf("invalid flag: %s\n\n", flag)
		Help("")
	}
}

func HelpHelp(verbose bool) {
	fmt.Printf("%-60s%s", "  [--help | -h] <flag>",
		"Show this help message if no flag is specified. Shows the specific help message if a flag is specified.\n")
	if verbose {
		fmt.Println("  example: gorn -h")
		fmt.Println("  example: gorn -h --naming-scheme (this will give more specific help on the naming scheme flag)")
		fmt.Println("  example: gorn -h --keep-ep-num (this will give more specific help on the keep-ep-num flag)")
	}
}

func HelpVersion(verbose bool) {
	fmt.Printf("%-60s%s", "  [--version | -v]",
		"Show version\n")
	if verbose {
		fmt.Println("  example: gorn -v")
		fmt.Println("  version:", version)
	}
}

func HelpRoot(verbose bool) {
	fmt.Printf("%-60s%s", "  [--root | -r] path/to/root",
		"Root directory containing series root and movie root\n\n")
	if verbose {
		fmt.Printf("  example: gorn -r /path/to/root\n\n")
		fmt.Println("  Root directory should contain series and movie roots where each root should contain series and movie entries respectively.")
		fmt.Println("  example root directory:")
		fmt.Println("  root")
		fmt.Println("  |__series")
		fmt.Println("  |  |__series title")
		fmt.Println("  |     |__media files, extra dirs, etc...")
		fmt.Println("  |")
		fmt.Println("  |__movies")
		fmt.Println("     |__movie title")
		fmt.Println("        |__media files, extra dirs, etc...")
	}
}

func HelpSeries(verbose bool) {
	fmt.Printf("%-60s%s", "  [--series | -s] path/to/series/root",
		"Series directory containing series entries\n\n")
	if verbose {
		fmt.Println("  example: gorn -s /path/to/series/root")
		fmt.Println("\n  Series directory should contain series entries.")
		fmt.Println("  example series directory:")
		fmt.Println("  series root")
		fmt.Println("  |__series title")
		fmt.Println("  |  |__media files, extra dirs, etc...")
		fmt.Println("  |")
		fmt.Println("  |__series entry 2")
		fmt.Println("  |  |__media files, extra dirs, etc...")
		fmt.Println("  |")
		fmt.Println("  |__series entry 3")
		fmt.Println("     |__media files, extra dirs, etc...")
	}
}

func HelpMovies(verbose bool) {
	fmt.Printf("%-60s%s", "  [--movies | -m] path/to/movies/root",
		"Movies directory containing movie entries\n\n")
	if verbose {
		fmt.Println("  example: gorn -m /path/to/movies/root")
		fmt.Println("\n  Movies directory should contain movie entries.")
		fmt.Println("  example movies directory:")
		fmt.Println("  movies root")
		fmt.Println("  |__movie title")
		fmt.Println("  |  |__media files, extra dirs, etc...")
		fmt.Println("  |")
		fmt.Println("  |__movie entry 2")
		fmt.Println("  |  |__media files, extra dirs, etc...")
		fmt.Println("  |")
		fmt.Println("  |__movie entry 3")
		fmt.Println("     |__media files, extra dirs, etc...")
	}
}

func HelpKEN(verbose bool) {
	fmt.Printf("%-60s%s", "  [--keep-ep-num | -ken] <all yes/no/default | var>",
		"Keep original episode numbers in filename based on common naming patterns\n\n")
	if verbose {
		fmt.Println("  common naming patterns taken into account are:")
		fmt.Println("    S01E02     |  S03.E04  | S05_E06 | S07-E08 | S09xE10 | S11 E12")
		fmt.Println("    01.02      |   03_04   |  05-06  |  07x08  | 09 10 ")
		fmt.Println("    Episode 01 | Episode02 |  EP03   |  EP-04  | E_05 | EP.06")
		fmt.Println("\n  '.', '-', 'x', '_', and ' ' are valid season-episode separators.")
		fmt.Println("  NOTE: This is not how the default naming scheme looks like in gorn. These common naming cases are just here to read the episode number from the filename.")
		fmt.Println("        second number is episode")
		fmt.Println("        if no common naming pattern is found, the file will not be renamed.")
		fmt.Println("\n  examples: gorn -ken all yes")
		fmt.Println("            gorn -ken all no")
		fmt.Println("            gorn -ken all default")
		fmt.Println("            gorn -ken var")
	}
}

func HelpSEN(verbose bool) {
	fmt.Printf("%-60s%s", "  [--starting-ep-num | -sen] <all int/default | var>",
		"Set the starting episode number in renaming.\n")
	if verbose {
		fmt.Println("\n  This can be useful if episodes are in absolute order but in different season directories for separation")
		fmt.Println("  User can specify different starting episode number for each of those seasons")
		fmt.Println("\n  examples: gorn -sen all 1")
		fmt.Println("            gorn -sen all 25")
		fmt.Println("            gorn -sen all default")
		fmt.Println("            gorn -sen var")
	}
}

func HelpS0(verbose bool) {
	fmt.Printf("%-60s%s", "  [--has-season-0 | -s0] <all yes/no/default | var>",
		"Treat extras/specials/OVA/etc directory as season 0\n")
	if verbose {
		fmt.Println("\n  Note that if this is set, there must be only one specials/extras/OVA directory under a series entry")
		fmt.Println("\n  This is more useful if specified at the series entry level by doing")
		fmt.Println("  'gorn -r path/to/root -s0 var'")
		fmt.Println("  This will let gorn prompt the user at: per series type level and per series entry level")
		fmt.Println("  if var is inputted at the per series type level, it will prompt the user at per series entry level which is where this flag will be most useful")
		fmt.Println("\n  examples: gorn -s0 all yes")
		fmt.Println("            gorn -s0 all no")
		fmt.Println("            gorn -s0 all default")
		fmt.Println("            gorn -s0 var")
	}
}

func HelpNS(verbose bool) {
	fmt.Printf("%-60s%s", "  [--naming-scheme | -ns] <naming-scheme>/default/var",
		"Change the naming scheme\n")
	if verbose {
		fmt.Println("\n  examples: gorn -ns default")
		fmt.Println(`            gorn -ns "S<season_num>E<episode_num> <parent: 1,5> <parent-parent: '_(\d+)_'> <p-3: 2,5> [<self: '\.(\w+)$'>]"`)
		fmt.Println("\n  Naming Scheme APIs:")
		fmt.Println("    1. <season_num>")
		fmt.Println("       represents the season number which is based on series type, and directory structure and naming")
		fmt.Println("       additional option for season num is padding with 0s")
		fmt.Println(`         "<season_num: 2>" pads the result to 2 digits`)
		fmt.Println(`         "<season_num: 3>" pads the result to 3 digits, etc...`)
		fmt.Println("\n    2. <episode_num>")
		fmt.Println("       represents the episode number which is either read from the filename or generated based on the `--keep-ep-nums` and `--starting-ep-num` flags")
		fmt.Println("       additional option for episode num is padding with 0s just like `<season_num>`")
		fmt.Println("\n    3. <parent> | <p>")
		fmt.Println("       represents the parent directory of the media file. if no option was specified, it will copy the whole name of the parent directory")
		fmt.Println(`       additional option is to select characters from the parent directory name.`)
		fmt.Println(`         range: "<parent: 0,3>" which will copy the first 4 characters of the parent directory name`)
		fmt.Println(`         regex: "<parent: 'S(\d+)'>" which will copy the capture group "(\d+)" that is prepended by "S" from the parent directory name`)
		fmt.Println()
		fmt.Println("         Notes on regex:")
		fmt.Println(`           it can only have one capture group per part`)
		fmt.Println(`           each part is separated by "|"`)
		fmt.Println(`             ie. "S(\d+)|E(\d+)" is valid. It has one capture group per part and has 2 parts`)
		fmt.Println(`             ie. "S(E|\d+)" has one capture group and one part.`)
		fmt.Println(`                 "|" inside parenthesis does not count as a part separator. only "|" outside parenthesis is part separator`)
		fmt.Println(`             ie. "'S(E)(\d+)|S(\d+)" is invalid since the first part has 2 capture groups, even if the second part has only 1 capture group`)
		fmt.Println()
		fmt.Println(`     another additional option is going above just the parent of the current directory.`)
		fmt.Println(`         "<parent-parent>" which will copy the parent of the parent directory`)
		fmt.Println(`         "<parent-parent: 0,4>" which will copy the first 4 characters of the parent of the parent directory`)
		fmt.Println(`         "<p>": short form. "<p>" is equivalent to "<parent>" in every way`)
		fmt.Println(`         "<p-2>": you can specify how much further up the directory tree you want to go by appending a number`)
		fmt.Println(`         "<p-2: _(\d+)_>" is equivalent to "<parent-parent: _(\d+)_>" in every way`)
		fmt.Println("\n    4. <self>")
		fmt.Println("       same as parent but instead of being based on the parent directory name, it is based on the name of the media file before renaming it")
		fmt.Println("       additional options are the same as well except for `<p-number>`. self has no short form")
	}
}
