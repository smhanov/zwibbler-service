
VERSION := 1.0
NAME := zwibbler_$(VERSION)

RPM_NAME := zwibbler-$(VERSION)-1.x86_64.rpm

zwibbler: *.go
	go build

$(NAME)_amd64.deb: zwibbler
	-rm $(NAME)_amd64.deb
	fpm -t deb -n zwibbler \
	--after-install after-install.sh \
	--before-remove before-remove.sh \
	-s dir \
	./zwibbler=/usr/bin/zwibbler ./zwibbler.conf=/etc/zwibbler.conf ./empty.db=/var/lib/zwibbler/zwibbler.db

$(RPM_NAME): zwibbler
	-rm $(RPM_NAME)
	fpm -t rpm -n zwibbler \
	--after-install after-install.sh \
	--before-remove before-remove.sh \
	-s dir \
	./zwibbler=/usr/bin/zwibbler ./zwibbler.conf=/etc/zwibbler.conf ./empty.db=/var/lib/zwibbler/zwibbler.db


deb: $(NAME)_amd64.deb
rpm: $(RPM_NAME)

