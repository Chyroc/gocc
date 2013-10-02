//Copyright 2013 Vastech SA (PTY) LTD
//
//   Licensed under the Apache License, Version 2.0 (the "License");
//   you may not use this file except in compliance with the License.
//   You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
//   Unless required by applicable law or agreed to in writing, software
//   distributed under the License is distributed on an "AS IS" BASIS,
//   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//   See the License for the specific language governing permissions and
//   limitations under the License.

/*
This package controls the generation of all parser-related code.
*/
package golang

import (
	"code.google.com/p/gocc/ast"
	"code.google.com/p/gocc/config"
	"code.google.com/p/gocc/parser/lrk/action"
	"code.google.com/p/gocc/parser/lrk/states"
	"code.google.com/p/gocc/parser/symbols"
)

func Gen(cfg config.Config, header string, prods []*ast.SyntaxBasicProd, symbols *symbols.Symbols, states *states.States, actions action.Actions) {
	genAction(cfg.OutDir())
	genActionTable(cfg.OutDir(), prods, symbols, states, actions)
	genErrors(cfg)
	genGotoTable(cfg.OutDir(), symbols, states)
	genParser(cfg, prods, symbols, states)
	genProductionsTable(cfg, header, prods, symbols, states)

	return
}
