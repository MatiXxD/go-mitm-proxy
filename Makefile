KEY_NAME = cert.key
CERT_NAME = cert.crt

gen_key:
	openssl genpkey -algorithm RSA -out $(KEY_NAME) -pkeyopt rsa_keygen_bits:2048

gen_cert:
	openssl req -x509 -new -nodes -key $(KEY_NAME) -sha256 -days 365 -out $(CERT_NAME) -subj "/C=RU/ST=Moscow/L=Moscow/O=Solist/OU=Ilya/CN=RootCA"

gen: gen_key gen_cert
	@mv $(KEY_NAME) certs/
	@mv $(CERT_NAME) certs/

.PHONY: clean
clean:
	@rm certs/*
