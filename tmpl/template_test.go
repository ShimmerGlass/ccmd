package tmpl

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNoVar(t *testing.T) {
	v := "dfghjklm"
	r, err := Exec(v, nil)
	require.Nil(t, err)
	require.Equal(t, v, r)
}

func TestOnlyVar(t *testing.T) {
	r, err := Exec("{var}", map[string]string{"var": "ok"})
	require.Nil(t, err)
	require.Equal(t, "ok", r)
}

func TestUndefinedVar(t *testing.T) {
	_, err := Exec("{var}", map[string]string{"nope": "ok"})
	require.NotNil(t, err)
}

func TestEnclosedVar(t *testing.T) {
	r, err := Exec("abc{var}abc", map[string]string{"var": "ok"})
	require.Nil(t, err)
	require.Equal(t, "abcokabc", r)
}

func TestStartVar(t *testing.T) {
	r, err := Exec("{var}abc", map[string]string{"var": "ok"})
	require.Nil(t, err)
	require.Equal(t, "okabc", r)
}

func TestEndVar(t *testing.T) {
	r, err := Exec("abc{var}", map[string]string{"var": "ok"})
	require.Nil(t, err)
	require.Equal(t, "abcok", r)
}

func TestEscape(t *testing.T) {
	r, err := Exec("{{var}}", nil)
	require.Nil(t, err)
	require.Equal(t, "{var}", r)
}
