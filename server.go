package main

import(
	"html/template"
	"net/http"
  "github.com/gorilla/mux"
  "github.com/gorilla/sessions"
  "strconv"
	"fmt"
  "time"
	"math/rand"
  "encoding/json"
  "io/ioutil"
  "strings"
)

//GLOBAL VARIABLES
//Session Storage for secure session variable handling using AES-256 Encryption
var store = sessions.NewCookieStore([]byte("30efe06c8b1c6a0cdca9298090d551cd86547a90435d15fd90fb3765af766072"))

//Custom Page Loaders
func loadPage(title string) (*Page, error){
	return &Page{Title: title, Body: "blank..."}, nil
}

func loadTimecard(tp timecardPage) (*timecardPage, error){
	return &tp, nil
}

func loadUsers(up userPage) (*userPage, error){
  return &up, nil
}

func loadAdmins(ap adminPage) (*adminPage, error){
  return &ap, nil
}

func loadTimecardEdit(tp timecardPage) (*timecardPage, error){
  return &tp, nil
}

func loadEditPunchPage(p Punches) (*Punches, error) {
  return &p, nil
}

func loadAddPunchPage(p Punches) (*Punches, error) {
  return &p, nil
}

func loadEditUserPage(p neicacUser) (*neicacUser, error) {
  return &p, nil
}

func loadEditAdminPage(p adminUser) (*adminUser, error) {
  return &p, nil
}
//END Custom Page Loaders

//GENERIC Page Loader Section
func indexViewHandler(w http.ResponseWriter, r *http.Request) {
	p, _ := loadPage("Test")
  renderTemplate(w, "./html/index", p)
}

func failedViewHandler(w http.ResponseWriter, r *http.Request) {
	p, _ := loadPage("Test")
  renderTemplate(w, "./html/failedLogin", p)
}

func adminViewHandler(w http.ResponseWriter, r *http.Request) {
	p, _ := loadPage("Test")
  renderTemplate(w, "./html/admin", p)
}

func adminLoginViewHandler(w http.ResponseWriter, r *http.Request) {
	p, _ := loadPage("Test")
  renderTemplate(w, "./html/adminLogin", p)
}

func punchInViewHandler(w http.ResponseWriter, r *http.Request) {
  p, _ := loadPage("Punch")
  renderTemplate(w, "./html/punchIn", p)
}

func punchOutViewHandler(w http.ResponseWriter, r *http.Request) {
  p, _ := loadPage("Punch")
  renderTemplate(w, "./html/punchOut", p)
}

func createAdminViewHandler(w http.ResponseWriter, r *http.Request) {
  p, _ := loadPage("Punch")
  renderTemplate(w, "./html/createAdmin", p)
}

func createAdminSuccessViewHandler(w http.ResponseWriter, r *http.Request) {
	p, _ := loadPage("Test")
  renderTemplate(w, "./html/createAdminSuccess", p)
}

func createAdminFailedViewHandler(w http.ResponseWriter, r *http.Request) {
	p, _ := loadPage("Test")
  renderTemplate(w, "./html/createAdminFailed", p)
}

func createUserViewHandler(w http.ResponseWriter, r *http.Request) {
  p, _ := loadPage("Punch")
  renderTemplate(w, "./html/createUser", p)
}

func createUserSuccessViewHandler(w http.ResponseWriter, r *http.Request) {
	p, _ := loadPage("Test")
  renderTemplate(w, "./html/createUserSuccess", p)
}

func createUserFailedViewHandler(w http.ResponseWriter, r *http.Request) {
	p, _ := loadPage("Test")
  renderTemplate(w, "./html/createUserFailed", p)
}
//End Generic Page Loader Section

