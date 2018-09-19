package main

import(
  "time"
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
  Current time.Time
  Sum time.Duration
  FormattedSum string
  Punchcard []Punches
  AllPunchcard []Pcard
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
  Pin int
  Punch time.Time
  FormattedPunch string
  Punchtype int
}
//END STRUCTS