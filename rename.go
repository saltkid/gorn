package main

type Rename interface {
	rename() error
}

type SeriesInfo struct {
	path            string
	series_type     string
	keep_ep_nums    bool
	starting_ep_num int
	seasons         map[int]string
	movies          []string
	has_season_0    bool
	extras_dirs     []string
}

type MovieInfo struct {
	path        string
	movie_type  string
	movies      []string
	extras_dirs []string
}

func (info *SeriesInfo) rename() error {
	return nil
}

func (info *MovieInfo) rename() error {
	return nil
}
