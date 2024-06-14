package main

import (
	"fmt"
	"github.com/george012/fltk_go"
	"log"
	"strings"
)

const (
	WIDGET_HEIGHT  = 25
	WIDGET_PADDING = 5
	WIDGET_WIDTH   = 80
)

var db = NewDataBase()

func init() {
	db.CreateUser("Emil", "Hans")
	db.CreateUser("Mustermann", "Max")
	db.CreateUser("Tisch", "Roman")
}

func main() {
	fltk_go.SetScheme("gtk+")

	win := fltk_go.NewWindow(
		WIDGET_WIDTH*4+WIDGET_PADDING*3,
		WIDGET_HEIGHT*6+WIDGET_PADDING*4)
	win.SetLabel("CURD")

	p := NewCrudPanel(win)
	p.Bind(db)

	win.End()
	win.Show()
	fltk_go.Run()
}

type CrudPanel struct {
	prefixInput  *fltk_go.Input
	listBrowser  *fltk_go.HoldBrowser
	nameInput    *fltk_go.Input
	surnameInput *fltk_go.Input
	createBtn    *fltk_go.Button
	updateBtn    *fltk_go.Button
	deleteBtn    *fltk_go.Button
}

func NewCrudPanel(win *fltk_go.Window) *CrudPanel {
	p := &CrudPanel{}

	win.Begin()

	col := fltk_go.NewFlex(WIDGET_PADDING, WIDGET_PADDING, win.W()-WIDGET_PADDING*2, win.H()-WIDGET_PADDING*2)
	col.SetGap(WIDGET_PADDING)

	{
		row := fltk_go.NewFlex(0, 0, 0, 0)
		col.Fixed(row, WIDGET_HEIGHT)
		row.SetType(fltk_go.ROW)
		fltk_go.NewBox(fltk_go.NO_BOX, 0, 0, 0, 0, "Filter prefix:").SetAlign(fltk_go.ALIGN_INSIDE | fltk_go.ALIGN_LEFT)
		p.prefixInput = fltk_go.NewInput(0, 0, 0, 0)
		fltk_go.NewBox(fltk_go.NO_BOX, 0, 0, 0, 0) // invisible
		fltk_go.NewBox(fltk_go.NO_BOX, 0, 0, 0, 0) // invisible
		row.End()
	}

	{
		row := fltk_go.NewFlex(0, 0, 0, 0*4)
		row.SetType(fltk_go.ROW)
		row.SetGap(WIDGET_PADDING)
		p.listBrowser = fltk_go.NewHoldBrowser(0, 0, 0, 0)
		{
			col := fltk_go.NewFlex(0, 0, 0, 0)
			{
				row := fltk_go.NewFlex(0, 0, 0, 0)
				col.Fixed(row, WIDGET_HEIGHT)
				row.SetType(fltk_go.ROW)
				fltk_go.NewBox(fltk_go.NO_BOX, 0, 0, 0, 0, "Name:").SetAlign(fltk_go.ALIGN_INSIDE | fltk_go.ALIGN_LEFT)
				p.nameInput = fltk_go.NewInput(0, 0, 0, 0)
				row.End()
			}
			{
				row := fltk_go.NewFlex(0, 0, 0, 0)
				col.Fixed(row, WIDGET_HEIGHT)
				row.SetType(fltk_go.ROW)
				fltk_go.NewBox(fltk_go.NO_BOX, 0, 0, 0, 0, "Surname:").SetAlign(fltk_go.ALIGN_INSIDE | fltk_go.ALIGN_LEFT)
				p.surnameInput = fltk_go.NewInput(0, 0, 0, 0)
				row.End()
			}
			col.End()
		}

		row.End()
	}

	{
		row := fltk_go.NewFlex(0, 0, 0, 0)
		col.Fixed(row, WIDGET_HEIGHT)
		row.SetType(fltk_go.ROW)
		row.SetGap(WIDGET_PADDING)
		p.createBtn = fltk_go.NewButton(0, 0, 0, 0)
		p.createBtn.SetLabel("Create")

		p.updateBtn = fltk_go.NewButton(0, 0, 0, 0)
		p.updateBtn.SetLabel("Update")

		p.deleteBtn = fltk_go.NewButton(0, 0, 0, 0)
		p.deleteBtn.SetLabel("Delete")
		fltk_go.NewBox(fltk_go.NO_BOX, 0, 0, 0, 0) // invisible
		row.End()
	}

	col.End()
	win.Resizable(col)
	return p
}

