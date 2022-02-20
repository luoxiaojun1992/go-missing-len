package pkg

import (
	"fmt"
	"go/ast"
	"go/token"
)

type MissingLenLinter struct {
	emptySliceMaps map[string]*sliceMapPos
	Hints []*Hint
}

type sliceMapPos struct {
	Pos token.Pos
	End token.Pos
}

type Hint struct {
	Pos token.Pos
	End token.Pos
	Category string
	Message string
	Suggestion string
}

func NewMissingLenLinter() *MissingLenLinter {
	linter := &MissingLenLinter{}
	linter.emptySliceMaps = make(map[string]*sliceMapPos)
	return linter
}

func (l *MissingLenLinter) addHint(newHint *Hint) {
	l.Hints = append(l.Hints, newHint)
}

func (l *MissingLenLinter) Reset() {
	l.emptySliceMaps = make(map[string]*sliceMapPos)
	l.Hints = nil
}

func (l *MissingLenLinter) Check(file *ast.File)  {
	ast.Walk(l, file)
}

func (l *MissingLenLinter) Visit(node ast.Node) ast.Visitor {
	if node != nil {
		switch n := node.(type) {
		case *ast.AssignStmt:
			switch asRh := n.Rhs[0].(type) {
			case *ast.CallExpr:
				callF := asRh.Fun.(*ast.Ident)
				if callF.Name == "make" {
					switch makeT := asRh.Args[0].(type) {
					case *ast.ArrayType:
						if makeT.Len == nil {
							switch initLen := asRh.Args[1].(type) {
							case *ast.BasicLit:
								if initLen.Value == "0" {
									switch asVar := n.Lhs[0].(type) {
									case *ast.Ident:
										l.emptySliceMaps[asVar.Name] = &sliceMapPos{
											Pos: asVar.Pos(),
											End: asVar.End(),
										}
									}
								}
							}
						}
					}
				}
			}
		case *ast.RangeStmt:
			switch rangeVar := n.X.(type) {
			case *ast.Ident:
				for _, rangeBodyStmt := range n.Body.List {
					switch rbs := rangeBodyStmt.(type) {
					case *ast.AssignStmt:
						switch rbsRh := rbs.Rhs[0].(type) {
						case *ast.CallExpr:
							callF := rbsRh.Fun.(*ast.Ident)
							if callF.Name == "append" {
								switch apVar := rbsRh.Args[0].(type) {
								case *ast.Ident:
									if emptySliceMap, ok := l.emptySliceMaps[apVar.Name]; ok {
										l.addHint(&Hint{
											Pos: emptySliceMap.Pos,
											End: emptySliceMap.End,
											Category: "missing-len",
											Message: fmt.Sprintf("Missing init len of slice[%s]", apVar.Name),
											Suggestion: fmt.Sprintf("May use len(%s)", rangeVar.Name),
										})
									}
								}
							}
						}
					}
				}
			}
		}
	}
	return l
}