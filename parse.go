package main

import "strings"

func parseItems(response string) (items []*item, unfinished bool) {
	unfinished = true
	lines := strings.Split(response, "\n")
	var currentItem *item

	for _, line := range lines {
		if path, ok := strings.CutPrefix(line, "=#=#= "); ok {
			path = strings.TrimSpace(path)
			if path == "END" {
				unfinished = false
				break
			}

			currentItem = &item{
				relPath: path,
			}
			items = append(items, currentItem)
		} else if currentItem != nil {
			currentItem.content = append(currentItem.content, []byte(line)...)
			currentItem.content = append(currentItem.content, '\n')
		}
	}

	return
}
