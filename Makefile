exe:
	rm -f ~/bin/gopsx
	rm -f `go env GOPATH`/bin/gopsx
	hash -r
	unalias -a gopsx
	go install

linked:
	rm -f ~/bin/gopsx.go
	rm -f `go env GOPATH`/bin/gopsx
	hash -r
	# cp gopsx.go ~/bin/gopsx
	ln -s `pwd`/gopsx.go ~/bin/gopsx.go
	chmod +x ~/bin/gopsx.go
	echo "alias gopsx=gopsx.go" > aliases && . aliases && rm aliases

script:
	rm -f ~/bin/gopsx.go
	rm -f `go env GOPATH`/bin/gopsx
	hash -r
	cp gopsx.go ~/bin/gopsx.go
	echo "alias gopsx=gopsx.go" > aliases && . aliases && rm aliases

run:
	go run gopsx.go
