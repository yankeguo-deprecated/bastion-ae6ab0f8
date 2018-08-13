package web

import (
	"github.com/novakit/nova"
	"github.com/novakit/view"
	"github.com/pkg/errors"
	"github.com/yankeguo/bastion/types"
	"golang.org/x/crypto/ssh"
	"strings"
)

func routeListKeys(c *nova.Context) (err error) {
	a, ks, v := authResult(c), keyService(c), view.Extract(c)
	var res1 *types.ListKeysResponse
	if res1, err = ks.ListKeys(c.Req.Context(), &types.ListKeysRequest{
		Account: a.User.Account,
	}); err != nil {
		return
	}
	v.Data["keys"] = res1.Keys
	v.DataAsJSON()
	return
}

func routeCreateKey(c *nova.Context) (err error) {
	a, ks, v := authResult(c), keyService(c), view.Extract(c)
	name := strings.TrimSpace(c.Req.FormValue("name"))
	var out ssh.PublicKey
	var cmt string
	if out, cmt, _, _, err = ssh.ParseAuthorizedKey([]byte(c.Req.FormValue("publicKey"))); err != nil {
		return
	}
	if len(name) == 0 {
		name = cmt
	}
	name = strings.TrimSpace(name)
	if len(name) == 0 {
		name = "NO NAME"
	}
	var res1 *types.CreateKeyResponse
	if res1, err = ks.CreateKey(c.Req.Context(), &types.CreateKeyRequest{
		Account:     a.User.Account,
		Fingerprint: ssh.FingerprintSHA256(out),
		Name:        name,
	}); err != nil {
		return
	}
	v.Data["key"] = res1.Key
	v.DataAsJSON()
	return
}

func routeDestroyKey(c *nova.Context) (err error) {
	a, ks, v := authResult(c), keyService(c), view.Extract(c)
	fingerprint := strings.TrimSpace(c.Req.FormValue("fingerprint"))
	var res1 *types.GetKeyResponse
	if res1, err = ks.GetKey(c.Req.Context(), &types.GetKeyRequest{
		Fingerprint: fingerprint,
	}); err != nil {
		return
	}
	if res1.Key.Account != a.User.Account {
		err = errors.New("not your key")
		return
	}
	if res1.Key.Source == types.KeySourceSandbox {
		err = errors.New("can not delete a sandbox key")
		return
	}
	if _, err = ks.DeleteKey(c.Req.Context(), &types.DeleteKeyRequest{
		Fingerprint: res1.Key.Fingerprint,
	}); err != nil {
		return
	}
	v.DataAsJSON()
	return
}
