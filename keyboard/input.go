package keyboard

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

const (
	balance = "balance"
	pay     = "pay"
	exit    = "exit"
)

var arrCommandsString = [...]string{balance, pay, exit}

func supportedCommandsComma() string {
	return strings.Join(arrCommandsString[:], ",")
}

type Keyboard struct {
	bm BalanceManager
	w  io.Writer
}

type BalanceManager interface {
	Balance() int
	Decrease(value int)
}

func New(bm BalanceManager, w io.Writer) *Keyboard {
	return &Keyboard{
		bm: bm,
		w:  w,
	}
}

func (k *Keyboard) Loop() error {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		line, err := reader.ReadString('\n')
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to scan from stdin, %w", err)
		}
		tokens := strings.SplitN(strings.TrimRight(line, "\n"), " ", 3)
		switch tokens[0] {
		case balance:
			fmt.Println(k.bm.Balance())
		case pay:
			if len(tokens) != 3 {
				fmt.Println("Pay command invalid, example 'pay Bob 10'")
				continue
			}
			n, err := strconv.ParseInt(tokens[2], 10, 64)
			if err != nil {
				fmt.Println("Only integer values allowed to specify in the 'pay' command")
				continue
			}
			if _, err := k.w.Write([]byte(tokens[1] + " " + tokens[2])); err != nil {
				fmt.Printf("failed to write message to peer, %s\n", err.Error())
				continue
			}
			k.bm.Decrease(int(n))
		case exit:
			fmt.Println("Goodbye!")
			return nil
		default:
			fmt.Printf("unknwon command, only %s supported\n", supportedCommandsComma())
		}
	}
	return nil
}
