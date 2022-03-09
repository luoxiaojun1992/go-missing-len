package pkg

import (
	"fmt"
	"go/ast"
	"go/token"
)

type MissingLenLinter struct {
	emptySliceMaps map[string]*slicePos
	emptyMapMaps   map[string]*mapPos
	Hints          []*Hint
}

type slicePos struct {
	Pos         token.Pos
	End         token.Pos
	InitLenNode *ast.BasicLit
}

type mapPos struct {
	Pos         token.Pos
	End         token.Pos
	InitLenNode *ast.BasicLit
}

type Hint struct {
	Pos        token.Pos
	End        token.Pos
	Category   string
	Message    string
	Suggestion string
}

func NewMissingLenLinter() *MissingLenLinter {
	linter := &MissingLenLinter{}
	linter.emptySliceMaps = make(map[string]*slicePos)
	linter.emptyMapMaps = make(map[string]*mapPos)
	return linter
}

func (l *MissingLenLinter) addHint(newHint *Hint) {
	l.Hints = append(l.Hints, newHint)
}

func (l *MissingLenLinter) Reset() {
	l.emptySliceMaps = make(map[string]*slicePos)
	l.emptyMapMaps = make(map[string]*mapPos)
	l.Hints = nil
}

func (l *MissingLenLinter) Check(file *ast.File) {
	ast.Walk(l, file)
}

func (l *MissingLenLinter) Visit(node ast.Node) ast.Visitor {
	if node != nil {
		switch n := node.(type) {
		case *ast.AssignStmt:
			switch asRh := n.Rhs[0].(type) {
			case *ast.CallExpr:
				callF, ok := asRh.Fun.(*ast.Ident)
				if !ok {
					break
				}
				if callF.Name == "make" {
					switch makeT := asRh.Args[0].(type) {
					case *ast.ArrayType:
						if makeT.Len == nil {
							switch initLen := asRh.Args[1].(type) {
							case *ast.BasicLit:
								if initLen.Value == "0" {
									l.addHint(&Hint{
										Pos:        initLen.Pos(),
										End:        initLen.End(),
										Category:   "missing-len",
										Message:    "Missing init len of slice",
										Suggestion: "Specific an init len of slice",
									})
									switch asVar := n.Lhs[0].(type) {
									case *ast.Ident:
										l.emptySliceMaps[asVar.Name] = &slicePos{
											Pos:         asVar.Pos(),
											End:         asVar.End(),
											InitLenNode: initLen,
										}
									}
								}
							}
						}
					case *ast.MapType:
						if len(asRh.Args) > 1 {
							switch initLen := asRh.Args[1].(type) {
							case *ast.BasicLit:
								if initLen.Value == "0" {
									l.addHint(&Hint{
										Pos:        initLen.Pos(),
										End:        initLen.End(),
										Category:   "missing-len",
										Message:    "Missing init len of map",
										Suggestion: "Specific an init len of map",
									})
									switch asVar := n.Lhs[0].(type) {
									case *ast.Ident:
										l.emptyMapMaps[asVar.Name] = &mapPos{
											Pos:         asVar.Pos(),
											End:         asVar.End(),
											InitLenNode: initLen,
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
				switch rangeVal := n.Value.(type) {
				case *ast.Ident:
					for _, rangeBodyStmt := range n.Body.List {
						switch rbs := rangeBodyStmt.(type) {
						case *ast.AssignStmt:
							switch rbsLh := rbs.Lhs[0].(type) {
							case *ast.IndexExpr:
								switch setVar := rbsLh.X.(type) {
								case *ast.Ident:
									if emptyMapPos, ok := l.emptyMapMaps[setVar.Name]; ok {
										switch setVal := rbs.Rhs[0].(type) {
										case *ast.Ident:
											if setVal.Name == rangeVal.Name {
												emptyMapPos.InitLenNode.Value = fmt.Sprintf("len(%s)", rangeVar.Name)
												l.addHint(&Hint{
													Pos:        emptyMapPos.Pos,
													End:        emptyMapPos.End,
													Category:   "missing-len",
													Message:    fmt.Sprintf("Missing init len of map[%s]", setVar.Name),
													Suggestion: fmt.Sprintf("May use len(%s)", rangeVar.Name),
												})
											}
										}
									}
								}
							}

							switch rbsRh := rbs.Rhs[0].(type) {
							case *ast.CallExpr:
								callF := rbsRh.Fun.(*ast.Ident)
								if callF.Name == "append" {
									switch apVar := rbsRh.Args[0].(type) {
									case *ast.Ident:
										if emptySlicePos, ok := l.emptySliceMaps[apVar.Name]; ok {
											switch apVal := rbsRh.Args[1].(type) {
											case *ast.Ident:
												if apVal.Name == rangeVal.Name {
													emptySlicePos.InitLenNode.Value = fmt.Sprintf("len(%s)", rangeVar.Name)
													l.addHint(&Hint{
														Pos:        emptySlicePos.Pos,
														End:        emptySlicePos.End,
														Category:   "missing-len",
														Message:    fmt.Sprintf("Missing init len of slice[%s]", apVar.Name),
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
			}
		}
	}
	return l
}
