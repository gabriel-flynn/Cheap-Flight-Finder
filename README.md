# Cheap-Flight-Finder

![cheap flight system design](./Cheap-Flight-Finder.png)

## Problem:
- I love to travel with friends but it can be exhausting trying to figure out whats a good deal (especially after factoring in baggage and other fees charged), finding flights that fit a certain criteria, etc. and I wanted an automated tool to take care of all that so all I need to do is wait to be notified by the tool and then decide if I want to book that flight or not. (And it's also a good excuse to learn more Golang and dive deeper into writing concurrent programs with it)

## Goals:
- This goal for this project is to be able to find cheap flights and have the ability to be really specific with the parameters I'm looking for. For example: I want to be notified of cheap flights to Vegas where I'd arrive Friday by X time then arrive back home by Sunday at Y time. 
- Make the Flight-Finder highly concurrent and fast (while obviously rate-limiting each specific API/website we're scraping to not overload the servers)


## Todo:
- [ ] Integrate Flight Providers:
    - [x] [Spirit Airlines](https://github.com/gabriel-flynn/Cheap-Flight-Finder/issues/1)
    - [ ] [Southwest Airlines](https://github.com/gabriel-flynn/Cheap-Flight-Finder/issues/3)
    - [ ] [Frontier Airlines](https://github.com/gabriel-flynn/Cheap-Flight-Finder/issues/4)
    - [ ] [JetBlue](https://github.com/gabriel-flynn/Cheap-Flight-Finder/issues/5)
- [] [Support round trip flights](https://github.com/gabriel-flynn/Cheap-Flight-Finder/issues/8)
- [ ] [Upload flight data to S3](https://github.com/gabriel-flynn/Cheap-Flight-Finder/issues/2)
- [x] [DynamoDB support](https://github.com/gabriel-flynn/Cheap-Flight-Finder/issues/6)
- [ ] Setup the python scraper with gRPC


## Possible Future Improvements:
- Turn the Cheap-Flight-Finder into a distributed application, it would be pretty similar to MapReduce -> Our Coordinator service would assign links to scrape (the Map step) to our worker nodes and our Reduce step would be where we run PySpark to find the cheapest flights
    - I think it would be really cool (and super cheap) to try to use EC2 spot instances to run this in a distributed manner (assuming AWS IPs aren't blocklisted by the different airlines/flight providers)

## Concerns:
- The data/flight prices we get from the flight providers might not be the most accurate depending on if the airlines use any kind of algorithm to determine a specific price per visitor. (Ex: Airline X sees my IP checking flights out pretty frequently so they assume I'm very likely to purchase a flight soon so they increase the flight price that I see by 15% since they think I'm very likely to buy it regardless of that price increase)