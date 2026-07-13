// ABOUTME: Public-CLI regression coverage for manuscript config directory selection.
// ABOUTME: Exercises single-input, multi-input, style, and precedence rules from issue #13.
package manuscript

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestManuscriptConfigDirectorySelection(t *testing.T) {
	root := testProjectRoot(t)
	folio := buildPublicFolio(t, root)

	t.Run("RT-13.1 relative single input uses input directory", func(t *testing.T) {
		project, home := newConfigProject(t)
		writeConfigInput(t, filepath.Join(project, "source", "chapter.md"), true)
		writeConfigFont(t, filepath.Join(project, "script.yaml"), "WrongCwdFont")
		writeConfigFont(t, filepath.Join(project, "source", "script.yaml"), "SingleRelativeFont")

		typst := runPublicFolio(t, folio, project, home, "source/chapter.md", "out.typ")
		assertContains(t, typst, `font: "SingleRelativeFont"`)
		assertNotContains(t, typst, `font: "WrongCwdFont"`)
	})

	t.Run("RT-13.2 absolute single input uses input directory", func(t *testing.T) {
		project, home := newConfigProject(t)
		input := filepath.Join(project, "source", "chapter.md")
		writeConfigInput(t, input, true)
		writeConfigFont(t, filepath.Join(project, "source", "script.yaml"), "SingleAbsoluteFont")

		typst := runPublicFolio(t, folio, project, home, input, "out.typ")
		assertContains(t, typst, `font: "SingleAbsoluteFont"`)
	})

	t.Run("RT-13.3 single input uses local US style config", func(t *testing.T) {
		project, home := newConfigProject(t)
		inputDir := filepath.Join(project, "source")
		writeConfigInput(t, filepath.Join(inputDir, "chapter.md"), true)
		writeFile(t, filepath.Join(inputDir, "script.yaml"), "wordcount: about 1,000 words\nfolio:\n  manuscript:\n    style: british\n    font: SingleBaseFont\n")
		writeConfigFont(t, filepath.Join(inputDir, "script-us.yaml"), "SingleUSFont")

		typst := runPublicFolio(t, folio, project, home, "--style", "us", "source/chapter.md", "out.typ")
		assertContains(t, typst, `font: "SingleUSFont"`)
		assertContains(t, typst, `#place(bottom + center`)
	})

	t.Run("RT-13.4 relative multi-input uses working directory", func(t *testing.T) {
		project, home := newConfigProject(t)
		writeSiblingInputs(t, project)
		writeConfigFont(t, filepath.Join(project, "script.yaml"), "MultiRelativeFont")

		typst := runPublicFolio(t, folio, project, home, "part1/ch01.md", "part2/ch02.md", "out.typ")
		assertContains(t, typst, `font: "MultiRelativeFont"`)
	})

	t.Run("RT-13.5 absolute multi-input uses working directory", func(t *testing.T) {
		project, home := newConfigProject(t)
		writeSiblingInputs(t, project)
		writeConfigFont(t, filepath.Join(project, "script.yaml"), "MultiAbsoluteFont")

		typst := runPublicFolio(t, folio, project, home,
			filepath.Join(project, "part1", "ch01.md"),
			filepath.Join(project, "part2", "ch02.md"),
			"out.typ")
		assertContains(t, typst, `font: "MultiAbsoluteFont"`)
	})

	t.Run("RT-13.6 nested multi-input still uses working directory", func(t *testing.T) {
		project, home := newConfigProject(t)
		writeConfigInput(t, filepath.Join(project, "chapters", "ch01.md"), true)
		writeConfigInput(t, filepath.Join(project, "chapters", "ch02.md"), false)
		writeConfigFont(t, filepath.Join(project, "script.yaml"), "NestedCwdFont")
		writeConfigFont(t, filepath.Join(project, "chapters", "script.yaml"), "WrongNestedFont")

		typst := runPublicFolio(t, folio, project, home, "chapters/ch01.md", "chapters/ch02.md", "out.typ")
		assertContains(t, typst, `font: "NestedCwdFont"`)
		assertNotContains(t, typst, `font: "WrongNestedFont"`)
	})

	t.Run("RT-13.7 multi-input uses working-directory US style config", func(t *testing.T) {
		project, home := newConfigProject(t)
		writeSiblingInputs(t, project)
		writeFile(t, filepath.Join(project, "script.yaml"), "wordcount: about 1,000 words\nfolio:\n  manuscript:\n    style: british\n    font: MultiBaseFont\n")
		writeConfigFont(t, filepath.Join(project, "script-us.yaml"), "MultiUSFont")

		typst := runPublicFolio(t, folio, project, home, "--style", "us", "part1/ch01.md", "part2/ch02.md", "out.typ")
		assertContains(t, typst, `font: "MultiUSFont"`)
		assertContains(t, typst, `#place(bottom + center`)
	})

	t.Run("RT-13.8 multi-input ignores input-adjacent configs", func(t *testing.T) {
		project, home := newConfigProject(t)
		writeSiblingInputs(t, project)
		writeConfigFont(t, filepath.Join(project, "script.yaml"), "SoleRequestFont")
		writeConfigFont(t, filepath.Join(project, "part1", "script.yaml"), "WrongPartOneFont")
		writeConfigFont(t, filepath.Join(project, "part2", "script.yaml"), "WrongPartTwoFont")

		typst := runPublicFolio(t, folio, project, home, "part1/ch01.md", "part2/ch02.md", "out.typ")
		assertContains(t, typst, `font: "SoleRequestFont"`)
		assertNotContains(t, typst, `font: "WrongPartOneFont"`)
		assertNotContains(t, typst, `font: "WrongPartTwoFont"`)
	})

	t.Run("RT-13.9 single-input precedence remains layered", func(t *testing.T) {
		project, home := newConfigProject(t)
		inputDir := filepath.Join(project, "source")
		writeConfigInput(t, filepath.Join(inputDir, "chapter.md"), true)
		writeConfigTitle(t, filepath.Join(home, ".config", "first-folio", "script.yaml"), "Global Title")
		writeConfigTitle(t, filepath.Join(inputDir, "script.yaml"), "Local Title")

		local := runPublicFolio(t, folio, project, home, "source/chapter.md", "local.typ")
		assertContains(t, local, "Local Title")
		assertNotContains(t, local, "Global Title")

		cli := runPublicFolio(t, folio, project, home, "--title", "CLI Title", "source/chapter.md", "cli.typ")
		assertContains(t, cli, "CLI Title")
		assertNotContains(t, cli, "Local Title")
	})

	t.Run("RT-13.10 multi-input precedence remains layered", func(t *testing.T) {
		project, home := newConfigProject(t)
		writeSiblingInputs(t, project)
		writeConfigTitle(t, filepath.Join(home, ".config", "first-folio", "script.yaml"), "Global Multi Title")
		writeConfigTitle(t, filepath.Join(project, "script.yaml"), "Local Multi Title")

		local := runPublicFolio(t, folio, project, home, "part1/ch01.md", "part2/ch02.md", "local.typ")
		assertContains(t, local, "Local Multi Title")
		assertNotContains(t, local, "Global Multi Title")

		cli := runPublicFolio(t, folio, project, home, "--title", "CLI Multi Title", "part1/ch01.md", "part2/ch02.md", "cli.typ")
		assertContains(t, cli, "CLI Multi Title")
		assertNotContains(t, cli, "Local Multi Title")
	})

	t.Run("RT-13.11 multi-input falls back to global and preset", func(t *testing.T) {
		project, home := newConfigProject(t)
		writeSiblingInputs(t, project)
		writeConfigTitle(t, filepath.Join(home, ".config", "first-folio", "script.yaml"), "Global Fallback Title")

		global := runPublicFolio(t, folio, project, home, "part1/ch01.md", "part2/ch02.md", "global.typ")
		assertContains(t, global, "Global Fallback Title")
		assertContains(t, global, `font: "Libertinus Serif"`)

		presetHome := filepath.Join(project, "empty-home")
		if err := os.MkdirAll(presetHome, 0o755); err != nil {
			t.Fatalf("creating empty HOME: %v", err)
		}
		preset := runPublicFolio(t, folio, project, presetHome, "part1/ch01.md", "part2/ch02.md", "preset.typ")
		assertContains(t, preset, "Source Title")
		assertContains(t, preset, `font: "Libertinus Serif"`)
	})
}

