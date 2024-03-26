# Capstone - Cocktails Recipes
A database recipes of drinks and cocktails.<br />
It includes online resources like cocktail images, thumbnails, and videos.

# System Features
- Runs a Rest API over a http.Server instance with chi.Mux support.
- The API configuration is set by system environment variables with the `CAPSTONE` word prefixed. If configuration is missed, set up the default one.
- Builds, launch, and test operations can be done by make statements. Moreover, it is useful for CI/CD implementations.
- Generates data files with proper permissions on the fly.
- All the data files using on test cases are created on the fly, as well as removed.
- Full logging support through the zerolog library.
- Clean Architecture implementation:
  - Entity
  - Controller
  - Service
  - Repository
- Patterns:
  - Singleton pattern
  - Dependency Injection
- Handling Errors:
  - Error propagation.
  - Custom errors.
  - The controller error handler determines the HTTP error code, status and full error description.

# How it works
You can retrieve all the cocktail recipes, a specific one by id and a filtered list.
- Filter name and filter values are case insensitive.
- Any wrong csv record is discarded.
- Update database from a public API.

# Cocktail recipes
You can get all the cocktail recipes or a filtered list of them.

To retrieve all the recipes:
```
http://localhost:8080/api/v0/cocktails
```
### Filtering recipes
You can get a filtered list of cocktail recipes. The following are the supported filters: 

Filtering recipes by ``ID``. Just numbers are allowed.
```
http://localhost:8080/api/v0/cocktail/id/1
```

Filtering recipes by ``NAME``. Possible values: adam, acapulco, after sex, americano, affair, acid, afternoon, etc.
```
http://localhost:8080/api/v0/cocktail/name/afterglow
```

Filtering recipes by ``ALCOHOLIC``. Possible values: alcoholic, non alcoholic, etc.
```
http://localhost:8080/api/v0/cocktail/alcoholic/alcoholic
```

Filtering recipes by ``CATEGORY``. Possible values: cocktail, shot, drink, coffee, tea, etc.
```
http://localhost:8080/api/v0/cocktail/category/ordinary
```

Filtering recipes by ``INGREDIENT``. Possible values: gin, amaretto, soda, vodka, applejack, scotch, sugar, lemonade etc. 
```
http://localhost:8080/api/v0/cocktail/ingredient/vodka
```

Filtering recipes by ``GLASS``. Possible values: cocktail, collins shot, martini, highball, old-fashioned, etc.
```
http://localhost:8080/api/v0/cocktail/glass/martini
```
### Retrieving Odd and Even recipes concurrently
We can get a filtered list of cocktail recipes by <b>even</b> and <b>odd</b> records based on the `ID` property from the database concurrently 
using the Worker-Pool pattern with following request format:
```
localhost:8080/api/v0/cocktails/{type}/{items}/{items-per-worker}"
```
Parameters description:
> <b>type:</b> is the number type. Only support "odd" or "even" </br>
<b>items:</b> is the amount of valid records to be processed. </br>
<b>items per worker:</b> is the amount of jobs that each worker performs. </br></br>
> *Note: the number of items per worker must not be higher than the items.*

To retrieve `Even` records:
```
localhost:8080/api/v0/cocktails/even/18/6
```

To retrieve `Odd` records:
```
localhost:8080/api/v0/cocktails/odd/18/6
```

# Administrative Tasks:
To update the database from the public API:
```
http://localhost:8080/api/v0/cocktail/updatedb
```

# Run
The API listens by default on the port `8080`. To run the API, execute the following command:
```
make run
```

# Testing
The controller, service, and repository layers implement unit tests with mock support.
All the data files required on the unit tests are created on the fly, as well as removed.<br />

To run tests of each layer of the clean architecture, use the following commands.
```
make test-controller
make test-service
make test-repository
```
Optional. Run all test layers of the architecture.
```
make test-clean-architecture     
```
Run all tests of the system.
```
make test-all     
```

