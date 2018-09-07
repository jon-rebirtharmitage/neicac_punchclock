package main

import(
	"html/template"
	"net/http"
  "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
  "github.com/gorilla/mux"
  "github.com/gorilla/sessions"
  "strconv"
	"fmt"
  "log"
  "time"
	"math/rand"
)

//GLOBAL VARIABLES
//Session Storage for secure session variable handling using AES-256 Encryption
var store = sessions.NewCookieStore([]byte("30efe06c8b1c6a0cdca9298090d551cd86547a90435d15fd90fb3765af766072"))


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
  Pin, Status int
  Current time.Time
  Punchcard []Punches
  AllPunchcard []Pcard
}

type neicacUser struct{
  Username  string
  Pin, Status int
}

type Pcard struct{
  Startdate time.Time
  Enddate time.Time
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

func punchInViewHandler(w http.ResponseWriter, r *http.Request) {
  p, _ := loadPage("Punch")
  renderTemplate(w, "./html/punchIn", p)
}

func punchOutViewHandler(w http.ResponseWriter, r *http.Request) {
  p, _ := loadPage("Punch")
  renderTemplate(w, "./html/punchOut", p)
}

func timecardViewHandler(w http.ResponseWriter, r *http.Request) {
  session, _ := store.Get(r, "neicac_punchcard")
  tp := timecardPage{}
  tp.Title = session.Values["username"].(string)
  tp.Body = session.Values["username"].(string)
  tp.Username = session.Values["username"].(string)
  tp.Pin = session.Values["pin"].(int)
  tp.Status = session.Values["status"].(int)
  tp.Current = time.Now()
  tp.Punchcard = getPunches(tp, 1)
  tp.Punchcard = filterPunchcardCurrent(tp)
  t := filterPunchcard(tp)
  for m := range t{
    fmt.Println(t[m])
  }
  for j := range tp.Punchcard{
    tp.Punchcard[j].FormattedPunch = tp.Punchcard[j].Punch.Format("Mon Jan _2 3:05 PM 2006")
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

//Goes through all punches to those that are within each range
func filterPunchcard(tp timecardPage) []Pcard {
  //Replace with global variable in database
  fromDate := time.Now().AddDate(0, 0, -4)
  //Replace above!!!
  pcard := []Pcard{}
  ctime := time.Now()
  dtime := time.Now()
  diff := time.Now().Sub(fromDate)
  c := int(diff.Hours())/24
  for j := 1; j < c; j++ {
    p := Pcard{}
    fmt.Println(-j)
    fmt.Println(-(j+1))
    ctime = ctime.AddDate(0, 0, -j)
    dtime = ctime.AddDate(0, 0, -(j+1))
    p.Startdate = dtime
    p.Enddate = ctime
    p.Title = "Test"
    for k := range tp.Punchcard{
      if inTimeSpan(dtime, ctime, tp.Punchcard[k].Punch){
        p.Punchf = append(p.Punchf, tp.Punchcard[k])
      }
    }
    pcard = append(pcard, p)
    fmt.Println(j)
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
  http.Redirect(w, r, "http://rebirtharmitage.com:8084", 302)
}

//Find the users information from the users table
func findUser(userPin int) bool{
  session, err := mgo.Dial("localhost:27017")
	if err != nil {
			panic(err)
	}
	defer session.Close()

	c := session.DB("neicac").C("test")

	result := neicacUser{}
	errorCheck := c.Find(bson.M{"pin":userPin}).One(&result)
	if errorCheck != nil{
    return false
	}
	return true
}

//Function access the current users pin to determine if they are signed in or out
//If they are in they get a sign out button and if out they are giving the sign in button
func checkStatus(userPin int) neicacUser{
  session, err := mgo.Dial("localhost:27017")
	if err != nil {
			panic(err)
	}
	defer session.Close()

	c := session.DB("neicac").C("test")

	result := neicacUser{}
	errorCheck := c.Find(bson.M{"pin":userPin}).One(&result)
	if errorCheck != nil{
    log.Fatal(errorCheck)
	}
	return result
}

//Function to enter a punch in for the user in the current session
func punchInUser(userID string){
  session, err := mgo.Dial("localhost:27017")
	if err != nil {
			panic(err)
	}
	defer session.Close()
  i, _ := strconv.Atoi(userID)
  
	c := session.DB("neicac").C("test")
  c.Update(bson.M{"pin": i}, bson.M{"$set": bson.M{"status": 1}})
  
  d := session.DB("neicac").C("timestamps")
  d.Insert(bson.M{"pin": i, "punch":time.Now(), "punchtype":0})
}

//Function to enter a punch in for the user in the current session
func punchOutUser(userID string){
  session, err := mgo.Dial("localhost:27017")
	if err != nil {
			panic(err)
	}
	defer session.Close()
  i, _ := strconv.Atoi(userID)

	c := session.DB("neicac").C("test")
  c.Update(bson.M{"pin": i}, bson.M{"$set": bson.M{"status": 0}})
  
  d := session.DB("neicac").C("timestamps")
  d.Insert(bson.M{"pin": i, "punch":time.Now(), "punchtype": 1})
}

//Get the users punches from the timestamp table
func getPunches(tp timecardPage, timecardType int) []Punches{
  session, err := mgo.Dial("localhost:27017")
	if err != nil {
			panic(err)
	}
	defer session.Close()

  c := session.DB("neicac").C("timestamps")
	result := []Punches{}
  iter := c.Find(bson.M{"pin":tp.Pin}).Limit(1000).Iter()
	err = iter.All(&result)
	if err != nil {
			log.Fatal(err)
	}
	return result
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
	http.Handle("/css/",http.StripPrefix("/css/", http.FileServer(http.Dir("./css"))))
	http.Handle("/fonts/",http.StripPrefix("/fonts/", http.FileServer(http.Dir("./fonts"))))
	http.Handle("/js/",http.StripPrefix("/js/", http.FileServer(http.Dir("./js"))))
	http.Handle("/vendor/",http.StripPrefix("/vendor/", http.FileServer(http.Dir("./vendor"))))
	http.Handle("/img/",http.StripPrefix("/img/", http.FileServer(http.Dir("./img"))))
  http.Handle("/", router)
	http.ListenAndServe(":8084", nil)
}