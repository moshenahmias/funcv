package funcv

import (
	"testing"
)

func TestConstSensitive(t *testing.T) {
	c := NewCommand("").AddConstant("test", false).MustCompile()

	fail := true

	_, err := c.Execute([]string{"test"}, func() {
		fail = false
	})

	if err != nil {
		t.Fatal(err)
	}

	if fail {
		t.Fatal("func not called")
	}
}

func TestConstInsensitive(t *testing.T) {
	c := NewCommand("").AddConstant("test", true).MustCompile()

	fail := true

	_, err := c.Execute([]string{"TeSt"}, func() {
		fail = false
	})

	if err != nil {
		t.Fatal(err)
	}

	if fail {
		t.Fatal("func not called")
	}
}

func TestStrVar(t *testing.T) {
	c := NewCommand("").AddStrVar("test", "").MustCompile()

	var v string

	_, err := c.Execute([]string{"xyz"}, func(a string) {
		v = a
	})

	if err != nil {
		t.Fatal(err)
	}

	if v != "xyz" {
		t.Fatal("wrong value", v)
	}
}

func TestIntVar(t *testing.T) {
	c := NewCommand("").AddIntVar("test", 10, "").MustCompile()

	var v int

	_, err := c.Execute([]string{"123"}, func(a int) {
		v = a
	})

	if err != nil {
		t.Fatal(err)
	}

	if v != 123 {
		t.Fatal("wrong value", v)
	}
}

func TestStrVarWithDefault(t *testing.T) {
	c := NewCommand("").AddStrVarWithDefault("test", "xyz", "").MustCompile()

	var v string

	_, err := c.Execute([]string{}, func(a string) {
		v = a
	})

	if err != nil {
		t.Fatal(err)
	}

	if v != "xyz" {
		t.Fatal("wrong value", v)
	}
}

func TestIntVarWithDefault(t *testing.T) {
	c := NewCommand("").AddIntVarWithDefault("test", 123, 10, "").MustCompile()

	var v int

	_, err := c.Execute([]string{}, func(a int) {
		v = a
	})

	if err != nil {
		t.Fatal(err)
	}

	if v != 123 {
		t.Fatal("wrong value", v)
	}
}

func TestStrFlag(t *testing.T) {
	c := NewCommand("").AddStrFlag("x", "xyz", "").MustCompile()

	var v string

	_, err := c.Execute([]string{"-x", "uvw"}, func(a string) {
		v = a
	})

	if err != nil {
		t.Fatal(err)
	}

	if v != "uvw" {
		t.Fatal("wrong value", v)
	}
}

func TestStrFlagDefault(t *testing.T) {
	c := NewCommand("").AddStrFlag("x", "xyz", "").MustCompile()

	var v string

	_, err := c.Execute([]string{}, func(a string) {
		v = a
	})

	if err != nil {
		t.Fatal(err)
	}

	if v != "xyz" {
		t.Fatal("wrong value", v)
	}
}

func TestIntFlag(t *testing.T) {
	c := NewCommand("").AddIntFlag("xx", 123, 10, "").MustCompile()

	var v int

	_, err := c.Execute([]string{"--xx", "456"}, func(a int) {
		v = a
	})

	if err != nil {
		t.Fatal(err)
	}

	if v != 456 {
		t.Fatal("wrong value", v)
	}
}

func TestIntFlagDefault(t *testing.T) {
	c := NewCommand("").AddIntFlag("xx", 123, 10, "").MustCompile()

	var v int

	_, err := c.Execute([]string{}, func(a int) {
		v = a
	})

	if err != nil {
		t.Fatal(err)
	}

	if v != 123 {
		t.Fatal("wrong value", v)
	}
}

func TestBoolFlag(t *testing.T) {
	c := NewCommand("").AddBoolFlag("b", "").MustCompile()

	var v bool

	_, err := c.Execute([]string{"-b"}, func(a bool) {
		v = a
	})

	if err != nil {
		t.Fatal(err)
	}

	if !v {
		t.Fatal("wrong value", v)
	}
}

func TestBoolFlagDefault(t *testing.T) {
	c := NewCommand("").AddBoolFlag("b", "").MustCompile()

	var v = true

	_, err := c.Execute([]string{}, func(a bool) {
		v = a
	})

	if err != nil {
		t.Fatal(err)
	}

	if v {
		t.Fatal("wrong value", v)
	}
}

func TestBoolFlagWithFalseParam(t *testing.T) {
	c := NewCommand("").AddBoolFlag("b", "").MustCompile()

	_, err := c.Execute([]string{"-b", "false"}, func(a bool) {
		if a {
			t.Fatal("wrong value", a)
		}
	})

	if err != nil {
		t.Fatal(err)
	}
}