//Timecard Viewer for users
func timecardViewHandler(w http.ResponseWriter, r *http.Request) {
  session, _ := store.Get(r, "neicac_punchcard")
  tp := timecardPage{}
  tp.Title = session.Values["username"].(string)
  tp.Body = session.Values["username"].(string)
  tp.Username = session.Values["username"].(string)
  tp.Pin = session.Values["pin"].(int)
  tp.Status = session.Values["status"].(int)
  tp.Fname = session.Values["fname"].(string)
  tp.Lname = session.Values["lname"].(string)
  tp.Current = time.Now()
  tp.Punchcard = getPunches(tp, 1)
  tp.AllPunchcard = filterPunchcard(tp)
  tp.Punchcard = filterPunchcardCurrent(tp)
  tp.Sum = findHoursWorked(tp)
  tp.FormattedSum = fmtDuration(tp.Sum)
  for m := range tp.AllPunchcard{
    fmt.Println(tp.AllPunchcard[m])
  }
  for j := range tp.Punchcard{
    tp.Punchcard[j].FormattedPunch = tp.Punchcard[j].Punch.Format("Mon Jan _2 3:04PM")
  }
  for k := range tp.AllPunchcard{
    for m := range tp.AllPunchcard[k].Punchf{
      tp.AllPunchcard[k].Punchf[m].FormattedPunch = tp.AllPunchcard[k].Punchf[m].Punch.Format("Mon Jan _2 3:05 PM 2006")
    }
  }
  p, _ := loadTimecard(tp)
  renderTimecard(w, "./html/timecard", p)
}

//Editable Time Card
func timecardEditViewHandler(w http.ResponseWriter, r *http.Request) {
  tp := timecardPage{}
  vars := mux.Vars(r)
  i, err := strconv.Atoi(vars["Value"])
  if err != nil {
    http.Redirect(w, r, "http://rebirtharmitage.com:8084/generalError", 301)
  }
  a := checkStatus(i)  
  tp.Title = a.Username
  tp.Body = a.Username
  tp.Username = a.Username
  tp.Pin = i
  tp.Status = a.Status
  tp.Fname = a.Fname
  tp.Lname = a.Lname
  tp.Department = a.Department
  tp.Current = time.Now()
  tp.Punchcard = getPunches(tp, 1)
  tp.AllPunchcard = filterPunchcard(tp)
  tp.Punchcard = filterPunchcardCurrent(tp)
  tp.Sum = findHoursWorked(tp)
  tp.FormattedSum = fmtDuration(tp.Sum)
  for m := range tp.AllPunchcard{
    fmt.Println(tp.AllPunchcard[m])
  }
  for j := range tp.Punchcard{
    tp.Punchcard[j].FormattedPunch = tp.Punchcard[j].Punch.Format("Mon Jan _2 3:04PM")
  }
  for k := range tp.AllPunchcard{
    for m := range tp.AllPunchcard[k].Punchf{
      tp.AllPunchcard[k].Punchf[m].FormattedPunch = tp.AllPunchcard[k].Punchf[m].Punch.Format("Mon Jan _2 3:04 PM")
    }
  }
  p, _ := loadTimecardEdit(tp)
  renderTimecardEdit(w, "./html/timecardEdit", p)
}

//HERE
func listUsersViewHandler(w http.ResponseWriter, r *http.Request) {
  session, _ := store.Get(r, "neicac_punchcard")
  up := userPage{}
  d := ""
  if session.Values["department"].(string) == "indirect" {
    d = "all"
    up.Title = "List of Users in all departments"
  }else{
    d = session.Values["department"].(string)
    up.Title = "List of Users in Department : " + session.Values["department"].(string)
  }
  up.Body = session.Values["department"].(string)
  up.Users = getUsers(d)
  p, _ := loadUsers(up)
  renderUsers(w, "./html/listUsers", p)
}

func listAdminsViewHandler(w http.ResponseWriter, r *http.Request) {
  session, _ := store.Get(r, "neicac_punchcard")
  ap := adminPage{}
  d := ""
  if session.Values["department"].(string) == "indirect" {
    d = "all"
    ap.Title = "List of Admins in all departments"
  }else{
    d = session.Values["department"].(string)
    ap.Title = "List of Admins in Department : " + session.Values["department"].(string)
  }
  ap.Body = session.Values["department"].(string)
  ap.Admins = getAdmins(d)
  p, _ := loadAdmins(ap)
  renderAdmins(w, "./html/listAdmins", p)
}

