package main

import (
)



//--------------------------------------------------------------
//func (cs Courses ) getsummarylist() http.Handler {
//   return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//    type PopupItems struct {
//      Id                 string          `json:"id"`
//      CourseTitle        string          `json:"course-title"`
//      Stars              int             `json:"stars"`
//      Duration           int             `json:"duration"`
//      Overview           string          `json:"overview"`
//      SelfpacedPrice     int             `json:"selfpacedprice"`
//      PublicPrice        int             `json:"publicprice"`
//      Courseicon         string          `json:"courseicon"`
//    }
//    (w).Header().Set("Access-Control-Allow-Headers","*")
//    (w).Header().Set("Access-Control-Allow-Origin", "*")
//    popi := PopupItems{}
//    popis := []PopupItems{}
//    var js []byte
//    var err error
//    //Iterate over all courses, Copy Id, Name, Stars, Duration, Overview, Price, and Courseicon
//    for _, ThisCourse := range cs.Cc {
//       fmt.Println("--------------------------------------------------")
//       fmt.Printf("_Course PopUp_  = %s, %s, %s, %s, %s, %s, %s\n", ThisCourse.Id, ThisCourse.CourseTitle, ThisCourse.Testimonials[0].Stars, ThisCourse.Duration, ThisCourse.Overview, ThisCourse.Price, ThisCourse.Courseicon)
//       popi.Id=ThisCourse.Id
//       popi.CourseTitle=ThisCourse.CourseTitle
//       popi.Stars=ThisCourse.Testimonials[0].Stars
//       popi.Duration=ThisCourse.Duration.Hours
////       popi.Overview=ThisCourse.Overview
//       popi.SelfpacedPrice=ThisCourse.Price.Selfpaced.PriceTags[0].price
//       popi.SelfpacedPrice=ThisCourse.Price.Public.PriceTags[0].price
//       popi.Courseicon=ThisCourse.Courseicon
//       popis = append(popis,popi)
//    }
//    //If no courses match, SEND THEM ALL! 
//       js, err = json.Marshal(popis)
//    if err != nil {
//       http.Error(w, err.Error(), http.StatusInternalServerError)
//       fmt.Printf("Error %s:\n", err)
//       return
//    }
//    w.Header().Set("Content-Type", "application/json")
//    w.Write(js)
//    return
//    })
//}





//func (cs Courses)  Select(id string) (Courses, error) {
//     log.Printf("WORKING: Looking for %s\n", id)
//     var c Courses
//		 for _, ThisCourse := range cs.Cc  {
//          if ThisCourse.Id == id  {
//              c.Cc = append(c.Cc, ThisCourse)
//              fmt.Printf("FOUND %d Record, returning: %s\n" , len(c.Cc), c.Cc[0].Id)
//              return c, nil
//          }
//      }
//     return c, errors.New(fmt.Sprintf("Course ID \"%s\" does NOT exist\n", id ))
//}






//func (cs Courses)  Search(ls string) (Courses, error) {
//     ls = strings.ToLower(ls)
//     fmt.Println("--------------------------------------------------")
//     log.Printf("SEARCH FUNC REPORTING: Searching for %s\n", ls)
//     var c Courses
//     i := 0
//     hits := 0
//     totalhits := 0
//     for _, ThisCourse := range cs.Cc {
//        hits = strings.Count( strings.ToLower(fmt.Sprintf("%v", cs.Cc[i])), ls )
//        totalhits = totalhits + hits
//        if  hits > 0 {
//            c.Cc = append(c.Cc, ThisCourse)
//            fmt.Printf("%s Course has %d hits\n", ThisCourse.Id, hits )
//        }
//        i++
//      }
//     if  totalhits == 0 {
//       return c, errors.New(fmt.Sprintf("No course contains any information regarding \"%s\"" , ls ))
//     }
//    return c, nil
//}



