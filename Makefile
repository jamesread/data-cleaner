default: service frontend

service:
	$(MAKE) -wC service default

protocol:
	$(MAKE) -wC proto buf

.PHONY: service protocol
