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
### 1. --help
Shows the simple help message. If user inputted a flag right after `--help`, it will show the detailed help message for that specific flag.

`-h` is the short form

Example: `gorn -h --naming-scheme`, `gorn --help --help`
### 2. --version
Shows the welcome message along with the version.

`-v` is the short form

Example: `gorn -v`, `gorn --version`
### 1. --keep-ep-num
By default, episode numbers are padded to 2 digits and will start at 01. These are automatically generated and renames the files based on natural sorting.

`--keep-ep-nums all no` is the default behavior if the flag is not present.

If `--keep-ep-nums` flag is present or user inputted `--keep-ep-nums all yes`, gorn will keep the original episode numbers in the filename based on common naming patterns.

if none was found in the filename, it will not rename for that specific file. This can be useful if you only have episodes that are canon, aka you don't have filler episodes, so you want to keep the episode number already in the filename.

If user inputted `--keep-ep-nums var`, gorn will ask for user input whether or not to keep the episode numbers again for each series entry.

`-ken` is the short form

### 2. --starting-ep-num
By default, episode numbers are padded to 2 digits and will start at 01. You can specify a different starting number to start at by `--starting-ep-num all <num>`.

If user inputted `--starting-ep-num var`, gorn will ask for user input again on what starting episode number to start at for each series entry.

`-sen` is the short form

### 3. --has-season-0
By default, the media files in specials/extras directory under a series entry are not renamed. `--has-season-0 all no` is the default behavior if the flag is not present. This will ignore the specials/extras directory.

If the flag `--has-season-0` or user inputted `--has-season-0 all yes`, gorn will rename the files in the specials/extras directory under a series entry, treating it as the *season 0* of the series entry.

*Note that there must be ***ONE*** special/extras directory under the series entry. If there are multiple, it won't rename the files and inform the user.*

If user inputted `--has-season-0 var`, gorn will ask for user input again on whether or not to rename the files in the specials/extras directory as season 0 under a series entry.

`-s0` is the short form

### 4. --naming-scheme
By default, gorn will rename the files differently based on the type of media. User can override this by `--naming-scheme all "<scheme>"` or `--naming-scheme var`.

`all "<scheme>"` overrides the naming scheme for all media files regardless of type (series only; movies will ignore these additional options)

`-ns` is the short form

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
below are guides on how to structure directories based on media type. provided also are the default naming schemes with a sample output. 

# Series / TV Shows
Series contain episodes which may be under a season. The filename of an episode number can be the ff:
1. `S01E01`, `S01 E01`, `S1E1`, `S100 E100`, `S01.E01`, `S01_E04`,  - *default for episodes*
2. `[0x1]`, `[00x11]` - *default for movies/specials in a series*
3. `Season 1 Episode 1`, `Season 1 Ep 1`
4. `EP08`, `E09`
5. your own custom naming scheme
    - `S<season_num>E<episode_num> - <parent-parent> <parent> something static`
    - output:
    - `S01E01 - Fruits Basket Season 1 something static`
    - `S01E02 - Fruits Basket Season 1 something static`

## current valid directory structures
### 1. single season no movie/s
directory input
```
<series root dir>
|__ <series name>
    |__ filename.mkv
    |__ filename2.mkv
    |__ ...
    |__ some other filename.mkv
```
sample output
```
Series
|__ Nichijou
    |__ S01E01 Nichijou.mkv
    |__ S01E02 Nichijou.mkv
    |__ ...
    |__ S01EXX Nichijou.mkv
```
default formatting
```
S<season_num>E<episode_num> <parent>
```

### 2. single season with movie/s
directory input
```
<series root dir>
|__ <series name>
    |__ <series name>
    |   |__ filename.mkv
    |   |__ filename2.mkv
    |   |__ ...
    |   |__ some other filename.mkv
    |
    |__ <movie name>
        |__ some filename.mkv
```
sample output
```
Series
|__ Neon Genesis Evangelion
    |__ Neon Genesis Evangelion
    |   |__ S01E01 Neon Genesis Evangelion.mkv
    |   |__ S01E02 Neon Genesis Evangelion.mkv
    |   |__ ...
    |   |__ S01EXX Neon Genesis Evangelion.mkv
    |
    |__ The End of Evangelion
        |__ Neon Genesis Evangelion The End of Evangelion [1x27].mkv
```
default formatting
```
episodes: S<season_num>E<episode_num> <parent>
movies: <parent-parent> <parent>
```
* note: `[1x27]` needs to be added manually since this **gorn** does not scrape data off tmdb/tvdb. 

