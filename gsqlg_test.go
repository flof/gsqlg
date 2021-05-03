package gsqlg

import (
	"log"
	"testing"
)

func TestModel(t *testing.T) {
	var (
		userTable           = &Table{TableName: "users", TableAlias: "u"}
		primaryGroupTable   = &Table{TableName: "groups", TableAlias: "prim_grp"}
		secondaryGroupTable = &Table{TableName: "groups", TableAlias: "sec_grp"}
	)

	s := Select{
		BaseSource: userTable,
		Joins: []Join{
			{
				JoinType: JoinTypeLeft,
				Source:   primaryGroupTable,
				Constraints: []JoinConstraint{
					{
						RefColumn:    SourceColumn{SourceRef: primaryGroupTable, Name: "id"},
						TargetColumn: SourceColumn{SourceRef: userTable, Name: "primary_group_id"},
					},
				},
			},
			{
				JoinType: JoinTypeLeft,
				Source:   secondaryGroupTable,
				Constraints: []JoinConstraint{
					{
						RefColumn:    SourceColumn{SourceRef: secondaryGroupTable, Name: "id"},
						TargetColumn: SourceColumn{SourceRef: userTable, Name: "secondary_group_id"},
					},
				},
			},
		},
		Columns: []SelectColumn{
			userTable.Column("id"),
			userTable.Column("first_name"),
			userTable.Column("last_name"),
			primaryGroupTable.Column("name"),
			secondaryGroupTable.Column("name"),
		},
	}

	gen := StatementGenerator{}

	stmt := gen.Select(s)
	log.Print(stmt.Text)
}

func TestModelMinimal(t *testing.T) {
	var (
		userTable = &Table{TableName: "users", TableAlias: "u"}
	)

	s := Select{
		BaseSource: userTable,
		Columns: []SelectColumn{
			userTable.Column("id"),
			userTable.Column("username"),
			RawSelectColumn{
				Raw: "NOW()",
				RawAlias: "created",
			},
		},
	}
	gen := StatementGenerator{}

	stmt := gen.Select(s)
	log.Print(stmt.Text)
}

func TestModelMinimal2(t *testing.T) {
	var (
		userTable      = &Table{TableName: "users", TableAlias: "u"}
		addressesTable = &Table{TableName: "addresses", TableAlias: "a"}
	)

	s := Select{
		BaseSource: userTable,
		Joins: []Join{
			{
				JoinType: JoinTypeInner,
				Source:   addressesTable,
				Constraints: []JoinConstraint{
					{
						RefColumn:    addressesTable.Column("user_id"),
						TargetColumn: userTable.Column("id"),
					},
				},
			},
		},
		Columns: []SelectColumn{
			userTable.Column("id"),
			userTable.Column("username"),
			addressesTable.Column("street"),
			addressesTable.Column("city"),
		},
	}
	gen := StatementGenerator{}

	stmt := gen.Select(s)
	log.Print(stmt.Text)
}
