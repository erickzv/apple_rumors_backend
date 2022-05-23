# Apple Rumors
This API scrapes Apple related rumors

Web App is live at https://apple-rumors.vercel.app

You can find the API live at https://apple-rumors.herokuapp.com/all_news

This API uses GO and the soup library

## API structure
`/all_news`
```json
{
  "https://url.com": [
    {
      "title": "Apple Rumors",
      "href": "/apple/rumor/path"
    },
    {
      "title": "Apple Rumors",
      "href": "/apple/rumor/path"
    }
  ],
  "https://url.com": [
    {
      "title": "Apple Rumors",
      "href": "/apple/rumor/path"
    }
  ]
}
```
Now append the href to the url and now you have a complete url

## Setup
Start by cloning the repo, once that done all you need to do is `go run main.go`

For the docker container run `docker build .`

This project already has the necessary files to be deployed in GCloud & Heroku

Made with ❤️ & GO
