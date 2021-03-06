package sqlparser

import (
	"github.com/knocknote/vitess-sqlparser/tidbparser/ast"
)

func convertFromCreateTableStmt(stmt *ast.CreateTableStmt, ddl *DDL) Statement {
	columns := []*ColumnDef{}
	for _, col := range stmt.Cols {
		columns = append(columns, &ColumnDef{
			Name:  col.Name.Name.String(),
			Type:  col.Tp.String(),
			Elems: col.Tp.Elems,
		})
	}
	constraints := []*Constraint{}
	for _, constraint := range stmt.Constraints {
		keys := []ColIdent{}
		for _, key := range constraint.Keys {
			keys = append(keys, NewColIdent(key.Column.Name.String()))
		}
		constraints = append(constraints, &Constraint{
			Type: ConstraintType(constraint.Tp),
			Name: constraint.Name,
			Keys: keys,
		})
	}
	return &CreateTable{
		DDL:         ddl,
		Columns:     columns,
		Constraints: constraints,
	}
}

func convertFromTruncateTableStmt(stmt *ast.TruncateTableStmt) Statement {
	return &TruncateTable{Table: TableName{Name: TableIdent{v: stmt.Table.Name.String()}}}
}

func convertTiDBStmtToVitessOtherAdmin(stmts []ast.StmtNode, admin *OtherAdmin) Statement {
	for _, stmt := range stmts {
		switch adminStmt := stmt.(type) {
		case *ast.TruncateTableStmt:
			return convertFromTruncateTableStmt(adminStmt)
		default:
			return admin
		}
	}
	return nil
}

func convertTiDBStmtToVitessDDL(stmts []ast.StmtNode, ddl *DDL) Statement {
	for _, stmt := range stmts {
		switch ddlStmt := stmt.(type) {
		case *ast.CreateTableStmt:
			return convertFromCreateTableStmt(ddlStmt, ddl)
		default:
			return ddl
		}
	}
	return nil
}
