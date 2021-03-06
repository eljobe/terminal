package terminal

import (
    "os"
    "reflect"
    "runtime"
    "testing"
    "strings"
)

var testCases = []struct {
    envVars []string
    r16 bool
    r256 bool
    r16million bool
}{
    {[]string{}, false, false, false},
    {[]string{"TERM=anything-256color"}, true, true, false},
    {[]string{"TERM=xterm-256"}, true, true, false},
    {[]string{"TERM=screen"}, true, false, false},
    {[]string{"TERM=some-screen-thing"}, false, false, false},
    {[]string{"TERM=xterm"}, true, false, false},
    {[]string{"TERM=some-xterm-thing"}, false, false, false},
    {[]string{"TERM=vt100"}, true, false, false},
    {[]string{"TERM=some-vt100-thing"}, false, false, false},
    {[]string{"TERM=ansi"}, true, false, false},
    {[]string{"TERM=some-ansi-thing"}, true, false, false},
    {[]string{"TERM=color"}, true, false, false},
    {[]string{"TERM=some-color-thing"}, true, false, false},
    {[]string{"TERM=cygwin"}, true, false, false},
    {[]string{"TERM=some-cygwin-thing"}, true, false, false},
    {[]string{"TERM=linux"}, true, false, false},
    {[]string{"TERM=some-linux-thing"}, true, false, false},
    {[]string{"TERM=dumb"}, false, false, false},
    {[]string{"TERM=dumb","TERM_PROGRAM=Apple_Terminal"}, true, true, false},
    {[]string{"TERM_PROGRAM=Hyper"}, true, true, true},
    {[]string{"TERM_PROGRAM=iTerm.app", "TERM_PROGRAM_VERSION=2.3"}, true, true, false},
    {[]string{"TERM_PROGRAM=iTerm.app", "TERM_PROGRAM_VERSION=3.1"}, true, true, true},
    {[]string{"TERM_PROGRAM=iTerm.app", "TERM_PROGRAM_VERSION=20.1"}, true, true, true},
    {[]string{"TERM=dumb","COLORTERM=256"}, true, true, false},
    {[]string{"COLORTERM=none"}, false, false, false},
    {[]string{"COLORTERM=basic"}, true, false, false},
    {[]string{"COLORTERM=4bit"}, true, false, false},
    {[]string{"COLORTERM=256"}, true, true, false},
    {[]string{"COLORTERM=8bit"}, true, true, false},
    {[]string{"COLORTERM=16m"}, true, true, true},
    {[]string{"COLORTERM=24bit"}, true, true, true},
    {[]string{"COLORTERM=Truecolor"}, true, true, true},
    {[]string{"COLORTERM=truecolor"}, true, true, true},
    {[]string{"COLORTERM=256", "LC_COLORTERM=16m"}, true, true, false},
    {[]string{"LC_COLORTERM=16m"}, true, true, true},
    {[]string{"USER_COLORTERM=256"}, true, true, false},
    {[]string{"LC_USER_COLORTERM=256"}, true, true, false},
    {[]string{"USER_COLORTERM=256", "LC_USER_COLORTERM=16m"}, true, true, false},
    {[]string{"COLORTERM=256", "USER_COLORTERM=4bit"}, true, false, false},
    {[]string{"COLORTERM=256", "USER_COLORTERM=16m"}, true, true, false},
    {[]string{"COLORTERM=256", "LC_USER_COLORTERM=4bit"}, true, false, false},
    {[]string{"COLORTERM=256", "LC_USER_COLORTERM=16m"}, true, true, false},
}

func TestColorSupport(t *testing.T) {
    for _, tc := range testCases {
        clearEnvironment()
        populateEnvironment(tc.envVars)
        cs := ColorSupport()
        pass := assertSupport(cs.Supports16Colors, tc.r16, t)
        pass = assertSupport(cs.Supports256Colors, tc.r256, t) && pass
        pass = assertSupport(cs.SupportsTrueColor, tc.r16million, t) && pass
        if !pass {
            t.Logf("Environment: %v", tc.envVars)
        }
    }
}

func assertSupport(check func() bool, expected bool, t *testing.T) bool {
    actual := check()
    if actual != expected {
        t.Errorf("%v: %t, wanted %t", funcName(check), actual, expected)
        return false
    }
    return true
}

func populateEnvironment(envVars []string) {
    for _, e := range envVars {
        parts := strings.Split(e, "=")
        k, v := parts[0], parts[1]
        os.Setenv(k, v)
    }
}

func clearEnvironment() {
    os.Unsetenv("TERM")
    os.Unsetenv("TERM_PROGRAM")
    os.Unsetenv("TERM_PROGRAM_VERSION")
    os.Unsetenv("COLORTERM")
    os.Unsetenv("USER_COLORTERM")
    os.Unsetenv("LC_COLORTERM")
    os.Unsetenv("LC_USER_COLORTERM")
}

func funcName(i interface{}) string {
    return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}
