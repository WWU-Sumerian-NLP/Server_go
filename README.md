# Server_go

This is our server that handles interactions between our dashboard, NLP libaries and other APIs.

## Description

Our server handles the relationship between the dashboard and our backend.  It imports our NLP libraries to use for each task and stores data in a relational database to handle the state of data between multiple services. </br>
![server_endpoints](https://github.com/WWU-Sumerian-NLP/images/blob/master/server_endpoints.png) </br>

This picture above shows the existing endpoints of our server. </br>

Below, shows the function for when the run entity extraction endpoint is called. We import our CDLI_Extractor module from our NLP library to perform this task. </br>
![running_libraries](https://github.com/WWU-Sumerian-NLP/images/blob/master/running_libraries_server.png) </br>

Finally, this shows the schema for our relational database </br>
![db_schema](https://github.com/WWU-Sumerian-NLP/images/blob/master/server_schema.png) </br>

## Authors
Hansel Guzman-Soto (https://www.linkedin.com/in/hansel-guzman-soto/)

