# chrome-web-store

A small utility to upload a packaged chrome extension to the chrome web store.

## Setup

1. ````go get github.com/macropodhq/chrome-web-store````

1. Create a Google Service worker account at https://console.developers.google.com/

1. Fetch the JSON for the service

1. Create a zip file of your extension to upload

## Example

```
chrome-web-store -d dist.zip -c credentials.json -e user@example.com -a appId
```

The service worker account you created has it's own email adderss, and won't allow you to upload. You'll need to set the email address used to login to the chrome web store
