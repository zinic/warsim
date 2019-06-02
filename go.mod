module github.com/zinic/warsim

go 1.12

require (
	github.com/BurntSushi/toml v0.3.1
	github.com/shurcooL/sanitized_anchor_name v1.0.0 // indirect
	golang.org/x/text v0.3.2
	gopkg.in/russross/blackfriday.v2 v2.0.1
)

replace gopkg.in/russross/blackfriday.v2 v2.0.1 => github.com/russross/blackfriday/v2 v2.0.1
