.PHONY: generate
generate:
	go run github.com/Khan/genqlient@48003b9627c3b2484701b331ee78eb46d7c53f84
	go build -o /tmp/fix-genqlient fix-genqlient.go
	/tmp/fix-genqlient internal/generated/graphql.go
	rm /tmp/fix-genqlient
