package gsqlg

import (
	"errors"
	"fmt"
)

var (
	ErrUnknownProjection = errors.New("unknown projection")
	ErrUnknownSource     = errors.New("unkown source")
)

type StatementGenerator struct {
}

type Statement struct {
	Text string
	Args []interface{}
}

func (s StatementGenerator) Select(sel Select) (Statement, error) {
	collector := &Collector{}

	err := writeSelect(collector, sel)
	if err != nil {
		return Statement{}, err
	}
	
	collector.AddGapless(";")

	return Statement{
		Text: collector.String(),
	}, nil
}

func writeSelect(collector *Collector, sel Select) error {
	collector.Add("select")

	for idx, proj := range sel.Projections {
		if idx > 0 {
			collector.AddGaplessLeft(",")
		}

		err := writeProjection(collector, proj)
		if err != nil {
			return err
		}
	}

	collector.Add("from")
	err := writeSourceDefinition(collector, sel.BaseSource)
	if err != nil {
		return err
	}

	for _, join := range sel.Joins {
		switch join.JoinType {
		case JoinTypeInner:
			collector.Add("inner")
		case JoinTypeLeft:
			collector.Add("left")
		case JoinTypeRight:
			collector.Add("right")
		case JoinTypeCross:
			collector.Add("cross")
		}

		collector.Add("join")

		err := writeSourceDefinition(collector, join.Source)
		if err != nil {
			return err
		}

		for idx, joinCon := range join.Constraints {
			if idx == 0 {
				collector.Add("on")
			} else {
				collector.Add("and")
			}

			collector.Add(joinCon.RefColumn.Source.GetAlias())
			collector.AddGapless(".")
			collector.Add(joinCon.RefColumn.Name)
			collector.Add("=")
			collector.Add(joinCon.TargetColumn.Source.GetAlias())
			collector.AddGapless(".")
			collector.Add(joinCon.TargetColumn.Name)
		}
	}

	return nil
}

func writeProjection(collector *Collector, projection Projection) error {
	switch p := projection.(type) {
	case SourceColumn:
		collector.Add(p.Source.GetAlias())
		collector.AddGapless(".")
		collector.Add(p.Name)
		collector.Add("as")
		collector.Add(fmt.Sprintf("%s_%s", p.Source.GetAlias(), p.Name))
	case RawSelectColumn:
		collector.Add(p.Raw)
		collector.Add("as")
		collector.Add(p.Alias)
	default:
		return ErrUnknownProjection
	}

	return nil
}

func writeSourceDefinition(collector *Collector, source Source) error {
	switch s := source.(type) {
	case *Table:
		collector.Add(s.Name)
		collector.Add(s.Alias)
	case *SubQuery:
		collector.AddGaplessRight("(")
		writeSelect(collector, s.Select)
		collector.AddGaplessLeft(")")
		collector.Add(s.Alias)
	default:
		return ErrUnknownSource
	}

	return nil
}