func (p *CrudPanel) Bind(db *DataBase) {
	p.prefixInput.SetCallbackCondition(fltk_go.WhenChanged)
	p.prefixInput.SetCallback(func() {
		p.ReflushList(db)
	})

	p.createBtn.SetCallback(func() {
		if p.nameInput.Value() != "" && p.surnameInput.Value() != "" {
			db.CreateUser(p.nameInput.Value(), p.surnameInput.Value())
			p.ReflushList(db)
		}
	})

	p.updateBtn.SetCallback(func() {
		line := p.listBrowser.Value()
		id := p.listBrowser.Data(line).(int)
		if db.UpdateUser(id, p.nameInput.Value(), p.surnameInput.Value()) {
			p.ReflushList(db)
		}
	})

	p.deleteBtn.SetCallback(func() {
		line := p.listBrowser.Value()
		id := p.listBrowser.Data(line).(int)
		if db.DeleteUser(id) {
			p.ReflushList(db)
		}
	})

	p.listBrowser.SetCallbackCondition(fltk_go.WhenChanged)
	p.listBrowser.SetCallback(func() {
		p.Update(db)
	})

	p.ReflushList(db)
}

func (p *CrudPanel) ReflushList(db *DataBase) {
	users := db.QueryUsers(p.prefixInput.Value())
	p.listBrowser.Clear()
	for _, user := range users {
		p.listBrowser.AddWithData(fmt.Sprintf("%s, %s", user.Name, user.Surname), user.ID)
	}

	p.Update(db)
}

func (p *CrudPanel) Update(db *DataBase) {
	if p.listBrowser.Value() == 0 {
		p.deleteBtn.Deactivate()
		p.updateBtn.Deactivate()
	} else {
		p.deleteBtn.Activate()
		p.updateBtn.Activate()
		id := p.listBrowser.Data(p.listBrowser.Value()).(int)
		user := db.FindUser(id)
		if user != nil {
			p.nameInput.SetValue(user.Name)
			p.surnameInput.SetValue(user.Surname)
		}
	}
}

type DataBase struct {
	UserList []*User
	lastID   int
}

func NewDataBase() *DataBase {
	db := &DataBase{}
	db.UserList = make([]*User, 0)
	db.lastID = 0
	return db
}

type User struct {
	ID      int
	Name    string
	Surname string
}

func (db *DataBase) CreateUser(name, surname string) int {
	user := &User{}
	user.Name = name
	user.Surname = surname
	db.UserList = append(db.UserList, user)
	user.ID = db.lastID + 1
	db.lastID++
	return user.ID
}

func (db *DataBase) FindUser(id int) *User {
	for _, user := range db.UserList {
		if user.ID == id {
			return user
		}
	}

	return nil
}

func (db *DataBase) UpdateUser(id int, name, surname string) bool {
	user := db.FindUser(id)
	if user == nil {
		log.Printf("fail to update: user with ID(%d) is not found", id)
		return false
	}

	user.Name = name
	user.Surname = surname
	return true
}

func (db *DataBase) DeleteUser(id int) bool {
	index := -1
	for index = range db.UserList {
		if db.UserList[index].ID == id {
			break
		}
	}

	if index < 0 {
		log.Printf("fail to delete: user with ID(%d) is not found", id)
		return false
	}

	db.UserList = append(db.UserList[:index], db.UserList[index+1:]...)
	return true
}

func (db *DataBase) QueryUsers(prefix string) []*User {
	if prefix == "" {
		return db.UserList
	}

	users := make([]*User, 0, len(db.UserList))
	for _, user := range db.UserList {
		if strings.HasPrefix(user.Surname, prefix) {
			users = append(users, user)
		}
	}

	return users
}
