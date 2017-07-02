package github

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"

	"github.com/fatih/color"
	"github.com/mitchellh/packer/common/json"
	"github.com/pkg/errors"
)

type githubRelease struct {
	ID     int64 `json:"id"`
	Assets []struct {
		BrowserDownloadURL string `json:"browser_download_url"`
		Name               string `json:"name"`
	} `json:"assets"`
	HtmlURL string `json:"html_url"`
	TagName string `json:"tag_name"`
}

func getLatestRelease(owner, repo string) (release githubRelease, err error) {
	client := http.DefaultClient
	resp, err := client.Get(fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", owner, repo))
	if err != nil {
		return
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	err = json.Unmarshal(data, &release)
	if err != nil {
		return
	}

	return
}

// IsLatestVersion checks the latest release on github and returns an error
// if it does not match version, error contains URL to download latest release
func IsLatestRelease(owner, repo, version string) (err error) {
	r, err := getLatestRelease(owner, repo)
	if err != nil {
		return
	}
	var v []string
	v = append(v, version, strings.TrimPrefix(r.TagName, "v"))
	sort.Strings(v)

	if v[len(v)-1] != version {
		err = errors.New(color.RedString("Running old version, newer version: %s - is available download from\n%s\n", r.TagName, r.HtmlURL))
	}

	return
}
