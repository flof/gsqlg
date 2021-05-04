package gsqlg

import "strings"

type Collector struct {
	s           strings.Builder
	noSpaceNext bool
}

func (c *Collector) Add(s string) {
	if c.s.Len() > 0 && !c.noSpaceNext {
		c.s.WriteString(" ")
	}

	c.s.WriteString(s)
	c.noSpaceNext = false
}

func (c *Collector) AddGapless(s string) {
	c.s.WriteString(s)
	c.noSpaceNext = true
}

func (c *Collector) AddGaplessLeft(s string) {
	c.s.WriteString(s)
	c.noSpaceNext = false
}

func (c *Collector) AddGaplessRight(s string) {
	if c.s.Len() > 0 && !c.noSpaceNext {
		c.s.WriteString(" ")
	}

	c.s.WriteString(s)
	c.noSpaceNext = true
}

func (c *Collector) String() string {
	return c.s.String()
}
