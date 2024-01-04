### current progress (to v1.0): see [this issue](https://github.com/saltkid/gorn/issues/1)
___ 
# Overview
Renames your movies and series based on directory naming and structure. Note that you still have to rename directories, just not the individual media files themselves. This is for easier metadata scraping when using jellyfin, kodi, plex, etc.

# Prerequisites
A directory to pass to the cli.
have any number of these (any combination of these will work too):
- a root directory containing series directories and/or movie directories
```
<root dir>
|__ <series root dir 1>
|   |__ <series entry 1>
|   |   |__ ...
|   |
|   |__ <series entry 2>
|       |__ ...
|
|__ <movie root dir 1>
    |__ <movie entry 1>
    |   |__ ...
    |
    |__ <movie entry 2>
        |__ ...

where ... may mean media files or subdirectories (like extras, specials, subs, etc)
```
- a series root directory containing series entries
```
<series root dir 1>
|__ <series entry 1>
|   |__ ...
|
|__ <series entry 2>
    |__ ...
```
- or a movie root directory containing movie entries
```
<movie root dir 1>
    |__ <movie entry 1>
    |   |__ ...
    |
    |__ <movie entry 2>
        |__ ...
```

# Usage
```
gorn --root path/to/root/dir
```
Renames all series and movies in the root directory based on directory naming and structure. The episode numbers per entry/season/part will be padded to 2 digits and will start at 01. movies will be named based on the directory they are in. naming scheme will vary depending on the type of media which gorn will detect:
- standalone movie, movie set
- single season, multiple seasons, and named parts/seasons. all with or without movies

```
gorn -r path/to/root/dir --root path/to/another/root/dir
```
User can specify multiple root directories to rename.

```
gorn --series path/to/series/root/dir --movies path/to/movies/root/dir -s path/to/another/series/root/dir -m path/to/another/movies/root/dir
```
User can specify series and movie root dirs separately, can specify only one of either, and can specify any number of dirs. Other than that, it shares the same default renaming behavior as specifying a root dir
___
## Additional Options
Some of these flags have default values if they are not specified.
This is different from the default value if a flag is specified but without a value.

- if `--keep-ep-nums` is not specified, it still has a default value (`all default`).
- if `--naming-scheme` is specified without a value, it will have a default value (`all yes`).

### 1. --help
- **short form:** `-h`
- **values:** `<other flags>`
- **default:** none

Shows the simple help message. If user inputted a flag right after `--help`, it will show the detailed help message for that specific flag.

`-h` is the short form

Example: `gorn -h --naming-scheme`, `gorn --help --help`
### 2. --version
- **short form:** `-v`
- **values:** none
- **default:** none

Shows the welcome message along with the version.

`-v` is the short form

Example: `gorn -v`, `gorn --version`

### 3. --options
- **short form:** `-o`
- **values:** none
- **default:** none

By default, **gorn** will populate flags 4-7 below with default values. This is on an *all media level*.

`--options` flag will ensure the flags are not populated. Since there will be no values, the user will be prompted to input them either during:
- *per series type level*
- *per series entry level* (if no input on per series type level)

`-o` is the short form

### 4. --keep-ep-num
- **short form:** `-ken`
- **values:** `all yes/no/default | var`
- **default if ommitted:** `all default`
- **default if no value:** `all yes`

By default, episode numbers are padded to 2 digits and will start at 01. These are automatically generated and renames the files based on natural sorting.

If this flag is set to `all yes`, **gorn** will keep the original episode numbers in the filename based on common naming patterns. If none was found in the filename, it will not rename for that specific file.

This can be useful if you only have episodes that are canon, aka you don't have filler episodes, so you want to keep the episode number already in the filename.

If this flag is set to `var`, **gorn** will ask for user input whether or not to keep the episode numbers again:
- *per series type level*
- *per series entry level*

### 5. --starting-ep-num
- **short form:** `-sen`
- **values:** `all <num>/default | var`
- **default if ommitted:** `all default`
- **default if no value:** `all 1`

By default, episode numbers are padded to 2 digits and will start at 01. You can specify a different starting number to start at by `--starting-ep-num all <num>`.

If this flag is set to `var`, **gorn** will ask for user input again on what starting episode number to start at:
- *per series type level*
- *per series entry level*

### 6. --has-season-0
- **short form:** `-s0`
- **values:** `all yes/no/default | var`
- **default if ommitted:** `all default`
- **default if no value:** `all yes`

By default, the media files in specials/extras directory under a series entry will be ignored and are not renamed.

