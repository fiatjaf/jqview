dist: deploy/linux/jqview deploy/windows/jqview.exe
	mkdir -p dist
	cd deploy/linux && tar -czvf jqview_linux.tar.gz jqview
	mv deploy/linux/jqview_linux.tar.gz dist/
	rm -f deploy/windows/jqview_windows.zip
	cd deploy/windows && zip jqview_windows *.exe
	mv deploy/windows/jqview_windows.zip dist/

deploy/linux/jqview: $(shell find . -name "*.go")
	qtdeploy -ldflags="-s -w" -fast build desktop github.com/fiatjaf/jqview

deploy/windows/jqview.exe: $(shell find . -name "*.go")
	qtdeploy -ldflags="-s -w" -docker build windows_64_static
