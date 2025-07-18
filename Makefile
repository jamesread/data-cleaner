default: service frontend

service:
	$(MAKE) -wC service default

proto: protocol

protocol:
	$(MAKE) -wC proto buf

frontend:
	$(MAKE) -wC frontend

.PHONY: service protocol frontend
