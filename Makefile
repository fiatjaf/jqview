dist: deploy/linux/jqview deploy/windows/jqview.exe
	mkdir -p dist
	cd deploy/linux && tar -czvf jqview_linux.tar.gz jqview
	cd deploy/windows && tar -czvf jqview_windows.tar.gz jqview.exe
	mv deploy/linux/jqview_linux.tar.gz dist/
	mv deploy/windows/jqview_windows.tar.gz dist/

deploy/linux/jqview: $(shell find . -name "*.go")
	qtdeploy -ldflags="-s -w" -fast build desktop github.com/fiatjaf/jqview

deploy/windows/jqview.exe: $(shell find . -name "*.go")
	qtdeploy -ldflags="-s -w" -docker build windows_64_static
