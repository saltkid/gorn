### current progress (to v1.0): see [this issue](https://github.com/saltkid/gorn/issues/1)
___
1. [Overview](#overview)
2. [Prerequisites](#prerequisites)
3. [Usage](#usage)
    1. [Optional Flags](#optional-flags)
___ 
# Overview
Renames your movies and series based on directory naming and structure. Note that you still have to rename directories, just not the individual media files themselves. This is for easier metadata scraping when using jellyfin, kodi, plex, etc.

# Prerequisites
Have at least one of any of these directories:
1. **root directory containing series roots and/or movie roots (subroots)**
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
*Note that there can be multiple series/movie subroots in the same root directory.*

2. **series subroot directory containing series entries**
```
<series root dir 1>
|__ <series entry 1>
|   |__ ...
|
|__ <series entry 2>
    |__ ...
```
3. **movie subroot directory containing movie entries**
```
<movie root dir 1>
    |__ <movie entry 1>
    |   |__ ...
    |
    |__ <movie entry 2>
        |__ ...
```
For a more detailed explanation of recommended directory structures, different series/movie types depending on structure, see [this wiki page](https://github.com/saltkid/gorn/wiki/Directory-Structure)
___
# Usage
To renames all series and movies in the root directory based on directory naming and structure:
```
gorn --root path/to/root/dir
```

User can specify multiple root/subroot directories to rename:
```
gorn -r path/to/root/dir --root path/to/another/root/dir
```

User can specify series and movie subroot dirs separately. User can also specify multiple subroot dirs. Other than that, it shares the same default renaming behavior as specifying a root dir
```
gorn --series path/to/series/subroot/dir
```
```
gorn --movies path/to/movies/subroot/dir
```
```
gorn -r path/to/root/dir -s path/to/another/series/subroot/dir -m path/to/another/movies/subroot/dir
```
___
## Optional Flags
These are the additional options that can be passed to the cli. For a more detailed explanation, see [this wiki page](https://github.com/saltkid/gorn/wiki/Usage#optional-flags)
1. `--help | -h`
    - **values:** `<other flags>`
2. `--version | -v`
    - **values:** none
3. `--options | -o`
    - **values:** none
4. `--keep-ep-num | -ken`
    - **values:** `all yes/no/default` or `var`
5. `--starting-ep-num | -sen`
    - **values:** `all <num>/default` or `var`
6. `--has-season-0 | -s0`
    - **values:** `all yes/no/default` or `var`
7. `--naming-scheme | -ns`
    - **values:** `all "<scheme>"/default` or `var`

### *scheme*
scheme can be composed of any character (as long as its a valid filename) and/or APIs enclosed in <> like:
- `S<season_num>E<episode_num>`
    - *output*: `S01E01`
- `S<season_num>E<episode_num> - <parent-parent> <parent> static text` 
    - *output*: `S01E01 - Fruits Basket Season 1 static text`

For more information, see [this wiki page](https://github.com/saltkid/gorn/wiki/Usage#naming-scheme-apis)
___

Credits: [@saltkid](https://github.com/saltkid)

License: MIT License
