package ast_test

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"github.com/mohammad-safakhou/diffmind/internal/extractor/ast"
)

func TestBinaryTransportStreamIsNotTypeScriptSource(t *testing.T) {
	dir:=t.TempDir();data:=[]byte{0x47,0x40,0,0x10,0xff,0xff}
	if e:=os.WriteFile(filepath.Join(dir,"video.ts"),data,0644);e!=nil {t.Fatal(e)}
	if f,e:=ast.ParseFile(context.Background(),dir,"video.ts");f!=nil||e==nil||!strings.Contains(e.Error(),"binary") {t.Fatalf("binary source parsed: %+v %v",f,e)}
	writeFile(t,dir,"server.ts",`const greeting = "سلام";`)
	idx:=buildIndex(t,dir);if len(idx.Files)!=1||len(idx.InputWarnings)!=1||idx.Files["server.ts"]==nil {t.Fatalf("text/source boundary incorrect: %+v / %+v",idx.Files,idx.InputWarnings)}
}
