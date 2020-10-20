package git

import (
	"bytes"
	"strings"
)

func ExtractRepositoryPath(gitUrl string) string {
	if bytes.ContainsRune([]byte(gitUrl), '@') {
		return parseAuthorizedRepositoryPath(gitUrl)
	}
	return parseHTTPRepositoryPath(gitUrl)
}

func parseAuthorizedRepositoryPath(path string) string {
	// http://vink:xxxxx@git.v2k.io/infra/account.git
	if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
		idx := strings.LastIndexByte(path, '@')
		path = path[idx+1:]
		index := strings.IndexByte(path, '/')
		path = path[index+1:]
	}
	if strings.HasPrefix(path, "ssh://") {
		replacer := strings.NewReplacer("ssh://", "", ".git", "", "\n", "")
		path = replacer.Replace(path)
		index := strings.IndexByte(path, '/')
		return path[index+1:]
	}
	{
		replacer := strings.NewReplacer(".git", "", "\n", "")
		path = replacer.Replace(path)
		index := strings.IndexByte(path, ':')
		return path[index+1:]
	}
}

func parseHTTPRepositoryPath(path string) string {
	replacer := strings.NewReplacer("https://", "", "http://", "", ".git", "", "\n", "")
	path = replacer.Replace(path)
	index := strings.IndexByte(path, '/')
	return path[index+1:]
}