### 3. multiple season no movie/s
directory input
```
<series root dir>
|__ <series name>
    |__ <season name>
    |   |__ filename.mkv
    |   |__ filename2.mkv
    |   |__ ...
    |   |__ some other filename.mkv
    |
    |__ <season name>
        |__ filename.mkv
        |__ filename2.mkv
        |__ ...
        |__ some other filename.mkv

```
sample output
```
Series
|__ Mob Psycho 100
    |__ Season 1
    |   |__ S01E01 Mob Psycho 100.mkv
    |   |__ S01E02 Mob Psycho 100.mkv
    |   |__ ...
    |   |__ S01EXX Mob Psycho 100.mkv
    |
    |__ Season 2
        |__ S02E01 Mob Psycho 100.mkv
        |__ S02E02 Mob Psycho 100.mkv
        |__ ...
        |__ S02EXX Mob Psycho 100.mkv
```
default formatting
```
episodes: S<season_num>E<episode_num> <parent-parent>
movies: <parent-parent> <parent>
```
### 4. multiple season with movie/s
directory input
```
<series root dir>
|__ <series name>
    |__ <special name>
    |   |__ filename.mkv
    |
    |__ <season name>
    |   |__ filename.mkv
    |   |__ filename2.mkv
    |   |__ ...
    |   |__ some other filename.mkv
    |
    |__ <season name>
        |__ filename.mkv
        |__ filename2.mkv
        |__ ...
        |__ some other filename.mkv

```
sample output
```
Series
|__ Fruits Basket
    |__ Prelude
    |   |__ Fruits Basket Prelude [0x1]
    |
    |__ Season 1
    |   |__ S01E01 Fruits Basket.mkv
    |   |__ S01E02 Fruits Basket.mkv
    |   |__ ...
    |   |__ S01EXX Fruits Basket.mkv
    |
    |__ Season 2
        |__ S02E01 Fruits Basket.mkv
        |__ S02E02 Fruits Basket.mkv
        |__ ...
        |__ S02EXX Fruits Basket.mkv
```
default formatting
```
episodes: S<season_num>E<episode_num> <parent-parent>
movies: <parent-parent> <parent>
```
* note: `[0x1]` needs to be added manually since this **gorn** does not scrape data off tmdb/tvdb.
### 5. named seasons with or without movies
* note: the `01. title` before the season name is important to determine order
    * `.` after digits can be `-` or `_`, and can be separated by spaces: `02 - title` `03__title`

directory input
```
<series root dir>
|__ <series name>
    |__ 01. <season name>
    |   |__ filename.mkv
    |   |__ filename2.mkv
    |   |__ ...
    |   |__ some other filename.mkv
    |
    |__ 02. <season name>
        |__ filename.mkv
        |__ filename2.mkv
        |__ ...
        |__ some other filename.mkv

```
sample output
```
Series
|__ JoJos Bizzare Adventure
    |__ 01. Phantom Blood
    |   |__ S01E01 JoJos Bizzare Adventure Phantom Blood.mkv
    |   |__ S01E02 JoJos Bizzare Adventure Phantom Blood.mkv
    |   |__ ...
    |   |__ S01EXX JoJos Bizzare Adventure Phantom Blood.mkv
    |
    |__ 01. Battle Tendency
    |   |__ S02E01 JoJos Bizzare Adventure Battle Tendency.mkv
    |   |__ S02E02 JoJos Bizzare Adventure Battle Tendency.mkv
    |   |__ ...
    |   |__ S02EXX JoJos Bizzare Adventure Battle Tendency.mkv
    |
    |__ 02. Stardust Crusaders
        |__ S03E01 JoJos Bizzare Adventure Stardust Crusaders.mkv
        |__ S03E02 JoJos Bizzare Adventure Stardust Crusaders.mkv
        |__ ...
        |__ S03EXX JoJos Bizzare Adventure Stardust Crusaders.mkv
```
default formatting
```
episodes: S<season_num>E<episode_num> <parent-parent> <parent>
movies: <parent-parent> <parent>
```

# Movies
Movies contain a movie file which may be under a movie set. The filename of a movie can be the ff

1. name of `parent_dir` which is most likely the title of the movie - *default for both standalone and movie sets*
2. your own custom naming scheme *(which may or may not be based on your parent directories)*
    - `<parent-parent> - <parent> something static`
    - `Rebuild of Evangelion - Evangelion 1.0 You are (Not) Alone something static`
    - `Rebuild of Evangelion - Evangelion 2.0 You can (Not) Advance something static`

## current valid directory structures
### 1. standalone movies
directory input
```
<movies root dir>
|__ <movie name>
    |__ filename.mkv
```
sample output
```
Movies
|__ Akira
    |__ Akira.mkv
```
default formatting
```
<parent>
```

### 2. movie sets
directory input
```
<movies root dir>
|__ <movie set name>
    |__ <movie name>
    |   |__ filename.mkv
    |
    |__ <movie name>
    |   |__ filename.mkv
    |
    |__ ...
    |
    |__ <movie name>
        |__ filename.mkv
```
sample output
```
Movies
|__ Rebuild of Evangelion
    |__ Evangelion 1.0 - You Are (Not) Alone
    |   |__ Evangelion 1.0 - You Are (Not) Alone.mkv
    |
    |__ Evangelion 2.0 - You Can (Not) Advance
    |   |__ Evangelion 2.0 - You Can (Not) Advance.mkv
    |
    |__ Evangelion 3.0 You Can (Not) Redo
    |   |__ Evangelion 3.0 You Can (Not) Redo.mkv
    |
    |__ Evangelion 3.0+1.0 Thrice Upon a Time
        |__ Evangelion 3.0+1.0 Thrice Upon a Time.mkv
```
default formatting
```
<parent>
```