# Rename movies and series (wip)
### current progress
- [x] read directories and files
- [ ] categorize series by type
- [ ] categorize movies by type
- [ ] rename series by type
- [ ] rename movies by type
- [ ] cli commands
- [ ] custom naming scheme
- [ ] add metadata (mediainfo)
- [ ] parallelize renaming


Renames your movies and series based on directory naming and structure.

you can choose to fully replace the filename
- `a random filename.mkv` --> `S01E01 <series title>.mkv`
- `another filename.mkv` --> `S01E02 <series title>.mkv`

or keep a part of the filename (ie episode title) while only replacing the season and episode number that is already present, just reformatting it. 
- `some series name season 1 episode 1 - episode title.mkv` --> `some series name S01E01 - episode title.mkv`
- `some series name season 1 episode 3 - episode title.mkv` --> `some series name S01E03 - episode title.mkv`

this can be useful if you only have episodes that are canon; aka you don't have filler episodes

# Series / TV Shows
Series contain episodes which may be under a season. The filename of an episode number can be the ff:
1. `S01E01`, `S01 E01`, `S1E1`, `S100 E100`, `S01.E01`, `S01_E04`,  - *default for episodes*
2. `[0x1]`, `[00x11]` - *default for movies/specials in a series*
3. `Season 1 Episode 1`, `Season 1 Ep 1`
4. `EP08`, `E09`
5. your own custom naming scheme *(which can be based on your parent directories)*
    - `S<season_num>E<episode_num> - <parent_parent_dir> <parent_dir> something static`
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
S<season_num>E<episode_num> <parent_dir>.<ext>
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
    |__ The End of Evangelion
        |__ Neon Genesis Evangelion The End of Evangelion [1x27].mkv
```
default formatting
```
episodes: S<season_num>E<episode_num> <parent_dir>.<ext>
movies: <parent_parent_dir> <parent_dir>.<ext>
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
    |__ Season 2
        |__ S02E01 Mob Psycho 100.mkv
        |__ S02E02 Mob Psycho 100.mkv
        |__ ...
        |__ S02EXX Mob Psycho 100.mkv
```
default formatting
```
episodes: S<season_num>E<episode_num> <parent_parent_dir>.<ext>
movies: <parent_parent_dir> <parent_dir>.<ext>
```
### 4. multiple season with movie/s
directory input
```
<series root dir>
|__ <series name>
    |__ <special name>
        |__ filename.mkv
    |__ <season name>
    |   |__ filename.mkv
    |   |__ filename2.mkv
    |   |__ ...
    |   |__ some other filename.mkv
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
        |__ Fruits Basket Prelude [0x1]
    |__ Season 1
    |   |__ S01E01 Fruits Basket.mkv
    |   |__ S01E02 Fruits Basket.mkv
    |   |__ ...
    |   |__ S01EXX Fruits Basket.mkv
    |__ Season 2
        |__ S02E01 Fruits Basket.mkv
        |__ S02E02 Fruits Basket.mkv
        |__ ...
        |__ S02EXX Fruits Basket.mkv
```
default formatting
```
episodes: S<season_num>E<episode_num> <parent_parent_dir>.<ext>
movies: <parent_parent_dir> <parent_dir>.<ext>
```
* note: `[0x1]` needs to be added manually since this **gorn** does not scrape data off tmdb/tvdb.
### 5. named seasons with or without movies
* note: the `01.` before the season name is important to determine order

directory input
```
<series root dir>
|__ <series name>
    |__ 01. <season name>
    |   |__ filename.mkv
    |   |__ filename2.mkv
    |   |__ ...
    |   |__ some other filename.mkv
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
    |__ 01. Battle Tendency
    |   |__ S02E01 JoJos Bizzare Adventure Battle Tendency.mkv
    |   |__ S02E02 JoJos Bizzare Adventure Battle Tendency.mkv
        |__ ...
    |   |__ S02EXX JoJos Bizzare Adventure Battle Tendency.mkv
    |__ 02. Stardust Crusaders
        |__ S03E01 JoJos Bizzare Adventure Stardust Crusaders.mkv
        |__ S03E02 JoJos Bizzare Adventure Stardust Crusaders.mkv
        |__ ...
        |__ S03EXX JoJos Bizzare Adventure Stardust Crusaders.mkv
```
default formatting
```
episodes: S<season_num>E<episode_num> <parent_parent_dir> <parent_dir>.<ext>
movies: <parent_parent_dir> <parent_dir>.<ext>
```