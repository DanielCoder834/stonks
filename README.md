# A Website that Presents Stock Data

Currently, the project is stil being developed as I need to set-up a better way to store and clean out older data in the database. 

## The Database 
The database is wrriten in Go-lang and acts like Redis database. Currently there are different methods in which you can store data from single hashmap layer, to BTrees to a multiple layer hashmap. I need a good way to delete older data as I want to prioritize utilizing my memory more effiecently. The majority of the main changes happen in the command handler file. 

## The API Caller 
Currently, the api caller is written in Python as many finiacial apis have better support for Python. The api currently uses Alpha Vantage and I hope to make the caller multi-threaded as I can start to use other apis such as Yahoo Finiance. The caller sends data to the database through its own Redis client, and gets data from the api every hour through a scheduler. 

## The Website 
The website currently has the ground-work ready but I want to focus on the backend logic before making the website look good. The website uses React as I hope to use libraries such as React-Router, and Charts.js. The website also has a Redis client with a scheduler in which it puts that data into a json file for Charts.js to utilize. 

