default: service frontend

service:
	cd service && make

protocol:
	buf generate
