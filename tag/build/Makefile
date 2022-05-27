# including only std functions
build:
	go build -o app ./

# including std and plus functions
build_plus:
	go build -tags plus -o app ./

# including std, plus and premium functions
build_premium:
	go build -tags premium -o app ./

# including all(std, plus, premium, ultimate) functions
build_ultimate:
	go build -tags "plus premium" -o app ./

clean:
	rm app
