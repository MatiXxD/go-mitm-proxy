DOCKER_COMPOSE_PATH = ./docker/docker-compose.yaml
PROJECT_NAME = mitmproxy

KEY_NAME = cert.key
CERT_NAME = cert.crt

# ==============================================================================
# DOCKER-COMPOSE
.PHONY: docker-compose-build
docker-compose-build:
	@docker compose -f $(DOCKER_COMPOSE_PATH) -p $(PROJECT_NAME) build

.PHONY: docker-compose-up
docker-compose-up:
	@docker compose -f $(DOCKER_COMPOSE_PATH) up -d

.PHONY: docker-compose-stop
docker-compose-stop:
	@docker compose -f $(DOCKER_COMPOSE_PATH) stop


# ==============================================================================
# CERT-GEN
.PHONY: gen_key
gen_key:
	openssl genpkey -algorithm RSA -out $(KEY_NAME) -pkeyopt rsa_keygen_bits:2048

.PHONY: gen_cert
gen_cert:
	openssl req -x509 -new -nodes -key $(KEY_NAME) -sha256 -days 365 -out $(CERT_NAME) -subj "/C=RU/ST=Moscow/L=Moscow/O=Solist/OU=Ilya/CN=RootCA"

# create certs folder before use
.PHONY: gen
gen: gen_key gen_cert
	@mv $(KEY_NAME) certs/
	@mv $(CERT_NAME) certs/

.PHONY: clean
clean:
	@rm certs/*