func buildPublicFolio(t *testing.T, root string) string {
	t.Helper()
	cmd := exec.Command("go", "build", "-o", filepath.Join(root, "bin", "folio-manuscript"), "./cmd/folio-manuscript")
	cmd.Dir = root
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("building manuscript helper: %v\n%s", err, string(output))
	}
	return filepath.Join(root, "bin", "folio")
}

func newConfigProject(t *testing.T) (string, string) {
	t.Helper()
	project := t.TempDir()
	home := filepath.Join(project, "home")
	if err := os.MkdirAll(home, 0o755); err != nil {
		t.Fatalf("creating isolated HOME: %v", err)
	}
	return project, home
}

func runPublicFolio(t *testing.T, folio string, project string, home string, args ...string) string {
	t.Helper()
	cmdArgs := append([]string{"manuscript"}, args...)
	cmd := exec.Command(folio, cmdArgs...)
	cmd.Dir = project
	cmd.Env = append(os.Environ(), "HOME="+home)
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("folio manuscript failed: %v\n%s", err, string(output))
	}
	return readFile(t, filepath.Join(project, args[len(args)-1]))
}

func writeSiblingInputs(t *testing.T, project string) {
	t.Helper()
	writeConfigInput(t, filepath.Join(project, "part1", "ch01.md"), true)
	writeConfigInput(t, filepath.Join(project, "part2", "ch02.md"), false)
}

func writeConfigInput(t *testing.T, path string, frontmatter bool) {
	t.Helper()
	content := "## Chapter\n\nBody text for configuration selection.\n"
	if frontmatter {
		content = "---\ntitle: Source Title\nauthor: Test Author\n---\n\n" + content
	}
	writeFile(t, path, content)
}

func writeConfigFont(t *testing.T, path string, font string) {
	t.Helper()
	writeFile(t, path, "folio:\n  manuscript:\n    font: "+font+"\n")
}

func writeConfigTitle(t *testing.T, path string, title string) {
	t.Helper()
	writeFile(t, path, "title: "+title+"\n")
}
