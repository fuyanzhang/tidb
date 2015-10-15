// Copyright 2015 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// See the License for the specific language governing permissions and
// limitations under the License.

package ast

var (
	_ StmtNode = &ExplainStmt{}
	_ StmtNode = &PrepareStmt{}
	_ StmtNode = &DeallocateStmt{}
	_ StmtNode = &ExecuteStmt{}
	_ StmtNode = &ShowStmt{}
	_ StmtNode = &BeginStmt{}
	_ StmtNode = &CommitStmt{}
	_ StmtNode = &RollbackStmt{}
	_ StmtNode = &UseStmt{}
	_ StmtNode = &SetStmt{}

	_ Node = &VariableAssignment{}
)

// AuthOption is used for parsing create use statement.
type AuthOption struct {
	// AuthString/HashString can be empty, so we need to decide which one to use.
	ByAuthString bool
	AuthString   string
	HashString   string
	// TODO: support auth_plugin
}

// ExplainStmt is a statement to provide information about how is SQL statement executed
// or get columns information in a table.
// See: https://dev.mysql.com/doc/refman/5.7/en/explain.html
type ExplainStmt struct {
	stmtNode

	Stmt DMLNode
}

// Accept implements Node Accept interface.
func (es *ExplainStmt) Accept(v Visitor) (Node, bool) {
	if !v.Enter(es) {
		return es, false
	}
	node, ok := es.Stmt.Accept(v)
	if !ok {
		return es, false
	}
	es.Stmt = node.(DMLNode)
	return v.Leave(es)
}

// PrepareStmt is a statement to prepares a SQL statement which contains placeholders,
// and it is executed with ExecuteStmt and released with DeallocateStmt.
// See: https://dev.mysql.com/doc/refman/5.7/en/prepare.html
type PrepareStmt struct {
	stmtNode

	InPrepare bool // true for prepare mode, false for use mode
	Name      string
	ID        uint32 // For binary protocol, there is no Name but only ID
	SQLStmt   Node   // The parsed statement from sql text with placeholder
}

// Accept implements Node Accept interface.
func (ps *PrepareStmt) Accept(v Visitor) (Node, bool) {
	if !v.Enter(ps) {
		return ps, false
	}
	node, ok := ps.SQLStmt.Accept(v)
	if !ok {
		return ps, false
	}
	ps.SQLStmt = node
	return v.Leave(ps)
}

// DeallocateStmt is a statement to release PreparedStmt.
// See: https://dev.mysql.com/doc/refman/5.7/en/deallocate-prepare.html
type DeallocateStmt struct {
	stmtNode

	Name string
	ID   uint32 // For binary protocol, there is no Name but only ID.
}

// Accept implements Node Accept interface.
func (ds *DeallocateStmt) Accept(v Visitor) (Node, bool) {
	if !v.Enter(ds) {
		return ds, false
	}
	return v.Leave(ds)
}

// ExecuteStmt is a statement to execute PreparedStmt.
// See: https://dev.mysql.com/doc/refman/5.7/en/execute.html
type ExecuteStmt struct {
	stmtNode

	Name      string
	ID        uint32 // For binary protocol, there is no Name but only ID
	UsingVars []ExprNode
}

// Accept implements Node Accept interface.
func (es *ExecuteStmt) Accept(v Visitor) (Node, bool) {
	if !v.Enter(es) {
		return es, false
	}
	for i, val := range es.UsingVars {
		node, ok := val.Accept(v)
		if !ok {
			return es, false
		}
		es.UsingVars[i] = node.(ExprNode)
	}
	return v.Leave(es)
}

// ShowStmt is a statement to provide information about databases, tables, columns and so on.
// See: https://dev.mysql.com/doc/refman/5.7/en/show.html
type ShowStmt struct {
	stmtNode

	Target int // Databases/Tables/Columns/....
	DBName string
	Table  *TableRef      // Used for showing columns.
	Column *ColumnRefExpr // Used for `desc table column`.
	Flag   int            // Some flag parsed from sql, such as FULL.
	Full   bool

	// Used by show variables
	GlobalScope bool
	Pattern     *PatternLikeExpr
	Where       ExprNode
}

// Accept implements Node Accept interface.
func (ss *ShowStmt) Accept(v Visitor) (Node, bool) {
	if !v.Enter(ss) {
		return ss, false
	}
	if ss.Table != nil {
		node, ok := ss.Table.Accept(v)
		if !ok {
			return ss, false
		}
		ss.Table = node.(*TableRef)
	}
	if ss.Column != nil {
		node, ok := ss.Column.Accept(v)
		if !ok {
			return ss, false
		}
		ss.Column = node.(*ColumnRefExpr)
	}
	if ss.Pattern != nil {
		node, ok := ss.Pattern.Accept(v)
		if !ok {
			return ss, false
		}
		ss.Pattern = node.(*PatternLikeExpr)
	}
	if ss.Where != nil {
		node, ok := ss.Where.Accept(v)
		if !ok {
			return ss, false
		}
		ss.Where = node.(ExprNode)
	}
	return v.Leave(ss)
}

// BeginStmt is a statement to start a new transaction.
// See: https://dev.mysql.com/doc/refman/5.7/en/commit.html
type BeginStmt struct {
	stmtNode
}

// Accept implements Node Accept interface.
func (bs *BeginStmt) Accept(v Visitor) (Node, bool) {
	if !v.Enter(bs) {
		return bs, false
	}
	return v.Leave(bs)
}

// CommitStmt is a statement to commit the current transaction.
// See: https://dev.mysql.com/doc/refman/5.7/en/commit.html
type CommitStmt struct {
	stmtNode
}

// Accept implements Node Accept interface.
func (cs *CommitStmt) Accept(v Visitor) (Node, bool) {
	if !v.Enter(cs) {
		return cs, false
	}
	return v.Leave(cs)
}

// RollbackStmt is a statement to roll back the current transaction.
// See: https://dev.mysql.com/doc/refman/5.7/en/commit.html
type RollbackStmt struct {
	stmtNode
}

// Accept implements Node Accept interface.
func (rs *RollbackStmt) Accept(v Visitor) (Node, bool) {
	if !v.Enter(rs) {
		return rs, false
	}
	return v.Leave(rs)
}

// UseStmt is a statement to use the DBName database as the current database.
// See: https://dev.mysql.com/doc/refman/5.7/en/use.html
type UseStmt struct {
	stmtNode

	DBName string
}

// Accept implements Node Accept interface.
func (us *UseStmt) Accept(v Visitor) (Node, bool) {
	if !v.Enter(us) {
		return us, false
	}
	return v.Leave(us)
}

// VariableAssignment is a variable assignment struct.
type VariableAssignment struct {
	node
	Name     string
	Value    ExprNode
	IsGlobal bool
	IsSystem bool
}

// Accept implements Node interface.
func (va *VariableAssignment) Accept(v Visitor) (Node, bool) {
	if !v.Enter(va) {
		return va, false
	}
	node, ok := va.Value.Accept(v)
	if !ok {
		return va, false
	}
	va.Value = node.(ExprNode)
	return v.Leave(va)
}

// SetStmt is the statement to set variables.
type SetStmt struct {
	stmtNode
	// Variables is the list of variable assignment.
	Variables []*VariableAssignment
}

// Accept implements Node Accept interface.
func (set *SetStmt) Accept(v Visitor) (Node, bool) {
	if !v.Enter(set) {
		return set, false
	}
	for i, val := range set.Variables {
		node, ok := val.Accept(v)
		if !ok {
			return set, false
		}
		set.Variables[i] = node.(*VariableAssignment)
	}
	return v.Leave(set)
}