func editPunchPageViewHandler(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Access-Control-Allow-Origin", "*")
  vars := mux.Vars(r)
  s := vars["Value"]
  fmt.Println(s)
  b := Punches{}
  a := loadPunch(s)
  if a == b {
    http.Redirect(w, r, "http://rebirtharmitage.com:8084/generalError", 301)
  }else{
    a.FormattedPunch = a.Punch.Format("2006-01-02 15:04")
    a.FormattedPunch = strings.Replace(a.FormattedPunch, " ", "T", -1)
    fmt.Println(a.FormattedPunch)
    p, _ := loadEditPunchPage(a)
    renderEditPunchPage(w, "./html/editPunchPage", p)
  }
}

//FUNCTIONAL CODE SECTION
//Contains functions to act on data outside of the web interfaces
func filterPunchcardCurrent(tp timecardPage) []Punches{
  dp := []Punches{}
  fromDate := time.Now().AddDate(0, 0, 0)
  fromDate = time.Date(fromDate.Year(), fromDate.Month(), fromDate.Day(), 0, 0, 0, 0, fromDate.Location())
  toDate := time.Now()
  for j := range tp.Punchcard{
    fmt.Println(tp.Punchcard[j].Punch)
    if inTimeSpan(fromDate, toDate, tp.Punchcard[j].Punch){
      dp = append(dp, tp.Punchcard[j])
    }
  }
  return dp
}

func findHoursWorked(tp timecardPage) time.Duration{
  a := time.Now()
  b := time.Now()
  q := a.Sub(b)
  lp := 0
  for k := range tp.Punchcard {
    //If the punch was a punchin
    if tp.Punchcard[k].Punchtype == 0 {
      lp = 0
      a = tp.Punchcard[k].Punch
    }else{
      lp = 1
      b = tp.Punchcard[k].Punch
      q = q + b.Sub(a)
    }
  }
  if lp == 0{
    q = q + time.Now().Sub(a)
  }
  return q
}

func findHoursWorkedAll(tp Pcard) time.Duration{
  a := time.Now()
  b := time.Now()
  q := a.Sub(b)
  for k := range tp.Punchf {
    //If the punch was a punchin
    if tp.Punchf[k].Punchtype == 0 {
      a = tp.Punchf[k].Punch
    }else{
      b = tp.Punchf[k].Punch
      q = q + b.Sub(a)
    }
  }
  return q
}

//Function to format time duration to readable and rounded numbers
func fmtDuration(d time.Duration) string {
    d = d.Round(time.Minute)
    h := d / time.Hour
    d -= h * time.Hour
    m := d / time.Minute
    return fmt.Sprintf("%02d:%02d", h, m)
}

//Goes through all punches to those that are within each range
func filterPunchcard(tp timecardPage) []Pcard {
  fromDate := time.Now().AddDate(0, 0, -14)
  pcard := []Pcard{}
  ctime := time.Now()
  ctime = time.Date(ctime.Year(), ctime.Month(), ctime.Day(), 0, 0, 0, 0, ctime.Location())
  dtime := time.Now()
  dtime = time.Date(dtime.Year(), dtime.Month(), dtime.Day(), 0, 0, 0, 0, dtime.Location())
  diff := time.Now().Sub(fromDate)
  c := int(diff.Hours())/24
  fmt.Println(c)
  for j := 0; j <= c; j++ {
    p := Pcard{}
    ctime = ctime.AddDate(0, 0, -j)
    dtime = ctime.AddDate(0, 0, -1)
    p.Startdate = dtime
    p.Enddate = ctime
    p.Title = "Punchcard from : " + dtime.Format("Mon Jan _2 2006")
    fmt.Println(j)
    for k := range tp.Punchcard{
      if inTimeSpan(dtime, ctime, tp.Punchcard[k].Punch){
        p.Punchf = append(p.Punchf, tp.Punchcard[k])
      }
    }
    if p.Punchf != nil{
      pcard = append(pcard, p)
    }
    ctime = time.Now()
    ctime = time.Date(ctime.Year(), ctime.Month(), ctime.Day(), 0, 0, 0, 0, ctime.Location())
    dtime = time.Now()
    dtime = time.Date(dtime.Year(), dtime.Month(), dtime.Day(), 0, 0, 0, 0, dtime.Location())
  }
  for m := range pcard{
    s := findHoursWorkedAll(pcard[m])
    t := fmtDuration(s)
    pcard[m].Sum = s
    pcard[m].FormattedSum = t
  }
  return pcard
}


