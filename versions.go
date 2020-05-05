package resource

import (
	"regexp"
	"strconv"
	"time"

	"github.com/google/go-github/github"
)

var defaultTagFilter = "^v?([^v].*)"

type versionParser struct {
	re *regexp.Regexp
}

func newVersionParser(filter string) (versionParser, error) {
	if filter == "" {
		filter = defaultTagFilter
	}
	re, err := regexp.Compile(filter)
	if err != nil {
		return versionParser{}, err
	}
	return versionParser{re: re}, nil
}

func (vp *versionParser) parse(tag string) string {
	matches := vp.re.FindStringSubmatch(tag)
	if len(matches) > 0 {
		return matches[len(matches)-1]
	}
	return ""
}

// getTimestamp returns the last time a give release was modified, including its
// assets.
func getTimestamp(release *github.RepositoryRelease) time.Time {
	var latestTime time.Time
	for _, asset := range release.Assets {
		if asset.CreatedAt != nil && asset.CreatedAt.After(latestTime) {
			latestTime = asset.CreatedAt.Time
		}
		if asset.UpdatedAt != nil && asset.UpdatedAt.After(latestTime) {
			latestTime = asset.UpdatedAt.Time
		}
	}
	if release.PublishedAt != nil && release.PublishedAt.After(latestTime) {
		latestTime = release.PublishedAt.Time
	} else if release.CreatedAt != nil && release.CreatedAt.After(latestTime) {
		latestTime = release.CreatedAt.Time
	}
	return latestTime
}

func versionFromRelease(release *github.RepositoryRelease) Version {
	v := Version{
		ID:        strconv.Itoa(*release.ID),
		Timestamp: getTimestamp(release),
	}
	if release.TagName != nil {
		v.Tag = *release.TagName
	}
	return v
}
