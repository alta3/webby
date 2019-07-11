{{ define "layout" }}
<!DOCTYPE html>
<html lang="en">

<head>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">
    <meta name="description" content="Alta3 a modern educational site template">
    <meta name="author" content="Ansonika">
    <title>Alta3 Research</title>

    <!-- Favicons-->
    <link rel="shortcut icon" href="/images/favicon.png" type="image/x-icon">
    <link rel="apple-touch-icon" type="image/x-icon" href="/img/apple-touch-icon-57x57-precomposed.png">
    <link rel="apple-touch-icon" type="image/x-icon" sizes="72x72" href="/img/apple-touch-icon-72x72-precomposed.png">
    <link rel="apple-touch-icon" type="image/x-icon" sizes="114x114" href="/img/apple-touch-icon-114x114-precomposed.png">
    <link rel="apple-touch-icon" type="image/x-icon" sizes="144x144" href="/img/apple-touch-icon-144x144-precomposed.png">

    <!-- BASE CSS -->
    <link href="/css/bootstrap.min.css" rel="stylesheet">
    <link href="/css/style.css" rel="stylesheet">
    <link href="/css/vendors.css" rel="stylesheet">
    <link href="/css/icon_fonts/css/all_icons.min.css" rel="stylesheet">
    <link href="/css/blog.css" rel="stylesheet">
    <!-- YOUR CUSTOM CSS -->
    <link href="/css/custom.css" rel="stylesheet">

</head>

<body>
         
     <div id="page">
         
     <header class="header menu_2">
         <div id="preloader"><div data-loader="circle-side"></div></div><!-- /Preload -->
         <div id="logo">
             <a href="/index.html"><img src="/images/Alta3-logo_large.png" width="300" height="100" data-retina="true" alt=""></a>
         </div>
         <ul id="top_menu">
             <li><a href="login.html" class="login">Login</a></li>
             <li><a href="#0" class="search-overlay-menu-btn">Search</a></li>
             <li class="hidden_tablet"><a href="/courses/courses-list.html" class="btn_1 rounded">Admission</a></li>
         </ul>
         <!-- /top_menu -->
         <a href="#menu" class="btn_mobile">
             <div class="hamburger hamburger--spin" id="hamburger">
                 <div class="hamburger-box">
                     <div class="hamburger-inner"></div>
                 </div>
             </div>
         </a>


         <nav id="menu" class="main-menu">
             <ul>



                 <li><span><a href="index.html">Home</a></span>
                 </li>
                 <li><span><a href="#0">Courses</a></span>
                     <ul>
                         <li><a href="/courses/courses-grid.html">Courses grid</a></li>
                         <li><a href="/courses/courses-grid-sidebar.html">Courses grid sidebar</a></li>
                         <li><a href="/courses/courses-list.html">Courses list</a></li>
                         <li><a href="/courses/courses-list-sidebar.html">Courses list sidebar</a></li>
                         <li><a href="/courses/course-detail.html">Course detail</a></li>
                         <li><a href="/courses/course-detail-2.html">Course detail working form</a></li>
                     </ul>
                 </li>
                 <li><span><a href="#0">Resources</a></span>
                     <ul>
                         <li><a href="about.html">About</a></li>
                         <li><a href="blog.html">Blog</a></li>
                         <li><a href="login.html">Login</a></li>
                         <li><a href="404.html">404 page</a></li>
                         <li><a href="faq.html">Faq</a></li>
                     </ul>
                 </li>
                 <li><span><a href="#0">Extra Pages</a></span>
                     <ul>
                         <li><a href="cart-1.html">Cart page 1</a></li>
                         <li><a href="cart-2.html">Cart page 2</a></li>
                         <li><a href="cart-3.html">Cart page 3</a></li>
                         <li><a href="pricing-tables.html">Responsive pricing tables</a></li>
                         <li><a href="/coming_soon/index.html">Coming soon</a></li>
                     </ul>
                 </li>

             </ul>
         </nav>

         <!-- Search Menu -->
         <div class="search-overlay-menu">
             <span class="search-overlay-close"><span class="closebt"><i class="ti-close"></i></span></span>
             <form role="search" id="searchform" method="get">
                 <input value="" name="q" type="search" placeholder="Search..." />
                 <button type="submit"><i class="icon_search"></i>
                 </button>
             </form>
         </div><!-- End Search Menu -->
     </header>
     <!-- /header -->



        
{{ template "body" }}     




     <footer>
         <div class="container margin_120_95">
             <div class="row">
                 <div class="col-lg-5 col-md-12 p-r-5">
                     <p><img src="/img/logo.png" width="149" height="42" data-retina="true" alt=""></p>
                     <p>Mea nibh meis philosophia eu. Duis legimus efficiantur ea sea. Id placerat tacimates definitionem sea, prima quidam vim no. Duo nobis persecuti cu. Nihil facilisi indoctum an vix, ut delectus expetendis vis.</p>
                     <div class="follow_us">
                         <ul>
                             <li>Follow us</li>
                             <li><a href="#0"><i class="ti-facebook"></i></a></li>
                             <li><a href="#0"><i class="ti-twitter-alt"></i></a></li>
                             <li><a href="#0"><i class="ti-google"></i></a></li>
                             <li><a href="#0"><i class="ti-pinterest"></i></a></li>
                             <li><a href="#0"><i class="ti-instagram"></i></a></li>
                         </ul>
                     </div>
                 </div>
                 <div class="col-lg-3 col-md-6 ml-lg-auto">
                     <h5>Useful links</h5>
                     <ul class="links">
                         <li><a href="#0">Admission</a></li>
                         <li><a href="#0">About</a></li>
                         <li><a href="#0">Login</a></li>
                         <li><a href="#0">Register</a></li>
                         <li><a href="#0">News &amp; Events</a></li>
                         <li><a href="#0">Contacts</a></li>
                     </ul>
                 </div>
                 <div class="col-lg-3 col-md-6">
                     <h5>Contact with Us</h5>
                     <ul class="contacts">
                         <li><a href="tel://17175664428"><i class="ti-mobile"></i> + 1 717 566 4428</a></li>
                         <li><a href="mailto:sales@alta3.com"><i class="ti-email"></i> sales@alta3.com</a></li>
                     </ul>
                     <div id="newsletter">
                     <h6>Newsletter</h6>
                     <div id="message-newsletter"></div>
                     <form method="post" action="/assets/newsletter.php" name="newsletter_form" id="newsletter_form">
                         <div class="form-group">
                             <input type="email" name="email_newsletter" id="email_newsletter" class="form-control" placeholder="Your email">
                             <input type="submit" value="Submit" id="submit-newsletter">
                         </div>
                     </form>
                     </div>
                 </div>
             </div>
             <!--/row-->
             <hr>
             <div class="row">
                 <div class="col-md-8">
                     <ul id="additional_links">
                         <li><a href="#0">Terms and conditions</a></li>
                         <li><a href="#0">Privacy</a></li>
                     </ul>
                 </div>
                 <div class="col-md-4">
                     <div id="copy">© 2017 Alta3</div>
                 </div>
             </div>
         </div>
     </footer>
     <!--/footer-->
     </div>
     <!-- page -->
     
     <!-- COMMON SCRIPTS -->



