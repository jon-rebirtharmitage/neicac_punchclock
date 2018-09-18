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

func loadPage(title string) (*Page, error){
	return &Page{Title: title, Body: "blank..."}, nil
}

func loadTimecard(tp timecardPage) (*timecardPage, error){
	return &tp, nil
}

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
    tp.Punchcard[j].FormattedPunch = tp.Punchcard[j].Punch.Format("Mon Jan _2 3:05 PM 2006")
  }
  for k := range tp.AllPunchcard{
    for m := range tp.AllPunchcard[k].Punchf{
      tp.AllPunchcard[k].Punchf[m].FormattedPunch = tp.AllPunchcard[k].Punchf[m].Punch.Format("Mon Jan _2 3:05 PM 2006")
    }
  }
  p, _ := loadTimecard(tp)
  renderTimecard(w, "./html/timecard", p)
}

func filterPunchcardCurrent(tp timecardPage) []Punches{
  dp := []Punches{}
  fromDate := time.Now().AddDate(0, 0, -1)
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



//Renders all html pages from the system
func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
  t, _ := template.ParseFiles(tmpl + ".html")
  t.Execute(w, p)
}

func renderTimecard(w http.ResponseWriter, tmpl string, p *timecardPage) {
  t, _ := template.ParseFiles(tmpl + ".html")
  t.Execute(w, p)
}

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
      //http.Redirect(w, r, "http://rebirtharmitage.com:8084/createAdminSuccess", 302)
    }else{
      fmt.Fprint(w, "User already Exists")
      //http.Redirect(w, r, "http://rebirtharmitage.com:8084/createAdminFailed", 302)
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
//   router.HandleFunc("/createUser", createUserViewHandler)
//   router.HandleFunc("/createUserProcess", createUserProcess).Methods("POST")
  router.HandleFunc("/createAdmin", createAdminViewHandler)
  router.HandleFunc("/createAdminFailed", createAdminFailedViewHandler)
  router.HandleFunc("/createAdminSuccess", createAdminSuccessViewHandler)
  router.HandleFunc("/createAdminProcess", createAdminProcess).Methods("POST")
  router.HandleFunc("/adminLogin", adminLoginViewHandler)
  router.HandleFunc("/adminLoginProcess/{Value}", adminLoginProcess)
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
