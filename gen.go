package main

//go:generate goapi-gen --package=spec --out .\internal\api\spec\planner-golang.spec.go .\internal\api\spec\planner-golang.spec.json
//go:generate tern migrate --migrations .\internal\pgstore\migrations\ --config .\internal\pgstore\migrations\tern.conf
//go:generate sqlc generate -f .\internal\pgstore\sqlc.yaml
