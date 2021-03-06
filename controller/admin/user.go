package admin

import (
    "github.com/roydong/notes/model"
)

type User struct {
    Base
}

func (c *User) Init() {
    c.Layout = "admin/layout"
}

func (c *User) Setting() {
    c.Base.Init()

    var u *model.User
    if u = c.User(); u == nil {
        panic("You have not signed in")
    }

    c.Render("admin/user/setting", u)
}

func (c *User) Signin() {
    form := new(model.UserForm)
    if c.Request.Method == "POST" {
        form.LoadData(c.Request)
        if form.Valid() {
            m := model.UserModel
            if user := m.FindByEmail(form.Email); user != nil &&
                user.CheckPasswd(form.Passwd) {
                c.Request.Session.Set("user", user, true)
                c.Redirect("/admin/setting", 302)
            }

            form.Message = "email or password error"
        }
    }

    c.Render("admin/user/signin", form)
}

func (c *User) Signup() {
    form := new(model.UserForm)
    if c.Request.Method == "POST" {
        form.LoadData(c.Request)
        if form.Valid() {
            m := model.UserModel
            if m.Exists(form.Email) {
                form.Message = "Email exists"
                goto RENDER
            }

            user := new(model.User)
            user.Name = form.Name
            user.Email = form.Email
            user.SetPasswd(form.Passwd)
            if m.Save(user) {
                c.Request.Session.Set("user", user, true)
                c.Redirect("/admin/setting", 302)
            }

            form.Message = "server error, could not save data"
        }
    }

RENDER:
    c.Render("admin/user/signup", form)
}
