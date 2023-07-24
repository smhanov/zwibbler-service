
VERSION := 10.0
NAME := zwibbler_$(VERSION)

RPM_NAME := zwibbler-$(VERSION)-1.x86_64.rpm

zwibbler: *.go
	go build -tags "$(TAGS)"

hae: *.go zwibserve/*.go
	make TAGS=hae

$(NAME)_amd64.deb: zwibbler *.sh
	-rm $(NAME)_amd64.deb
	fpm -t deb -n zwibbler \
	-v $(VERSION) \
	--after-install after-install.sh \
	--before-remove before-remove.sh \
	-s dir \
	./zwibbler=/usr/bin/zwibbler ./zwibbler.conf=/etc/zwibbler.conf

$(RPM_NAME): zwibbler
	-rm $(RPM_NAME)
	fpm -t rpm -n zwibbler \
	-v $(VERSION) \
	--after-install after-install.sh \
	--before-remove before-remove.sh \
	-s dir \
	./zwibbler=/usr/bin/zwibbler ./zwibbler.conf=/etc/zwibbler.conf

deb: $(NAME)_amd64.deb
rpm: $(RPM_NAME)

clean:
	-rm $(NAME)_amd64.deb $(RPM_NAME) zwibbler
