sync:
	git fetch upstream
	git checkout master
	git merge upstream/master

gofancyimports:
	find . -name \*.go -print0 | xargs -0 gofancyimports fix -w -l github.com/paullesiak
