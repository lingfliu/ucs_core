package conn

type ConnMsg struct {
}

type Coder struct {
	on_decode func(msg *ConnMsg)
}

func (c *Coder) _task_decode() {
	msg := &ConnMsg{}

	if c.on_decode != nil {
		c.on_decode(msg)
	}
}

func (c *Coder) PutBytes(bs []byte) {
}

func (c *Coder) Encode() {
}

type ConnCli struct {
	coder Coder
	conn  Conn
}

func (c *ConnCli) Connect() {
}

func (c *ConnCli) Disconnect() {

}

func (c *ConnCli) _task_read() {
	for {
		bf := c.conn.ReadToBuff()
		c.coder.PutBytes(bf)
	}
}
