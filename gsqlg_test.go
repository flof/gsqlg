package gsqlg

import (
	"log"
	"testing"
)

func TestModel(t *testing.T) {
	var (
		userTable           = &Table{Name: "users", Alias: "u"}
		primaryGroupTable   = &Table{Name: "groups", Alias: "prim_grp"}
		secondaryGroupTable = &Table{Name: "groups", Alias: "sec_grp"}
	)

	s := Select{
		BaseSource: userTable,
		Joins: []Join{
			{
				JoinType: JoinTypeLeft,
				Source:   primaryGroupTable,
				Constraints: []JoinConstraint{
					{
						RefColumn:    SourceColumn{Source: primaryGroupTable, Name: "id"},
						TargetColumn: SourceColumn{Source: userTable, Name: "primary_group_id"},
					},
				},
			},
			{
				JoinType: JoinTypeLeft,
				Source:   secondaryGroupTable,
				Constraints: []JoinConstraint{
					{
						RefColumn:    SourceColumn{Source: secondaryGroupTable, Name: "id"},
						TargetColumn: SourceColumn{Source: userTable, Name: "secondary_group_id"},
					},
				},
			},
		},
		Projections: []Projection{
			userTable.Column("id"),
			userTable.Column("first_name"),
			userTable.Column("last_name"),
			primaryGroupTable.Column("name"),
			secondaryGroupTable.Column("name"),
		},
	}

	gen := StatementGenerator{}

	stmt, err := gen.Select(s)
	if err != nil {
		log.Fatal(err)
	}

	log.Print(stmt.Text)
}

func TestModelMinimal(t *testing.T) {
	var (
		userTable = &Table{Name: "users", Alias: "u"}
	)

	s := Select{
		BaseSource: userTable,
		Projections: []Projection{
			userTable.Column("id"),
			userTable.Column("username"),
			RawSelectColumn{
				Raw:   "NOW()",
				Alias: "created",
			},
		},
	}
	gen := StatementGenerator{}

	stmt, err := gen.Select(s)
	if err != nil {
		log.Fatal(err)
	}

	log.Print(stmt.Text)
}

func TestModelMinimal2(t *testing.T) {
	var (
		userTable      = &Table{Name: "users", Alias: "u"}
		addressesTable = &Table{Name: "addresses", Alias: "a"}
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
		Projections: []Projection{
			userTable.Column("id"),
			userTable.Column("username"),
			addressesTable.Column("street"),
			addressesTable.Column("city"),
		},
	}
	gen := StatementGenerator{}

	stmt, err := gen.Select(s)
	if err != nil {
		log.Fatal(err)
	}

	log.Print(stmt.Text)
}

func TestSubQuery(t *testing.T) {
	var (
		userTable      = &Table{Name: "users", Alias: "u"}
		addressesTable = &Table{Name: "addresses", Alias: "a"}
	)

	addresses := Select{
		BaseSource: addressesTable,
		Projections: []Projection{
			addressesTable.Column("user_id"),
			addressesTable.Column("street"),
			addressesTable.Column("city"),
		},
	}

	addressesSubQuery := &SubQuery{Select: addresses, Alias: "a"}

	s := Select{
		BaseSource: userTable,
		Joins: []Join{
			{
				JoinType: JoinTypeInner,
				Source:   addressesSubQuery,
				Constraints: []JoinConstraint{
					{
						RefColumn:    addressesTable.Column("user_id"),
						TargetColumn: userTable.Column("id"),
					},
				},
			},
		},
		Projections: []Projection{
			userTable.Column("id"),
			userTable.Column("username"),
			addressesSubQuery.Column("street"),
			addressesSubQuery.Column("city"),
		},
	}
	gen := StatementGenerator{}

	stmt, err := gen.Select(s)
	if err != nil {
		log.Fatal(err)
	}

	log.Print(stmt.Text)
}
