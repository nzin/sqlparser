package parser

type SqlSelect struct {
	Table  string
	Column string
}

type Sql struct {
	Select []SqlSelect
	From   []string
	Where  string
}
