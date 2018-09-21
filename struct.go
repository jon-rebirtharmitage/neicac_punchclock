package main

import(
  "time"
  "gopkg.in/mgo.v2/bson"
)

//STUCTS
/*
TYPE : Page
struct for use with HTTP/TEMPLATE to display web pages.  Webpages internal data is stored here.
*/
type Page struct {
	Title string
	Body  string
}

type UserPage struct {
	Title string
	Body  string
  Content []string
}

type timecardPage struct {
	Title string
	Body  string
  Username string
  Fname, Lname string
  Pin, Status int
  Department string
  Current time.Time
  Sum time.Duration
  FormattedSum string
  Punchcard []Punches
  AllPunchcard []Pcard
}

type timecardReport struct {
  Startdate string `json:"startdate"`
  Enddate   string `json:"enddate"`
  Pin       string `json:"pin"`
}

type userPage struct {
  Title string
  Body string
  Users []neicacUser
}

type adminPage struct {
  Title string
  Body string
  Admins []adminUser
}

type neicacUser struct{
  Username  string
  Fname, Lname string
  Department string
  Pin, Status int
}

type newUser struct {
	Firstname       string `json:"firstname"`
	Lastname        string `json:"lastname"`
	Pin             string `json:"pin"`
	Department      string `json:"Department"`
}

type adminUser struct{
  Username string
  Fname string
  Lname string
  Password string
  Department string
  Role string
  Auth bool
}

type newAdmin struct {
	Firstname       string `json:"firstname"`
	Lastname        string `json:"lastname"`
	Username        string `json:"username"`
	Password        string `json:"Password"`
	PasswordConfirm string `json:"PasswordConfirm"`
	Department      string `json:"Department"`
}

type Pcard struct{
  Startdate time.Time
  Enddate time.Time
  Sum time.Duration
  FormattedSum string
  Title string
  Punchf []Punches
}

type Punches struct{
  Id bson.ObjectId      `json:"id" bson:"_id,omitempty"`
  Pin int               `json:"pin"`
  Punch time.Time       `json:"punch"`
  FormattedPunch string `json:"formattedPunch"`
  Punchtype int         `json:"punchtype"`
}

type EPunch struct{
  Id string             `json:"id"`
  Pin int            `json:"pin"`
  Punch string          `json:"punch"`
  FormattedPunch string `json:"formattedPunch"`
  Punchtype int      `json:"punchtype"`
}
//END STRUCTS