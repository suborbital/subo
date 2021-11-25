now = $(shell date +'%Y-%m-%dT%TZ')
commit = $(shell git rev-parse HEAD)
var_path = github.com/suborbital/subo/subo/release
RELEASE_FLAGS = "-X $(var_path).CommitHash=$(commit)\
 -X $(var_path).BuildTime=$(now)"
