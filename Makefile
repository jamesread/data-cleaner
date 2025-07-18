default: service frontend

service:
	$(MAKE) -wC service default

proto: protocol

protocol:
	$(MAKE) -wC proto buf

.PHONY: service protocol
