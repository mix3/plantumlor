package plantuml

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type Format int

const (
	PNG Format = iota
	SVG
	EPS
	PDF
	VDX
	XMI
	HTML
	TXT
	UTXT
)

func (f Format) String() string {
	switch f {
	case PNG:
		return "png"
	case SVG:
		return "svg"
	case EPS:
		return "eps"
	case PDF:
		return "pdf"
	case VDX:
		return "vdx"
	case XMI:
		return "xmi"
	case HTML:
		return "html"
	case TXT:
		return "txt"
	case UTXT:
		return "utxt"
	default:
		return "Unknown"
	}
}

func (f Format) Flag() string {
	return fmt.Sprintf("-t%s", f)
}

type PlantUML string

var NilPlantUML = PlantUML("")

func NewPlantUML(path string) (PlantUML, error) {
	if path == "" {
		path = command("which", "plantuml")
	}

	if path == "" {
		return NilPlantUML, fmt.Errorf("%s: command not found", path)
	}

	out, err := exec.Command(path, "-version").Output()
	if err != nil {
		return NilPlantUML, err
	}

	part := make([]byte, 8)
	copy(part, out)

	if string(part) != "PlantUML" {
		return NilPlantUML, fmt.Errorf("%s: not plantuml command", path)
	}

	return PlantUML(path), nil
}

func command(base string, args ...string) string {
	out, _ := exec.Command(base, args...).Output()
	return strings.TrimRight(string(out), "\n")
}

func (p PlantUML) Transfer(input string, format Format) ([]byte, error) {
	os.Setenv("JAVA_TOOL_OPTIONS", "-Djava.awt.headless=true")
	cmd := exec.Command(string(p), "-charset", "utf8", "-q", "-p", format.Flag())
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}
