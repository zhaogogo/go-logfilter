package grok

import (
	"bufio"
	"fmt"
	"github.com/rs/zerolog/log"
	"io"
	"net/http"
	"os"
	"path"
	"regexp"
	"strings"
)

type Grok struct {
	p           *regexp.Regexp
	subexpNames []string
	ignoreBlank bool

	patterns     map[string]string
	patternPaths []string
}

func getFiles(filepath string) ([]string, error) {
	if strings.HasPrefix(filepath, "http://") || strings.HasPrefix(filepath, "https://") {
		return []string{filepath}, nil
	}

	fi, err := os.Stat(filepath)
	if err != nil {
		return nil, err
	}

	if !fi.IsDir() {
		return []string{filepath}, nil
	}

	f, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}

	list, err := f.Readdir(-1)
	f.Close()

	if err != nil {
		return nil, err
	}
	files := make([]string, 0)
	for _, l := range list {
		if l.Mode().IsRegular() {
			files = append(files, path.Join(filepath, l.Name()))
		}
	}
	return files, nil
}

func (g *Grok) loadPattern(filename string) {
	var r *bufio.Reader
	if strings.HasPrefix(filename, "http://") || strings.HasPrefix(filename, "https://") {
		resp, err := http.Get(filename)
		if err != nil {
			log.Fatal().Err(err).Msg("load pattern http failed")
		}
		defer resp.Body.Close()
		r = bufio.NewReader(resp.Body)
	} else {
		f, err := os.Open(filename)
		if err != nil {
			log.Fatal().Err(err).Msg("load pattern file failed")
		}
		r = bufio.NewReader(f)
	}
	for {
		line, isPrefix, err := r.ReadLine()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal().Err(err).Msgf("read pattenrs failed from %s", filename)
		}
		if isPrefix {
			log.Fatal().Msgf("readline prefix from %s", filename)
		}
		if len(line) == 0 || line[0] == '#' {
			continue
		}
		ss := strings.SplitN(string(line), " ", 2)
		if len(ss) != 2 {
			log.Fatal().Msgf("splited `%s` length !=2", string(line))
		}
		g.patterns[ss[0]] = ss[1]
	}
}

func (g *Grok) loadPatterns() {
	for _, path := range g.patternPaths {
		files, err := getFiles(path)
		if err != nil {
			log.Fatal().Err(err).Msg("build g filter")
		}
		for _, file := range files {
			g.loadPattern(file)
		}
	}
	log.Info().Msgf("gork plugin load patterns:%s", g.patterns)
}

func (g *Grok) translateMatchPattern(match string) string {
	// (:)非捕获组
	p := regexp.MustCompile(`%{\w+?(:\w+?)?}`)
	var r string = ""
	for {
		r = p.ReplaceAllStringFunc(match, g.replaceFunc)
		if r == match {
			return r
		}
		match = r
	}
}

func (g *Grok) replaceFunc(s string) string {
	p := regexp.MustCompile(`%{(\w+?)(?::(\w+?))?}`)

	rst := p.FindAllStringSubmatch(s, -1)
	if len(rst) != 1 {
		log.Fatal().Msgf("sub match in `%s` != 1", s)
	}
	if pattern, ok := g.patterns[rst[0][1]]; ok {
		if rst[0][2] == "" {
			return fmt.Sprintf("(%s)", pattern)
		} else {
			return fmt.Sprintf("(?P<%s>%s)", rst[0][2], pattern)
		}
	} else {
		log.Fatal().Msgf("`%s` could not be found", rst[0][1])
		return ""
	}
}

func (g *Grok) grok(input string) map[string]string {
	rst := make(map[string]string)
	for i, substring := range g.p.FindStringSubmatch(input) {
		if g.subexpNames[i] == "" {
			continue
		}
		if g.ignoreBlank && substring == "" {
			continue
		}
		rst[g.subexpNames[i]] = substring
	}
	return rst
}

func NewGrok(match string, patternPaths []string, ignoreBlank bool) *Grok {
	g := &Grok{
		patternPaths: patternPaths,
		patterns:     make(map[string]string),
		ignoreBlank:  ignoreBlank,
	}
	g.loadPatterns()

	finalPattern := g.translateMatchPattern(match)
	log.Info().Msgf("final match pattern: %s", finalPattern)
	p := regexp.MustCompile(finalPattern)
	g.p = p

	g.subexpNames = p.SubexpNames()
	return g
}
