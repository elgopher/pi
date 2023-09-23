package p8

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type Parser struct {
	currentSectionName string
}

func (p *Parser) Parse(file string) (File, error) {
	f, err := os.Open(file)
	if err != nil {
		return File{}, err
	}

	scanner := bufio.NewScanner(bufio.NewReader(f))
	var version int
	if version, err = p.readHeader(scanner); err != nil {
		return File{}, err
	}

	var sections []Section

	for {
		section, err := p.readSection(scanner)
		sections = append(sections, section)
		if err == io.EOF {
			break
		} else if err != nil {
			return File{}, err
		}
	}

	return File{
		Version:  version,
		Sections: sections,
	}, nil
}

func (p *Parser) readHeader(scanner *bufio.Scanner) (version int, err error) {
	if !scanner.Scan() {
		return 0, fmt.Errorf("error reading first line from p8 file header: %w", scanner.Err())
	}

	firstLine := scanner.Text()

	if firstLine != "pico-8 cartridge // http://www.pico-8.com" {
		return 0, fmt.Errorf("input file is not p8 cartridge file. Header expected")
	}

	if !scanner.Scan() {
		return 0, fmt.Errorf("error reading second line from p8 file header: %w", scanner.Err())
	}

	versionLine := scanner.Text()
	if !strings.HasPrefix(versionLine, "version ") {
		return 0, fmt.Errorf("input file is not p8 cartridge file. Version in header expected, but not found")
	}

	versionString := strings.SplitN(versionLine, " ", 2)[1]

	version, err = strconv.Atoi(versionString)
	if err != nil {
		return 0, fmt.Errorf("input file is not p8 cartridge file. Version number in header expected, but found %s", versionString)
	}

	return version, nil
}

func (p *Parser) readSection(scanner *bufio.Scanner) (Section, error) {
	if p.currentSectionName == "" {
		if !scanner.Scan() {
			return Section{}, fmt.Errorf("error reading section name")
		}
		p.currentSectionName = scanner.Text()
	}

	sectionName := p.currentSectionName

	var lines []string
	for {
		scanOk := scanner.Scan()
		line := scanner.Text()
		if p.isSectionName(line) {
			p.currentSectionName = line
			return Section{
				Name:  sectionName,
				Lines: lines,
			}, nil
		}
		lines = append(lines, line)
		if !scanOk {
			if scanner.Err() == nil {
				return Section{
					Name:  sectionName,
					Lines: lines,
				}, io.EOF
			}
			return Section{}, fmt.Errorf("error readling section line: %w", scanner.Err())
		}
	}

}

func (p *Parser) isSectionName(line string) bool {
	return line == "__sfx__" || line == "__gfx__" || line == "__lua__" || line == "__music__"
}

type File struct {
	Version  int
	Sections []Section
}

type Section struct {
	Name  string
	Lines []string
}
