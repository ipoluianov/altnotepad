package editor

import (
	"strings"
	"testing"
)

func names(doc *Document) string {
	var s []string
	for _, f := range FunctionList(doc) {
		s = append(s, f.Name)
	}
	return strings.Join(s, ",")
}

func TestFunctionList(t *testing.T) {
	cases := []struct{ lang, text, want string }{
		{"c", "#include <x.h>\n// int fake(void) {\nstatic int counter;\nint main(int argc, char **argv) {\n  if (x) {\n    foo(1);\n  }\n}\nvoid Foo::bar() const\n{\n}\n", "main,Foo::bar"},
		{"go", "package x\n\nfunc (s *S) Method(a int) error {\n}\n\nfunc Free() {}\ntype T struct {\n}\n", "Method,Free,T"},
		{"python", "class A:\n    def f(self):\n        pass\n\nasync def g():\n    pass\n", "A,f,g"},
		{"markdown", "# Title\ntext\n## Sub\n```\n# not a heading? \n```\n", "Title,Sub,not a heading?"},
		{"javascript", "function a() {}\nconst b = (x) => x\nclass C {\n  method(x) {\n    if (x) {\n    }\n  }\n}\n", "a,b,C,method"},
	}
	for _, c := range cases {
		doc := NewDocument()
		doc.SetText([]byte(c.text))
		doc.SetLanguage(LanguageByID(c.lang))
		if got := names(doc); got != c.want {
			t.Errorf("%s: %q, want %q", c.lang, got, c.want)
		}
	}
}