//RENDER SECTION
//Renders all html pages from the system
func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
  t, _ := template.ParseFiles(tmpl + ".html")
  t.Execute(w, p)
}

func renderTimecard(w http.ResponseWriter, tmpl string, p *timecardPage) {
  t, _ := template.ParseFiles(tmpl + ".html")
  t.Execute(w, p)
}

func renderUsers(w http.ResponseWriter, tmpl string, p *userPage) {
  t, _ := template.ParseFiles(tmpl + ".html")
  t.Execute(w, p)
}

func renderAdmins(w http.ResponseWriter, tmpl string, p *adminPage) {
  t, _ := template.ParseFiles(tmpl + ".html")
  t.Execute(w, p)
}

func renderTimecardEdit(w http.ResponseWriter, tmpl string, p *timecardPage) {
  t, _ := template.ParseFiles(tmpl + ".html")
  t.Execute(w, p)
}

func renderEditPunchPage(w http.ResponseWriter, tmpl string, p *Punches) {
  t, _ := template.ParseFiles(tmpl + ".html")
  t.Execute(w, p)
}

func renderAddPunchPage(w http.ResponseWriter, tmpl string, p *Punches) {
  t, _ := template.ParseFiles(tmpl + ".html")
  t.Execute(w, p)
}

func renderEditUserPage(w http.ResponseWriter, tmpl string, p *neicacUser) {
  t, _ := template.ParseFiles(tmpl + ".html")
  t.Execute(w, p)
}

func renderEditAdminPage(w http.ResponseWriter, tmpl string, p *adminUser) {
  t, _ := template.ParseFiles(tmpl + ".html")
  t.Execute(w, p)
}
//END RENDER SECTION

//Created unique session ID to insure the user is in the current session
//When entering in data about logging in or logging out
func CreateSessionID() (string){
	source := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKMNOPQRSTUVWXYZ"
	rand.Seed(time.Now().UnixNano())
	s := ""
	for i := 0; i < 24; i++{
		s = s + string(source[rand.Intn(50)])
	}
	return s
}

//Checks if date is within specified range
func inTimeSpan(start, end, check time.Time) bool {
	return check.After(start) && check.Before(end)
}

