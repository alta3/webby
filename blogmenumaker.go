package main

import (
  "fmt"
)


func (b Blogs) blogmenumaker() BlogMenus  {
	  var existing         bool
	  var blogmenus        []BlogsByCategory
		var blogsbycategory  BlogsByCategory
    var blogs            Blogs
	  var categories       []string
		// Iterate over all blogs, and derive a list of unique categories
		for _, thisblog := range  b {
        //Iterate over array of categories
				existing = false
				for _,  thiscategory := range categories {
								if thiscategory == thisblog.Category {
							  existing = true
					      }
			  }
        if existing == false {
          categories = append (categories, thisblog.Category)
		    }
    }
    //At this stage, a list of unique categories has been gathered,
		//so build the blogmenu
    // Interate over each category 
		for _, thiscategory := range categories {
				fmt.Printf("\"%s\"\n",thiscategory)
				//Iternate over every blog for that category
        for _, thisblog := range b {
                if thisblog.Category == thiscategory {
							      blogs = append( blogs, thisblog)
										fmt.Printf("  - %s\n", thisblog.Title)
                }
         }
				 blogsbycategory.BlogCategory = thiscategory
				 blogsbycategory.Blogs = blogs
				 blogs = nil
				 blogmenus = append(blogmenus,blogsbycategory)
    }
		return blogmenus
}


