package star

import "strings"

func getKindFilenamesFromProject(prefix string) []string {
	file, scanner := fileScanner("etc/Project.xml")
	defer file.Close()
	var files []string
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if !strings.HasPrefix(line, prefix) {
			continue
		}
		line = strings.TrimSuffix(strings.TrimPrefix(line, prefix), `"`)
		files = append(files, line)
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
	return files
}

func getDtaParmFiles() []string {
	return getKindFilenamesFromProject(`DtaParm="file://`)
}

func getServiceFiles() []string {
	return getKindFilenamesFromProject(`Service="file://`)
}

func getFormatFiles() []string {
	return getKindFilenamesFromProject(`Format="file://`)
}
