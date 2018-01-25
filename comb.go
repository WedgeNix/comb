package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/WedgeNix/comb/lib"
)

type Choice int

const (
	AddInput Choice = 19 // ^S
	CalcComb Choice = 24 // ^X
)

var (
	intexp  = regexp.MustCompile(`[-+]?\d+(,\d+)*`)
	numexp  = regexp.MustCompile(`[+-]?\d+(?:,\d+)*(?:[.]\d+)?`)
	goalexp = regexp.MustCompile(`(?i)\b(goal|aim|fin|final|finish|sum|end|stop|val|value|get|grab|num|number|last|here|this|add|to|towards|there)\b`)
	fileexp = regexp.MustCompile(`"[^"]+"`)

	rdr = bufio.NewReader(os.Stdin)

	sprint     = fmt.Sprint
	itoa       = strconv.Itoa
	atoi       = strconv.Atoi
	repeat     = strings.Repeat
	parsefloat = strconv.ParseFloat
)

func main() {
	var (
		allgoals,
		allpools [][]float64
		ctrlx []int
	)

	choice := AddInput

	for {
		switch choice {
		case AddInput:
			addinput(&allgoals, &allpools)
		case CalcComb:
			if err := params(2, ctrlx, 1, len(allgoals), 1, len(allpools)); err != nil {
				print(err)
				break
			}
			gi, pi := ctrlx[0]-1, ctrlx[1]-1
			for _, goal := range allgoals[gi] {
				print("goal " + ftoa(goal) + "  " + sprint(lib.Find(goal, allpools[pi])))
			}
		}

		ctrlxstr := rng(allgoals) + rng(allpools)
		spaces := repeat(" ", len(ctrlxstr))

		var ctrlxfn string
		if len(allgoals) > 0 && len(allpools) > 0 {
			ctrlxfn = "ctrl-X" + ctrlxstr + "  find goals using pool"
		}

		print(
			"ctrl-S"+spaces+"  add input",
			ctrlxfn,
			"",
			"goals  "+sprint(allgoals),
			"pools  "+sprint(allpools),
		)

		choice, ctrlx = atoiParams(choose())
	}
}

func params(cnt int, args []int, fromto ...int) error {
	if L := len(args); L != cnt || len(fromto)/2 != L {
		return errors.New("bad parameters")
	}
	for a, arg := range args {
		from, to := fromto[a<<1], fromto[a<<1+1]
		if arg < from {
			return errors.New("bad arguments  " + itoa(arg) + "<" + itoa(from))
		} else if arg > to {
			return errors.New("bad arguments  " + itoa(arg) + ">" + itoa(to))
		}
	}
	return nil
}

func ftoa(f float64) string {
	return sprint(f)
}
func atoiParams(c Choice, s string) (Choice, []int) {
	return c, atois(s)
}
func atois(s string) []int {
	var ints []int
	for _, intstr := range intexp.FindAllString(s, -1) {
		if n, err := atoi(intstr); err == nil {
			ints = append(ints, n)
		}
	}
	return ints
}
func atofs(s string) []float64 {
	var nums []float64
	for _, numstr := range numexp.FindAllString(s, -1) {
		if n, err := parsefloat(numstr, 64); err == nil {
			nums = append(nums, n)
		}
	}
	return nums
}

func rng(f [][]float64) string {
	L := len(f)
	switch L {
	case 0:
		return ""
	case 1:
		return "  1"
	default:
		return "  1-" + strconv.Itoa(L)
	}
}

func choose() (Choice, string) {
	s, _ := rdr.ReadString('\n')
	s = clean(s)
	if len(s) > 0 {
		return Choice(int(s[0])), s
	}
	return -1, s
}

func addinput(ag, ap *[][]float64) {
	var (
		allgoals,
		allpools,
		allgoalsf,
		allpoolsf [][]float64
		poolbuf,
		poolbuff []float64
	)

	print(
		"paste input/drop file",
		"",
		"ctrl-S  continue",
	)

NextLine:
	for cmd := 0; cmd != 19; {
		in, _ := rdr.ReadString('\n')
		in = clean(in)

		if len(in) > 0 {
			cmd = int(in[0])

			for _, file := range fileexp.FindAllString(in, -1) {
				b, err := ioutil.ReadFile(unquote(file))
				if err != nil {
					println(err)
					continue NextLine
				}

				for _, lnb := range bytes.Split(b, []byte("\n")) {
					scan(string(lnb), &allgoalsf, &allpoolsf, &poolbuff)
				}
				flushpool(&allpoolsf, &poolbuff)
			}

			scan(in, &allgoals, &allpools, &poolbuf)
		}
	}
	flushpool(&allpools, &poolbuf)

	allgoals = append(allgoalsf, allgoals...)
	allpools = append(allpoolsf, allpools...)

	*ag, *ap = append(*ag, allgoals...), append(*ap, allpools...)
}

func clean(s string) string {
	return strings.NewReplacer("\r\n", "", "\r", "", "\n", "").Replace(strings.Trim(s, " \t"))
}

func unquote(s string) string {
	return strings.Trim(s, "\"")
}

func print(lns ...interface{}) { println(lns...) }
func println(lns ...interface{}) {
	buf := "\n"
	for _, ln := range lns {
		buf += ">>> " + fmt.Sprintln(ln)
	}
	fmt.Println(buf)
}

func flushpool(ap *[][]float64, pb *[]float64) {
	allpools := *ap
	poolbuf := *pb
	if len(poolbuf) > 0 {
		allpools = append(allpools, poolbuf)
		poolbuf = nil
	}
	*ap = allpools
	*pb = poolbuf
}

func scan(in string, ag, ap *[][]float64, pb *[]float64) {
	allgoals := *ag
	allpools := *ap
	poolbuf := *pb
	nums := atofs(in)
	if len(nums) > 0 {
		if goalexp.MatchString(in) {
			allgoals = append(allgoals, nums)
			flushpool(&allpools, &poolbuf)
		} else {
			poolbuf = append(poolbuf, nums...)
		}
	}
	*ag = allgoals
	*ap = allpools
	*pb = poolbuf
}
