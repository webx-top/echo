package echo

func (c *XContext) SetTransaction(t Transaction) {
	c.transaction = NewTransaction(t)
}

func (c *XContext) Transaction() Transaction {
	return c.transaction.Transaction
}

func (c *XContext) Begin() error {
	return c.transaction.Begin(c)
}

func (c *XContext) Rollback() error {
	return c.transaction.Rollback(c)
}

func (c *XContext) Commit() error {
	return c.transaction.Commit(c)
}

func (c *XContext) End(succeed bool) error {
	return c.transaction.End(c, succeed)
}