func TestBoolFlagWithTrueParam(t *testing.T) {
	c := NewCommand("").AddBoolFlag("b", "").MustCompile()

	_, err := c.Execute([]string{"-b", "true"}, func(a bool) {
		if !a {
			t.Fatal("wrong value", a)
		}
	})

	if err != nil {
		t.Fatal(err)
	}
}

func TestBoolFlagWithoutParam(t *testing.T) {
	c := NewCommand("").AddBoolFlag("b", "").MustCompile()

	_, err := c.Execute([]string{"-b"}, func(a bool) {
		if !a {
			t.Fatal("wrong value", a)
		}
	})

	if err != nil {
		t.Fatal(err)
	}
}

func TestStrFlagWithoutParam(t *testing.T) {
	c := NewCommand("").AddStrFlag("s", "xyz", "").MustCompile()

	if _, err := c.Execute([]string{"-s"}, nil); err == nil {
		t.Fatal(err)
	}
}

func TestIntFlagWithoutParam(t *testing.T) {
	c := NewCommand("").AddIntFlag("i", 123, 10, "").MustCompile()

	if _, err := c.Execute([]string{"-i"}, nil); err == nil {
		t.Fatal(err)
	}
}

func TestBadFlagName000(t *testing.T) {
	c := NewCommand("").AddStrFlag("x", "xyz", "").MustCompile()
	_, err := c.Execute([]string{"x", "uvw"}, nil)

	if err == nil {
		t.Fatal("wrong flag name passed")
	}
}

func TestBadFlagName001(t *testing.T) {
	c := NewCommand("").AddStrFlag("x", "xyz", "").MustCompile()
	_, err := c.Execute([]string{"--x", "uvw"}, nil)

	if err == nil {
		t.Fatal("wrong flag name passed")
	}
}

func TestBadFlagName002(t *testing.T) {
	c := NewCommand("").AddStrFlag("xx", "xyz", "").MustCompile()
	_, err := c.Execute([]string{"xx", "uvw"}, nil)

	if err == nil {
		t.Fatal("wrong flag name passed")
	}
}

func TestBadFlagName003(t *testing.T) {
	c := NewCommand("").AddStrFlag("xx", "xyz", "").MustCompile()
	_, err := c.Execute([]string{"-xx", "uvw"}, nil)

	if err == nil {
		t.Fatal("wrong flag name passed")
	}
}

func TestCombined(t *testing.T) {
	c := NewCommand("").
		AddConstant("test", false).
		AddStrFlag("x", "xxx", "").
		AddIntFlag("y", 111, 10, "").
		AddBoolFlag("z", "").
		AddStrVar("v1", "").
		AddIntVar("v2", 10, "").
		AddStrVarWithDefault("v3", "v3def", "").
		AddIntVarWithDefault("v4", 444, 10, "").
		MustCompile()

	args := []string{"test", "-x", "xxxx", "-y", "1111", "-z", "111", "222"}

	_, err := c.Execute(args, func(x string, y int, z bool, v1 string, v2 int, v3 string, v4 int) {
		if x != "xxxx" {
			t.Fatal("wrong x", x)
		}

		if y != 1111 {
			t.Fatal("wrong y", y)
		}

		if !z {
			t.Fatal("wrong z", z)
		}

		if v1 != "111" {
			t.Fatal("wrong v1", v1)
		}

		if v2 != 222 {
			t.Fatal("wrong v2", v2)
		}

		if v3 != "v3def" {
			t.Fatal("wrong v3", v3)
		}

		if v4 != 444 {
			t.Fatal("wrong v4", v4)
		}
	})

	if err != nil {
		t.Fatal(err)
	}
}

func TestMissingFuncParams(t *testing.T) {
	c := NewCommand("").
		AddConstant("test", false).
		AddStrFlag("x", "xxx", "").
		AddIntFlag("y", 111, 10, "").
		AddBoolFlag("z", "").
		AddStrVar("v1", "").
		AddIntVar("v2", 10, "").
		AddStrVarWithDefault("v3", "v3def", "").
		AddIntVarWithDefault("v4", 444, 10, "").
		MustCompile()

	args := []string{"test", "-x", "xxxx", "-y", "1111", "-z", "111", "222"}

	_, err := c.Execute(args, func(x string, y int, z bool, v1 string, v2 int) {

	})

	if err == nil {
		t.Fatal("missing params passed")
	}
}

