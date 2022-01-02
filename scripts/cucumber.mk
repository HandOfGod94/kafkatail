cucumber-setup:
	ruby --version
	bundle install
	bundle binstub cucumber --path bin

cucumber-test:
	bin/cucumber
