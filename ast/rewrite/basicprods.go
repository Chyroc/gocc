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

package rewrite

import (
	"code.google.com/p/gocc/ast"
	"fmt"
)

type augId map[string]int

func (this augId) Next(prodId string) int {
	if num, exist := this[prodId]; exist {
		this[prodId] = num + 1
	} else {
		this[prodId] = 1
	}
	return this[prodId]
}

func BasicProds(prods ast.SyntaxProdList) []*ast.SyntaxBasicProd {
	bprods := make([]*ast.SyntaxBasicProd, 0, 256)
	augIds := make(augId)
	for _, prod := range prods {
		for _, body := range prod.SyntaxExpression {
			bprods = append(bprods, rewriteSyntaxProd(string(prod.Id.Lit), augIds, body)...)
		}
	}
	return bprods
}

func rewriteSyntaxProd(prodId string, idx augId, expr ...*ast.SyntaxBody) (basicProds []*ast.SyntaxBasicProd) {
	basicProds = make([]*ast.SyntaxBasicProd, 0, 8)
	for _, body := range expr {
		basicProds = append(basicProds, rewriteSyntaxProdBody(string(prodId), idx, body)...)
	}
	return
}

/*
Return the basic productions for (prodId, body) and the next name index for prodId
*/
func rewriteSyntaxProdBody(prodId string, idx augId, body *ast.SyntaxBody) []*ast.SyntaxBasicProd {
	basicProds := make([]*ast.SyntaxBasicProd, 0, 8)
	prod := &ast.SyntaxBasicProd{
		Id:     prodId,
		Error:  body.Error,
		Terms:  make(ast.SyntaxTerms, 0, len(body.Terms)),
		Action: body.Action,
	}
	basicProds = append(basicProds, prod)
	for ti, term := range body.Terms {
		if term.Basic() {
			prod.Terms = append(prod.Terms, term)
		} else {
			newProdId := augmentedProdId(prodId, idx)
			prod.Terms = append(prod.Terms, ast.SyntaxProdId(newProdId))

			switch t := term.(type) {
			case ast.SyntaxGroupExpression:
				basicProds = append(basicProds, rewriteSyntaxProd(newProdId, idx, t...)...)
			case ast.SyntaxOptionalExpression:
				basicProds = append(basicProds, rewriteSyntaxProd(prodId, idx, bodyWithoutTerm(body, ti))...)
				basicProds = append(basicProds, rewriteSyntaxProd(newProdId, idx, t...)...)
			case ast.SyntaxRepeatedExpression:
				basicProds = append(basicProds, rewriteSyntaxProd(prodId, idx, bodyWithoutTerm(body, ti))...)
				repTermId := repPid(newProdId)
				basicProds = append(basicProds, repProds(newProdId, repTermId)...)
				for _, t1 := range t {
					prods := rewriteSyntaxProd(repTermId, idx, t1)
					basicProds = append(basicProds, prods...)
				}
			default:
				panic(fmt.Sprintf("Unexpected type of non-basic term: %T", term))
			}
		}
	}
	return basicProds
}

func repProds(repPid, repTid string) (prods []*ast.SyntaxBasicProd) {
	return []*ast.SyntaxBasicProd{
		{
			Id:     repPid,
			Error:  false,
			Terms:  ast.SyntaxTerms{ast.SyntaxProdId(repTid)},
			Action: "[]interface{}{X[0]}, nil",
		},
		{
			Id:     repPid,
			Error:  false,
			Terms:  ast.SyntaxTerms{ast.SyntaxProdId(repPid), ast.SyntaxProdId(repTid)},
			Action: "append(X[0].([]interface{}), X[1]), nil",
		},
	}
}

func repPid(pid string) string {
	return fmt.Sprintf("%s_RepTerm", pid)
}

func bodyWithoutTerm(body *ast.SyntaxBody, term int) *ast.SyntaxBody {
	newTerms := make(ast.SyntaxTerms, 0, len(body.Terms)-1)
	for i, t := range body.Terms {
		if i != term {
			newTerms = append(newTerms, t)
		}
	}
	return &ast.SyntaxBody{
		Error:  body.Error,
		Terms:  newTerms,
		Action: body.Action,
	}
}

func augmentedProdId(prodId string, index augId) string {
	return fmt.Sprintf("%s~%d", prodId, index.Next(prodId))
}
