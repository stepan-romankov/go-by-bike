###### Prerequisites:
* Docker
* Docker Compose

###### Notes:
* APP and TESTS uses separate databases
  
###### How to start:
* All in one: `docker-compose -f stack.yml up --build`
* Only tests: `docker-compose -f stack.yml up tests`
* Only app: `docker-compose -f stack.yml up app`

###### How to see logs:
* For tests: docker-compose -f stack.yml logs -f tests
* For app: docker-compose -f stack.yml logs -f app