If the flag is set to `all yes`, **gorn** will rename the files in the specials/extras directory under a series entry, treating it as the *season 0* of the series entry.

*Note that if the flag is set to `all yes`, there must be ***ONE*** special/extras directory under the series entry. If there are multiple, it won't rename the files and inform the user.*

*if the flag is set to `all no`, there can be any number of specials/extras directories*

If this flag is set to `var`, **gorn** will ask for user input again on whether or not to rename the files in the specials/extras directory as season 0 under a series entry:
- *per series type level*
- *per series entry level*

### 7. --naming-scheme
- **short form:** `-ns`
- **values:** `all "<scheme>"/default | var`
- **default if ommitted:** `all default`
- **default if no value:** none

By default, **gorn** will rename the files differently based on the type of media. User can override this by `--naming-scheme all "<scheme>"` or `--naming-scheme var`.

`all "<scheme>"` overrides the naming scheme for all media files regardless of type (series only; movies will ignore these additional options)

If the flag is set to `var`, **gorn** will ask for user input again on the naming scheme:
- *per series type level*
- *per series entry level*

### *scheme*

scheme can be composed of any character (as long as its a valid filename) and/or APIs enclosed in <> like:
- `S<season_num>E<episode_num>`
    - *output*: `S01E01`
- `S<season_num>E<episode_num> - <parent-parent> <parent> static text` 
    - *output*: `S01E01 - Fruits Basket Season 1 static text`

## Naming Scheme APIs
Current APIs are:
1. `<season_num>`
    - represents the season number which is based on series type, and directory structure and naming
    - additional option for season num is padding with 0s
        - `<season_num: 2>` which pads the result to 2 digits
        - `<season_num: 3>` which pads the result to 3 digits
        - etc ...

2. `<episode_num>`
    - represents the episode number which is either read from the filename or generated based on the `--keep-ep-nums` and `--starting-ep-num` flags
    - additional option for episode num is padding with 0s just like `<season_num>`

3. `<parent>`
    - represents the parent directory of the media file. if no option was specified, it will copy the whole name of the parent directory

    - additional option for parent is to select the range of characters from the parent directory name
    - it can be:
        - a range of two numbers like `<parent: 0,3>`
            - `<parent: 0,3>` which will copy the first 4 characters of the parent directory name
        - a single number like `<parent: 4>`
            - `<parent: 4>` which will copy the 5th character of the parent directory name
        - a regex expression enclosed in single quotes like `<parent: 'S(\d+)'>`
            - `<parent: 'S(\d+)'>` which will copy the capture group `(\d+)` that is prepended by `S` from the parent directory name. Notes:
                1. it can only have one capture group per part
                2. each part is separated by `|`
                3. ie. `S(\d+)|E(\d+)` is valid. It has one capture group per part and has 2 parts
                4. ie. `S(E|\d+)` has one capture group and one part. `|` inside parenthesis does not count as a part separator. only `|` outside parenthesis is part separator
                5. ie. `'S(E)(\d+)|S(\d+)` is invalid since the first part has 2 capture groups, even if the second part has only 1

    - another additional option is going above just the parent of the current directory.
        - `<parent-parent>` which will copy the parent of the parent directory

        - `<parent-parent: 0,4>` which will copy the first 4 characters of the parent of the parent directory

        - `<parent-parent-parent>` which will copy the parent of the parent of the parent directory

        - `<p>`: short form. `<p>` is equivalent to `<parent>` in every way

        - `<p-2>`: you can specify how much further up the directory tree you want to go by appending a number

        - `<p-2>` is equivalent to `<parent-parent>` in every way

4. `<self>`
    - same as parent but instead of being based on the parent directory name, it is based on the name of the media file before renaming it
    - additional options are the same as well except for `<p-number>`. self has no short form

___
# Root Directory Structure Overview

Root directories should contain series roots and/or movie roots (let's call these subroots). Each subroot should contain series and movie entries respectively.`

sample root directory
```
<root dir>
|__ <series subroot>
|   |__ <series entry>
|   |   |__ ...
|   |
|   |__ <series entry>
|       |__ ...
|
|__ <movie subroot>
|   |__ <movie entry>
|   |   |__ ...
|   |
|   |__ <movie entry>
|       |__ ...
|
|__ <movie subroot>
    |__ <movie entry>
        |__ ...
```
*where `...` may mean media files or subdirectories like extras, specials, subs, etc*

For more information about Subroot (series/movies) Directory Structures, see [this wiki page](https://github.com/saltkid/gorn/wiki/Directory-Structure)