//Login function uses post from index page to validate user login into their account
func LoginProcess(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Access-Control-Allow-Origin", "*")
  vars := mux.Vars(r)
  i, err := strconv.Atoi(vars["Value"])
  if err != nil {
    http.Redirect(w, r, "http://rebirtharmitage.com:8084/failedLogin", 301)
  }
  session, err := store.Get(r, "neicac_punchcard")
    if err != nil {
      http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
  a := findUser(i)
  if a == true {
    b := checkStatus(i)
    session.Values["session"] = "test"
    session.Values["userID"] = vars["Value"]
    session.Values["status"] = b.Status
    session.Values["authenticated"] = true
    session.Values["username"] = b.Username
    session.Values["pin"] = b.Pin
    session.Values["fname"] = b.Fname
    session.Values["lname"] = b.Lname
    session.Save(r, w)
    if b.Status == 0 {
      http.Redirect(w, r, "http://rebirtharmitage.com:8084/punchIn", 302)
    }else{
      http.Redirect(w, r, "http://rebirtharmitage.com:8084/punchOut", 302)
    }
  }else{
    http.Redirect(w, r, "http://rebirtharmitage.com:8084/failedLogin", 301)
  }
}

//Login to the admin system
func adminLoginProcess(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Access-Control-Allow-Origin", "*")
  vars := mux.Vars(r)
  s := strings.Split(vars["Value"], "`")
  uname, pass := s[0], s[1]
  session, err := store.Get(r, "neicac_punchcard")
    if err != nil {
      http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
  a, b := adminLoginTest(uname, pass)
  fmt.Println(a)
  if b {
    session.Values["session"] = "admin"
    session.Values["userID"] = uname
    session.Values["department"] = a.Department
    session.Values["authenticated"] = true
    session.Save(r, w)
    http.Redirect(w, r, "http://rebirtharmitage.com:8084/admin", 302)
  } else {
    session.Options.MaxAge = -1
    session.Save(r, w)
    http.Redirect(w, r, "http://rebirtharmitage.com:8084/failedLogin", 301)
  }
}

//Login to the admin system
func createAdminProcess(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Access-Control-Allow-Origin", "*")
  var m newAdmin
  b, _ := ioutil.ReadAll(r.Body)
  json.Unmarshal(b, &m)
  fmt.Println(m.Username)
  if m.Password != m.PasswordConfirm {
    fmt.Fprint(w, "Passwords did not match")
  }else{
    c, _ := adminLoginTest(m.Username, m.Password)
    if c.Username == ""{
      p := createAdmin(m)
      if p {
        fmt.Fprint(w, "User Account Created")
      }else{
        fmt.Fprint(w, "User already Exists")
      }
    }else{
      fmt.Fprint(w, "User already Exists")
    }
    
  }
}

//Login to the admin system
func createUserProcess(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Access-Control-Allow-Origin", "*")
  var m newUser
  b, _ := ioutil.ReadAll(r.Body)
  json.Unmarshal(b, &m)
  fmt.Println(m.Pin)
  i, err := strconv.Atoi(m.Pin)
  if err != nil {
    fmt.Fprint(w, "Pin was Invalid")
  }else{
    c := findUser(i)
    if c != true{
      p := createUser(m)
      if p {
        fmt.Fprint(w, "User Account Created")
      }else{
        fmt.Fprint(w, "Pin already Exists")
      }
    }else{
      fmt.Fprint(w, "Pin already Exists")
   }
  }
}

//Login to the admin system
func logoffProcess(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Access-Control-Allow-Origin", "*")
  session, err := store.Get(r, "neicac_punchcard")
    if err != nil {
      http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
  session.Options.MaxAge = -1
  session.Save(r, w)
  http.Redirect(w, r, "http://rebirtharmitage.com:8084", 301)
}

//Punch In Process page used to enter in punch in and change client status in system
func punchInProcess(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Access-Control-Allow-Origin", "*")
  session, _ := store.Get(r, "neicac_punchcard")
  userID := session.Values["userID"].(string)
  if session.Values["authenticated"].(bool) {
    if session.Values["status"] == 0 {
      punchInUser(userID)
    }
  } 
  session.Options.MaxAge = -1
  session.Save(r, w)
  http.Redirect(w, r, "http://rebirtharmitage.com:8084", 302)
}

//Punch Out Process page used to enter punch out, create total time since last punch and client status in system
func punchOutProcess(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Access-Control-Allow-Origin", "*")
  session, _ := store.Get(r, "neicac_punchcard")
  userID := session.Values["userID"].(string)
  if session.Values["authenticated"].(bool) {
    if session.Values["status"] == 1 {
      punchOutUser(userID)
    }
  }
  session.Options.MaxAge = -1
  session.Save(r, w)
  http.Redirect(w, r, "http://rebirtharmitage.com:8084", 302)
}

func masterPunchIn(w http.ResponseWriter, r *http.Request){
  vars := mux.Vars(r)
  punchInUser(vars["Value"])
  http.Redirect(w, r, "http://rebirtharmitage.com:8084/listUsers", 302)
}

func masterPunchOut(w http.ResponseWriter, r *http.Request){
  vars := mux.Vars(r)
  punchOutUser(vars["Value"])
  http.Redirect(w, r, "http://rebirtharmitage.com:8084/listUsers", 302)
}

func deletePunch(w http.ResponseWriter, r *http.Request){
  vars := mux.Vars(r)
  a := vars["Value"]
  fmt.Println(a)
  deletePunchMongo(a)
  http.Redirect(w, r, "http://rebirtharmitage.com:8084/admin", 302)
}

func editPunchProcess(w http.ResponseWriter, r *http.Request){
  w.Header().Set("Access-Control-Allow-Origin", "*")
  var p EPunch
  b, _ := ioutil.ReadAll(r.Body)
  json.Unmarshal(b, &p)
  if p.Pin == 0 {
    fmt.Fprint(w, "Punch was not edited.")
  }else{
    layout := "2006-01-02T15:04"
    t, err := time.Parse(layout, p.FormattedPunch)
    
    if err != nil {
        fmt.Println(err)
    }
    t = t.Add(5*time.Hour)
    updatePunch(p, t)
    fmt.Fprint(w, "Punch was edited.")
  }
}

func addPunchProcess(w http.ResponseWriter, r *http.Request){
  w.Header().Set("Access-Control-Allow-Origin", "*")
  var p EPunch
  var pe Punches
  b, _ := ioutil.ReadAll(r.Body)
  json.Unmarshal(b, &p)
  if p.Pin == 0 {
    fmt.Fprint(w, "Punch was not added.")
  }else{
    layout := "2006-01-02T15:04"
    t, err := time.Parse(layout, p.FormattedPunch)
    
    if err != nil {
        fmt.Println(err)
    }
    t = t.Add(5*time.Hour)
    pe.Punch = t
    pe.Pin = p.Pin
    pe.Punchtype = p.Punchtype
    addPunch(pe)
    fmt.Fprint(w, "Punch was added.")
  }
}

func addPunchViewHandler(w http.ResponseWriter, r *http.Request){
  w.Header().Set("Access-Control-Allow-Origin", "*")
  vars := mux.Vars(r)
  i, err := strconv.Atoi(vars["Value"])
  if err != nil {
    http.Redirect(w, r, "http://rebirtharmitage.com:8084/generalError", 302)
  }else{
    m := Punches{}
    m.Pin = i
    p, _ := loadAddPunchPage(m)
    renderEditPunchPage(w, "./html/addPunchPage", p)
  }
}


func editUserPageViewHandler(w http.ResponseWriter, r *http.Request){
  w.Header().Set("Access-Control-Allow-Origin", "*")
  vars := mux.Vars(r)
  i, err := strconv.Atoi(vars["Value"])
  if err != nil {
    http.Redirect(w, r, "http://rebirtharmitage.com:8084/generalError", 302)
  }else{
    m := checkStatus(i)
    p, _ := loadEditUserPage(m)
    renderEditUserPage(w, "./html/editUserPage", p)
  }
}

func editUserProcess(w http.ResponseWriter, r *http.Request){
  w.Header().Set("Access-Control-Allow-Origin", "*")
  var p newUser
  var a neicacUser
  b, _ := ioutil.ReadAll(r.Body)
  json.Unmarshal(b, &p)
  i, err := strconv.Atoi(p.Pin)
  if err != nil {
    fmt.Fprint(w, "Pin was Invalid")
  }else{
    c := findUser(i)
    if c {
        a.Pin = i
        a.Fname = p.Firstname
        a.Lname = p.Lastname
        a.Department = p.Department
        updateUser(a)
        fmt.Fprint(w, "User account has been edited.")
      }else{
        fmt.Fprint(w, "Pin already Exists")
     }
  }
}

func editAdminPageViewHandler(w http.ResponseWriter, r *http.Request){
  w.Header().Set("Access-Control-Allow-Origin", "*")
  vars := mux.Vars(r)
  a := vars["Value"]
  m, n := getAdmin(a)
  if n {
      p, _ := loadEditAdminPage(m)
      renderEditAdminPage(w, "./html/editAdminPage", p)
  }else{
    http.Redirect(w, r, "http://rebirtharmitage.com:8084/generalError", 302)
  }
}

//Edit Existing Admin User
func editAdminProcess(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Access-Control-Allow-Origin", "*")
  var m newAdmin
  b, _ := ioutil.ReadAll(r.Body)
  json.Unmarshal(b, &m)
  if m.Password != m.PasswordConfirm {
    fmt.Fprint(w, "Passwords did not match")
  }else{
    c, _ := adminLoginTest(m.Username, m.Password)
    if c.Username != ""{
      updateAdmin(m)
      fmt.Fprint(w, "User Account Edited")
    }else{
      fmt.Fprint(w, "User already Exists")
    }
    
  }
}


/*
	Start of ROUTER Section
*/
var router = mux.NewRouter()

func main() {
  router.HandleFunc("/", indexViewHandler)
  router.HandleFunc("/loginProcess/{Value}", LoginProcess)
  router.HandleFunc("/punchIn", punchInViewHandler)
  router.HandleFunc("/punchInProcess", punchInProcess)
  router.HandleFunc("/punchOut", punchOutViewHandler)
  router.HandleFunc("/punchOutProcess", punchOutProcess)
  router.HandleFunc("/failedLogin", failedViewHandler)
  router.HandleFunc("/timecard", timecardViewHandler)
  router.HandleFunc("/admin", adminViewHandler)
  router.HandleFunc("/createUser", createUserViewHandler)
  router.HandleFunc("/createUserProcess", createUserProcess).Methods("POST")
  router.HandleFunc("/createUserFailed", createUserFailedViewHandler)
  router.HandleFunc("/createUserSuccess", createUserSuccessViewHandler)
  router.HandleFunc("/listUsers", listUsersViewHandler)
  router.HandleFunc("/createAdmin", createAdminViewHandler)
  router.HandleFunc("/createAdminFailed", createAdminFailedViewHandler)
  router.HandleFunc("/createAdminSuccess", createAdminSuccessViewHandler)
  router.HandleFunc("/createAdminProcess", createAdminProcess).Methods("POST")
  router.HandleFunc("/adminLogin", adminLoginViewHandler)
  router.HandleFunc("/adminLoginProcess/{Value}", adminLoginProcess)
  router.HandleFunc("/listAdmins", listAdminsViewHandler)
  router.HandleFunc("/timecardEdit/{Value}", timecardEditViewHandler)
  router.HandleFunc("/masterPunchIn/{Value}", masterPunchIn)
  router.HandleFunc("/masterPunchOut/{Value}", masterPunchOut)
  router.HandleFunc("/deletePunch/{Value}", deletePunch)
  router.HandleFunc("/addPunchPage/{Value}", addPunchViewHandler)
  router.HandleFunc("/addPunchProcess", addPunchProcess).Methods("POST")
  router.HandleFunc("/editPunchPage/{Value}", editPunchPageViewHandler)
  router.HandleFunc("/editPunchProcess", editPunchProcess).Methods("POST")
  router.HandleFunc("/editUserPage/{Value}", editUserPageViewHandler)
  router.HandleFunc("/editUserProcess", editUserProcess).Methods("POST")
  router.HandleFunc("/editAdminPage/{Value}", editAdminPageViewHandler)
  router.HandleFunc("/editAdminProcess", editAdminProcess).Methods("POST")
  router.HandleFunc("/logout", logoffProcess)
	http.Handle("/css/",http.StripPrefix("/css/", http.FileServer(http.Dir("./css"))))
  http.Handle("/jQueryAssets/",http.StripPrefix("/jQueryAssets/", http.FileServer(http.Dir("./JQueryAssets"))))
	http.Handle("/fonts/",http.StripPrefix("/fonts/", http.FileServer(http.Dir("./fonts"))))
	http.Handle("/js/",http.StripPrefix("/js/", http.FileServer(http.Dir("./js"))))
	http.Handle("/vendor/",http.StripPrefix("/vendor/", http.FileServer(http.Dir("./vendor"))))
	http.Handle("/img/",http.StripPrefix("/img/", http.FileServer(http.Dir("./img"))))
  http.Handle("/", router)
	http.ListenAndServe(":8084", nil)
}
