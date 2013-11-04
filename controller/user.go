package controller

import (
    "net/http"
    "github.com/roydong/potato"
    "github.com/roydong/topic/model"
)

type User struct {
    Base
}

func (c *User) Setting() {
    var u *model.User
    if u = c.User(); u == nil {
        potato.Panic(http.StatusUnauthorized, "You have not signed in")
    }

    c.Title = "Admin - Setting"
    c.Render("user/setting", u)
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
                c.Redirect("/setting")
                return
            }

            form.Message = "email or password wrong"
        }
    }

    c.Title = "Admin - Sign in"
    c.Render("user/signin", form)
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
            potato.L.Println(form, user)
            user.Email = form.Email
            user.SetPasswd(form.Passwd)
            if m.Save(user) {
                c.Request.Session.Set("user", user, true)
                c.Redirect("/setting")
                return
            }

            form.Message = "server error, could not save data"
        }
    }

    RENDER:
        c.Title = "Admin - Sign up"
        c.Render("user/signup", form)
}
