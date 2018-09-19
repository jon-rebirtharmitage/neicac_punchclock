package main

import(
  "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
  "strconv"
  "time"
  "log"
)


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

//Get the users punches from the timestamp table
func getUsers(d string) []neicacUser{
  session, err := mgo.Dial("localhost:27017")
	if err != nil {
			panic(err)
	}
	defer session.Close()

  c := session.DB("neicac").C("test")
	result := []neicacUser{}
  if d == "all" {
    iter := c.Find(bson.M{}).Limit(1000).Iter()
    err = iter.All(&result)
    if err != nil {
      log.Fatal(err)
    }
  }else{
    iter := c.Find(bson.M{"department":d}).Limit(1000).Iter()
    err = iter.All(&result)
    if err != nil {
      log.Fatal(err)
    }
  }

	return result
}

//Get the users punches from the timestamp table
func getAdmins(d string) []adminUser{
  session, err := mgo.Dial("localhost:27017")
	if err != nil {
			panic(err)
	}
	defer session.Close()

  c := session.DB("neicac").C("admin")
	result := []adminUser{}
  if d == "all" {
    iter := c.Find(bson.M{}).Limit(1000).Iter()
    err = iter.All(&result)
    if err != nil {
      log.Fatal(err)
    }
  }else{
    iter := c.Find(bson.M{"department":d}).Limit(1000).Iter()
    err = iter.All(&result)
    if err != nil {
      log.Fatal(err)
    }
  }

	return result
}

//Get the users punches from the timestamp table
func adminLoginTest(uname string, pass string) (adminUser, bool){
  session, err := mgo.Dial("localhost:27017")
	if err != nil {
			panic(err)
	}
	defer session.Close()
  key := []byte("0b6c92070071d82753b6f747")
  c := session.DB("neicac").C("admin")
	result := adminUser{}
  errorCheck := c.Find(bson.M{"username":uname}).One(&result)
	if errorCheck != nil{
    return result, false
	}
  if decrypt(key, result.Password) == pass{
    return result, true
  }
  return result, false
}

//Get the users punches from the timestamp table
func createAdmin(n newAdmin) bool {
  session, err := mgo.Dial("localhost:27017")
	if err != nil {
			panic(err)
	}
	defer session.Close()
  key := []byte("0b6c92070071d82753b6f747")
  c := session.DB("neicac").C("admin")
  m := adminUser{}
  m.Username = n.Username
  m.Fname = n.Firstname
  m.Lname = n.Lastname
  m.Password = encrypt(key, n.Password)
  m.Department = n.Department
  m.Role = "Admin"
  m.Auth = false
  c.Insert(bson.M{"username":m.Username,"fname":m.Fname,"lname":m.Lname,"password":m.Password,"department":m.Department,"role":m.Role,"auth":m.Auth})
  return true
}

//Takes the user information and creates a new one and inserts into the user database
func createUser(n newUser) bool {
  session, err := mgo.Dial("localhost:27017")
	if err != nil {
			panic(err)
	}
	defer session.Close()
  c := session.DB("neicac").C("test")
  m := neicacUser{}
  i, _ := strconv.Atoi(n.Pin)
  m.Fname = n.Firstname
  m.Lname = n.Lastname
  m.Pin = i
  m.Department = n.Department
  m.Status = 0
  c.Insert(bson.M{"fname":m.Fname,"lname":m.Lname,"pin":m.Pin,"department":m.Department,"status":m.Status})
  return true
}