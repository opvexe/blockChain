package main

func main() {
	bc := NewBlockChain()
	defer bc.db.Close()
	c:=NewCmd(bc)
	c.Run()
}
