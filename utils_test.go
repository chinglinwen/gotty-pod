package main

import (
	"fmt"
	"testing"
)

func TestWalk(t *testing.T) {
	list, err := Walk("/data/fluentd")
	if err != nil {
		t.Error("walk err", err)
		return
	}
	fmt.Println(list)
}

func TestListPods(t *testing.T) {
	list, err := listpods()
	if err != nil {
		t.Error("walk err", err)
		return
	}
	fmt.Println(list)
}

var gitlist = "H4sIAAAAAAAA/3RW7W7cuA59lXkB1fcWuED+NneQTdPtomi6wAKLxYCWaIuxTLr6GI/z9Av5K+Nx+mfMcw5JURJH0t/QdcXz8YEYWOMhoyCJTQ/nCVlzih50cwI2XshsSZJwAK8tnbFoUomeMeIbBRyodKjEmZUzohv0Fc10CQFPAf2ZdPYHN0TSQWUG/VbtvICOdMaTBue22jU4lFQbiFBgdKvdoQ/CEHKsPiFHL92wqg1UDSiPhsKhdKIbbYG4GH8VOPTxoKUYEvdI2crxg6oTGcwwRGAD3oQMIoZIXM+8blTqag+TI3HMNWKfQY/YuEF57MTHOVRbFSz4NX5weD1MMlej9uNvizgOVznpl00qHo8/8v6g39LPxyeCB+D6HmSr/IF9+OblBXXcCl8Qu09aS+Ib4RHkMcFPAt7x3+W2ngcHHCxx/bvs/I+fUrTPxy+7iGD33s/Hr/BEsKv/8fgdS4j4NXfGqGjMa11oR7pRFsQAKS3tRnQCrLxwXcJOGxIr6GjD3akh8c7rF8khGYq7FEZ6dgJmS1YqlzJVcquMM/gwDfLhdhBTKdG3zN2H2zJNpd6bzzzuO7R0uKNriNjDsOF+m7hP3z5v514jx1MYQsT21GO5Ee3/3k0frPTmZjJjHjXlmQSSUPxAaP/s7tPwRj0fT/8HekhPtOHuge5B/iJYGjLzj8dTNZ53BC5E8biVcuO5tfHmRJ85JD8ekZv098+n/M/ZkF+f7k9T43rhqJBNUX1EFSJE0lekNWo8RNW4Pm++i6UclYdaHHBdREvejNs3E1raVnhB0yymFiloKVVp3y5kj6WKIi4sxJVWglwoB3TEtWqJaVFaIP74n//eLfiVLok3HuOHB+BLAlbLHfGumO+K94Rd3QvxQly/kMfVY624Ha+OqyIXe7kBZgimfdPWdLdlvtpUgvyC3Vad16lOwC8Eq78jYGMK4rOQxtO4f8grD3QykmvcKXPEzvHQ2a4IDTmnWugmZMHjaBk8zxsZOjL5bzN+xibaMnOLbMlpRWaus+J0WjNV6P2wgNxrZ+qKgTi32Gg78HBGp/JlO2fKtE1SwmTWEC16VYo0Iw42NThaZ+rUsjsZl/l6YpNtVXtJXWHdzyt0eX3x+ZjwEEn4RriChl5tusIp5KmPp8e1V6JXm6dx6cBHRl8sxmG61IsqsESqhgU7qVVHTuJCaOGKatUKU5Q1bNqYGeRpt+CbsBA2tYvZsPQOTb06Tx/FNfFlzTa/I9agtyfVzCwvqs7BWmpN0UGpNKnmLqj8/FiU6C835RvsnLwTeW43gW+PtJv4Tbm9+CZe+Ux9udQ+JeyRXy1y7YiLHVFLjMPWxQNW1KhKfA/5yEvRbhyc1MThrPd5lJM6HCxG4domKkjCaZBkQX4m4H/+BQAA//8BAAD//7hfUhnqCgAA"

func TestUnCompress(t *testing.T) {
	b, err := UnCompress(gitlist)
	if err != nil {
		t.Error("uncompress err", err)
		return
	}
	fmt.Println(b)
}

func TestCompress(t *testing.T) {
	a := Compress("aa")
	b, _ := UnCompress(a)
	if b != "aa" {
		t.Error("error")
	}
}