func TestNotAFunc(t *testing.T) {
	c := NewCommand("").
		AddConstant("test", false).
		AddStrFlag("x", "xxx", "").
		AddIntFlag("y", 111, 10, "").
		AddBoolFlag("z", "").
		AddStrVar("v1", "").
		AddIntVar("v2", 10, "").
		AddStrVarWithDefault("v3", "v3def", "").
		AddIntVarWithDefault("v4", 444, 10, "").
		MustCompile()

	args := []string{"test", "-x", "xxxx", "-y", "1111", "-z", "111", "222"}

	_, err := c.Execute(args, "")

	if err == nil {
		t.Fatal("not a func passed")
	}
}

func TestGroup000(t *testing.T) {

	cmd0, cmd1, cmd2 := false, false, false

	var grp Group

	if err := NewCommand("").AddConstant("cmd0", false).ToGroup(&grp, func() {
		cmd0 = true
	}); err != nil {
		t.Fatal(err)
	}

	if err := NewCommand("").AddConstant("cmd1", false).ToGroup(&grp, func() {
		cmd1 = true
	}); err != nil {
		t.Fatal(err)
	}

	if err := NewCommand("").AddConstant("cmd2", false).ToGroup(&grp, func() {
		cmd2 = true
	}); err != nil {
		t.Fatal(err)
	}

	if cmd0 || cmd1 || cmd2 {
		t.FailNow()
	}

	if grp.ExecuteAll([]string{"notcmd"}) != 0 || cmd0 || cmd1 || cmd2 {
		t.FailNow()
	}

	if grp.ExecuteAll([]string{"cmd0"}) != 1 || !cmd0 || cmd1 || cmd2 {
		t.FailNow()
	}

	if grp.ExecuteAll([]string{"cmd1"}) != 1 || !cmd0 || !cmd1 || cmd2 {
		t.FailNow()
	}

	if grp.ExecuteAll([]string{"cmd2"}) != 1 || !cmd0 || !cmd1 || !cmd2 {
		t.FailNow()
	}
}

func TestGroup001(t *testing.T) {

	cmd0, cmd1, cmd2 := false, false, false
	var grp Group

	if err := NewCommand("").AddConstant("cmd", false).ToGroup(&grp, func() {
		cmd0 = true
	}); err != nil {
		t.Fatal(err)
	}

	if err := NewCommand("").AddConstant("cmd", false).ToGroup(&grp, func() {
		cmd1 = true
	}); err != nil {
		t.Fatal(err)
	}

	if err := NewCommand("").AddConstant("cmd", false).ToGroup(&grp, func() {
		cmd2 = true
	}); err != nil {
		t.Fatal(err)
	}

	if cmd0 || cmd1 || cmd2 {
		t.FailNow()
	}

	if grp.ExecuteAll([]string{"notcmd"}) != 0 || cmd0 || cmd1 || cmd2 {
		t.FailNow()
	}

	if grp.ExecuteAll([]string{"cmd"}) != 3 || !cmd0 || !cmd1 || !cmd2 {
		t.FailNow()
	}
}

func TestGroup002(t *testing.T) {

	cmd0, cmd1, cmd2 := false, false, false
	var grp Group

	if err := NewCommand("").AddConstant("cmd0", false).ToGroup(&grp, func() {
		cmd0 = true
	}); err != nil {
		t.Fatal(err)
	}

	if err := NewCommand("").AddConstant("cmd1", false).ToGroup(&grp, func() {
		cmd1 = true
	}); err != nil {
		t.Fatal(err)
	}

	if err := NewCommand("").AddConstant("cmd2", false).ToGroup(&grp, func() {
		cmd2 = true
	}); err != nil {
		t.Fatal(err)
	}

	if cmd0 || cmd1 || cmd2 {
		t.FailNow()
	}

	if grp.ExecuteFirst([]string{"notcmd"}) >= 0 || cmd0 || cmd1 || cmd2 {
		t.FailNow()
	}

	if grp.ExecuteFirst([]string{"cmd0"}) != 0 || !cmd0 || cmd1 || cmd2 {
		t.FailNow()
	}

	cmd0, cmd1, cmd2 = false, false, false

	if grp.ExecuteFirst([]string{"cmd1"}) != 1 || cmd0 || !cmd1 || cmd2 {
		t.FailNow()
	}

	cmd0, cmd1, cmd2 = false, false, false

	if grp.ExecuteFirst([]string{"cmd2"}) != 2 || cmd0 || cmd1 || !cmd2 {
		t.FailNow()
	}
}

func TestCompatibilityDegree000(t *testing.T) {
	c := NewCommand("").AddConstant("test", false).MustCompile()
	n, err := c.Execute([]string{"test"}, nil)

	if err != nil {
		t.Fatal(err)
	}

	if n != 1 {
		t.Fatal("n =", n)
	}
}