<script>
function myFunction() {
    var t = document.createTextNode("This is a paragraph.");

    newcourse = document.createElement("div");
    newcourse.className = "row no-gutters";

    plork = document.createElement("div");
    plork.className = "box_list wow";

    l3 = document.createElement("div");
    l3.className = "col-lg-5";

    l4 = document.createlement("figure");
    l4.classname = "block-reveal";

    r1 = document.createElement("div");
    r1.classname = "block-horizontal";

    l2 = document.createElement("img");
    l2.classname = "img-fluid";

    r2 = document.createElement("div");
    r2.className = "col-lg-7";

    r3 = document.createlement("figure");
    r3.classname = "wrapper";


    newcourse.appendChild("plork");
    document.getElementById("mycourselist").appendChild("newcourse");






}
</script>


<script>

var BLOCKS_PER_CHART = 4;
function addcourse() {
   chartContainer = "margin_60_35"
   listcourses(chartContainer)
}

function listcourses(chartContainer) {
  var container = document.createElement("div");
  var text = "Hello World!";
  var blockDiv, textSpan;  // used in the for loop

  container.className = "box_list wow";
  document.getElementById(chartContainer.replace("#", "")).appendChild(container);

  for(var i = 0; i < BLOCKS_PER_CHART; i++) {
    blockDiv = document.createElement("div");
    blockDiv.className = "block";
    textSpan = document.createElement("span");
    textSpan.append(text);  // see note about browser compatibility
    blockDiv.append(textSpan);
    container.append(blockDiv);
  }
}    

</script>
    <script src="/js/jquery-2.2.4.min.js"></script>
    <script src="/js/common_scripts.js"></script>
    <script src="/js/main.js"></script>
    <script src="/assets/validate.js"></script>
    <script src="/assets/courselist.js"></script>
     
</body>
</html>
{{ end }}
