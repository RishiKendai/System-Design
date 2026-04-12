package main

import (
	"fmt"

	"github.com/lld/url-shortner/url"
	"github.com/lld/url-shortner/user"
)

func main() {
	service := url.NewURLShortnerService(10)
	// user
	user1 := user.NewUser("John Doe", "john.doe@example.com")
	url1 := service.GenerateShortCode("https://www.google.com", user1.ID, "", nil)
	service.PrintUrl(url1.ShortCode)
	fmt.Println("Get URL: ", service.GetURL(url1.ShortCode))
	service.PrintUrl(url1.ShortCode)

	url2 := service.GenerateShortCode("https://www.google.com", user1.ID, "google", nil)
	service.PrintUrl(url2.ShortCode)
	fmt.Println("Get URL: ", service.GetURL(url2.ShortCode))
	fmt.Println("Get URL: ", service.GetURL(url2.ShortCode))
	fmt.Println("Get URL: ", service.GetURL(url2.ShortCode))
	service.PrintUrl(url2.ShortCode)

}