func TestCompatibilityDegree001(t *testing.T) {
	c := NewCommand("").AddConstant("const0", false).AddConstant("const1", false).MustCompile()
	n, err := c.Execute([]string{}, nil)

	if err == nil {
		t.FailNow()
	}

	if n != 0 {
		t.Fatal("n =", n)
	}
}

func TestCompatibilityDegree002(t *testing.T) {
	c := NewCommand("").AddConstant("const0", false).AddIntVar("var", 10, "").MustCompile()
	n, err := c.Execute([]string{"const0"}, nil)

	if err == nil {
		t.FailNow()
	}

	if n != 1 {
		t.Fatal("n =", n)
	}
}

func TestCompatibilityDegree003(t *testing.T) {
	c := NewCommand("").AddIntVar("var", 10, "").MustCompile()
	n, err := c.Execute([]string{"asd"}, nil)

	if err == nil {
		t.FailNow()
	}

	if n != 0 {
		t.Fatal("n =", n)
	}
}

func TestCompatibilityDegree004(t *testing.T) {
	c := NewCommand("").AddIntVar("var", 10, "").MustCompile()
	n, err := c.Execute([]string{"123"}, nil)

	if err != nil {
		t.Fatal(err)
	}

	if n != 1 {
		t.Fatal("n =", n)
	}
}

func TestCompatibilityDegree005(t *testing.T) {
	c := NewCommand("").AddBoolFlag("b", "").MustCompile()
	n, err := c.Execute([]string{"-b", "true"}, nil)

	if err != nil {
		t.Fatal(err)
	}

	if n != 2 {
		t.Fatal("n =", n)
	}
}

func TestStrVariadic(t *testing.T) {
	c := NewCommand("").AddConstant("test", false).AddStrVariadic("params", "").MustCompile()

	fail := true

	n, err := c.Execute([]string{"test", "cr0", "cr1", "cr2"}, func(params ...string) {
		fail = false
		if len(params) != 3 || params[0] != "cr0" || params[1] != "cr1" || params[2] != "cr2" {
			t.Fatal("params = ", params)
		}
	})

	if err != nil {
		t.Fatal(err)
	}

	if fail {
		t.Fatal("func not called")
	}

	if n != 4 {
		t.Fatal("n =", n)
	}
}

func TestIntVariadic(t *testing.T) {
	c := NewCommand("").AddConstant("test", false).AddIntVariadic("params", 10, "").MustCompile()

	fail := true

	n, err := c.Execute([]string{"test", "123", "456", "789"}, func(params ...int) {
		fail = false
		if len(params) != 3 || params[0] != 123 || params[1] != 456 || params[2] != 789 {
			t.Fatal("params = ", params)
		}
	})

	if err != nil {
		t.Fatal(err)
	}

	if fail {
		t.Fatal("func not called")
	}

	if n != 4 {
		t.Fatal("n =", n)
	}
}

func TestVariadic000(t *testing.T) {
	c := NewCommand("").AddConstant("test", false).AddStrVariadic("params", "").MustCompile()

	fail := true

	n, err := c.Execute([]string{"test"}, func(params ...string) {
		fail = false
		if len(params) != 0 {
			t.Fatal("len(params) =", len(params))
		}
	})

	if err != nil {
		t.Fatal(err)
	}

	if fail {
		t.Fatal("func not called")
	}

	if n != 1 {
		t.Fatal("n =", n)
	}
}

func TestVariadic001(t *testing.T) {
	c := NewCommand("").AddConstant("test", false).AddStrVar("s", "").AddStrVariadic("params", "").MustCompile()

	fail := true

	n, err := c.Execute([]string{"test", "***"}, func(s string, params ...string) {
		fail = false

		if s != "***" {
			t.Fatal("s =", s)
		}

		if len(params) != 0 {
			t.Fatal("len(params) =", len(params))
		}
	})

	if err != nil {
		t.Fatal(err)
	}

	if fail {
		t.Fatal("func not called")
	}

	if n != 2 {
		t.Fatal("n =", n)
	}
}

func TestVariadic002(t *testing.T) {
	c := NewCommand("").AddConstant("test", false).AddStrVar("s", "").AddStrVariadic("params", "").MustCompile()

	fail := true

	n, err := c.Execute([]string{"test", "***", "$$$"}, func(s string, params ...string) {
		fail = false

		if s != "***" {
			t.Fatal("s =", s)
		}

		if len(params) != 1 || params[0] != "$$$" {
			t.Fatal("params =", params)
		}
	})

	if err != nil {
		t.Fatal(err)
	}

	if fail {
		t.Fatal("func not called")
	}

	if n != 3 {
		t.Fatal("n =", n)
	}
}
