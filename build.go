package main

import (
	"flag"
	"fmt"
)

type buildCmd struct {
	output string
}

func (cmd *buildCmd) Name() string {
	return "build"
}

func (cmd *buildCmd) Desc() string {
	return "builds output"
}

func (cmd *buildCmd) Help() string {
	return `Usage: build [flags]:

The build command converts the content files in the publish directory
into static output files using the transformations specified in the
tmpl directory.`
}

func (cmd *buildCmd) Flags(fset *flag.FlagSet) {
	fset.StringVar(&cmd.output, "o", "build", "output directory name relative to project root")
}

func (cmd *buildCmd) Run(args []string) error {
	root, ok := getRoot()
	if !ok {
		return noRootErr
	}

	//publish := filepath.Join(root, "publish")
	//output := filepath.Join(root, cmd.output)

	tmpl, err := loadTmpl(root)
	if err != nil {
		return fmt.Errorf("failed to load templates: %v", err)
	}

	for t, tmpl := range tmpl {
		fmt.Printf("%v: %#v\n", t, tmpl)
	}
	return nil

	//var meta map[string]interface{}
	//return walk(publish, func(p string, fi os.FileInfo) error {
	//	if fi.IsDir() {
	//		return nil
	//	}

	//	in, err := os.Open(filepath.Join(publish, p))
	//	if err != nil {
	//		return fmt.Errorf("failed to open %q: %v", p, err)
	//	}
	//	defer in.Close()

	//	meta = make(map[string]interface{}, len(meta))
	//	inr, err := readMeta(in, &meta)
	//	if err != nil {
	//		return fmt.Errorf("failed to load meta from %q: %v", p, err)
	//	}

	//	tmplType, ok := meta["type"].(string)
	//	if !ok {
	//		return fmt.Errorf("no type in %q", p)
	//	}

	//	t, ok := tmpl[tmplType]
	//	if !ok {
	//		return fmt.Errorf("unknown type %q in %q", tmplType, p)
	//	}

	//	var content strings.Builder
	//	err = t.tmpl.Execute(&content, map[string]interface{}{
	//		"tmpl": t,
	//		"meta": meta,
	//	})
	//	if err != nil {
	//		return fmt.Errorf("failed to execute %q: %v", tmplType, err)
	//	}

	//	err = os.MkdirAll(filepath.Join(output, filepath.Dir(p)), 0755)
	//	if err != nil {
	//		return fmt.Errorf("failed to create directory for %q: %v", p, err)
	//	}

	//	op := strings.TrimSuffix(p, ".md") + filepath.Ext(tmplType)
	//	out, err := os.OpenFile(
	//		filepath.Join(output, op),
	//		os.O_WRONLY|os.O_CREATE|os.O_TRUNC,
	//		0644,
	//	)
	//	if err != nil {
	//		return fmt.Errorf("failed to create %q: %v", p, err)
	//	}
	//	defer out.Close()

	//	return nil
	//})
}
