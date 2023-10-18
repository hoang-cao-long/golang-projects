pre-commit:
	pre-commit run --all-files

fsm: finite-state-machine
	(cd ./finite-state-machine/ && make $(filter-out $@, $(MAKECMDGOALS)))

lgorm: learn-gorm
	(cd ./learn-gorm/ && make $(filter-out $@, $(MAKECMDGOALS)))

lgraphql: learn-graphql
	(cd ./learn-graphql/ && make $(filter-out $@, $(MAKECMDGOALS)))

lgrpc: learn-grpc
	(cd ./learn-grpc/ && make $(filter-out $@, $(MAKECMDGOALS)))
	
rss: rss-services
	(cd ./rss-services/ && make $(filter-out $@, $(MAKECMDGOALS)))