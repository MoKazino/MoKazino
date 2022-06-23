build_ui:
	rm -rf .tmp_build_ui || true
	git clone --depth=1 --branch=master https://github.com/MoKazino/ui .tmp_build_ui
	cd .tmp_build_ui && make build
build_go: 
	echo > ./.tmp_build_ui/build/imports.txt
	sh -c 'cd .tmp_build_ui/build && for i in $(shell cd .tmp_build_ui/build && find . -type f); do echo "//go:embed $$i" >> imports.txt; done'
	rm -rf .tmp_build || true
	mkdir .tmp_build
	cp vendor web main.go go.mod go.sum .tmp_build -r
	cp .tmp_build_ui/build/* .tmp_build/web -r
	sed 's/var files embed.FS//' -i .tmp_build/web/embed_files.go
	cat .tmp_build/web/imports.txt >> .tmp_build/web/embed_files.go
	echo "var files embed.FS" >> .tmp_build/web/embed_files.go
	sed 's/\.\///g' -i .tmp_build/web/embed_files.go
	cd .tmp_build && go build -v
	mv .tmp_build/xmrdice .

build: clean build_ui build_go

clean:
	rm -rf .tmp_build
	rm -rf .tmp_build_ui

