## Unofficial client for ntfy

Homepage: https://ntfy.sh

Use: 

```go

    import "github.com/mlplabs/common-go-pkg/pkg/clients/ntfy"
....

	var (
		yourAppName = "my-app"
		ntfyBaseUrl = "https://ntfy-base-url"
		ntfyToken   = "access_token_from_ntfy_ui"
	)
    
    // create client
	ntfyClient := ntfy.NewHttpClient(
		yourAppName,
		&ntfy.Config{
			BaseURL: ntfyBaseUrl,
			Token: ntfyToken},
		)

	// required
	msg := ntfy.Message{
		Topic: "my-topic",
		Text:  "Message text",
	}

	// optional
	msg.SetPriority(ntfy.HighPriority).
		SetTitle("Message with high priority").
		SetTags("tag1, tag2")
    msg.SetClick("https://target") 
	
	// send message
	ntfyClient.SendMessage(ctx, &msg)



```