var userinfo = "eyJpZCI6NzUsInVzZXJuYW1lIjoid2VuemhlbmdsaW4iLCJlbWFpbCI6IndlbnpoZW5nbGluQGhhb2RhaS5uZXQiLCJuYW1lIjoid2VuemhlbmdsaW4iLCJzdGF0ZSI6ImFjdGl2ZSIsImNyZWF0ZWRfYXQiOiIyMDE4LTEyLTEwVDAzOjExOjQzLjQ3WiIsImJpbyI6IiIsImxvY2F0aW9uIjoiIiwicHVibGljX2VtYWlsIjoiIiwic2t5cGUiOiIiLCJsaW5rZWRpbiI6IiIsInR3aXR0ZXIiOiIiLCJ3ZWJzaXRlX3VybCI6IiIsIm9yZ2FuaXphdGlvbiI6IiIsImV4dGVybl91aWQiOiIiLCJwcm92aWRlciI6IiIsInRoZW1lX2lkIjoxLCJsYXN0X2FjdGl2aXR5X29uIjoiMjAxOS0wMi0xOCIsImNvbG9yX3NjaGVtZV9pZCI6MSwiaXNfYWRtaW4iOnRydWUsImF2YXRhcl91cmwiOiJodHRwczovL3d3dy5ncmF2YXRhci5jb20vYXZhdGFyLzhhOGYwNzkxMDZkZDFlNTcxYmJkYjQ2ZDA0OThhYWJlP3M9ODBcdTAwMjZkPWlkZW50aWNvbiIsImNhbl9jcmVhdGVfZ3JvdXAiOnRydWUsImNhbl9jcmVhdGVfcHJvamVjdCI6dHJ1ZSwicHJvamVjdHNfbGltaXQiOjEwMDAwMCwiY3VycmVudF9zaWduX2luX2F0IjoiMjAxOS0wMi0xOFQwODowNjo1MC41MjNaIiwibGFzdF9zaWduX2luX2F0IjoiMjAxOS0wMi0xNFQwNjoxNzoyNC44NzZaIiwiY29uZmlybWVkX2F0IjoiMjAxOC0xMi0xMFQwMzoxMTo0My40MjNaIiwidHdvX2ZhY3Rvcl9lbmFibGVkIjpmYWxzZSwiaWRlbnRpdGllcyI6W3sicHJvdmlkZXIiOiJsZGFwbWFpbiIsImV4dGVybl91aWQiOiJjbj13ZW56aGVuZ2xpbixvdT3mioDmnK/kuK3lv4Msb3U9cGVvcGxlLGRjPWhhb2RhaSxkYz1uZXQifV0sImV4dGVybmFsIjpmYWxzZSwicHJpdmF0ZV9wcm9maWxlIjpmYWxzZSwic2hhcmVkX3J1bm5lcnNfbWludXRlc19saW1pdCI6MH0="

/*
{"id":75,"username":"wenzhenglin","email":"wenzhenglin@haodai.net","name":"wenzhenglin","state":"active","created_at":"2018-12-10T03:11:43.47Z","bio":"","location":"","public_email":"","skype":"","linkedin":"","twitter":"","website_url":"","organization":"","extern_uid":"","provider":"","theme_id":1,"last_activity_on":"2019-02-18","color_scheme_id":1,"is_admin":true,"avatar_url":"https://www.gravatar.com/avatar/8a8f079106dd1e571bbdb46d0498aabe?s=80\u0026d=identicon","can_create_group":true,"can_create_project":true,"projects_limit":100000,"current_sign_in_at":"2019-02-18T08:06:50.523Z","last_sign_in_at":"2019-02-14T06:17:24.876Z","confirmed_at":"2018-12-10T03:11:43.423Z","two_factor_enabled":false,"identities":[{"provider":"ldapmain","extern_uid":"cn=wenzhenglin,ou=技术中心,ou=people,dc=haodai,dc=net"}],"external":false,"private_profile":false,"shared_runners_minutes_limit":0}
*/
func TestParseUser(t *testing.T) {
	name, id, err := ParseUserInfo(userinfo)
	if err != nil {
		t.Error("parse user err", err)
		return
	}
	if id != 75 {
		t.Error("id error, want 75, got: ", id)
		return
	}
	fmt.Println("got ", name, id)
}
