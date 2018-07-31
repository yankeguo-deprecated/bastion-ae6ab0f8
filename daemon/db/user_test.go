package db

import "testing"

func TestUser_BeforeSave(t *testing.T) {
	db, err := Open(temporaryFile())
	db.LogMode(true)
	defer db.Close()
	if err != nil {
		t.Fatal(err)
	}
	if err = db.Migrate(); err != nil {
		t.Fatal(err)
	}
	u := &User{
		Account: "test+user",
		IsAdmin: 1,
	}
	u.SetPassword("qwerty")
	if err = db.Create(u).Error; err == nil {
		t.Fatal("failed to stop bad account")
	}
}
