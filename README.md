# webby


Return a list of all blog.id: 
`http://ultracoolness.com:56644/api/v1/blog/search/*`   (Just add an asterisk)
Example:
http://ultracoolness.com:56644/api/v1/blog/search/*`
Returns:  

```
["Law","BulletProofCode","Max"]  
```

Search for a specific blog:  
http://ultracoolness.com:56644/api/v1/blog/search/<search-string> <-- case insensitive
Example:  
`http://ultracoolness.com:56644/api/v1/blog/search/cat`  
Returns:

```
["Max"]
```
 
Return a specific blog detail using blog.id:  
`http://ultracoolness.com:56644/api/v1/blog/id/<blog.id>` <-- case insensitive  
Example: 
`http://ultracoolness.com:56644/api/v1/blog/id/MaX`

```
{
  "Id": "Max",
  "Title": "Office Cat",
  "Date": "March 22",
  "Author": "Kat Tastrophy",
  "Content": "\u003cp\u003eThe Office Cat\u003c/p\u003e\n\n\u003cp\u003eMax the cat\nthe office cat\nhe\u0026rsquo;s very fat.\nWE all know that\nmostly dog, but really fat\nexpects us all his head to pat\nHe eats his food out of a vat\nthen goes outside to eat a rat\u003c/p\u003e\n\n\u003cp\u003eMax the cat\nthe office cat\nwears many hats.\nthe biggest cat I\u0026rsquo;ve ever seen\nsmiles at me, awkwardly\nMax likes to meow and to purr\nMax likes to growl\nand walk like a sow.\u003c/p\u003e\n\n\u003cp\u003eMax the cat\nthe office cat\nwears many hats.\nusurps my desk\nhijacks the carpet\ntakes over the keyboard\nhow is there an argument\u003c/p\u003e\n\n\u003cp\u003eMax the cat\nthe office cat\nwears many hats\nand that in fact\nevery day I walk into the office\nIf I\u0026rsquo;m grumpy,upset,angry or sad at\nsomething tells me that day ameliorated Max\u003c/p\u003e\n\n\u003cp\u003eThe office cat known as Max\nMighty Max at that\nwill not be denied or held back\nas Max paws open the sliding door\nonce upstairs and gaining tracks\nOh, he is there to put back Max\u003c/p\u003e\n\n\u003cp\u003esay goodbye to the office cat\nrare sighting seen, but Max will be back\n\\to work, which \u0026ldquo;what is work?\u0026rdquo;\nEverything is done, where is the perk?\u003c/p\u003e\n\n\u003cp\u003eThe Office Cat\u003c/p\u003e\n"
}
```
